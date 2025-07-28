#!/bin/bash


set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Health score
HEALTH_SCORE=0
MAX_SCORE=100

# yassified output idc i like the emojis
print_status() {
    local status=$1
    local message=$2
    local points=$3
    
    if [ "$status" = "PASS" ]; then
        echo -e "${GREEN}✅ PASS${NC} $message ${CYAN}(+$points pts)${NC}"
        HEALTH_SCORE=$((HEALTH_SCORE + points))
    elif [ "$status" = "FAIL" ]; then
        echo -e "${RED}❌ FAIL${NC} $message ${CYAN}(0 pts)${NC}"
    elif [ "$status" = "WARN" ]; then
        echo -e "${YELLOW}⚠️  WARN${NC} $message ${CYAN}(+$((points/2)) pts)${NC}"
        HEALTH_SCORE=$((HEALTH_SCORE + points/2))
    else
        echo -e "${BLUE}ℹ️  INFO${NC} $message"
    fi
}

# I feel like theres a better way to extract owner and repo from GitHub URL...
extract_repo_info() {
    local url=$1
    # Remove trailing slash and .git
    url=$(echo "$url" | sed 's|/$||' | sed 's|\.git$||')
    
    # Extract owner/repo using regex
    if [[ $url =~ github\.com/([^/]+)/([^/]+) ]]; then
        REPO_OWNER="${BASH_REMATCH[1]}"
        REPO_NAME="${BASH_REMATCH[2]}"
    else
        echo -e "${RED}Error: Invalid GitHub URL format${NC}"
        exit 1
    fi
}

# GitHub API calls
github_api() {
    local endpoint=$1
    local url="https://api.github.com/repos/$REPO_OWNER/$REPO_NAME$endpoint"
    
    # Add GitHub token if available (for higher rate limits)
    if [ ! -z "$GITHUB_TOKEN" ]; then
        curl -s -H "Authorization: token $GITHUB_TOKEN" "$url"
    else
        curl -s "$url"
    fi
}

# if file exists in repo
file_exists() {
    local file=$1
    local response=$(github_api "/contents/$file")
    
    if echo "$response" | grep -q '"name"'; then
        return 0
    else
        return 1
    fi
}

