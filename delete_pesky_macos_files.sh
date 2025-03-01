#!/bin/sh

# Delete macOS resource fork files and .DS_Store files
find "$(pwd)" -type f \( -name '._*' -o -name '.DS_Store' \) -exec rm -v {} \;

echo "All macOS metadata files (.DS_Store and ._*) have been removed."
