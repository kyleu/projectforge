package project

type Info struct {
	Org             string `json:"org,omitempty"`
	AuthorName      string `json:"authorName,omitempty"`
	AuthorEmail     string `json:"authorEmail,omitempty"`
	License         string `json:"license,omitempty"`
	Bundle          string `json:"bundle,omitempty"`
	SigningIdentity string `json:"signingIdentity,omitempty"`
	Homepage        string `json:"homepage,omitempty"`
	Sourcecode      string `json:"sourcecode,omitempty"`
	Summary         string `json:"summary,omitempty"`
	Description     string `json:"description,omitempty"`
}
