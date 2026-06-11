package main

import (
	"context"
	"fmt"
	"net/url"
	"os/signal"
	"reflect"
	"syscall"

	"github.com/go-co-op/gocron/v2"
	"go.uber.org/zap"

	"detector/internal/infrastructure/inspector/composite"
	"detector/internal/infrastructure/inspector/http"
	"detector/internal/infrastructure/inspector/ping"
	repository "detector/internal/infrastructure/repository/memory"
	inspectionservice "detector/internal/inspection/application/service"
	"detector/internal/inspection/domain/inspector"
	application2 "detector/internal/report/application"
	service2 "detector/internal/report/application/service"
	"detector/internal/scheduler/application"
	"detector/internal/scheduler/application/submitter"

	"detector/internal/route/application"
	"detector/internal/route/application/dto"
)

func main() {
	repo := repository.NewMemoryRouteRepository()
	service := routeapplication.NewRouteService(repo)
	u := url.URL{Scheme: "https", Host: "example.com"}
	route, _ := service.Add(routedto.AddRouteCommand{URL: u})
	cfgPing := ping.NewInspectorConfig()
	cfgHttp := http.NewInspectorConfig()
	pingInspector := ping.NewInspector(*cfgPing)
	httpInspector := http.NewInspector(*cfgHttp)

	compositeInspector := composite.NewCompositeInspector(composite.InspectorConfig{
		Inspectors: map[string]inspector.Inspector{
			"ping": pingInspector,
			"http": httpInspector,
		},
	})

	bridge := createBridge()
	err := bridge.Register(route.ID, compositeInspector)

	if err != nil {
		fmt.Println(err)
	}
	logger, _ := zap.NewDevelopment()
	defer func(logger *zap.Logger) {
		_ = logger.Sync()
	}(logger)
	processor := application2.NewReportProcessor(logger)
	reportsaver := repository.NewReportSaver()
	reportService := service2.NewReportService(processor, reportsaver)
	printSubmitter := submitter.NewReportSubmitter(reportService)
	cronJob := gocron.CronJob("*/1 * * * *", false)

	sch := application.NewScheduler(service, bridge, printSubmitter, cronJob, logger)
	err = sch.Start()
	if err != nil {
		logger.Error(err.Error())
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	<-ctx.Done()

	defer func(sch *application.Scheduler) {
		err := sch.Stop()
		if err != nil {
			logger.Error(err.Error())
		}
	}(sch)
}

func createBridge() *inspectionservice.RouteInspectorBridge {
	repo := repository.NewMemoryRouteMethodRepository()

	registry := inspector.NewFactoryRegistry()
	registry.Register("ping", reflect.TypeFor[*ping.Inspector](), &ping.InspectorFactory{})
	registry.Register("http", reflect.TypeFor[*http.Inspector](), &http.InspectorFactory{})
	registry.Register("composite", reflect.TypeFor[*composite.Inspector](), &composite.InspectorFactory{})

	riBridge := inspectionservice.NewRouteInspectorBridge(&registry, repo)

	return &riBridge
}
