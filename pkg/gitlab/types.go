package gitlab

// Status values of a pipeline
const (
	StatusRunning  = "running"
	StatusPending  = "pending"
	StatusSuccess  = "success"
	StatusFailed   = "failed"
	StatusCanceled = "canceled"
	StatusKipped   = "skipped"
)

// Repo holds all relevant infos for a GitLab repo
type Repo struct {
	ID       int
	Name     string
	FullPath string
}

// RepoStatus holds the state of the last pipeline for each ref
type RepoStatus map[string]PipelineStatus

type PipelineStatus struct {
	ID       int
	Status   string
	Username string
	UserID   int
	Ref      string
	SHA      string
}
