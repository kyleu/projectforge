package tui

import (
	"projectforge.dev/projectforge/app/doctor"
)

type Config struct {
	projectsLoading bool
	projectsErr     error

	doctorChecks  doctor.Checks
	doctorResults map[string]*doctor.Result
	doctorLoading bool
	doctorRunning bool
	doctorErr     error
}
