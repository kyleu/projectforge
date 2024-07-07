package git

import "github.com/samber/lo"

type Action struct {
	Key         string `json:"key"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type Actions []*Action

var (
	ActionStatus     = &Action{Key: "status", Title: "Status", Description: "Returns the status of the repository"}
	ActionCreateRepo = &Action{Key: "createrepo", Title: "Create Repo", Description: "Creates a git repo and makes an initial commit"}
	ActionMagic      = &Action{Key: "magic", Title: "Magic", Description: "Does everything it can to bring the repo up to date (stash, pull, pop, commit, push)"}
	ActionFetch      = &Action{Key: "fetch", Title: "Fetch", Description: "Fetches the latest changes from the repository"}
	ActionCommit     = &Action{Key: "commit", Title: "Commit", Description: "Adds all files, commits with the provided message"}
	ActionPull       = &Action{Key: "pull", Title: "Pull", Description: "Pulls pending commits from upstream"}
	ActionPush       = &Action{Key: "push", Title: "Push", Description: "Pushes pending commits to the remote"}
	ActionReset      = &Action{Key: "reset", Title: "Reset", Description: "Resets all local changes; be careful"}
	ActionBranch     = &Action{Key: "branch", Title: "Branch", Description: "Switch to a new branch"}
	ActionUndoCommit = &Action{Key: "undocommit", Title: "Undo", Description: "Removes the most recent commit, keeping all local changes"}
	ActionOutdated   = &Action{Key: "outdated", Title: "Outdated", Description: "Finds commits since last tag"}
	ActionHistory    = &Action{Key: "history", Title: "History", Description: "Visualize the git history"}

	allActions = Actions{
		ActionStatus, ActionCreateRepo, ActionMagic, ActionFetch, ActionCommit,
		ActionPull, ActionPush, ActionReset, ActionBranch, ActionUndoCommit, ActionOutdated, ActionHistory,
	}
)

func ActionStatusFromString(key string) *Action {
	if key == "" {
		return ActionStatus
	}
	return lo.FindOrElse(allActions, nil, func(act *Action) bool {
		return act.Key == key
	})
}

func (r *Result) Actions() Actions {
	ret := Actions{ActionStatus}
	if r.Status == "no repo" {
		return append(ret, ActionCreateRepo)
	}
	ret = append(ret, ActionFetch)
	if dirty := r.DataStringArray("dirty"); len(dirty) > 0 {
		ret = append(ret, ActionCommit)
	}
	if r.DataInt("commitsAhead") > 0 {
		ret = append(ret, ActionPush)
	}
	if r.DataInt("commitsBehind") > 0 {
		ret = append(ret, ActionPull)
	}
	ret = append(ret, ActionReset, ActionUndoCommit, ActionMagic, ActionHistory)
	return ret
}
