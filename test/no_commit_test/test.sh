#!/bin/bash

#FIXME
git status &> /dev/null

#saving the initial state of the repository
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

#should be still on master
if [[ `git symbolic-ref HEAD` != "refs/heads/master" ]] ; then
	echo "Should be still on master"
	exit 1
fi