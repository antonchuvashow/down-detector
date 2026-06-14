package apisuperset

import (
	"os"
	"strings"
)

func DashboardsFromEnv() []Dashboard {
	rawDashboards := os.Getenv("SUPERSET_DASHBOARDS")
	if rawDashboards != "" {
		return parseDashboards(rawDashboards)
	}

	return []Dashboard{
		{
			Name: getEnv("SUPERSET_DASHBOARD_NAME", "Main"),
			ID:   getEnv("SUPERSET_DASHBOARD_ID", "8f9771f2-2c8a-4e7e-8f10-a5ecd7d39255"),
		},
	}
}

func parseDashboards(rawDashboards string) []Dashboard {
	dashboards := make([]Dashboard, 0)

	for _, rawDashboard := range strings.Split(rawDashboards, ",") {
		rawDashboard = strings.TrimSpace(rawDashboard)
		if rawDashboard == "" {
			continue
		}

		name, id, ok := strings.Cut(rawDashboard, ":")
		if !ok {
			continue
		}

		name = strings.TrimSpace(name)
		id = strings.TrimSpace(id)
		if name == "" || id == "" {
			continue
		}

		dashboards = append(dashboards, Dashboard{Name: name, ID: id})
	}

	return dashboards
}

func getEnv(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}
