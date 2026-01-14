#!/bin/bash

# --- Installer for the Cooler Tool (Go Version) ---

echo "--- Cooler Tool Installer ---"
echo "This script will install the necessary dependencies for the Cooler tool."
echo
echo "It will perform the following actions:"
echo "1. Install the Go compiler."
echo "2. Build the 'cooler' executable from the source code."
echo "3. Install Node.js (v18) and npm for the Gemini CLI."
echo "4. Install the Google Gemini CLI via npm."
echo "5. Install system utilities: 'lm-sensors' and 'cpufrequtils'."
echo "6. Configure hardware sensors."
echo
echo "This will install packages on your system and requires sudo privileges."
read -p "Do you want to continue? (y/n) " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]
then
    echo "Installation cancelled."
    exit 1
fi

# --- 1. Install Go ---
if command -v go &> /dev/null
then
    echo "Go is already installed. Skipping."
else
    echo "Installing Go compiler..."
    sudo apt-get update
    sudo apt-get install -y golang-go
    echo "Go installation complete."
fi

# --- 2. Build the executable ---
echo "Building the Cooler executable with optimizations..."
# Build flags explanation:
# -ldflags="-s -w": Strip symbol table and debug info for smaller binary
# -trimpath: Remove file system paths from compiled binary for reproducibility
go build -ldflags="-s -w" -trimpath -o cooler ./cmd/cooler
if [ $? -eq 0 ]; then
    echo "Build successful."
else
    echo "Build failed. Please check for Go installation and source code errors."
    exit 1
fi

# --- 3. Install Node.js and npm ---
if command -v node &> /dev/null
then
    echo "Node.js is already installed. Skipping."
else
    echo "Installing Node.js v18..."
    curl -fsSL https://deb.nodesource.com/setup_18.x | sudo -E bash -
    sudo apt-get install -y nodejs
    echo "Node.js installation complete."
fi

# --- 4. Install Gemini CLI ---
if command -v gemini &> /dev/null
then
    echo "Gemini CLI is already installed. Skipping."
else
    echo "Installing Gemini CLI..."
    sudo npm install -g @google/gemini-cli
    echo "Gemini CLI installation complete."
    echo "IMPORTANT: You may need to run 'gemini' once manually to log in."
fi

# --- 5. Install System Utilities ---
echo "Installing lm-sensors and cpufrequtils..."
sudo apt-get install -y lm-sensors cpufrequtils

# --- 6. Configure Sensors ---
echo "Configuring hardware sensors..."
sudo sensors-detect --auto
sudo service kmod start

echo
echo "--- Installation Complete! ---"
echo "You can now run the tool with: ./cooler"
echo "(You might need to start a new terminal session for some commands to be available.)"