package opentracing_interceptors 

import (
	"fmt"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"github.com/opentracing/opentracing-go/ext"
	"google.golang.org/grpc/metadata"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc/grpclog"
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

func TraceClient(ctx context.Context, tracer opentracing.Tracer, operation string) (opentracing.Span, context.Context){
	span := tracer.StartSpan(operation)
	md, ok := metadata.FromContext(ctx)
	if !ok{
		md = metadata.MD{}
	}
	err := tracer.Inject(span.Context(),opentracing.TextMap,metadataReaderWriter{md})
	if err != nil{
		grpclog.Fatalf("Failed to create tracer %v",err)
	}

	return span, metadata.NewContext(ctx, md)
}

func TraceServer(tracer opentracing.Tracer) (grpc.ServerOption){
	return grpc.UnaryInterceptor(
		func (ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler,)(resp interface{}, err error){
			md, ok := metadata.FromContext(ctx)
			if !ok {
				md = metadata.MD{}
			}	
			sctx, err := tracer.Extract(opentracing.TextMap,metadataReaderWriter{md})
			span := tracer.StartSpan(info.FullMethod, ext.RPCServerOption(sctx))
			ext.Component.Set(span, "grpc")

			if err != nil {
				ext.Error.Set(span, true)
			}
			defer span.Finish()
			ctx = opentracing.ContextWithSpan(ctx,span)

			return handler(ctx,req)
			})
}

