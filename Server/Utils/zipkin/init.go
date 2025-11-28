package zipkin

import (
	zipkin "github.com/openzipkin/zipkin-go"
	"github.com/openzipkin/zipkin-go/reporter"
	reporterHTTP "github.com/openzipkin/zipkin-go/reporter/http"
	"log"
)

type Tracing struct {
	Tracer   *zipkin.Tracer
	Reporter *reporter.Reporter
}

func CreateTracing(serviceName, zipkinURL, hostPort string) *Tracing {

	// 1. 创建 Reporter：专门负责把 span 上报给 Zipkin
	reporter := reporterHTTP.NewReporter(zipkinURL)

	// 2. 创建 Endpoint：告诉 Zipkin “我是谁”
	endpoint, err := zipkin.NewEndpoint(serviceName, hostPort)
	if err != nil {
		log.Fatalf("unable to create local endpoint: %+v\n", err)
	}

	// 3. 创建 Tracer：用来创建 span、管理链路关系
	tracer, err := zipkin.NewTracer(
		reporter,
		zipkin.WithLocalEndpoint(endpoint),
		zipkin.WithSharedSpans(false), // server/client 分别生成 span（最佳实践）
	)
	if err != nil {
		log.Fatalf("unable to create tracer: %+v\n", err)
	}

	// 4. 打包返回
	return &Tracing{
		Tracer:   tracer,
		Reporter: &reporter,
	}
}
