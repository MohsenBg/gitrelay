package git

import (
	"context"
	"fmt"

	"github.com/google/go-github/v85/github"
	"golang.org/x/oauth2"
)

const refPrefix = "refs/heads/"

// NewClient initializes a new GitClient using an OAuth2 token.
func NewClient(ctx context.Context, token string) *GitClient {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	return &GitClient{
		client: github.NewClient(tc),
	}
}

// ValidateToken verifies if the provided token is valid by fetching the user profile.
func (gc *GitClient) ValidateToken(ctx context.Context) error {
	_, _, err := gc.client.Users.Get(ctx, "")
	if err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidToken, err)
	}
	return nil
}

// GetOwner retrieves and caches the authenticated user's login name thread-safely.
func (gc *GitClient) GetOwner(ctx context.Context) (string, error) {
	gc.mu.RLock()
	if gc.owner != "" {
		defer gc.mu.RUnlock()
		return gc.owner, nil
	}
	gc.mu.RUnlock()

	gc.mu.Lock()
	defer gc.mu.Unlock()

	if gc.owner != "" {
		return gc.owner, nil
	}

	user, _, err := gc.client.Users.Get(ctx, "")
	if err != nil {
		return "", wrap("github: get user profile", err)
	}

	gc.owner = user.GetLogin()
	return gc.owner, nil
}

// CreateRepo creates a new repository. It returns ErrRepoExists if it already exists.
func (gc *GitClient) CreateRepo(ctx context.Context, name string, private bool) error {
	repo := &github.Repository{
		Name:    github.Ptr(name),
		Private: github.Ptr(private),
	}

	_, _, err := gc.client.Repositories.Create(ctx, "", repo)
	if err != nil {
		return wrap("github: create repo", err)
	}
	return nil
}

// DeleteRepo removes a repository. Returns ErrRepoNotFound if it doesn't exist.
func (gc *GitClient) DeleteRepo(ctx context.Context, repo string) error {
	owner, err := gc.GetOwner(ctx)
	if err != nil {
		return err
	}

	_, err = gc.client.Repositories.Delete(ctx, owner, repo)
	if err != nil {
		if isNotFound(err) {
			return ErrRepoNotFound
		}
		return wrap("github: delete repo", err)
	}
	return nil
}

// CreateBranch creates a new branch. Returns ErrBranchExists if already present.
func (gc *GitClient) CreateBranch(ctx context.Context, repo, branchName, fromBranch string) error {
	owner, err := gc.GetOwner(ctx)
	if err != nil {
		return err
	}

	newRefPath := fmt.Sprintf("%s%s", refPrefix, branchName)
	_, _, err = gc.client.Repositories.GetBranch(ctx, owner, repo, newRefPath, 0)
	if err == nil {
		return ErrBranchExists
	}

	baseRef, _, err := gc.client.Repositories.GetBranch(ctx, owner, repo, fromBranch, 0)
	if err != nil {
		if isNotFound(err) {
			return ErrBranchNotFound
		}
		return wrap("github: get base branch", err)
	}

	ref := github.CreateRef{
		Ref: newRefPath,
		SHA: baseRef.Commit.GetSHA(),
	}

	_, _, err = gc.client.Git.CreateRef(ctx, owner, repo, ref)
	if err != nil {
		return wrap("github: create branch ref", err)
	}
	return nil
}

// DeleteBranch removes a branch. Returns ErrBranchNotFound if it doesn't exist.
func (gc *GitClient) DeleteBranch(ctx context.Context, repo, branchName string) error {
	owner, err := gc.GetOwner(ctx)
	if err != nil {
		return err
	}

	ref := fmt.Sprint(refPrefix, branchName)
	_, err = gc.client.Git.DeleteRef(ctx, owner, repo, ref)
	if err != nil {
		if isNotFound(err) {
			return ErrBranchNotFound
		}
		return wrap("github: delete branch", err)
	}
	return nil
}

// GetFile retrieves the content of a file from a specific branch.
func (gc *GitClient) GetFile(ctx context.Context, repo, path, branch string) ([]byte, error) {
	owner, err := gc.GetOwner(ctx)
	if err != nil {
		return nil, err
	}

	opts := &github.RepositoryContentGetOptions{Ref: branch}
	fileContent, _, _, err := gc.client.Repositories.GetContents(ctx, owner, repo, path, opts)
	if err != nil {
		if isNotFound(err) {
			return nil, ErrFileNotFound
		}
		return nil, wrap("github: get file content", err)
	}

	content, err := fileContent.GetContent()
	if err != nil {
		return nil, wrap("github: decode file content", err)
	}

	return []byte(content), nil
}

// CreateFile adds a new file. Returns ErrFileExists if the file already exists at that path.
func (gc *GitClient) CreateFile(ctx context.Context, repo, path, message string, content []byte, branch string) error {
	owner, err := gc.GetOwner(ctx)
	if err != nil {
		return err
	}

	// Check if file already exists
	_, _, _, err = gc.client.Repositories.GetContents(ctx, owner, repo, path, &github.RepositoryContentGetOptions{Ref: branch})
	if err == nil {
		return ErrFileExists
	}

	opts := &github.RepositoryContentFileOptions{
		Message: github.Ptr(message),
		Content: content,
		Branch:  github.Ptr(branch),
	}

	_, _, err = gc.client.Repositories.CreateFile(ctx, owner, repo, path, opts)
	if err != nil {
		return wrap("github: create file", err)
	}
	return nil
}

// UpdateFile modifies an existing file. Returns ErrFileNotFound if the file does not exist.
func (gc *GitClient) UpdateFile(ctx context.Context, repo, path, message string, content []byte, branch string) error {
	owner, err := gc.GetOwner(ctx)
	if err != nil {
		return err
	}

	// We must get the current file's SHA to update it
	file, _, _, err := gc.client.Repositories.GetContents(ctx, owner, repo, path, &github.RepositoryContentGetOptions{Ref: branch})
	if err != nil {
		if isNotFound(err) {
			return ErrFileNotFound
		}
		return wrap("github: get file for update", err)
	}

	opts := &github.RepositoryContentFileOptions{
		Message: github.Ptr(message),
		Content: content,
		SHA:     file.SHA,
		Branch:  github.Ptr(branch),
	}

	_, _, err = gc.client.Repositories.UpdateFile(ctx, owner, repo, path, opts)
	if err != nil {
		return wrap("github: update file", err)
	}
	return nil
}
