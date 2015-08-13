prefix = /usr/local
bindir = $(prefix)/bin
datarootdir = $(prefix)/share
mandir = $(datarootdir)/man
docdir = $(datarootdir)/doc/git-doc
htmldir = $(docdir)

cat_scripts = cat $(2) $(3) $(4) $(5) > bin/$(1) \
	&& chmod +x bin/$(1)

generate_docs = asciidoc --out-file=doc/$(1).html src/doc/$(1).1.txt \
    && a2x -f manpage src/doc/$(1).1.txt --no-xmllint --destination-dir=doc

install_docs = cp -f doc/$(1).1 $(mandir)/man1/$(1).1 \
	&& cp -f doc/$(1).html $(htmldir)

build:
	@mkdir -p bin
	$(call cat_scripts,git-octopus,src/lib/common,src/lib/git-merge-octopus-fork.sh,src/git-octopus)
	$(call cat_scripts,git-conflict,src/lib/common,src/lib/hash-conflict,src/git-conflict)
	$(call cat_scripts,git-apply-conflict-resolution,src/lib/common,src/lib/hash-conflict,src/git-apply-conflict-resolution)
	@echo 'Build success'

build-docs:
	mkdir -p doc
	$(call generate_docs,git-octopus)
	$(call generate_docs,git-conflict)

install-bin: build
	@cp -f bin/git-octopus $(bindir) && echo 'Installing $(bindir)/git-octopus'
	@cp -f bin/git-conflict $(bindir) && echo 'Installing $(bindir)/git-conflict'
	@cp -f bin/git-apply-conflict-resolution $(bindir) && echo 'Installing $(bindir)/git-apply-conflict-resolution'

install-docs:
	@echo 'Installing documentation'
	@mkdir -p $(htmldir)
	$(call install_docs,git-octopus)
	$(call install_docs,git-conflict)

install: install-bin install-docs

uninstall:
	rm $(bindir)/git-octopus
	rm $(bindir)/git-conflict
	rm $(bindir)/git-apply-conflict-resolution
	rm $(mandir)/man1/git-octopus.1
	rm $(htmldir)/git-octopus.html