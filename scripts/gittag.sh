#!/bin/bash

semver=$(echo "$MESSAGE" | grep -i "release-" | cut -d '-' -f 2 | head -n 1)
if [ "$semver" == "" ]; then
  semver=$(git describe --tags --abbrev=0)
fi

echo "$semver"
