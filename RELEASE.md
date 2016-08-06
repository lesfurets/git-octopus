#Source

##Bump the version
This is not automated. Search references of the current version number in the `src` folder and change them.

##Generate the documentation
Note that even if there's no documentation change, this is still needed because the version number is visible in the man doc.

Run `make build-docs install-docs` (requires asciidoc to be installed). Make sure the documentation looks good.

##Commit & push
Commit the version bump and the generated documentation.

Push to `master` and `gh-pages` branches.

##Tag and patch note
Create a [release on github](https://github.com/lesfurets/git-octopus/releases). Write a patch note and publish it. This will tag the current master.

#Homebrew
`git-octopus` is part of the homebrew/core tap. Retrieve the tap installation path with
```
brew tap-info homebrew/core
```
git-octopus's formula is in `Formula/git-octopus.rb`

Here's the guidelines : 
* [CONTRIBUTING.md of homebrew/core tap](https://github.com/homebrew/homebrew-core/blob/master/.github/CONTRIBUTING.md)
* [Formula Cookbook](https://github.com/Homebrew/brew/blob/master/share/doc/homebrew/Formula-Cookbook.md)