package email

import (
	"fmt"
	"notification/internal/models"
)

// MockEmailSender is a mock implementation of the Sender interface.
// Used for testing purposes.
type MockEmailSender struct{}

func NewMockEmailSender() Sender {
	return &MockEmailSender{}
}

func (s *MockEmailSender) Send(msg models.EmailMessage) error {
	fmt.Printf("Sending email:\n\tTo: %s\n\tSubject: %s\n\tMessage: %s\n", msg.To, msg.Subject, msg.Message)
	return nil
}
