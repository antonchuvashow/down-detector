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
	"detector/internal/infrastructure/repository/clickhouse/chconn"
	"detector/internal/infrastructure/repository/clickhouse/chreport"
	"detector/internal/infrastructure/repository/postgres"
	"detector/internal/infrastructure/repository/postgres/pgroute"
	"detector/internal/infrastructure/repository/postgres/pgrouteassignment"
	inspectorapp "detector/internal/inspector/application"
	inspector "detector/internal/inspector/domain"
	reportapp "detector/internal/report/application"
	"detector/internal/report/domain"
	"detector/internal/route/application"
	"detector/internal/route/application/dto"
	route "detector/internal/route/domain"
	"detector/internal/scheduler/application"
	"detector/internal/scheduler/application/submitter"
)

func main() {
	postgresConfig := postgres.ConfigFromEnv()
	postgresConn, err := postgres.New(postgresConfig)
	if err != nil {
		panic(err)
	}

	reportSubmitterDescriptor := report.Descriptor{
		Source:    report.SourceTypeInspector,
		Latitude:  55.160023,
		Longitude: 61.401998,
		IP:        net.ParseIP(""),
		Platform:  report.PlatformTypeIOS,
	}

	routeRepo := pgroute.NewRepository(postgresConn)
	routeService := routeapp.NewService(routeRepo)
	bridge := createBridge(postgresConn)

	registerRoute("", routeService, routeRepo, bridge)
	registerRoute("", routeService, routeRepo, bridge)
	// registerRoute("", routeService, routeRepo, bridge)

	logger, _ := zap.NewDevelopment()
	defer func(logger *zap.Logger) {
		_ = logger.Sync()
	}(logger)
	processor := reportapp.NewPrintProcessor(logger)
	cfg := chconn.Config{
		Addr:     "localhost:9000",
		Database: "analytics",
		Username: "dev",
		Password: "password",
	}

	conn, err := chconn.New(cfg)
	if err != nil {
		panic(err)
	}

	reportsaver := chreport.NewRepository(conn, logger)
	reportService := reportapp.NewService(processor, reportsaver)
	printSubmitter := submitter.NewReportSubmitter(reportService, reportSubmitterDescriptor)

	cronJob := gocron.CronJob("*/10 * * * * *", true)

	sch := schedulerapp.NewScheduler(routeService, bridge, printSubmitter, cronJob, logger)
	err = sch.Start()
	if err != nil {
		logger.Error(err.Error())
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	<-ctx.Done()

	defer func(sch *schedulerapp.Scheduler) {
		err := sch.Stop()
		if err != nil {
			logger.Error(err.Error())
		}
	}(sch)
}

func registerRoute(domain string, routeService *routeapp.Service, routeRepo *pgroute.Repository, bridge *inspectorapp.RouteInspectorBridge) {
	u := url.URL{Host: domain}
	routes, err := routeRepo.Search(u)
	if err != nil {
		panic(err)
	}
	var rt route.Route
	if len(routes) == 0 {
		rt, _ = routeService.Add(routedto.AddCommand{URL: u})
	} else if len(routes) == 1 {
		rt = routes[0]
	} else {
		panic("too many routes")
	}

	cfgPing := ping.NewInspectorConfig()
	cfgPing.Interval = new(time.Millisecond * 100)
	cfgPing.PingCount = new(10)
	// cfgHttp := http.NewInspectorConfig()
	pingInspector := ping.NewInspector(*cfgPing)
	// httpInspector := http.NewInspector(*cfgHttp)

	compositeInspector := composite.NewInspector(composite.InspectorConfig{
		Inspectors: map[string]inspector.Inspector{
			"ping": pingInspector,
			// "http": httpInspector,
		},
	})

	err = bridge.Register(rt.ID, compositeInspector)
	if err != nil {
		panic(err)
	}
}

func createBridge(postgresConn *sql.DB) *inspectorapp.RouteInspectorBridge {
	repo := pgrouteassignment.NewRepository(postgresConn)

	registry := inspector.NewFactoryRegistry()
	registry.Register("ping", reflect.TypeFor[*ping.Inspector](), &ping.InspectorFactory{})
	registry.Register("http", reflect.TypeFor[*http.Inspector](), &http.InspectorFactory{})
	registry.Register("composite", reflect.TypeFor[*composite.Inspector](), &composite.InspectorFactory{})

	riBridge := inspectorapp.NewRouteInspectorBridge(&registry, repo)

	return &riBridge
}
