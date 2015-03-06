#git-octopus
>git-octopus - extends git-merge with branch naming patterns.

##Installation
Clone this repository and add the bin/ directory to the PATH.

##Usage
The main purpose of the command is to merge branches based on naming patterns
```
git octopus features/*
```
Given a set of branches features/something1 features/something2 features/something3, this command will do
```
git merge features/something1 features/something2 features/something3
```

See git-octopus' [documentation](http://lesfurets.github.io/git-octopus/doc/git-octopus.html) for more information.

##Conflicts detection
If the merge fails, the command will start a diagnostic phase that will try to merge each branch separated from the others into the current branch to find out which one is in conflict.

Note that when an octopus fails, it doesn't necessarily mean that a given branch has conflicts with the current branch, it could actually be with any other one. This means that the diagnostic may not find any conflicts. 

##Continuous integration
git-octopus is meant to be used in a continuous integration build flow. The goal is to merge all the branches the developpers are working on. 

The implementation in a command line job is pretty straight forward
```
#Usually, the CI fetches only the branch the job has been started on, so we need to fetch the rest
git fetch -p

#Lets merge feature branches and the master from origin
git octopus origin/features/* origin/master

if [ $? -eq 0 ] ; then
  git push origin +HEAD:octopus
else
  exit 1
fi
```
