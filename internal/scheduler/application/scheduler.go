package application

import (
	"fmt"

	"github.com/go-co-op/gocron/v2"
	"go.uber.org/zap"

	inspectionservice "detector/internal/inspection/application/service"
	routeapplication "detector/internal/route/application"
)

type Scheduler struct {
	routeService              *routeapplication.RouteService
	routeInspectorBridge      *inspectionservice.RouteInspectorBridge
	inspectionResultSubmitter InspectionResultSubmitter
	cronJob                   gocron.JobDefinition
	cronScheduler             gocron.Scheduler
	logger                    *zap.Logger
}

func NewScheduler(service *routeapplication.RouteService,
	bridge *inspectionservice.RouteInspectorBridge,
	submitter InspectionResultSubmitter,
	cronJob gocron.JobDefinition,
	logger *zap.Logger) *Scheduler {
	return &Scheduler{
		routeService:              service,
		routeInspectorBridge:      bridge,
		inspectionResultSubmitter: submitter,
		cronJob:                   cronJob,
		logger:                    logger,
	}
}

func (s *Scheduler) Start() error {
	var err error
	s.cronScheduler, err = gocron.NewScheduler()
	if err != nil {
		return fmt.Errorf("error creating cron scheduler: %w", err)
	}

	_, err = s.cronScheduler.NewJob(
		s.cronJob,
		gocron.NewTask(s.tick),
	)
	if err != nil {
		return fmt.Errorf("error creating inspection job: %w", err)
	}
	s.cronScheduler.Start()
	return nil
}

func (s *Scheduler) Stop() error {
	err := s.cronScheduler.Shutdown()
	if err != nil {
		return fmt.Errorf("shutdown cron job: %w", err)
	}
	return nil
}

func (s *Scheduler) tick() {
	routes, err := s.routeService.GetAllRoutes()
	if err != nil {
		s.logger.Error("error getting routes:", zap.Error(err))
		return
	}
	for _, route := range routes {
		insp, err := s.routeInspectorBridge.FindInspector(route.ID)
		if err != nil {
			s.logger.Error("error finding inspector:", zap.Error(err))
			continue
		}
		result, err := insp.Inspect(route)
		if err != nil {
			s.logger.Error("error inspecting route:", zap.Error(err))
			continue
		}
		err = s.inspectionResultSubmitter.Submit(result)
		if err != nil {
			s.logger.Error("error submitting inspection result:", zap.Error(err))
		}
	}
}
