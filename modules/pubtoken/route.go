package pubtoken

const (
	basePath   = "v1/pubtoken"
	detailPath = basePath + "/:id"
)

func (module *PubTokenModule) registerRoutes() {
	module.app.Get(basePath, module.jwtService.GetHandler(), module.userMiddleware.CanAccess, module.controller.handleList)
	module.app.Post(basePath, module.jwtService.GetHandler(), module.userMiddleware.CanAccess, module.controller.handleCreate)
	module.app.Get(detailPath, module.jwtService.GetHandler(), module.userMiddleware.CanAccess, module.controller.handleDetail)
	module.app.Put(detailPath, module.jwtService.GetHandler(), module.userMiddleware.CanAccess, module.controller.handleUpdate)
	module.app.Delete(detailPath, module.jwtService.GetHandler(), module.userMiddleware.CanAccess, module.controller.handleDelete)
}
