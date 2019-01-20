package monitor

import (
	"fmt"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/apinnecke/gitlab-razer-device-build-monitor/pkg/config"
	"github.com/xanzy/go-gitlab"
)

var (
	addIncludeSubgroups = gitlab.OptionFunc(func(req *http.Request) error {
		v := req.URL.Query()
		v.Add("include_subgroups", "true")
		req.URL.RawQuery = v.Encode()
		return nil
	})
)

type fetcherClient interface {
}

// RepoFetcher fetches a list of repositories from GitLab
type RepoFetcher struct {
	logger *logrus.Entry
	client *gitlab.Client
	config *config.Config
}

// NewRepoFetcher returns a new RepoFetcher instance
func NewRepoFetcher(logger *logrus.Entry, client *gitlab.Client, config *config.Config) (*RepoFetcher, error) {
	return &RepoFetcher{
		logger: logger,
		client: client,
		config: config,
	}, nil
}

func (f *RepoFetcher) getCurrentUserName() (string, error) {
	user, _, err := f.client.Users.CurrentUser()
	if err != nil {
		return "", fmt.Errorf("failed to fetch currently authenticated user: %v", err)
	}

	return user.Username, nil
}

// Fetch fetches a list of accessible repos within the groups set in config file
func (f *RepoFetcher) Fetch() ([]string, error) {
	var repos []string

	opt := &gitlab.ListGroupProjectsOptions{
		MinAccessLevel: gitlab.AccessLevel(gitlab.DeveloperPermissions),
		ListOptions: gitlab.ListOptions{
			PerPage: 100,
		},
	}

	for _, group := range f.config.Groups {
		projects, _, err := f.client.Groups.ListGroupProjects(group, opt, addIncludeSubgroups)
		if err != nil {
			return []string{}, fmt.Errorf("failed to fetch GitLab projects or group %q: %v", group, err)
		}

		for _, p := range projects {
			repos = append(repos, p.PathWithNamespace)
		}

	}

	return repos, nil
}
