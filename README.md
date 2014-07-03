git-octopus
===========

Scripts to work with a lot of branches without suffering the merge hell. Continuous Merge capacity for the build system.

I'm working on publishing our scripts to get feedback from other teams, if you can help for review before publishing I'll be glad. I'm having hard time to explain this simply, and will re-work it until it is good enough. Probably with a sample project.
So you're feedbacks are really welcome.

One Sentence description :
* By continuously merging features in progress, we allow each development to be simply based on master with early detection of conflicts with other features, and many options for solving it. 

The minimal to know :
* Any new feature is a branch named feature/XXX-xxx_name where XXX-xxx is a Jira ID
* Oneline  :
     1/ take the list of branches named "features/xxx" => <BRANCH_LIST>
     2/ run "git merge -s octopus master <BRANCH_LIST>" (built-in in git !)
     3/ if it's ok, then you know that all the branches are compatible, if not, then you have the conflicting files, and probably the branche causing the issue. The octopus scripting we're using is mainly there to add some tooling to find the branches in conflicts.
* the you are safe to work totally independently each feature, you'll be warned if problems will occur later.

From the developer side :
* Create a feature branch
* Develop on it
* Push stable steps on the "origin"
* The CI will merge all branches together in a "features-octopus" tag and warn if conflicts appears
* The this tag can be compiled and deployed on a staging server to test like with a "develop" branch
* When the feature is finished, the branch can be safely merge in the master (using a release branch or not)


This allows to :
* Detect conflict between branches before merge (coders can run the script locally as pre-commit)
* Have a virtual "develop" branch (in git-flow fashion) easy to repare
* Ensure each branch is compatible with others
* Merge only branches that are ready for production
* Mave many strategies to solve conflicts (as no manual merge as been pushed)

One more thing:
* This enforces really INVEST tasks, and empowers the good pratices of Kanban !

Solving Conflicts: 
The whole idea (and discovery) of this branching model is to detect branches conflicts and leave many options for solving, and where manual merge is the last recommended solution (because it links the 2 features together).
A. EASY and SAFE resolutions:
1/ When conflict is due to adding code at the same place, move insertion point (helps git to resolve)
2/ Remove this branch from the server, or rollback your commit (easy rollback, with no impact on other developpers)
3/ Cherry-pick a small refactoring from one branch to another (mikado strategy)
4/ Merge the 2 branches (if really linked), and ship the result together
... (I'll probably add others there)

B. ADVANCED resolutions
1/ Rebase a branch on the conflicting one, so you can manually resolve the conflit, but then 1 branch should be released before the other one : THEY ARE NOW LINKED
2/ Push a commit on the master (mikado on the master)
... (I'll probably add others there)

Mainly here is what I plan to publish :
* octopus.sh : the merge script + --tag option (to create the tag) + --analyse option (to find the conflicts)

I have a lot more ideas explaining how it helps us, but this is what I guess is the minimum viable knownledge.

After I'll have:
* octopus-check.sh : the list of feature branches and numbers of commits in the branch + list of issues involved (helps to have a view on current status).
* octopus-remote-run.sh : a script to locally build the branch, perform an octopus-merge and build the result, to avoir pushing bad commits (golden octopus !)
