package internal

import (
	"context"
	"net/http"
	"zygarde/lib/httpx"
)

type CellClient struct {
	client *httpx.Client
}

func NewCellClient(baseURL string) *CellClient {
	c := httpx.NewClient(baseURL)

	return &CellClient{
		client: c,
	}
}

func (c *CellClient) Ping(ctx context.Context, path string) (*PongResponse, error) {
	return httpx.Do[PongResponse](ctx, c.client, httpx.Request{
		Method:         http.MethodGet,
		Path:           path,
		ExpectedStatus: 200,
	})
}
