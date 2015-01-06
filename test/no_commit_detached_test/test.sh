#!/bin/bash
git status &> /dev/null

git checkout --detach

sha1=`git rev-parse HEAD`

git octopus -n features/*

#should be back to HEAD
if [[ `git rev-parse HEAD` != $sha1 ]] ; then
	exit 1
fi

#repository should be clean
if [[ -n `git diff-index HEAD` ]] ; then
	exit 1
fi

#should be still detached
git symbolic-ref HEAD &> /dev/null

if [ $? -eq 0 ] ; then
	echo "Repository should remains detached"
	exit 1
fi