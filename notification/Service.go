package notification

type Service interface {
	PushCreateNotification(key string, description string, flagURL string) error
	PushDeleteNotification(key string, description string) error
}
