package services_test

import (
	"context"
	"fmt"

	"jwt-auth/internal/domain/entities"
	"jwt-auth/internal/domain/repositories"
)

// Mock user repository
type mockUserRepository struct {
	users map[string]*entities.User
}

func newMockUserRepository() repositories.UserRepository {
	return &mockUserRepository{
		users: make(map[string]*entities.User),
	}
}

func (r *mockUserRepository) Create(ctx context.Context, user *entities.User) error {
	user.ID = len(r.users) + 1
	r.users[user.Email] = user
	return nil
}

func (r *mockUserRepository) GetByEmail(ctx context.Context, email string) (*entities.User, error) {
	if user, ok := r.users[email]; ok {
		return user, nil
	}
	return nil, fmt.Errorf("user not found")
}

func (r *mockUserRepository) GetByID(ctx context.Context, id int) (*entities.User, error) {
	for _, user := range r.users {
		if user.ID == id {
			return user, nil
		}
	}
	return nil, fmt.Errorf("user not found")
}

func (r *mockUserRepository) Update(ctx context.Context, user *entities.User) error {
	if existing, ok := r.users[user.Email]; ok {
		if existing.ID == user.ID {
			r.users[user.Email] = user
			return nil
		}
	}
	return fmt.Errorf("user not found")
}

func (r *mockUserRepository) Delete(ctx context.Context, id int) error {
	for email, user := range r.users {
		if user.ID == id {
			delete(r.users, email)
			return nil
		}
	}
	return fmt.Errorf("user not found")
}

// Mock email service
type mockEmailService struct{}

func newMockEmailService() *mockEmailService {
	return &mockEmailService{}
}

// Implement the EmailService interface
func (m *mockEmailService) SendEmail(ctx context.Context, to, subject, body string) error {
	// No-op/mock
	return nil
}
