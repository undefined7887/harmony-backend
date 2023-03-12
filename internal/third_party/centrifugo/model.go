package centrifugo

type Request[P any] struct {
	Method string `json:"method"`
	Params P      `json:"params"`
}

type Response[R any] struct {
	Result R `json:"result"`

	Error *ResponseError `json:"error,omitempty"`
}

func NewResponse[T any](result T) Response[T] {
	return Response[T]{
		Result: result,
	}
}

type EmptyResponse struct {
	Error ResponseError `json:"error"`
}

type ResponseError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
