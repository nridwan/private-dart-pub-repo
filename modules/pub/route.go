package pub

const (
	basePath          = "v1/pub/api/packages"
	versionListPath   = basePath + "/:package"
	versionDetailPath = versionListPath + "/versions/:version"
)

func (module *PubModule) registerRoutes() {
	module.app.Get(versionListPath, module.jwtService.GetOptionalHandler(), module.controller.handleVersionList)
	module.app.Get(versionListPath+"/versions/new", module.jwtService.GetHandler(), module.middleware.CanAccess, module.middleware.CanWrite, module.controller.handleGetUploadUrl)
	module.app.Post(versionListPath+"/versions/newUpload", module.jwtService.GetHandler(), module.middleware.CanAccess, module.middleware.CanWrite, module.controller.handleDoUpload)
	module.app.Get(versionListPath+"/versions/newUploadFinish", module.jwtService.GetHandler(), module.middleware.CanAccess, module.middleware.CanWrite, module.controller.handleFinishUpload)
	module.app.Get(versionDetailPath, module.jwtService.GetOptionalHandler(), module.controller.handleVersionDetail)
}
