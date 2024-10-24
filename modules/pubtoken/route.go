package pubtoken

const (
	basePath   = "v1/pubtoken"
	detailPath = basePath + "/:id"
)

func (module *UserModule) registerRoutes() {
	module.app.Get(basePath, module.jwtService.GetHandler(), module.userMiddleware.CanAccess, module.userMiddleware.IsAdmin, module.controller.handleList)
	module.app.Post(basePath, module.jwtService.GetHandler(), module.userMiddleware.CanAccess, module.userMiddleware.IsAdmin, module.controller.handleCreate)
	module.app.Get(detailPath, module.jwtService.GetHandler(), module.userMiddleware.CanAccess, module.userMiddleware.IsAdmin, module.controller.handleDetail)
	module.app.Delete(detailPath, module.jwtService.GetHandler(), module.userMiddleware.CanAccess, module.userMiddleware.IsAdmin, module.controller.handleDelete)
}
