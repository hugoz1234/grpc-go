package util 

import (
	"time"
	"fmt"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"github.com/opentracing/opentracing-go"
	"github.com/lightstep/lightstep-tracer-go"
)

type metadataReaderWriter struct {
	metadata.MD
}

func (w metadataReaderWriter) Set(key, val string) {
	fmt.Println(key, " -> ", val)
	w.MD[key] = append(w.MD[key], val)
}

func (w metadataReaderWriter) ForeachKey(handler func(key, val string) error) error {
	for k, vals := range w.MD {
		for _, v := range vals {
			if dk, dv, err := metadata.DecodeKeyValue(k, v); err == nil {
				if err = handler(dk, dv); err != nil {
					return err
				}
			} else {
				return err
			}
		}
	}
	return nil
}

func InitSpan(tracer opentracing.Tracer, ctx context.Context, operation string) (opentracing.Span, context.Context){
	span := tracer.StartSpanWithOptions(opentracing.StartSpanOptions{
		OperationName: operation,
	    StartTime: time.Now(),
	})

	span.LogEvent(operation+"_called")
	md := metadata.New(make(map[string]string))
	//inject span rep into context 
	tracer.Inject(span,opentracing.TextMap,metadataReaderWriter{md})

	return span, metadata.NewContext(context.Background(), md)
}

func Setup(accessToken string) (opentracing.Tracer){
	lightstepTracer := lightstep.NewTracer(lightstep.Options{
        AccessToken: accessToken,
    })

    return lightstepTracer
}

func SetupServerInterceptor(tracer opentracing.Tracer) (grpc.ServerOption){
	return grpc.UnaryInterceptor(
		func (ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler,)(resp interface{}, err error){
			md, _ := metadata.FromContext(ctx)			

			span, err := tracer.Join(info.FullMethod,opentracing.TextMap,metadataReaderWriter{md})
			span.LogEvent(info.FullMethod+"_called")

			defer span.FinishWithOptions(opentracing.FinishOptions{FinishTime: time.Now()})
			ctx = opentracing.ContextWithSpan(ctx,span)

			return handler(ctx,req)
			})
}

func FinishSpan(span opentracing.Span){
	span.FinishWithOptions(opentracing.FinishOptions{FinishTime: time.Now()})
}

