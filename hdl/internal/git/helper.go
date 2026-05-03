package git

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/google/go-github/v85/github"
)

func isNotFound(err error) bool {
	var e *github.ErrorResponse
	if errors.As(err, &e) {
		return e.Response.StatusCode == http.StatusNotFound
	}
	return false
}

func wrap(action string, err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", action, err)
}
