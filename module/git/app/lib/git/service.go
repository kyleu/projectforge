package git

const (
	noUpdates = "no updates"
)

type Service struct {
	Key  string `json:"key"`
	Path string `json:"path,omitzero"`
}

func NewService(key string, path string) *Service {
	return &Service{Key: key, Path: path}
}
