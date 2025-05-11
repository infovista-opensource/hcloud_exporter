#!/bin/sh
# Script to configure a repository-specific Git user
# Usage: ./setup-git-user.sh "Your Name" "your.email@example.com"

# Exit if any command fails
set -e

if [ "$#" -ne 2 ]; then
    echo "Usage: $0 \"Your Name\" \"your.email@example.com\""
    exit 1
fi

# Set repository-specific Git user
git config --local user.name "$1"
git config --local user.email "$2"

# Verify the changes
echo "\nRepository-specific Git user configured:"
git config --local --get user.name
git config --local --get user.email

echo "\nThese settings will only apply to the hcloud_exporter repository."
echo "Your global Git settings remain unchanged for other repositories."
