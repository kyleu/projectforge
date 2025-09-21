package git

import (
	"strconv"

	"github.com/samber/lo"

	"{{{ .Package }}}/app/util"
)

var ResultFields = []string{"Project", "Status", "Data", "Error"}

type Result struct {
	Project string        `json:"project,omitzero"`
	Status  string        `json:"status,omitzero"`
	Data    util.ValueMap `json:"data,omitzero"`
	Error   string        `json:"error,omitzero"`
}

func NewResult(prj string, status string, data util.ValueMap) *Result {
	return &Result{Project: prj, Status: status, Data: data}
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

func (r *Result) String() string {
	ret := r.Status
	if branch := r.Data.GetStringOpt("branch"); branch != "" {
		ret += "[" + branch + "]"
	}
	ahead := r.Data.GetStringOpt("commitsAhead")
	behind := r.Data.GetStringOpt("commitsBehind")
	if ahead != "" || behind != "" {
		ret += " ("
		if ahead != "" {
			ret += ahead + " ahead"
		}
		if behind != "" {
			if ahead != "" {
				ret += ", "
			}
			ret += behind + " behind"
		}
		ret += ")"
	}
	return ret
}

func (r *Result) Strings() []string {
	return []string{r.Project, r.Status, r.Data.String(), r.Error}
}

func (r *Result) ToCSV() ([]string, [][]string) {
	return ResultFields, [][]string{r.Strings()}
}

type Results []*Result

func (r Results) Get(key string) *Result {
	return lo.FindOrElse(r, nil, func(x *Result) bool {
		return x.Project == key
	})
}

func (r Results) ToCSV() ([]string, [][]string) {
	return ResultFields, lo.Map(r, func(x *Result, _ int) []string {
		return x.Strings()
	})
}
