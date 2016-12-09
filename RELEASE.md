# Source

## Bump the version
This is not automated. Search references of the current version number in the `src` folder and change them.

## Generate the documentation
Note that even if there's no documentation change, this is still needed because the version number is visible in the man doc.

Run `make build-docs install-docs` (requires asciidoc to be installed). Make sure the documentation looks good.

## Commit & push
Commit the version bump and the generated documentation.

Push to `master` and `gh-pages` branches.

## Tag and patch note
Create a [release on github](https://github.com/lesfurets/git-octopus/releases). Write a patch note and publish it. This will tag the current master.

# Homebrew
`git-octopus` is part of the homebrew/core tap. Retrieve the tap installation path with
```
brew tap-info homebrew/core
```
git-octopus's formula is in `Formula/git-octopus.rb`

Here's the guidelines : 
* [CONTRIBUTING.md of homebrew/core tap](https://github.com/homebrew/homebrew-core/blob/master/.github/CONTRIBUTING.md)
* [Formula Cookbook](https://github.com/Homebrew/brew/blob/master/share/doc/homebrew/Formula-Cookbook.md)

# RPM - Fedora/RHEL
## Prerequisites

Install mock and add current user to mock group:

    sudo yum install mock
    sudo usermod -a -G mock $(id -u -n)
    
Clone the spec file project (the project contains the spec file needed to build the rpm)

    git clone https://github.com/danoliv/git-octopus-spec.git

## Update the spec file

**You should modify the spec file following the official guidelines: [https://fedoraproject.org/wiki/Packaging:Guidelines](https://fedoraproject.org/wiki/Packaging:Guidelines)**

Check the last available version of git-octopus from the official repo: [https://github.com/lesfurets/git-octopus/releases/latest](https://github.com/lesfurets/git-octopus/releases/latest)

Update the Version tag of the spec file to match the latest version, set the Release number to 1 for a new version, increase it if it is only a packaging modification:

    Name:   	git-octopus
    Version:	1.4
    Release:	1%{?dist}
    Summary:	Git commands for continuous delivery
    
Update the changelog in the spec files:

    %changelog
    * Tue Dec 06 2016 Andrea Baita <andrea@baita.pro> - 1.4-2
    - added documentation build, updated build requires
     
    * Wed Nov 30 2016 Andrea Baita <andrea@baita.pro> - 1.4-1
    - Packaging of version 1.4.
     
    * Thu Nov 17 2016 Xavier Bachelot <xavier@bachelot.org> - 1.3-1
    - Initial package.
    
## Build the RPM

retrieve the tarball, the file will be put into `~/rpmbuild/SOURCES` (will create a directory if not exists)

    spectool -g -R git-octopus.spec
    
build the source rpm, the file will be put into `~/rpmbuild/SRPMS/` (will create a directory if not exists):

    rpmbuild -bs git-octopus.spec
    
finally build the rpm, by passing the src.rpm file created in the previous step, indicate a configuration to use from `/etc/mock/`, without the path prefix and the .cfg suffix:

    mock -r <configuration> ~/rpmbuild/SRPMS/git-octopus-<version>.<release>.src.rpm

example:
    
    mock -r epel-6-x86_64 ~/rpmbuild/SRPMS/git-octopus-1.4-2.el6.src.rpm

the results will be usually available in the directory: `/var/lib/mock/<configuration>/result` (check the output of mock command)

**If the rpm build fails the spec there could have been some incompatible modification on the code, the spec file should be updated accordingly.**

## Test the new package

check the rpm by compiling and installing in the local machine

    mvn clean install
    sudo yum install <specify local rpm>
    
check the git octopus version

    git octopus -v

## Make a pull request

Please make a pull request following the github guide: https://guides.github.com/activities/forking/