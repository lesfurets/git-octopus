#!/bin/bash

#FIXME
git status &> /dev/null

git branch -D features/feat2 features/feat3
git branch features/feat0

git octopus features/*
merged=`git branch --merged`
if [[ $merged != *feat0* ]] ; then
	exit 1
fi
if [[ $merged != *feat1* ]] ; then
	exit 1
fi
if [[ $merged != *master* ]] ; then
	exit 1
fi
