#!/bin/bash

#FIXME
git status &> /dev/null

mergeBase=$(git merge-base --all HEAD features/feat1)
git read-tree -um --aggressive $mergeBase HEAD features/feat1
git merge-index -o -q git-merge-one-file -a

git store-conflict