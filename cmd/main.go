package main

import (
	"context"
	"database/sql"
	"net"
	"net/url"
	"os/signal"
	"reflect"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/go-co-op/gocron/v2"
	"go.uber.org/zap"

	mygin "detector/internal/infrastructure/api/gin"
	apiinspector "detector/internal/infrastructure/api/inspector"
	apireport "detector/internal/infrastructure/api/report"
	apiroute "detector/internal/infrastructure/api/route"
	apisuperset "detector/internal/infrastructure/api/superset"
	"detector/internal/infrastructure/client/superset"
	"detector/internal/infrastructure/inspector/composite"
	"detector/internal/infrastructure/inspector/http"
	"detector/internal/infrastructure/inspector/ping"
	"detector/internal/infrastructure/repository/clickhouse/chconn"
	"detector/internal/infrastructure/repository/clickhouse/chreport"
	"detector/internal/infrastructure/repository/clickhouse/chroute"
	"detector/internal/infrastructure/repository/postgres"
	"detector/internal/infrastructure/repository/postgres/pgroute"
	"detector/internal/infrastructure/repository/postgres/pgrouteassignment"
	inspectorapp "detector/internal/inspector/application"
	inspector "detector/internal/inspector/domain"
	reportapp "detector/internal/report/application"
	"detector/internal/report/domain"
	"detector/internal/route/application"
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
	eventChannel := make(chan routeapp.Event)

	routeService := routeapp.NewService(routeRepo, eventChannel)
	bridge := createBridge(postgresConn)

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
	chrouteListener := chroute.NewEventListener(conn, logger, eventChannel)
	go chrouteListener.Listen()

	reportService := reportapp.NewService(processor, reportsaver)
	printSubmitter := submitter.NewReportSubmitter(reportService, reportSubmitterDescriptor)

	cronJob := gocron.CronJob("*/10 * * * * *", true)

	sch := schedulerapp.NewScheduler(routeService, bridge, printSubmitter, cronJob, logger)
	err = sch.Start()
	if err != nil {
		logger.Error(err.Error())
	}
	parsedurl, err := url.Parse("http://localhost:8088")
	if err != nil {
		panic(err)
	}
	supersetConfig := superset.Config{
		BaseURL:       *parsedurl,
		AdminUser:     "admin",
		AdminPassword: "admin",
	}

	guestDescriptor := superset.GuestDescriptor{
		Username:  "guest",
		Firstname: "guest",
		Lastname:  "guest",
	}

	handlers := mygin.Handlers{
		Route:     apiroute.NewHandler(routeService),
		Inspector: apiinspector.NewHandler(bridge),
		Report:    apireport.NewHandler(reportService),
		Superset: apisuperset.NewHandler(superset.NewClient(supersetConfig, logger), guestDescriptor, []apisuperset.Dashboard{
			{
				Name: "Main",
				ID:   "8f9771f2-2c8a-4e7e-8f10-a5ecd7d39255",
			},
		}, logger),
	}

	srvCfg := mygin.Config{
		Port:     5436,
		Mode:     gin.TestMode,
		Handlers: handlers,
	}
	srv := mygin.NewServer(srvCfg)
	err = srv.Start()
	if err != nil {
		logger.Error(err.Error())
	}
	defer srv.Shutdown()

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

func createBridge(postgresConn *sql.DB) *inspectorapp.RouteInspectorBridge {
	repo := pgrouteassignment.NewRepository(postgresConn)

	registry := inspector.NewFactoryRegistry()
	registry.Register("ping", reflect.TypeFor[*ping.Inspector](), &ping.InspectorFactory{})
	registry.Register("http", reflect.TypeFor[*http.Inspector](), &http.InspectorFactory{})
	registry.Register("composite", reflect.TypeFor[*composite.Inspector](), &composite.InspectorFactory{})

	riBridge := inspectorapp.NewRouteInspectorBridge(&registry, repo)

	return &riBridge
}
