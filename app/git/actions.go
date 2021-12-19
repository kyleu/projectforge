package git

type Action struct {
	Key         string `json:"key"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type Actions []*Action

var (
	ActionStatus     = &Action{Key: "status", Title: "Status", Description: "Returns the status of the repository"}
	ActionCreateRepo = &Action{Key: "createrepo", Title: "Create Repo", Description: "Creates a git repo and makes an initial commit"}
	ActionMagic      = &Action{Key: "magic", Title: "Magic", Description: "Does everything it can to bring the repo up to date"}
	ActionCommit     = &Action{Key: "commit", Title: "Commit", Description: "Adds all files, commits with the provided message"}

	AllActions = Actions{ActionStatus, ActionCreateRepo, ActionMagic, ActionCommit}
)

