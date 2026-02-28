package grep

import (
	"path/filepath"
	"strconv"
	"strings"

	"{{{ .Package }}}/app/lib/search/result"
)

type Match struct {
	File    string `json:"file"`
	Offset  int    `json:"offset"`
	LineNum int    `json:"lineNum"`
	Text    string `json:"text"`
	Match   string `json:"match"`
	Start   int    `json:"start"`
	End     int    `json:"end"`
}

type Matches []*Match

type Response struct {
	Matches       Matches  `json:"matches"`
	Request       *Request `json:"request,omitzero"`
	BytesSearched int      `json:"bytesSearched,omitzero"`
	ElapsedNanos  int      `json:"elapsedNanos"`
	Debug         any      `json:"debug,omitzero"`
}

type ToSearchResultOptions struct {
	Type           string
	ID             string
	Title          string
	URL            string
	Icon           string
	MatchKeyPrefix string
	TermsLower     []string
	Limit          int
}

func (r *Response) ToSearchResult(opts *ToSearchResultOptions) *result.Result {
	if r == nil {
		return nil
	}
	if len(r.Matches) == 0 {
		return nil
	}
	if opts.Limit <= 0 {
		opts.Limit = 50
	}

	matches := result.Matches{}
	seen := map[string]bool{}
	for _, m := range r.Matches {
		if m == nil {
			continue
		}
		line := strings.TrimSpace(strings.TrimRight(m.Text, "\n"))
		if len(opts.TermsLower) > 1 {
			ll := strings.ToLower(line)
			ok := true
			for _, t := range opts.TermsLower[1:] {
				if t == "" {
					continue
				}
				if !strings.Contains(ll, t) {
					ok = false
					break
				}
			}
			if !ok {
				continue
			}
		}

		k := filepath.ToSlash(strings.TrimSpace(m.File))
		if m.LineNum > 0 {
			k = k + ":" + strconv.Itoa(m.LineNum)
		}
		if k == "" {
			continue
		}
		k = opts.MatchKeyPrefix + k
		if seen[k] {
			continue
		}
		seen[k] = true

		matches = append(matches, &result.Match{Key: k, Value: line})
		if len(matches) >= opts.Limit {
			break
		}
	}
	if len(matches) == 0 {
		return nil
	}
	matches.Sort()

	return &result.Result{Type: opts.Type, ID: opts.ID, Title: opts.Title, URL: opts.URL, Icon: opts.Icon, Matches: matches}
}
