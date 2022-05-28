package controller

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/controller/cutil"
	"projectforge.dev/projectforge/app/git"
	"projectforge.dev/projectforge/views/verror"
	"projectforge.dev/projectforge/views/vgit"
)

var (
	messageArg = &cutil.Arg{Key: "message", Title: "Message", Description: "The message to used for the commit"}
	dryRunArg  = &cutil.Arg{Key: "dryRun", Title: "Dry Run", Description: "Runs without any destructive operations", Type: "bool", Default: "true"}

	gitBranchArgs = cutil.Args{{Key: "name", Title: "Branch Name", Description: "The name to used for the new branch"}}
	gitCommitArgs = cutil.Args{messageArg}
	gitMagicArgs  = cutil.Args{messageArg, dryRunArg}
)

func GitAction(rc *fasthttp.RequestCtx) {
	a, _ := cutil.RCRequiredString(rc, "act", false)
	act("git.action."+a, rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		key, err := cutil.RCRequiredString(rc, "key", false)
		if err != nil {
			return "", err
		}
		prj, err := as.Services.Projects.Get(key)
		if err != nil {
			return "", errors.Wrap(err, "unable to load project")
		}
		var bc = []string{"projects", prj.Key, "Git"}
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
			if len(argRes.Missing) > 0 {
				ps.Data = argRes
				url := fmt.Sprintf("/git/%s/commit", prj.Key)
				return render(rc, as, &verror.Args{URL: url, Directions: "Enter your commit message", ArgRes: argRes}, ps, bc...)
			}
			result, err = as.Services.Git.Commit(ps.Context, prj, argRes.Values["message"], ps.Logger)
		case git.ActionUndoCommit.Key:
			result, err = as.Services.Git.UndoCommit(ps.Context, prj, ps.Logger)
		case git.ActionPull.Key:
			result, err = as.Services.Git.Pull(ps.Context, prj, ps.Logger)
		case git.ActionPush.Key:
			result, err = as.Services.Git.Push(ps.Context, prj, ps.Logger)
		case git.ActionBranch.Key:
			argRes := cutil.CollectArgs(rc, gitBranchArgs)
			if len(argRes.Missing) > 0 {
				url := fmt.Sprintf("/git/%s/branch", prj.Key)
				ps.Data = argRes
				return render(rc, as, &verror.Args{URL: url, Directions: "Enter your branch name", ArgRes: argRes}, ps, bc...)
			}
			result, err = as.Services.Git.Commit(ps.Context, prj, argRes.Values["message"], ps.Logger)
		case git.ActionMagic.Key:
			argRes := cutil.CollectArgs(rc, gitMagicArgs)
			if len(argRes.Missing) > 0 {
				url := fmt.Sprintf("/git/%s/magic", prj.Key)
				ps.Data = argRes
				return render(rc, as, &verror.Args{URL: url, Directions: "Enter your commit message", ArgRes: argRes}, ps, bc...)
			}
			dryRun := argRes.Values["dryRun"] == "true"
			result, err = as.Services.Git.Magic(ps.Context, prj, argRes.Values["message"], dryRun, ps.Logger)
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
		return render(rc, as, &vgit.Result{Action: actn, Result: result}, ps, bc...)
	})
}
