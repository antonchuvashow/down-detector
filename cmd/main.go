package main

import (
	"context"
	"database/sql"
	"net"
	"net/url"
	"os/signal"
	"reflect"
	"syscall"
	"time"

	"github.com/go-co-op/gocron/v2"
	"go.uber.org/zap"

	"detector/internal/infrastructure/inspector/composite"
	"detector/internal/infrastructure/inspector/http"
	"detector/internal/infrastructure/inspector/ping"
	"detector/internal/infrastructure/repository/clickhouse/connection"
	"detector/internal/infrastructure/repository/clickhouse/report"
	connection2 "detector/internal/infrastructure/repository/postgres/connection"
	route2 "detector/internal/infrastructure/repository/postgres/route"
	"detector/internal/infrastructure/repository/postgres/routemethod"
	inspectionservice "detector/internal/inspection/application/service"
	"detector/internal/inspection/domain/inspector"
	application2 "detector/internal/report/application"
	service2 "detector/internal/report/application/service"
	"detector/internal/report/domain"
	routedomain "detector/internal/route/domain"
	"detector/internal/scheduler/application"
	"detector/internal/scheduler/application/submitter"

	"detector/internal/route/application"
	"detector/internal/route/application/dto"
)

func main() {
	postgresConfig := connection2.ConfigFromEnv()
	postgresConn, err := connection2.New(postgresConfig)
	if err != nil {
		panic(err)
	}

	reportSubmitterDescriptor := domain.Descriptor{
		Source:    domain.SourceTypeInspector,
		Latitude:  55.160023,
		Longitude: 61.401998,
		IP:        net.ParseIP(""),
		Platform:  domain.PlatformTypeIOS,
	}

	routeRepo := route2.NewPostgresRouteRepository(postgresConn)
	routeService := routeapplication.NewRouteService(routeRepo)
	bridge := createBridge(postgresConn)

	registerRoute("", routeService, routeRepo, bridge)
	registerRoute("", routeService, routeRepo, bridge)
	// registerRoute("", routeService, routeRepo, bridge)

	logger, _ := zap.NewDevelopment()
	defer func(logger *zap.Logger) {
		_ = logger.Sync()
	}(logger)
	processor := application2.NewReportProcessor(logger)
	cfg := connection.Config{
		Addr:     "localhost:9000",
		Database: "analytics",
		Username: "dev",
		Password: "password",
	}

	conn, err := connection.New(cfg)
	if err != nil {
		panic(err)
	}

	reportsaver := report.NewRepository(conn, logger)
	reportService := service2.NewReportService(processor, reportsaver)
	printSubmitter := submitter.NewReportSubmitter(reportService, reportSubmitterDescriptor)

	cronJob := gocron.CronJob("*/10 * * * * *", true)

	sch := application.NewScheduler(routeService, bridge, printSubmitter, cronJob, logger)
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

func registerRoute(domain string, routeService *routeapplication.RouteService, routeRepo *route2.Repository, bridge *inspectionservice.RouteInspectorBridge) {
	u := url.URL{Host: domain}
	routes, err := routeRepo.Search(u)
	if err != nil {
		panic(err)
	}
	var route routedomain.Route
	if len(routes) == 0 {
		route, _ = routeService.Add(routedto.AddRouteCommand{URL: u})
	} else if len(routes) == 1 {
		route = routes[0]
	} else {
		panic("too many routes")
	}

	cfgPing := ping.NewInspectorConfig()
	cfgPing.Interval = new(time.Millisecond * 100)
	cfgPing.PingCount = new(10)
	// cfgHttp := http.NewInspectorConfig()
	pingInspector := ping.NewInspector(*cfgPing)
	// httpInspector := http.NewInspector(*cfgHttp)

	compositeInspector := composite.NewCompositeInspector(composite.InspectorConfig{
		Inspectors: map[string]inspector.Inspector{
			"ping": pingInspector,
			// "http": httpInspector,
		},
	})

	err = bridge.Register(route.ID, compositeInspector)
	if err != nil {
		panic(err)
	}
}

func createBridge(postgresConn *sql.DB) *inspectionservice.RouteInspectorBridge {
	repo := routemethod.NewPostgresRouteMethodRepository(postgresConn)

	registry := inspector.NewFactoryRegistry()
	registry.Register("ping", reflect.TypeFor[*ping.Inspector](), &ping.InspectorFactory{})
	registry.Register("http", reflect.TypeFor[*http.Inspector](), &http.InspectorFactory{})
	registry.Register("composite", reflect.TypeFor[*composite.Inspector](), &composite.InspectorFactory{})

	riBridge := inspectionservice.NewRouteInspectorBridge(&registry, repo)

	return &riBridge
}
