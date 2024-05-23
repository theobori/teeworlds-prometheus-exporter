package http

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

var (
	// Default HTTP client
	httpDefaultClient = http.Client{Timeout: 10 * time.Second}
)

// Perform a HTTP get request, then fill the
// `target` struct with the body JSON reponse
func HTTPGetJson(ctx context.Context, c *http.Client, url string, target any) error {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return fmt.Errorf("creating request: GET %q: %v", url, err)
	}

	r, err := c.Do(req)
	if err != nil {
		return fmt.Errorf("request: %v", err)
	}

	defer r.Body.Close()

	if r.StatusCode != http.StatusOK {
		return fmt.Errorf("response status: %s", r.Status)
	}

	err = json.NewDecoder(r.Body).Decode(target)
	if err != nil {
		return err
	}

	return nil
}
