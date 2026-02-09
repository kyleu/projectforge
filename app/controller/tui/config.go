package tui

import (
	"projectforge.dev/projectforge/app/doctor"
	"projectforge.dev/projectforge/app/project/action"
)

type Config struct {
	projectsLoading bool
	projectsErr     error
	projectKey      string

	projectActionRunning bool
	projectActionErr     error
	projectActionResults map[string]*action.Result

	doctorChecks  doctor.Checks
	doctorResults map[string]*doctor.Result
	doctorLoading bool
	doctorRunning bool
	doctorErr     error
}
