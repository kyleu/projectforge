package git

const (
	ok        = "OK"
	noUpdates = "no updates"
)

type Service struct{}

func NewService() *Service {
	return &Service{}
}
