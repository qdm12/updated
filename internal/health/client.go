package health

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

func IsClientMode(args []string) bool {
	return len(args) > 1 && args[1] == "healthcheck"
}

type Client interface {
	Query(ctx context.Context) error
}

type client struct {
	*http.Client
}

func NewClient() Client {
	const timeout = 5 * time.Second
	return &client{
		Client: &http.Client{Timeout: timeout},
	}
}

var (
	ErrReadHealthcheckBody = errors.New("cannot read healthcheck response body")
	ErrBadStatusCode       = errors.New("bad HTTP status code")
)

// Query sends an HTTP request to the other instance of
// the program, and to its internal healthcheck server.
func (c *client) Query(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://127.0.0.1:9999", nil)
	if err != nil {
		return err
	}
	resp, err := c.Do(req)
	if err != nil {
		return err
	} else if resp.StatusCode == http.StatusOK {
		return nil
	}

	b, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return fmt.Errorf("%w: %s: %s", ErrReadHealthcheckBody, resp.Status, err)
	}

	return fmt.Errorf("%w: %s: %s", ErrBadStatusCode, resp.Status, string(b))
}
