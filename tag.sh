#!/bin/bash

git checkout master || exit -1
git pull || exit -1

# Add tag before merging from master
tag=$1

echo "Add tag $tag to master"

git tag -a $tag -m "New version ${tag}"
git push origin $tag
