package centrifugo

import (
	"context"
	"fmt"
	"github.com/undefined7887/harmony-backend/internal/config"
	"github.com/undefined7887/harmony-backend/internal/util/http"

	"github.com/go-resty/resty/v2"
)

type Client struct {
	config *config.Centrifugo
	client *resty.Client
}

func NewClient(config *config.Centrifugo) *Client {
	client := resty.New().
		SetHeader(httputil.HeaderAuthorization, fmt.Sprintf("apikey %s", config.ApiKey)).
		SetBaseURL(config.ApiAddress)

	return &Client{
		config: config,
		client: client,
	}
}

const (
	publishMethod = "publish"
)

type PublishRequest struct {
	Channel string `json:"channel"`
	Data    any    `json:"data"`
}

type PublishResponse struct {
	Offset int64  `json:"offset"`
	Epoch  string `json:"epoch"`
}

func (c *Client) Publish(ctx context.Context, channel string, data any) (*PublishResponse, error) {
	resp, err := c.client.R().
		SetContext(ctx).
		SetBody(&request[PublishRequest]{
			Method: publishMethod,
			Params: PublishRequest{
				Channel: channel,
				Data:    data,
			},
		}).
		SetResult(&response[PublishResponse]{}).
		Post("")
	if err != nil {
		return nil, fmt.Errorf("centrifugo: %v", err)
	}

	result, err := handleResponse[PublishResponse](resp)
	if err != nil {
		return nil, fmt.Errorf("centrifugo: %v", err)
	}

	return result, nil
}

func handleResponse[R any](resp *resty.Response) (*R, error) {
	if resp.IsError() {
		return nil, &HttpError{
			StatusCode: resp.StatusCode(),
		}
	}

	result := resp.Result().(*response[R])

	// Centrifugo can send errors with '200 OK' status
	if result.Error.Code > 0 {
		return nil, &ApiError{
			Code:    result.Error.Code,
			Message: result.Error.Message,
		}
	}

	return &result.Result, nil
}
