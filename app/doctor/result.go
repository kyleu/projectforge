package doctor

type Result struct {
	Check    *Check   `json:"-"`
	Key      string   `json:"key"`
	Title    string   `json:"title"`
	Status   string   `json:"status,omitempty"`
	Summary  string   `json:"summary,omitempty"`
	Errors   Errors   `json:"errors,omitempty"`
	Duration int      `json:"duration,omitempty"`
	Solution []string `json:"solution,omitempty"`
	Logs     []string `json:"logs,omitempty"`
}

func NewResult(check *Check, key string, title string, summary string) *Result {
	return &Result{Check: check, Key: key, Title: title, Status: "OK", Summary: summary}
}

func (p *Result) AddLog(msg string) *Result {
	p.Logs = append(p.Logs, msg)
	return p
}

func (p *Result) WithError(err *Error) *Result {
	p.Status = "error"
	p.Errors = append(p.Errors, err)
	return p
}

func (p *Result) AddSolution(msg string) *Result {
	p.Solution = append(p.Solution, msg)
	return p
}

type Results []*Result
