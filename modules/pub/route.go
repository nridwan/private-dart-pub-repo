package pub

const (
	basePath            = "v1/pub/api/packages"
	versionListPath     = basePath + "/:package"
	versionDetailPath   = versionListPath + "/versions/:version"
	getUploadUrlPath    = basePath + "/versions/new"
	uploadUrlPath       = basePath + "/versions/newUpload"
	finishUploadUrlPath = basePath + "/versions/newUploadFinish"
)

func (module *PubModule) registerRoutes() {
	module.app.Get(getUploadUrlPath, module.jwtService.GetHandler(), module.middleware.CanAccess, module.middleware.CanWrite, module.controller.handleGetUploadUrl)
	module.app.Post(uploadUrlPath, module.jwtService.GetHandler(), module.middleware.CanAccess, module.middleware.CanWrite, module.controller.handleDoUpload)
	module.app.Get(finishUploadUrlPath, module.jwtService.GetHandler(), module.middleware.CanAccess, module.middleware.CanWrite, module.controller.handleFinishUpload)
	module.app.Get(versionListPath, module.jwtService.GetOptionalHandler(), module.controller.handleVersionList)
	module.app.Get(versionDetailPath, module.jwtService.GetOptionalHandler(), module.controller.handleVersionDetail)
}
