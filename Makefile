prefix = /usr/local
bindir = $(prefix)/bin
datarootdir = $(prefix)/share
mandir = $(datarootdir)/man
docdir = $(datarootdir)/doc/git-doc
htmldir = $(docdir)

CAT_SCRIPTS = @`cat $^ > bin/$@ && chmod +x bin/$@`

build: setup git-octopus git-conflict git-apply-conflict-resolution
	@echo 'Build success'

setup: 
	@mkdir -p bin

git-octopus: src/lib/common src/lib/git-merge-octopus-fork.sh src/git-octopus
	${CAT_SCRIPTS}
	
git-conflict: src/lib/common src/lib/hash-conflict src/git-conflict
	${CAT_SCRIPTS}

git-apply-conflict-resolution: src/lib/common src/lib/hash-conflict src/git-apply-conflict-resolution 
	${CAT_SCRIPTS}

install-bin: build
	@cp -f bin/git-octopus $(bindir) && echo 'Installing $(bindir)/git-octopus'
	@cp -f bin/git-conflict $(bindir) && echo 'Installing $(bindir)/git-conflict'
	@cp -f bin/git-apply-conflict-resolution $(bindir) && echo 'Installing $(bindir)/git-apply-conflict-resolution'

install-docs:
	@echo 'Installing documentation'
	@cp -f man/man1/git-octopus.1 $(mandir)/man1/git-octopus.1
	@mkdir -p $(htmldir)
	@cp -f doc/git-octopus.html $(htmldir)

install: install-bin install-docs

uninstall:
	rm $(bindir)/git-octopus
	rm $(bindir)/git-conflict
	rm $(bindir)/git-apply-conflict-resolution
	rm $(mandir)/man1/git-octopus.1
	rm $(htmldir)/git-octopus.html