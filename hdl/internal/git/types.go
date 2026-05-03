package git

import (
	"sync"

	"github.com/google/go-github/v85/github"
)

type GitClient struct {
	client *github.Client
	owner  string
	mu     sync.RWMutex
}
