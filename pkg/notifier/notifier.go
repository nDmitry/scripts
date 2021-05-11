package notifier

type Notifier interface {
	Notify(text string) error
}

type NotifierImpl struct {
	Notifier
}
