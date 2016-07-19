package opentracing_interceptors

import (
  "golang.org/x/net/context"
  "github.com/opentracing/opentracing-go"

  "testing"
  opint "opentracing_interceptors"
)

func TestTraceClient(t *testing.T){
    tracer := opentracing.GlobalTracer()
    span, md := opint.TraceClient(context.Background(), tracer, "TEST_OPERATION")

    if span == nil || md ==nil {
      t.Errorf("nil span/md")
    }
}

func TestTraceServer(t *testing.T){
  tracer := opentracing.GlobalTracer()

  if serverOption := opint.TraceServer(tracer); serverOption == nil{
      t.Errorf("nil serverOption")
  }
}
