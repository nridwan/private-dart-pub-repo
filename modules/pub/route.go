package pub

const (
	basePath            = "v1/pub/packages"
	apiPath             = "v1/pub/api/packages"
	versionListPath     = apiPath + "/:package"
	versionDetailPath   = versionListPath + "/versions/:version"
	getUploadUrlPath    = apiPath + "/versions/new"
	uploadUrlPath       = apiPath + "/versions/newUpload"
	finishUploadUrlPath = apiPath + "/versions/newUploadFinish"
	downloadPath        = basePath + "/:package/versions/:version.tar.gz"

	queryPackageListPath   = "v1/pub/query/packages"
	queryPackageUpdatePath = queryPackageListPath + "/:package"
	queryVersionListPath   = queryPackageUpdatePath + "/versions"
	queryVersionDetailPath = queryVersionListPath + "/:version"
)

func (module *PubModule) registerRoutes() {
	module.app.Get(getUploadUrlPath, module.jwtService.GetHandler(), module.middleware.CanAccess,
		module.middleware.CanWrite, module.controller.handleGetUploadUrl)
	module.app.Post(uploadUrlPath, module.jwtService.GetHandler(), module.middleware.CanAccess,
		module.middleware.CanWrite, module.controller.handleDoUpload)
	module.app.Get(finishUploadUrlPath, module.jwtService.GetHandler(), module.middleware.CanAccess,
		module.middleware.CanWrite, module.controller.handleFinishUpload)
	module.app.Get(versionListPath, module.jwtService.GetOptionalHandler(), module.controller.handleVersionList)
	module.app.Get(versionDetailPath, module.jwtService.GetOptionalHandler(), module.controller.handleVersionDetail)
	module.app.Get(downloadPath, module.jwtService.GetOptionalHandler(), module.controller.handleDownloadPath)

	module.app.Get(queryPackageListPath, module.jwtService.GetOptionalHandler(), module.controller.handleQueryPackageList)
	module.app.Put(queryPackageUpdatePath, module.jwtService.GetHandler(), module.userMiddleware.CanAccess,
		module.userMiddleware.IsAdmin, module.controller.handleQueryPackageUpdate)
	module.app.Get(queryVersionListPath, module.jwtService.GetOptionalHandler(), module.controller.handleQueryVersionList)
	module.app.Get(queryVersionDetailPath, module.jwtService.GetOptionalHandler(), module.controller.handleQueryVersionDetail)
}
