#!/bin/bash

tag=$1

echo "Remove tag $tag from master"

git tag -d $tag
git push origin :refs/tags/$tag
