package fireblocksdk

import (
	"time"
)

type PostRequestOption struct {
	idempotencyKey string
}

type GetRequestOption struct {
	paged bool
}

func WithIdempotencyKey(idempotencyKey string) func(*PostRequestOption) {
	return func(option *PostRequestOption) {
		option.idempotencyKey = idempotencyKey
	}
}

func WithPaged(paged bool) func(*GetRequestOption) {
	return func(option *GetRequestOption) {
		option.paged = paged
	}
}

type SDKOptions struct {
	timeoutMilliseconds time.Duration
	auth                IAuthProvider
	client              IAPIClient
}

type FireblocksSDK struct {
	baseURL string
	client  IAPIClient
	auth    IAuthProvider
}

func WithAuthProvider(auth IAuthProvider) func(o *SDKOptions) {
	return func(o *SDKOptions) {
		o.auth = auth
	}
}

func WithAPIClient(client IAPIClient) func(o *SDKOptions) {
	return func(o *SDKOptions) {
		o.client = client
	}
}

func WithTimout(timeout time.Duration) func(o *SDKOptions) {
	return func(o *SDKOptions) {
		o.timeoutMilliseconds = timeout
	}
}

func CreateSDK(apikey, privateKey string, baseURL string, opts ...func(o *SDKOptions)) (*FireblocksSDK, error) {
	opt := &SDKOptions{}

	for _, o := range opts {
		o(opt)
	}

	if opt.auth == nil {
		provider, err := NewAuthProvider(apikey, privateKey)
		if err != nil {
			return nil, err
		}

		opt.auth = provider
	}

	if opt.client == nil {
		opt.client = NewAPIClient(opt.auth, baseURL)
	}

	sdk := &FireblocksSDK{
		baseURL: baseURL,
		client:  opt.client,
		auth:    opt.auth,
	}

	return sdk, nil
}