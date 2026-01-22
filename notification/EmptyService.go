package notification

type EmptyService struct{}

var _ Service = (*EmptyService)(nil)

func NewEmptyService() *EmptyService {
	return &EmptyService{}
}

func (s *EmptyService) PushCreateNotification(key string, description string, flagURL string) error {
	return nil
}

func (s *EmptyService) PushDeleteNotification(key string, description string) error {
	return nil
}
