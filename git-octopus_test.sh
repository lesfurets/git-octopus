#!/bin/bash
git commit --allow-empty -m "first commit"
head=$(git rev-parse HEAD)
git update-ref refs/heads/test1 $head
git update-ref refs/remotes/origin/test1 $head
git update-ref refs/remotes/origin/test2 $head
