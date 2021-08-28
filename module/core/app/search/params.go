package search

type Params struct {
	Q string `json:"q"`
}

func (r *Params) String() string {
	return r.Q
}
