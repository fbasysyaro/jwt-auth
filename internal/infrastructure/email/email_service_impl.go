package email

import (
	"context"
	"fmt"
	"jwt-auth/internal/domain/services"
)

type emailServiceImpl struct{}

func NewEmailService() services.EmailService {
	return &emailServiceImpl{}
}

func (e *emailServiceImpl) SendEmail(ctx context.Context, to, subject, body string) error {
	fmt.Printf("Sending email to %s: %s - %s\n", to, subject, body)
	return nil
}