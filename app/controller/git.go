package controller

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"
	"projectforge.dev/app/git"
	"projectforge.dev/views/verror"
	"projectforge.dev/views/vgit"

	"projectforge.dev/app"
	"projectforge.dev/app/controller/cutil"
)

var (
	gitCommitArgs = cutil.Args{{Key: "message", Title: "Message", Description: "The message to used for the commit"}}
	gitBranchArgs = cutil.Args{{Key: "name", Title: "Branch Name", Description: "The name to used for the new branch"}}
)

func GitAction(rc *fasthttp.RequestCtx) {
	a, _ := RCRequiredString(rc, "act", false)
	act("git.action."+a, rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		key, err := RCRequiredString(rc, "key", false)
		if err != nil {
			return "", err
		}
		prj, err := as.Services.Projects.Get(key)
		if err != nil {
			return "", errors.Wrap(err, "unable to load project")
		}
		var result *git.Result
		switch a {
		case git.ActionStatus.Key, "":
			// runs by default
		case git.ActionCreateRepo.Key:
			result, err = as.Services.Git.CreateRepo(prj)
		case git.ActionMagic.Key:
			result, err = as.Services.Git.Magic(prj)
		case git.ActionFetch.Key:
			result, err = as.Services.Git.Fetch(prj)
		case git.ActionCommit.Key:
			argRes := cutil.CollectArgs(rc, gitCommitArgs)
			if len(argRes.Missing) > 0 {
				url := fmt.Sprintf("/git/%s/commit", prj.Key)
				ps.Data = argRes
				return render(rc, as, &verror.Args{URL: url, Directions: "Enter your commit message", ArgRes: argRes}, ps, "projects", prj.Key, "Git")
			}
			result, err = as.Services.Git.Commit(prj, argRes.Values["message"])
		case git.ActionBranch.Key:
			argRes := cutil.CollectArgs(rc, gitBranchArgs)
			if len(argRes.Missing) > 0 {
				url := fmt.Sprintf("/git/%s/branch", prj.Key)
				ps.Data = argRes
				return render(rc, as, &verror.Args{URL: url, Directions: "Enter your branch name", ArgRes: argRes}, ps, "projects", prj.Key, "Git")
			}
			result, err = as.Services.Git.Commit(prj, argRes.Values["message"])
		default:
			err = errors.Errorf("unhandled action [%s]", a)
		}
		if err != nil {
			return "", err
		}
		statusResult, _ := as.Services.Git.Status(prj)
		if result == nil {
			result = statusResult
		} else {
			result.Data = result.Data.Merge(statusResult.Data)
		}
		ps.Data = result
		return render(rc, as, &vgit.Result{Result: result}, ps, "projects", prj.Key, "Git")
	})
}
