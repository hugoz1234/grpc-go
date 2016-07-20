# OpenTracing Interceptor API for Unary RPCs

This API enforces the standards specified by the [OpenTracing project](http://opentracing.io).
Any tracing implementation that follows the OpenTracing spec can use this API.

##### Usage  with gRPC's [route_guide](https://github.com/grpc/grpc-go/tree/master/examples/route_guide) example

Server Interceptor: simply call `TraceServer(...)', passing in your OpenTracing tracer.

```go
	import opint ".../opentracing_interceptors"
	...
	func main(){
		...
		var opts []grpc.ServerOption
		...
		var tracer = opentracing.InitGlobalTracer(
				//tracing impl specific
				some_tracing_impl.New(...),
			     )
		opts = append(opts, opint.TraceServer(tracer))
		...
	}
```
Client Interceptor: The functional equivalent of an interceptor in the client can be achieved with two lines of code within a unary RPC func body.

```go
    import opint ".../opentracing_interceptors"
    ...
    var tracer opentracing.Tracer //tracer must have global scope
    ...

    func printFeature(client pb.RouteGuideClient, point *pb.Point) {
        ...
	    span, ctx := opint.ClientTrace(context.Background(), tracer, "printFeature")
	    defer span.Finish()
        ...
	    feature, err := client.GetFeature(ctx, point) //calls the server
        ...
     }
 ```