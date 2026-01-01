# protonats

This monorepo contains shared definitions used by the Go and Java implementation of protonats.
ProtoNATS is a protoc plugin that generates server and client code for NATS microservices, in the given target language.

Currently only Go is supported, but support for Java is on the way.

This repository will only be used after the refactor is done, in the meanwhile please refer to the [Go repository](https://github.com/d0x7/protoc-gen-go-nats) for the Go implementation.

## Interceptors

> [!WARNING]
> Modifying Headers in Interceptors requires special attention.

When writing a client interceptor, you are working with the network layer, not the context layer.

As the caller (i.e. your app), you need to use `protonats.WithOutgoingHeader(ctx, ...)` (and adjacent methods, starting with `WithOutgoingHeader*` or `SetOutgoingHeader*`) to set initial metadata.  
This metadata will then be available from within client interceptors, to be set as headers.

Once the interceptor is called, the context has already been evaluated for headers, and any changes to the context inside the interceptor will not be visible to the network layer, but they will be visible to subsequent interceptors in the chain.

So for example, you can modify the context in an interceptor to pass data to subsequent interceptors, but if you want to add headers to be sent over the network, you need to use the provided `nats.Header` map.

## License

[MIT](LICENSE) © 2026 Dorian Heinrichs
