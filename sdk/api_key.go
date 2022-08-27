package sendgrid

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sendgrid/rest"
)

// APIKey is a Sendgrid API key.
type APIKey struct {
	ID     string   `json:"api_key_id,omitempty"`
	APIKey string   `json:"api_key,omitempty"`
	Name   string   `json:"name,omitempty"`
	Scopes []string `json:"scopes,omitempty"`
}

func parseAPIKey(respBody string) (*APIKey, RequestError) {
	var body APIKey
	if err := json.Unmarshal([]byte(respBody), &body); err != nil {
		return nil, RequestError{
			StatusCode: http.StatusInternalServerError,
			Err:        fmt.Errorf("failed parsing API key: %w", err),
		}
	}

	return &body, RequestError{StatusCode: http.StatusOK, Err: nil}
}

func parseAPIKeys(respBody string) ([]APIKey, RequestError) {
	var body []APIKey
	if err := json.Unmarshal([]byte(respBody), &body); err != nil {
		return nil, RequestError{
			StatusCode: http.StatusInternalServerError,
			Err:        fmt.Errorf("failed parsing API key: %w", err),
		}
	}

	return body, RequestError{StatusCode: http.StatusOK, Err: nil}
}

// CreateAPIKey creates an APIKey and returns it.
func (c *Client) CreateAPIKey(req *APIKey) (*APIKey, RequestError) {
	if req.Name == "" {
		return nil, RequestError{
			StatusCode: http.StatusInternalServerError,
			Err:        ErrNameRequired,
		}
	}

	respBody, statusCode, err := c.Post("POST", "/api_keys", req)
	if err != nil {
		return nil, RequestError{
			StatusCode: http.StatusInternalServerError,
			Err:        fmt.Errorf("failed creating API key: %w", err),
		}
	}

	if statusCode >= http.StatusMultipleChoices {
		return nil, RequestError{
			StatusCode: statusCode,
			Err:        fmt.Errorf("%w, status: %d, response: %s", ErrFailedCreatingAPIKey, statusCode, respBody),
		}
	}

	return parseAPIKey(respBody)
}

// ReadAPIKey retreives an APIKey and returns it.
func (c *Client) ReadAPIKey(id string) (*APIKey, RequestError) {
	if id == "" {
		return nil, RequestError{
			StatusCode: http.StatusInternalServerError,
			Err:        ErrAPIKeyIDRequired,
		}
	}

	respBody, _, err := c.Get("GET", "/api_keys/"+id)
	if err != nil {
		return nil, RequestError{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	return parseAPIKey(respBody)
}

func (c *Client) ReadAPIKeys() ([]APIKey, RequestError) {
	respBody, _, err := c.Get("GET", "/api_keys")
	if err != nil {
		return nil, RequestError{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	return parseAPIKeys(respBody)
}

// UpdateAPIKey edits an APIKey and returns it.
func (c *Client) UpdateAPIKey(id string, req *APIKey) (*APIKey, RequestError) {
	if id == "" {
		return nil, RequestError{
			StatusCode: http.StatusInternalServerError,
			Err:        ErrAPIKeyIDRequired,
		}
	}

	var method rest.Method
	if len(req.Scopes) > 0 {
		method = rest.Put
	} else {
		method = rest.Patch
	}

	respBody, _, err := c.Post(method, "/api_keys/"+id, req)
	if err != nil {
		return nil, RequestError{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	return parseAPIKey(respBody)
}

// DeleteAPIKey deletes an APIKey.
func (c *Client) DeleteAPIKey(id string) (bool, *RequestError) {
	if id == "" {
		return false, &RequestError{
			StatusCode: http.StatusInternalServerError,
			Err:        ErrAPIKeyIDRequired,
		}
	}

	responseBody, statusCode, err := c.Get("DELETE", "/api_keys/"+id)
	if err != nil {
		return false, &RequestError{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	if statusCode >= http.StatusMultipleChoices && statusCode != http.StatusNotFound {
		return false, &RequestError{
			StatusCode: statusCode,
			Err:        fmt.Errorf("%w, status: %d, response: %s", ErrFailedDeletingAPIKey, statusCode, responseBody),
		}
	}

	return true, nil
}
