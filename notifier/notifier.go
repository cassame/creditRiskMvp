package notifier

import (
	"credit-risk-mvp/internal/domain"
	"log"
)

type Notifier interface {
	Notify(app domain.Application, status string) error
}
type LogNotifier struct{}

func (n LogNotifier) Notify(app domain.Application, status string) error {
	log.Printf("notify: phone=%s status=%s\n", app.Phone, status)
	return nil
}
