PREFIX = /usr/local
BINDIR = $(PREFIX)/bin
MANDIR = $(PREFIX)/share/man
HTMLDIR = `git --html-path`

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
	@cp -f bin/git-octopus $(BINDIR) && echo 'Installing $(BINDIR)/git-octopus'
	@cp -f bin/git-conflict $(BINDIR) && echo 'Installing $(BINDIR)/git-conflict'
	@cp -f bin/git-apply-conflict-resolution $(BINDIR) && echo 'Installing $(BINDIR)/git-apply-conflict-resolution'

install-docs:
	@echo 'Installing documentation'
	@cp -f man/man1/git-octopus.1 $(MANDIR)/man1/git-octopus.1
	@mkdir -p $(HTMLDIR)
	@cp -f doc/git-octopus.html $(HTMLDIR)

install: install-bin install-docs

uninstall:
	rm $(BINDIR)/git-octopus
	rm $(BINDIR)/git-conflict
	rm $(BINDIR)/git-apply-conflict-resolution
	rm $(MANDIR)/man1/git-octopus.1
	rm $(HTMLDIR)/git-octopus.html