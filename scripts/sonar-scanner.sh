#!/bin/bash
set -e

# SonarScanner installation and execution script
# This script downloads SonarScanner if not present and runs the analysis

SONAR_SCANNER_VERSION=6.2.1.4610-linux-x64
SONAR_SCANNER_HOME=$HOME/.sonar/sonar-scanner-$SONAR_SCANNER_VERSION
SONAR_SCANNER_BIN=$HOME/.sonar/sonar-scanner-$SONAR_SCANNER_VERSION/bin/sonar-scanner

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
    echo "📥 Downloading SonarScanner CLI $SONAR_SCANNER_VERSION..."
    mkdir -p $HOME/.sonar
    
    # Download and extract
    curl -L --output $HOME/.sonar/$ZIP_FILE $DOWNLOAD_URL
    
    echo "📦 Extracting SonarScanner..."
    unzip -o $HOME/.sonar/$ZIP_FILE -d $HOME/.sonar/
    
    # Clean up zip
    rm $HOME/.sonar/$ZIP_FILE
    
    echo "✅ SonarScanner installed at: $SONAR_SCANNER_BIN"
else
    echo "✅ SonarScanner already installed"
fi

# Check if project key is configured
PROJECT_KEY=$(grep "^sonar.projectKey=" sonar-project.properties | cut -d'=' -f2)
if [ "$PROJECT_KEY" = "YOUR_GITHUB_USERNAME_YOUR_PROJECT_KEY" ]; then
    echo ""
    echo "❌ Error: Project key not configured!"
    echo ""
    echo "Please update sonar-project.properties with your SonarCloud credentials:"
    echo "  1. Go to https://sonarcloud.io"
    echo "  2. Create a project"
    echo "  3. Copy your project key"
    echo "  4. Update sonar-project.properties:"
    echo "     sonar.projectKey=YOUR_ACTUAL_PROJECT_KEY"
    echo "     sonar.organization=YOUR_ACTUAL_ORGANIZATION"
    echo ""
    exit 1
fi

# Run SonarScanner
echo "🔍 Running SonarScanner..."
$SONAR_SCANNER_BIN \
    -Dsonar.projectKey=$PROJECT_KEY \
    -Dsonar.sources=. \
    -Dsonar.host.url=https://sonarcloud.io \
    -Dsonar.token=$SONAR_TOKEN \
    -Dsonar.go.coverage.reportPaths=coverage.out

echo ""
echo "✅ Analysis complete!"
echo "📊 View results at: https://sonarcloud.io/dashboard?id=$PROJECT_KEY"
