#!/bin/bash

#FIXME
git status &> /dev/null

git checkout features/feat1
git config octopus.pattern features/*
git config --add octopus.pattern master

/home/git-octopus
merged=`git branch --merged`
if [[ $merged != *feat1* ]] ; then
	exit 1
fi
if [[ $merged != *feat2* ]] ; then
	exit 1
fi
if [[ $merged != *feat3* ]] ; then
	exit 1
fi
if [[ $merged != *master* ]] ; then
	exit 1
fi