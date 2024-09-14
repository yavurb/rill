#!/bin/bash

version_type=$(echo "$MESSAGE" | grep -o -i "\[MAJOR\]\|\[MINOR\]\|\[PATCH\]" | tr -d '[]' | tr '[:upper:]' '[:lower:]')
if [ "$version_type" == "" ]; then
  echo "Error: No version type found in the message."
  exit 1
fi

echo "$version_type"
