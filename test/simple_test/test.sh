#!/bin/bash
git status
/home/git-octopus.sh features/*
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