package centrifugo

type request[P any] struct {
	Method string `json:"method"`
	Params P      `json:"params"`
}

type response[R any] struct {
	Result R `json:"result"`

	Error responseError `json:"error"`
}

type emptyResponse struct {
	Error responseError `json:"error"`
}

type responseError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
