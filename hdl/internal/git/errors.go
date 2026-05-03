package git

import "errors"

var (
	ErrInvalidToken   = errors.New("git: invalid github token provided")
	ErrRepoNotFound   = errors.New("git: repository not found")
	ErrRepoExists     = errors.New("git: repository already exists")
	ErrBranchNotFound = errors.New("git: branch not found")
	ErrBranchExists   = errors.New("git: branch already exists")
	ErrFileNotFound   = errors.New("git: file not found")
	ErrFileExists     = errors.New("git: file already exists")
)
