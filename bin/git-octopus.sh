#!/bin/bash
usage() {
cat <<EOF
Usage:
    [OPTION..] [<refspec>]

    Exécute le merge octopus des branches remotes/origin/features origin/master et <refspec> en effectuant un checkout de <refspec> au préalable.
    La commande commence par faire un fetch -p.
    si <refspec> n'est pas précisé, utilise HEAD.
    La commande laisse le repo sur la branche <refspec> à la fin de son exécution.

OPTION :
    -a, --analyse exécute un merge octopus partiel entre <refspec>, origin/master et chaque feature branch individuellement pour détecter avec quelle branch se produit le conflit.
    --push push le commit resultant du merge (en cas de succes) sur origin octopus-features
    -s, --stop-on-failure Stops leaving the merge unresolved
    -h, --help, -help, help
    --porcelain retourne uniquement le SHA1 de la branche octopus-features
EOF
}

line_break(){
    log "------------------------------------------------------------------------"
}

log(){
    if [ $porcelain -eq 0 ]
    then
        echo "$1"
    fi
}

log_porcelain(){
    if [ $porcelain -eq 1 ]
    then
        : $1
    fi
}

startPoint(){
    line_break
    log "Stoping..."
    line_break
    log "HEAD -> $triggeredBranch"
    git reset -q --hard
    git checkout -q $triggeredBranch
}

trap "startPoint && exit 1;" SIGINT SIGQUIT

remote=`git config remote.origin.url | sed -e 's/.*git\/\(.*\)/\1/'`
triggeredBranch=`git symbolic-ref HEAD`
triggeredBranch=${triggeredBranch#refs/heads/}
octopusBranchName="octopus-features"
doPush=0
doAnalyse=0
doStopOnFailure=0
porcelain=0

for param in $@
do
  if [ $param = '--help' ] || [ $param = '-h' ] || [ $param = '-help' ] || [ $param = 'help' ]
  then
    usage
    exit 0
  else
    if [ $param = '--push' ]
    then
      doPush=1
    else
      if [ $param = '-a' ] || [ $param = '--analyse' ]
      then
        doAnalyse=1
      else
        if [ $param = '-s' ] || [ $param = '--stop-on-failure' ]
        then
          doStopOnFailure=1
        else
          if [ $param = '--porcelain' ]
          then
            porcelain=1
          else
            triggeredBranch=$param
          fi
        fi
      fi
    fi
  fi
done

line_break
log "Infos"
line_break
log "Directory : `pwd`"
log "Repository : $remote"
log "Branch : $triggeredBranch"
log "Commit : `git rev-parse HEAD`"
log "Push resulting merge : $doPush"
log "Analyse : $doAnalyse"
log "StopOnFailure : $doStopOnFailure"

line_break

if [[ -n `git status --porcelain` ]]
then
    git status
    log
    log "Le repo doit être clean"
    log
    line_break
    log "OCTOPUS FAILED"
    line_break
    exit 1
fi
log "Fetching ..."
line_break
#fetch de toutes les features
git fetch -p
line_break

features+="$triggeredBranch origin/master "
for branch in $(git for-each-ref --format="%(refname)" refs/remotes/origin/features)
do
    features+="$branch "
done

log "Branches à merger :"
log
BRANCHF="%-10s%-10s%s\n"
log "$( printf $BRANCHF id tag branch )"
for branch in $features
do
    log "$( printf "$BRANCHF" "`git rev-parse --short $branch`" "$( git describe --tags $branch | cut -d '-' -f 1 )" "${branch#refs/}" )"
done

line_break
log "Merge octopus"
line_break

git checkout -q --detach
git merge -q --no-edit $features

if [ $? -eq 0 ]
then
    if [ $doPush -eq 1 ]
    then
        line_break
        log "Push du merge sur octopus"
        line_break
        git push -f origin HEAD:$octopusBranchName
    else
        log
        log "merge créé : `git rev-parse HEAD`"
        log_porcelain "`git rev-parse HEAD`"
        git checkout -q $triggeredBranch
        line_break
    fi
    log "OCTOPUS SUCCESS"
    line_break
else
    line_break
    if [ $doStopOnFailure -eq 1 ]
    then
        log "OCTOPUS FAILED"
        log "   Stopped on failure, the merge is left unresolved."
        log "   To return to original state : git checkout -f $triggeredBranch"
        line_break
        exit 1
    fi
    if [ $doAnalyse -eq 1 ]
    then
        git checkout -q --detach -f $triggeredBranch
        git reset HEAD
        log "Recherche des branches en conflict avec $triggeredBranch..."
        line_break
        for branch in $features
        do
            if [ "$branch" != "$triggeredBranch" ] && [ "$branch" != "origin/master" ]
            then
                log "merge partiel :"
                log "`git rev-parse --short $triggeredBranch`  $triggeredBranch"
                log "`git rev-parse --short $branch`  ${branch#refs/remotes/}"
                log "`git rev-parse --short origin/master`  origin/master"
                log
                #bug de git ? il faut 3 branches en parametre pour que ca se passe bien
                git merge --no-commit -q $branch origin/master
                if [ $? -eq 0 ]
                then
                    log "##teamcity[message text='Merge ${triggeredBranch#refs/heads/} ${branch#refs/remotes/origin/} origin/master: OK' ]"
                    log "SUCCESS"
                else
                    git diff
                    log "##teamcity[message errorDetails='Merge ${triggeredBranch#refs/heads/} ${branch#refs/remotes/origin/} origin/master: FAILED' status='ERROR']"
                    log "FAILED"
                    conflicts+="$branch "
                fi
                git checkout -q --detach -f $triggeredBranch
                line_break
            fi
        done
        if [ -z "$conflicts" ]; then
            log "Aucun conflits trouvé ! regarde sur teamcity c'est peut être pas toi qui casse l'octopus ..."
        else
            log "Branches en conflits avec ${triggeredBranch#refs/}"
            for branch in $conflicts
            do
                log "  ${branch#refs/}"
            done
        fi
        line_break
    else
        log "use git octopus --analyse/-a pour détecter avec quelle(s) branche(s) se produit le conflit"
        line_break
    fi

    log "OCTOPUS FAILED"
    line_break

    git checkout -q -f $triggeredBranch
    git reset HEAD
    exit 1
fi
