DESTDIR = /usr/local/bin
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

install: build
	@cp -f bin/git-octopus $(DESTDIR) && echo 'Installing $(DESTDIR)/git-octopus'
	@cp -f bin/git-conflict $(DESTDIR) && echo 'Installing $(DESTDIR)/git-conflict'
	@cp -f bin/git-apply-conflict-resolution $(DESTDIR) && echo 'Installing $(DESTDIR)/git-apply-conflict-resolution'

uninstall:
	rm $(DESTDIR)/git-octopus
	rm $(DESTDIR)/git-conflict
	rm $(DESTDIR)/git-apply-conflict-resolution