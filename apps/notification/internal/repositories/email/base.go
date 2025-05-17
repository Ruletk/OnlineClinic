package email

import "notification/internal/models"

// Sender interface defines the methods for sending emails.
type Sender interface {
	Send(message models.EmailMessage) error
}
