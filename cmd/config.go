package main

import (
	"net"
	"os"
	"strconv"

	"github.com/go-co-op/gocron/v2"

	"detector/internal/report/domain"
)

func reportSubmitterDescriptorFromEnv() report.Descriptor {
	return report.Descriptor{
		Source:    report.SourceType(getEnv("REPORT_SUBMITTER_SOURCE", string(report.SourceTypeInspector))),
		Latitude:  getEnvFloat("REPORT_SUBMITTER_LATITUDE", 55.160023),
		Longitude: getEnvFloat("REPORT_SUBMITTER_LONGITUDE", 61.401998),
		IP:        net.ParseIP(getEnv("REPORT_SUBMITTER_IP", "")),
		Platform:  report.PlatformType(getEnv("REPORT_SUBMITTER_PLATFORM", string(report.PlatformTypeIOS))),
	}
}

func schedulerCronJobFromEnv() gocron.JobDefinition {
	return gocron.CronJob(getEnv("SCHEDULER_CRON", "*/10 * * * * *"), true)
}

func getEnv(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}

func getEnvFloat(key string, fallback float64) float64 {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	parsed, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return fallback
	}

	return parsed
}
