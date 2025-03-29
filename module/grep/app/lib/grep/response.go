package grep

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
	Request       *Request `json:"request,omitempty"`
	BytesSearched int      `json:"bytesSearched,omitempty"`
	ElapsedNanos  int      `json:"elapsedNanos"`
	Debug         any      `json:"debug,omitempty"`
}