# health check 
check_repo_health() {
    echo -e "${PURPLE}🔍 Repository Health Check for: $REPO_OWNER/$REPO_NAME${NC}"
    echo -e "${BLUE}================================================${NC}"
    echo ""
    
    REPO_DATA=$(github_api "")
    
    if echo "$REPO_DATA" | grep -q '"message": "Not Found"'; then
        echo -e "${RED}Error: Repository not found or not accessible${NC}"
        exit 1
    fi
    
    # some grep practice 
    DESCRIPTION=$(echo "$REPO_DATA" | grep -o '"description": "[^"]*"' | cut -d'"' -f4)
    LANGUAGE=$(echo "$REPO_DATA" | grep -o '"language": "[^"]*"' | cut -d'"' -f4)
    STARS=$(echo "$REPO_DATA" | grep -o '"stargazers_count": [0-9]*' | cut -d':' -f2 | tr -d ' ')
    FORKS=$(echo "$REPO_DATA" | grep -o '"forks_count": [0-9]*' | cut -d':' -f2 | tr -d ' ')
    ISSUES=$(echo "$REPO_DATA" | grep -o '"open_issues_count": [0-9]*' | cut -d':' -f2 | tr -d ' ')
    HAS_WIKI=$(echo "$REPO_DATA" | grep -o '"has_wiki": [a-z]*' | cut -d':' -f2 | tr -d ' ')
    CREATED_AT=$(echo "$REPO_DATA" | grep -o '"created_at": "[^"]*"' | cut -d'"' -f4)
    UPDATED_AT=$(echo "$REPO_DATA" | grep -o '"updated_at": "[^"]*"' | cut -d'"' -f4)
    
    echo -e "${CYAN}📊 Basic Repository Information:${NC}"
    echo "   Language: ${LANGUAGE:-"Not specified"}"
    echo "   Stars: $STARS | Forks: $FORKS | Open Issues: $ISSUES"
    echo "   Created: $CREATED_AT"
    echo "   Last Updated: $UPDATED_AT"
    echo ""
    
    echo -e "${BLUE}🔎 Checking Essential Files...${NC}"
    if file_exists "README.md" || file_exists "README.rst" || file_exists "README.txt" || file_exists "README"; then
        print_status "PASS" "README file exists" 15
    else
        print_status "FAIL" "No README file found" 15
    fi
    
    if file_exists "LICENSE" || file_exists "LICENSE.md" || file_exists "LICENSE.txt" || file_exists "COPYING"; then
        print_status "PASS" "LICENSE file exists" 15
    else
        print_status "FAIL" "No LICENSE file found" 15
    fi
    
    if file_exists ".gitignore"; then
        print_status "PASS" ".gitignore file exists" 10
    else
        print_status "WARN" "No .gitignore file found" 10
    fi
    
    if [ ! -z "$DESCRIPTION" ] && [ "$DESCRIPTION" != "null" ]; then
        print_status "PASS" "Repository has description" 5
    else
        print_status "FAIL" "Repository lacks description" 5
    fi
    
    echo ""
    echo -e "${BLUE}🚀 Checking CI/CD Setup...${NC}"
    ci_found=false
    
    if file_exists ".github/workflows" 2>/dev/null; then
        print_status "PASS" "GitHub Actions workflow found" 10
        ci_found=true
    fi
    
    for ci_file in ".travis.yml" ".circleci/config.yml" "Jenkinsfile" ".gitlab-ci.yml" "azure-pipelines.yml"; do
        if file_exists "$ci_file"; then
            print_status "PASS" "CI/CD configuration found ($ci_file)" 10
            ci_found=true
            break
        fi
    done
    
    if [ "$ci_found" = false ]; then
        print_status "WARN" "No CI/CD configuration found" 10
    fi
    
    echo ""
    echo -e "${BLUE}👥 Checking Community Health...${NC}"
    if file_exists "CONTRIBUTING.md" || file_exists "CONTRIBUTING.rst" || file_exists ".github/CONTRIBUTING.md"; then
        print_status "PASS" "Contributing guidelines exist" 5
    else
        print_status "WARN" "No contributing guidelines found" 5
    fi
    
    if file_exists "CODE_OF_CONDUCT.md" || file_exists ".github/CODE_OF_CONDUCT.md"; then
        print_status "PASS" "Code of Conduct exists" 5
    else
        print_status "WARN" "No Code of Conduct found" 5
    fi
    
    if file_exists ".github/ISSUE_TEMPLATE" || file_exists ".github/issue_template.md"; then
        print_status "PASS" "Issue templates exist" 5
    else
        print_status "WARN" "No issue templates found" 5
    fi
    
    echo ""
    echo -e "${BLUE}🔧 Language-Specific Checks...${NC}"
    case $LANGUAGE in
        "JavaScript")
            if file_exists "package.json"; then
                print_status "PASS" "package.json exists" 10
            else
                print_status "FAIL" "package.json missing for JavaScript project" 10
            fi
            ;;
        "Python")
            if file_exists "requirements.txt" || file_exists "setup.py" || file_exists "pyproject.toml" || file_exists "Pipfile"; then
                print_status "PASS" "Python dependencies file exists" 10
            else
                print_status "WARN" "No Python dependencies file found" 10
            fi
            ;;
        "Go")
            if file_exists "go.mod"; then
                print_status "PASS" "go.mod exists" 10
            else
                print_status "WARN" "go.mod missing for Go project" 10
            fi
            ;;
        "Java")
            if file_exists "pom.xml" || file_exists "build.gradle"; then
                print_status "PASS" "Build configuration exists" 10
            else
                print_status "WARN" "No build configuration found" 10
            fi
            ;;
        *)
            print_status "INFO" "Language-specific checks skipped for $LANGUAGE" 0
            HEALTH_SCORE=$((HEALTH_SCORE + 5)) # Don't penalize for unknown languages
            ;;
    esac
    

    echo ""
    echo -e "${BLUE}📈 Activity Analysis...${NC}"
    
    COMMITS_DATA=$(github_api "/commits?per_page=10")
    RECENT_COMMITS=$(echo "$COMMITS_DATA" | grep -o '"date": "[^"]*"' | head -5 | wc -l)
    
    if [ "$RECENT_COMMITS" -gt 0 ]; then
        print_status "PASS" "Repository has recent commits" 10
    else
        print_status "WARN" "No recent commits found" 10
    fi
    
    if [ "$ISSUES" -lt 50 ]; then
        print_status "PASS" "Manageable number of open issues ($ISSUES)" 5
    elif [ "$ISSUES" -lt 100 ]; then
        print_status "WARN" "High number of open issues ($ISSUES)" 5
    else
        print_status "FAIL" "Very high number of open issues ($ISSUES)" 5
    fi
    
    echo ""
    echo -e "${PURPLE}================================================${NC}"
    echo -e "${CYAN}🏥 FINAL HEALTH SCORE: $HEALTH_SCORE/$MAX_SCORE${NC}"
    
    if [ "$HEALTH_SCORE" -ge 80 ]; then
        echo -e "${GREEN}🌟 EXCELLENT HEALTH - This repository follows best practices!${NC}"
    elif [ "$HEALTH_SCORE" -ge 60 ]; then
        echo -e "${YELLOW}👍 GOOD HEALTH - Room for some improvements${NC}"
    elif [ "$HEALTH_SCORE" -ge 40 ]; then
        echo -e "${YELLOW}⚠️  FAIR HEALTH - Several areas need attention${NC}"
    else
        echo -e "${RED}🚨 POOR HEALTH - Significant improvements needed${NC}"
    fi
    
    echo -e "${PURPLE}================================================${NC}"
}

if [ $# -eq 0 ]; then
    echo "Usage: $0 <github_repo_url>"
    echo "Example: $0 https://github.com/microsoft/vscode"
    exit 1
fi

extract_repo_info "$1"

check_repo_health
