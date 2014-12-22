#!/bin/bash
asciidoc --out-file=doc/git-octopus.html doc/git-octopus.1.txt
a2x -f manpage doc/git-octopus.1.txt --verbose --no-xmllint --destination-dir=man/man1/
