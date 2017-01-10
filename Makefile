prefix = /usr/local
bindir = $(prefix)/bin
datarootdir = $(prefix)/share
mandir = $(datarootdir)/man
man1dir = $(mandir)/man1
docdir = $(datarootdir)/doc/git-doc
htmldir = $(docdir)

cat_scripts = cat $(2) $(3) $(4) $(5) > bin/$(1) \
	&& chmod +x bin/$(1)

generate_docs = asciidoc -d manpage --out-file=doc/$(1).html src/doc/$(1).1.txt \
    && a2x -f manpage src/doc/$(1).1.txt --no-xmllint --destination-dir=doc

install_docs = cp -f doc/$(1).1 $(man1dir)/$(1).1 \
	&& cp -f doc/$(1).html $(htmldir)

fmt:
	gofmt -w **/*.go

build:
	@mkdir -p bin
	$(call cat_scripts,git-octopus,src/lib/common,src/lib/git-merge-octopus-fork.sh,src/git-octopus)
	$(call cat_scripts,git-conflict,src/lib/common,src/lib/hash-conflict,src/git-conflict)
	$(call cat_scripts,git-apply-conflict-resolution,src/lib/common,src/lib/hash-conflict,src/git-apply-conflict-resolution)
	@echo 'Build success'

build-docs:
	@mkdir -p doc
	$(call generate_docs,git-octopus)
	$(call generate_docs,git-conflict)

install-bin: build
	@mkdir -p $(bindir)
	@cp -f bin/git-octopus $(bindir) && echo 'Installing $(bindir)/git-octopus'
	@cp -f bin/git-conflict $(bindir) && echo 'Installing $(bindir)/git-conflict'
	@cp -f bin/git-apply-conflict-resolution $(bindir) && echo 'Installing $(bindir)/git-apply-conflict-resolution'

install-docs:
	@echo 'Installing documentation'
	@mkdir -p $(htmldir)
	@mkdir -p $(man1dir)
	$(call install_docs,git-octopus)
	$(call install_docs,git-conflict)

install: install-bin install-docs

uninstall:
	rm $(bindir)/git-octopus
	rm $(bindir)/git-conflict
	rm $(bindir)/git-apply-conflict-resolution
	rm $(man1dir)/git-octopus.1
	rm $(man1dir)/git-conflict.1
	rm $(htmldir)/git-octopus.html
	rm $(htmldir)/git-conflict.html

go-build = GOOS=$(1) GOARCH=$(2) go build -o git-octopus-$(1)-$(2)-2.0.beta1

go-cross-compile:
	$(call go-build,darwin,386)
	$(call go-build,darwin,amd64)
	$(call go-build,freebsd,386)
	$(call go-build,freebsd,amd64)
	$(call go-build,linux,386)
	$(call go-build,linux,amd64)
	$(call go-build,windows,386)
	$(call go-build,windows,amd64)
