package user

import (
	"time"

	"github.com/gofiber/fiber/v2/middleware/limiter"
)

const (
	basePath   = "v1/users"
	detailPath = basePath + "/:id"
)

func (module *UserModule) registerRoutes() {
	publicRateLimiter := limiter.New(limiter.Config{
		Expiration:        60 * time.Second,
		LimiterMiddleware: limiter.SlidingWindow{},
	})

	module.app.Post(basePath+"/login", publicRateLimiter, module.controller.handleLogin)
	module.app.Get(basePath+"/profile", module.jwtService.GetHandler(), module.Middleware.CanAccess, module.controller.handleProfile)
	module.app.Post(basePath+"/refresh", module.jwtService.GetHandler(), module.Middleware.CanRefresh, module.controller.handleRefresh)
	module.app.Post(basePath+"/forgot-password/otp", publicRateLimiter, module.controller.handleForgotOtp)
	module.app.Post(basePath+"/forgot-password/create-password", publicRateLimiter, module.controller.handleForgotCreatePassword)
	module.app.Get(basePath, module.jwtService.GetHandler(), module.Middleware.CanAccess, module.Middleware.IsAdmin, module.controller.handleList)
	module.app.Post(basePath, module.jwtService.GetHandler(), module.Middleware.CanAccess, module.Middleware.IsAdmin, module.controller.handleCreate)
	module.app.Get(detailPath, module.jwtService.GetHandler(), module.Middleware.CanAccess, module.Middleware.IsAdmin, module.controller.handleDetail)
	module.app.Put(detailPath, module.jwtService.GetHandler(), module.Middleware.CanAccess, module.Middleware.IsAdmin, module.controller.handleUpdate)
	module.app.Delete(detailPath, module.jwtService.GetHandler(), module.Middleware.CanAccess, module.Middleware.IsAdmin, module.controller.handleDelete)
}
