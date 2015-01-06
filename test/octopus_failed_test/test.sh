#!/bin/bash

#FIXME
git status &> /dev/null

sha1=`git rev-parse HEAD`

git octopus -n features/*

if [ $? -ne 1 ] ; then
	exit 1
fi

#should be back to HEAD
if [[ `git rev-parse HEAD` != $sha1 ]] ; then
	exit 1
fi

#repository should be clean
if [[ -n `git diff-index HEAD` ]] ; then
	exit 1
fi