#!/bin/bash

#FIXME
git status &> /dev/null

#saving the initial state of the repository
sha1=`git rev-parse HEAD`

git octopus

if [ $? -ne 0 ] ; then
	exit 1
fi

if [[ `git rev-parse HEAD` != $sha1 ]] ; then
	echo "should stayed at HEAD"
	exit 1
fi

if [[ -n `git diff-index HEAD` ]] ; then
	echo "repository should be clean"
	exit 1
fi

if [[ `git symbolic-ref HEAD` != "refs/heads/master" ]] ; then
	echo "Should be still on master"
	exit 1
fi