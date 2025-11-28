package protonats

type ClientOptions interface {
	AddUnaryInterceptor(UnaryClientInterceptor)
	ChainUnaryInterceptors(...UnaryClientInterceptor)
}

type ClientOption func(options ClientOptions)

func WithUnaryInterceptor(interceptor UnaryClientInterceptor) ClientOption {
	return func(options ClientOptions) {
		options.AddUnaryInterceptor(interceptor)
	}
}

func WithChainUnaryInterceptor(interceptors ...UnaryClientInterceptor) ClientOption {
	return func(options ClientOptions) {
		options.ChainUnaryInterceptors(interceptors...)
	}
}
