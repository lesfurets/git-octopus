CAT_SCRIPTS = @`cat $^ > bin/$@ && chmod +x bin/$@`

cat: setup git-octopus git-conflict git-apply-conflict-resolution

setup: 
	@mkdir -p bin

git-octopus: src/lib/common src/lib/git-merge-octopus-fork.sh src/git-octopus
	${CAT_SCRIPTS}
	
git-conflict: src/lib/common src/lib/hash-conflict src/git-conflict
	${CAT_SCRIPTS}

git-apply-conflict-resolution: src/lib/common src/lib/hash-conflict src/git-apply-conflict-resolution 
	${CAT_SCRIPTS}
