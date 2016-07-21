package opentracing_interceptors

import (
  "golang.org/x/net/context"
  "github.com/opentracing/opentracing-go"
  "github.com/opentracing/opentracing-go/mocktracer"

  "github.com/stretchr/testify/assert"

  "testing"
  opint "opentracing_interceptors"
)

func TestTraceClient(t *testing.T){
    tracer := mocktracer.New()
    opentracing.InitGlobalTracer(tracer)
    span, md := opint.TraceClient(context.Background(), tracer, "TEST_OPERATION")

    if span == nil || md ==nil {
      t.Errorf("nil span/md")
    }
    span.Finish()

    spans := tracer.FinishedSpans()
    assert.Equal(t, 1, len(spans))
    assert.Equal(t, spans[0].Tracer(), tracer)
    assert.Equal(t,spans[0].OperationName, "TEST_OPERATION")

}

func TestTraceServer(t *testing.T){
  opentracing.InitGlobalTracer(mocktracer.New())
  tracer := opentracing.GlobalTracer()

  if serverOption := opint.TraceServer(tracer); serverOption == nil{
      t.Errorf("nil serverOption")
  }

}
