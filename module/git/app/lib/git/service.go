package git

const (
	ok        = "OK"
	noUpdates = "no updates"
)

type Service struct {
	Key  string `json:"key"`
	Path string `json:"path,omitempty"`
}

func NewService(key string, path string) *Service {
	return &Service{Key: key, Path: path}
}
