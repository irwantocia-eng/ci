#!/bin/bash
set -e

# SonarScanner installation and execution script
# Supports both local SonarQube and SonarCloud

# Load environment variables from .env if present
if [ -f .env ]; then
    export $(cat .env | grep -v '^#' | xargs)
fi

SONAR_SCANNER_VERSION=6.2.1.4610-linux-x64
SONAR_SCANNER_HOME=$HOME/.sonar/sonar-scanner-$SONAR_SCANNER_VERSION
SONAR_SCANNER_BIN=$HOME/.sonar/sonar-scanner-$SONAR_SCANNER_VERSION/bin/sonar-scanner

# Default values
SONAR_HOST_URL=${SONAR_HOST_URL:-http://localhost:9000}
SONAR_PROJECT_KEY=${SONAR_PROJECT_KEY:-ci}
SONAR_TOKEN=${SONAR_TOKEN:-}

# Detect OS
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

if [ "$OS" = "linux" ]; then
    SONAR_SCANNER_BIN=$HOME/.sonar/sonar-scanner-$SONAR_SCANNER_VERSION/bin/sonar-scanner
    ZIP_FILE="sonar-scanner-cli-$SONAR_SCANNER_VERSION-linux-x64.zip"
    DOWNLOAD_URL="https://binaries.sonarsource.com/Distribution/sonar-scanner-cli/$ZIP_FILE"
elif [ "$OS" = "darwin" ]; then
    SONAR_SCANNER_BIN=$HOME/.sonar/sonar-scanner-$SONAR_SCANNER_VERSION-macosx/bin/sonar-scanner
    ZIP_FILE="sonar-scanner-cli-$SONAR_SCANNER_VERSION-macosx.zip"
    DOWNLOAD_URL="https://binaries.sonarsource.com/Distribution/sonar-scanner-cli/$ZIP_FILE"
else
    echo "Unsupported OS: $OS"
    exit 1
fi

# Check if SonarScanner is already installed
if [ ! -f "$SONAR_SCANNER_BIN" ]; then
    echo "Downloading SonarScanner CLI $SONAR_SCANNER_VERSION..."
    mkdir -p $HOME/.sonar
    
    # Download and extract
    curl -L --output $HOME/.sonar/$ZIP_FILE $DOWNLOAD_URL
    
    echo "Extracting SonarScanner..."
    unzip -o $HOME/.sonar/$ZIP_FILE -d $HOME/.sonar/
    
    # Clean up zip
    rm $HOME/.sonar/$ZIP_FILE
    
    echo "SonarScanner installed at: $SONAR_SCANNER_BIN"
else
    echo "SonarScanner already installed"
fi

# Check if token is configured
if [ -z "$SONAR_TOKEN" ]; then
    echo ""
    echo "Error: SONAR_TOKEN not configured!"
    echo ""
    echo "Please set SONAR_TOKEN environment variable or add it to .env file:"
    echo "  1. Start SonarQube: docker compose up -d"
    echo "  2. Go to http://localhost:9000"
    echo "  3. Login and generate a token: My Account → Security → Generate Token"
    echo "  4. Create .env file with:"
    echo "     SONAR_TOKEN=your_token_here"
    echo ""
    exit 1
fi

# Check if project key is configured
if [ "$SONAR_PROJECT_KEY" = "YOUR_GITHUB_USERNAME_YOUR_PROJECT_KEY" ]; then
    echo ""
    echo "Error: Project key not configured!"
    echo ""
    echo "Please update .env or set SONAR_PROJECT_KEY with your project key"
    echo ""
    exit 1
fi

# Run SonarScanner
echo "Running SonarScanner..."
echo "  Host: $SONAR_HOST_URL"
echo "  Project: $SONAR_PROJECT_KEY"

$SONAR_SCANNER_BIN \
    -Dsonar.projectKey=$SONAR_PROJECT_KEY \
    -Dsonar.sources=. \
    -Dsonar.host.url=$SONAR_HOST_URL \
    -Dsonar.token=$SONAR_TOKEN \
    -Dsonar.go.coverage.reportPaths=coverage.out \
    -Dsonar.go.golangci-lint.reportPaths=golangci-lint-report.xml

echo ""
echo "Analysis complete!"
echo "View results at: $SONAR_HOST_URL/dashboard?id=$SONAR_PROJECT_KEY"
