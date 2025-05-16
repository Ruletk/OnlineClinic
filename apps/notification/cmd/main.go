package main

import (
	"fmt"
	"notification/internal/repositories/email"
)

func main() {
	emailSender := email.NewMockEmailSender()

	fmt.Println("Hello World!", emailSender)
}
