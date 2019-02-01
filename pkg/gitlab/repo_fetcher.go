package gitlab

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

type groupsClient interface {
	ListGroupProjects(gid interface{}, opt *gitlab.ListGroupProjectsOptions, options ...gitlab.OptionFunc) ([]*gitlab.Project, *gitlab.Response, error)
}

// RepoFetcher fetches a list of repositories from GitLab
type RepoFetcher struct {
	logger *logrus.Entry
	client groupsClient
	config *config.Config
}

// NewRepoFetcher returns a new RepoFetcher instance
func NewRepoFetcher(logger *logrus.Entry, client groupsClient, config *config.Config) (*RepoFetcher, error) {
	return &RepoFetcher{
		logger: logger,
		client: client,
		config: config,
	}, nil
}

// Fetch fetches a list of accessible repos within the groups set in config file
func (f *RepoFetcher) GetProjectsWithAtLeastDevAccess() ([]Repo, error) {
	var repos []Repo

	opt := &gitlab.ListGroupProjectsOptions{
		MinAccessLevel: gitlab.AccessLevel(gitlab.DeveloperPermissions),
		ListOptions: gitlab.ListOptions{
			PerPage: 100,
		},
	}

	f.logger.Debugf("Fetching gitlab repos for %d groups (%s) ...", len(f.config.Groups), f.config.Groups)

	for _, group := range f.config.Groups {
		projects, _, err := f.client.ListGroupProjects(group, opt, addIncludeSubgroups)
		if err != nil {
			return []Repo{}, fmt.Errorf("failed to fetch GitLab projects or group %q: %v", group, err)
		}

		for _, p := range projects {
			repos = append(repos, Repo{
				ID:       p.ID,
				Name:     p.Name,
				FullPath: p.PathWithNamespace,
			})
		}

	}

	f.logger.Debugf("Fetching gitlab repos done. Got %d repos.", len(repos))

	return repos, nil
}
