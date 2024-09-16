package monitor

import (
	"github.com/gofiber/fiber/v2/middleware/monitor"
)

func (module *MonitorModule) registerRoutes() {
	module.app.Get("/monitor", monitor.New(monitor.Config{Title: module.config.Getenv("APP_CODE", "App") + " Monitor"}))
}
