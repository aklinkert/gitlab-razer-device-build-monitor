package gitlab

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/xanzy/go-gitlab"
)

type usersClient interface {
	CurrentUser(options ...gitlab.OptionFunc) (*gitlab.User, *gitlab.Response, error)
}

// UserFetcher fetches the users details from the GitLab API
type UserFetcher struct {
	logger *logrus.Entry
	client usersClient
}

// NewUserFetcher returns a new UserFetcher instance
func NewUserFetcher(logger *logrus.Entry, client usersClient) (*UserFetcher, error) {
	return &UserFetcher{
		logger: logger,
		client: client,
	}, nil
}

// GetCurrentUserName fetches the current username from the GitLab API
func (u *UserFetcher) GetCurrentUserName() (string, error) {
	user, _, err := u.client.CurrentUser()
	if err != nil {
		return "", fmt.Errorf("failed to fetch currently authenticated user: %v", err)
	}

	return user.Username, nil
}
