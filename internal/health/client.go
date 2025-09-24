package health

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

// IsClientMode checks whether the program is run in client mode
// to query the internal healthcheck server of another instance
// of the program.
func IsClientMode(args []string) bool {
	return len(args) > 1 && args[1] == "healthcheck"
}

// Client queries the internal healthcheck server of another instance
// of the program.
type Client struct {
	*http.Client
}

// NewClient creates a new health client with a timeout.
func NewClient() *Client {
	const timeout = 5 * time.Second
	return &Client{
		Client: &http.Client{Timeout: timeout},
	}
}

var (
	ErrReadHealthcheckBody = errors.New("cannot read healthcheck response body")
	ErrBadStatusCode       = errors.New("bad HTTP status code")
)

// Query sends an HTTP request to the other instance of
// the program, and to its internal healthcheck server.
func (c *Client) Query(ctx context.Context) error {
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
	if err != nil {
		_ = resp.Body.Close()
		return fmt.Errorf("%w: %s: %w", ErrReadHealthcheckBody, resp.Status, err)
	}

	err = resp.Body.Close()
	if err != nil {
		return fmt.Errorf("closing body: %w", err)
	}

	return fmt.Errorf("%w: %s: %s", ErrBadStatusCode, resp.Status, string(b))
}
