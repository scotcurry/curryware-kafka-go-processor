#!/bin/bash
VERSION_FILE="version.txt"
if [ -f "$VERSION_FILE" ]; then
    VERSION=$(cat "$VERSION_FILE")
    # Extract the build number part, increment it, and replace it
    BUILD_NUM=$(echo "$VERSION" | awk -F. '{print $NF}')
    NEW_BUILD_NUM=$((BUILD_NUM+1))
    NEW_VERSION="1.0.0-build.$NEW_BUILD_NUM"
    echo "$NEW_VERSION" > "$VERSION_FILE"
else
    echo "1.0.0-build.1" > "$VERSION_FILE"
fi