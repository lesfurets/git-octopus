#git-octopus
>The continuous merge workflow is meant for continuous integration/delivery and is based on feature branching. git-octopus provides git commands to implement it.

##Installation

###Requirements
Requires git >= 1.8
You need to have a command `shasum` in the path. This is the case on most unix based systems. If you're on Windows/Cygwin you may have to install it.

###Homebrew
If you know and use [Homebrew](http://brew.sh), you just need to do:
```bash
brew update
brew install git-octopus
```

###From sources
Download the latest [release](https://github.com/lesfurets/git-octopus/releases/latest) or clone this repository. Go into the directory and type
```bash
make install
```

Make sure the installation works
```bash
git octopus -v
```

## What you'll find
Two additionnal git commands : 

### git octopus
Extends `git merge` with branch naming patterns. For instance
```
git octopus features/*
```
Merges all branches named features/ into the current branch.
See [git-octopus(1)](http://lesfurets.github.io/git-octopus/doc/git-octopus.html).

### git conflict
Allows you to record conflicts resolutions that `git octopus` can reuse.
Conflicts resolutions are standard refs so they can be pushed/fetched.
See the conflicts management section bellow and [git-conflict(1)](http://lesfurets.github.io/git-octopus/doc/git-conflict.html).

##The Continuous Merge

### What is it all about ?
Feature branching and continuous integration don't live well together. The idea of this project is to reconcile those two by using the mighty power of git.

I gave a talk about why and how to use it at Devoxx France 2015, but it's in french ;) https://www.parleys.com/tutorial/le-continuous-merge-chez-lesfurets-com

###The branching model
The simpliest form of the model is to have a mainline branch, let's call it `master`, and feature branches on top of that master. In a continuous delivery workflow you won't need more than that. 

* The `master` branch, or however you call it, is in a ready-to-ship state. Nobody commits on it.
* A feature branch is a change, as small as possible, that can bring the `master` from a ready-to-ship state to an other.

This means that all the work is done in feature branches. Don't be afraid to have many, one branch per developer is fine. Keep feature branches independent from each other, that's the key for having a fluent delivery pipe.

###The workflow
`git octopus` allows you to merge all you feature branches together at any moment so you can have an assembly of all the work that is going on and finally do a continuous integration job on that merge. here's how it works : 

A developer pushes a change on his feature branch. There is a job in your continuous integration system that will trigger and do this bash command :

```bash
git octopus origin/features/* origin/master && git push origin +HEAD:octopus
```
This job computes a merge with all feature branches and the master, and then pushes the result on a branch `octopus` on origin. 
The new merge commit on `octopus` will now trigger an other job that will build and deploy this merge on your test servers etc ...
Note that the octopus merge is not kept in any history line. The next push on any feature branch will trigger the build of a new merge that will be forced push again on `octopus`.

Once a feature branch is validated on your test environment, you can merge it on master.

### Managing conflicts
If `git-octopus` fails, it will do a diagnostic sequence to figure out the conflict precisely. It can lead to two cases : 

* A conflict has been found

	1. Ask yourself if you could avoid that conflict. Rewriting the history is possible as long as you're alone working on the branch. 

	2. Use [git-conflict](http://lesfurets.github.io/git-octopus/doc/git-conflict.html) to record a resolution and push it to origin. See the documentation for more details.

	3. Consider to remove one of the conflicting branches from the continuous integration (I.E. rename the branch so it won't get caught in the merge) and wait for the other to be merged in `master`. Then you'll be able to update and resolve the conflict.

	4. Rebase one branch on top of the other (depending on which one you want to ship first). This has to be the last resort because you'll loose branches independency.

* No conflict found

	1. Someone else might breaks the merge, look at previous octopus job executions.

	2. You felt in a complex case. There are ongoing works to prevent that from happening but for the moment this might happen. Don't hesitate to open an issue !

## Community

We have a [Google Group](https://groups.google.com/forum/#!forum/git-octopus), feel free to come and discuss with us. You can also send an email to git-octopus@googlegroups.com.

[![Analytics](https://ga-beacon.appspot.com/UA-79856083-1/README.md?pixel&useReferrer)](https://github.com/igrigorik/ga-beacon)
