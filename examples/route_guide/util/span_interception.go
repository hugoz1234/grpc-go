package util 

import (
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"github.com/opentracing/opentracing-go"
	"github.com/lightstep/lightstep-tracer-go"
)

func Inject(op string) (opentracing.Span, context.Context) {
	//create span
	span := opentracing.GlobalTracer().StartSpanWithOptions(opentracing.StartSpanOptions{
		OperationName: op,
	    StartTime: time.Now(),
	})
	span.LogEvent("printFeature_called")
	//inject span rep into context metadata
	values := make(map[string]string)
	values["op"] = op

	return span, metadata.NewContext(context.Background(), metadata.New(values))
}

func Setup(accessToken string) (opentracing.Tracer){
	lightstepTracer := lightstep.NewTracer(lightstep.Options{
        AccessToken: accessToken,
    })

    return lightstepTracer
}

func join(ctx context.Context) (opentracing.Span, error){
	mp, _ := metadata.FromContext(ctx)
	span := opentracing.GlobalTracer().StartSpanWithOptions(opentracing.StartSpanOptions{
		OperationName: mp["op"][0],
	    StartTime: time.Now(),
	})
	return span, nil
}

func SetupServerInterceptor() (grpc.ServerOption){
	return grpc.UnaryInterceptor(
		func (ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler,)(resp interface{}, err error){
			span, err := join(ctx)
			span.LogEvent("GetFeature_called")

			defer span.FinishWithOptions(opentracing.FinishOptions{FinishTime: time.Now()})
			ctx = opentracing.ContextWithSpan(ctx,span)

			return handler(ctx,req)
			})
}

func FinishSpan(span opentracing.Span){
	span.FinishWithOptions(opentracing.FinishOptions{FinishTime: time.Now()})
}


