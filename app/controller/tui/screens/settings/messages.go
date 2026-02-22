package settings

import "projectforge.dev/projectforge/app/lib/task"

type errMsg struct{ err error }

type taskRunMsg struct {
	result *task.Result
	err    error
}
