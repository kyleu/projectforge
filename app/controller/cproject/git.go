package cproject

import (
	"fmt"
	"strconv"

	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/controller"
	"projectforge.dev/projectforge/app/controller/cutil"
	"projectforge.dev/projectforge/app/project/git"
	"projectforge.dev/projectforge/app/util"
	"projectforge.dev/projectforge/views/vgit"
	"projectforge.dev/projectforge/views/vpage"
)

var (
	authorsArg = &cutil.Arg{Key: "authors", Title: "Authors", Description: "Limit to a set of author emails (comma-separated)"}
	branchArg  = &cutil.Arg{Key: "name", Title: "Branch Name", Description: "The name to used for the new branch"}
	dryRunArg  = &cutil.Arg{Key: "dryRun", Title: "Dry Run", Description: "Runs without any destructive operations", Type: "bool", Default: util.BoolTrue}
	limitArg   = &cutil.Arg{Key: "limit", Title: "Limit", Description: "Limits the results to, at most, this amount", Type: "number", Default: "100"}
	messageArg = &cutil.Arg{Key: "message", Title: "Message", Description: "The message to used for the commit"}
	pathArg    = &cutil.Arg{Key: "path", Title: "Path", Description: "Limits the results to the provided path (leave blank for all)"}
	sinceArg   = &cutil.Arg{Key: "since", Title: "Since", Description: "Limit to a date range", Type: "datetime"}

	gitBranchArgs  = cutil.Args{branchArg}
	gitCommitArgs  = cutil.Args{messageArg}
	gitHistoryArgs = cutil.Args{pathArg, sinceArg, authorsArg, limitArg}
	gitMagicArgs   = cutil.Args{messageArg, dryRunArg}
)

func GitAction(rc *fasthttp.RequestCtx) {
	a, _ := cutil.RCRequiredString(rc, "act", false)
	controller.Act("git.action."+a, rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		key, err := cutil.RCRequiredString(rc, "key", false)
		if err != nil {
			return "", err
		}
		prj, err := as.Services.Projects.Get(key)
		if err != nil {
			return "", errors.Wrap(err, "unable to load project")
		}
		bc := []string{"projects", prj.Key, "Git"}
		var result *git.Result
		actn := git.ActionStatusFromString(a)
		switch a {
		case git.ActionStatus.Key, "":
			// runs by default
		case git.ActionCreateRepo.Key:
			result, err = as.Services.Git.CreateRepo(ps.Context, prj, ps.Logger)
		case git.ActionFetch.Key:
			result, err = as.Services.Git.Fetch(ps.Context, prj, ps.Logger)
		case git.ActionCommit.Key:
			argRes := cutil.CollectArgs(rc, gitCommitArgs)
			if argRes.HasMissing() {
				ps.Data = argRes
				url := fmt.Sprintf("/git/%s/commit", prj.Key)
				return controller.Render(rc, as, &vpage.Args{URL: url, Directions: "Enter your commit message", ArgRes: argRes}, ps, bc...)
			}
			result, err = as.Services.Git.Commit(ps.Context, prj, argRes.Values.GetStringOpt("message"), ps.Logger)
		case git.ActionUndoCommit.Key:
			result, err = as.Services.Git.UndoCommit(ps.Context, prj, ps.Logger)
		case git.ActionPull.Key:
			result, err = as.Services.Git.Pull(ps.Context, prj, ps.Logger)
		case git.ActionPush.Key:
			result, err = as.Services.Git.Push(ps.Context, prj, ps.Logger)
		case git.ActionOutdated.Key:
			result, err = as.Services.Git.Outdated(ps.Context, prj, ps.Logger)
		case git.ActionHistory.Key:
			argRes := cutil.CollectArgs(rc, gitHistoryArgs)
			if argRes.HasMissing() {
				ps.Data = argRes
				url := fmt.Sprintf("/git/%s/history", prj.Key)
				return controller.Render(rc, as, &vpage.Args{URL: url, Directions: "Choose your options", ArgRes: argRes}, ps, bc...)
			}
			path := argRes.Values.GetStringOpt("paths")
			since, _ := util.TimeFromString(argRes.Values.GetStringOpt("since"))
			authors := util.StringSplitAndTrim(argRes.Values.GetStringOpt("authors"), ",")
			commit := string(rc.URI().QueryArgs().Peek("commit"))
			limit, _ := strconv.ParseInt(argRes.Values.GetStringOpt("limit"), 10, 32)
			hist := &git.HistoryResult{Path: path, Since: since, Authors: authors, Commit: commit, Limit: int(limit)}
			result, err = as.Services.Git.History(ps.Context, prj, hist, ps.Logger)
		case git.ActionBranch.Key:
			argRes := cutil.CollectArgs(rc, gitBranchArgs)
			if argRes.HasMissing() {
				url := fmt.Sprintf("/git/%s/branch", prj.Key)
				ps.Data = argRes
				return controller.Render(rc, as, &vpage.Args{URL: url, Directions: "Enter your branch name", ArgRes: argRes}, ps, bc...)
			}
			result, err = as.Services.Git.Commit(ps.Context, prj, argRes.Values.GetStringOpt("message"), ps.Logger)
		case git.ActionMagic.Key:
			argRes := cutil.CollectArgs(rc, gitMagicArgs)
			if argRes.HasMissing() {
				url := fmt.Sprintf("/git/%s/magic", prj.Key)
				ps.Data = argRes
				return controller.Render(rc, as, &vpage.Args{URL: url, Directions: "Enter your commit message", ArgRes: argRes}, ps, bc...)
			}
			dryRun := argRes.Values.GetStringOpt("dryRun") == util.BoolTrue
			result, err = as.Services.Git.Magic(ps.Context, prj, argRes.Values.GetStringOpt("message"), dryRun, ps.Logger)
		default:
			err = errors.Errorf("unhandled action [%s]", a)
		}
		if err != nil {
			return "", err
		}
		statusResult, err := as.Services.Git.Status(ps.Context, prj, ps.Logger)
		if err != nil {
			return "", errors.Wrap(err, "unable to pull repo status")
		}

		if result == nil {
			result = statusResult
		} else {
			result.Data = result.Data.Merge(statusResult.Data)
		}
		ps.Data = result
		return controller.Render(rc, as, &vgit.Result{Action: actn, Result: result}, ps, bc...)
	})
}
