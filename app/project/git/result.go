package git

import (
	"strconv"

	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/project"
	"projectforge.dev/projectforge/app/util"
)

var ResultFields = []string{"Project", "Status", "Data", "Error"}

type Result struct {
	Project *project.Project `json:"-"`
	Status  string           `json:"status"`
	Data    util.ValueMap    `json:"data"`
	Error   string           `json:"error"`
}

func NewResult(prj *project.Project, status string, data util.ValueMap) *Result {
	return &Result{Project: prj, Status: status, Data: data}
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

func (r *Result) History() *HistoryResult {
	if r.Data == nil {
		return nil
	}
	if x, ok := r.Data["history"]; ok {
		h := &HistoryResult{}
		_ = util.CycleJSON(x, h)
		return h
	}
	return nil
}

func (r *Result) CleanData() util.ValueMap {
	return lo.OmitByKeys(r.Data, []string{"branch", "dirty", "status"})
}

func (r *Result) DataString(k string) string {
	if r.Data == nil {
		return ""
	}
	return r.Data.GetStringOpt(k)
}

func (r *Result) DataInt(k string) int {
	ret, _ := strconv.ParseInt(r.DataString(k), 10, 32)
	return int(ret)
}

func (r *Result) DataStringArray(k string) []string {
	if r.Data == nil {
		return nil
	}
	ret, _ := r.Data.GetStringArray(k, true)
	return ret
}

func (r *Result) Strings() []string {
	return []string{r.Project.Key, r.Status, r.Data.String(), r.Error}
}

func (r *Result) ToCSV() ([]string, [][]string) {
	return ResultFields, [][]string{r.Strings()}
}

type Results []*Result

func (r Results) Get(key string) *Result {
	return lo.FindOrElse(r, nil, func(x *Result) bool {
		return x.Project.Key == key
	})
}

func (r Results) ToCSV() ([]string, [][]string) {
	return ResultFields, lo.Map(r, func(x *Result, _ int) []string {
		return x.Strings()
	})
}
