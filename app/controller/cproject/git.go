package cproject

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/controller"
	"projectforge.dev/projectforge/app/controller/cutil"
	"projectforge.dev/projectforge/app/lib/git"
	"projectforge.dev/projectforge/app/util"
	"projectforge.dev/projectforge/views/vgit"
	"projectforge.dev/projectforge/views/vpage"
)

var (
	authorsArg = &util.FieldDesc{Key: "authors", Title: "Authors", Description: "Limit to a set of author emails (comma-separated)"}
	branchArg  = &util.FieldDesc{Key: "name", Title: "Branch Name", Description: "The name to used for the new branch"}
	cmdArg     = &util.FieldDesc{Key: "cmd", Title: "Command", Description: "The command to run"}
	dirArg     = &util.FieldDesc{Key: "dir", Title: "Directory", Description: "The directory to run the command in"}
	dryRunArg  = &util.FieldDesc{Key: "dryRun", Title: "Dry Run", Description: "Runs without any destructive operations", Type: "bool", Default: util.BoolTrue}
	limitArg   = &util.FieldDesc{Key: "limit", Title: "Limit", Description: "Limits the results to, at most, this amount", Type: "number", Default: "100"}
	messageArg = &util.FieldDesc{Key: "message", Title: "Message", Description: "The message to used for the commit"}
	pathArg    = &util.FieldDesc{Key: "path", Title: "Path", Description: "Limits the results to the provided path (leave blank for all)"}
	sinceArg   = &util.FieldDesc{Key: "since", Title: "Since", Description: "Limit to a date range", Type: "datetime"}

	gitBranchArgs  = util.FieldDescs{branchArg}
	gitCommitArgs  = util.FieldDescs{messageArg}
	gitHistoryArgs = util.FieldDescs{pathArg, sinceArg, authorsArg, limitArg}
	gitMagicArgs   = util.FieldDescs{messageArg, dryRunArg}
	gitCustomArgs  = util.FieldDescs{cmdArg, dirArg}
)

func GitAction(w http.ResponseWriter, r *http.Request) {
	a, _ := cutil.PathString(r, "act", false)
	controller.Act("git.action."+a, w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
		key, err := cutil.PathString(r, "key", false)
		if err != nil {
			return "", err
		}
		prj, err := as.Services.Projects.Get(key)
		if err != nil {
			return "", errors.Wrap(err, "unable to load project")
		}
		bc := []string{"projects", prj.Key, "Git**git"}
		var result *git.Result
		actn := git.ActionStatusFromString(a)
		gs := git.NewService(prj.Key, prj.Path)
		switch a {
		case git.ActionStatus.Key, "":
			// runs by default
		case git.ActionCreateRepo.Key:
			result, err = gs.CreateRepo(ps.Context, ps.Logger)
		case git.ActionFetch.Key:
			result, err = gs.Fetch(ps.Context, ps.Logger)
		case git.ActionCommit.Key:
			argRes := util.FieldDescsCollect(r, gitCommitArgs)
			if argRes.HasMissing() {
				ps.SetTitleAndData("Git Commit", argRes)
				url := fmt.Sprintf("/git/%s/commit", prj.Key)
				return controller.Render(r, as, &vpage.Args{URL: url, Directions: "Enter your commit message", Results: argRes}, ps, bc...)
			}
			result, err = gs.Commit(ps.Context, argRes.Values.GetStringOpt("message"), ps.Logger)
		case git.ActionUndoCommit.Key:
			result, err = gs.UndoCommit(ps.Context, ps.Logger)
		case git.ActionPull.Key:
			result, err = gs.Pull(ps.Context, ps.Logger)
		case git.ActionPush.Key:
			result, err = gs.Push(ps.Context, ps.Logger)
		case git.ActionReset.Key:
			result, err = gs.Reset(ps.Context, ps.Logger)
		case git.ActionOutdated.Key:
			result, err = gs.Outdated(ps.Context, ps.Logger)
		case git.ActionHistory.Key:
			argRes := util.FieldDescsCollect(r, gitHistoryArgs)
			if argRes.HasMissing() {
				ps.SetTitleAndData("Git History", argRes)
				url := fmt.Sprintf("/git/%s/history", prj.Key)
				return controller.Render(r, as, &vpage.Args{URL: url, Directions: "Choose your options", Results: argRes}, ps, bc...)
			}
			path := argRes.Values.GetStringOpt("paths")
			since, _ := util.TimeFromString(argRes.Values.GetStringOpt("since"))
			authors := util.StringSplitAndTrim(argRes.Values.GetStringOpt("authors"), ",")
			commit := cutil.QueryStringString(r, "commit")
			limit, _ := strconv.ParseInt(argRes.Values.GetStringOpt("limit"), 10, 32)
			args := &git.HistoryArgs{Path: path, Since: since, Authors: authors, Commit: commit, Limit: int(limit)}
			result, err = gs.History(ps.Context, args, ps.Logger)
		case git.ActionBranch.Key:
			argRes := util.FieldDescsCollect(r, gitBranchArgs)
			if argRes.HasMissing() {
				url := fmt.Sprintf("/git/%s/branch", prj.Key)
				ps.SetTitleAndData("Git Branch", argRes)
				return controller.Render(r, as, &vpage.Args{URL: url, Directions: "Enter your branch name", Results: argRes}, ps, bc...)
			}
			result, err = gs.Commit(ps.Context, argRes.Values.GetStringOpt("message"), ps.Logger)
		case git.ActionMagic.Key:
			argRes := util.FieldDescsCollect(r, gitMagicArgs)
			if argRes.HasMissing() {
				url := fmt.Sprintf("/git/%s/magic", prj.Key)
				ps.SetTitleAndData("Git Magic!", argRes)
				warning := "Are you sure you'd like to commit and push?"
				page := &vpage.Args{URL: url, Directions: "Enter your commit message", Results: argRes, Warning: warning}
				return controller.Render(r, as, page, ps, bc...)
			}
			dryRun := argRes.Values.GetStringOpt("dryRun") == util.BoolTrue
			result, err = gs.Magic(ps.Context, argRes.Values.GetStringOpt("message"), dryRun, ps.Logger)
		default:
			err = errors.Errorf("unhandled action [%s]", a)
		}
		if err != nil {
			return "", err
		}
		statusResult, err := gs.Status(ps.Context, ps.Logger)
		if err != nil {
			return "", errors.Wrap(err, "unable to pull repo status")
		}

		if result == nil {
			result = statusResult
		} else {
			result.Data = result.Data.Merge(statusResult.Data)
		}
		ps.Data = result
		page := &vgit.Result{Action: actn, Result: result, URL: prj.WebPath(), Title: prj.Title(), Icon: prj.IconSafe()}
		return controller.Render(r, as, page, ps, bc...)
	})
}
