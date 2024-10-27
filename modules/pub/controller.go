package pub

import (
	"private-pub-repo/modules/app"
	"private-pub-repo/modules/pubtoken"
	"private-pub-repo/utils"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type pubController struct {
	service         PubService
	responseService app.ResponseService
	validator       *validator.Validate
	middleware      pubtoken.PubTokenJwtMiddleware
}

func newPubController(service PubService, responseService app.ResponseService, validator *validator.Validate, middleware pubtoken.PubTokenJwtMiddleware) *pubController {
	return &pubController{
		service:         service,
		responseService: responseService,
		validator:       validator,
		middleware:      middleware,
	}
}

// handlers start

func (controller *pubController) handleVersionList(ctx *fiber.Ctx) error {
	packageName := ctx.Params("package")

	publicOnly := !utils.HasJwt(ctx) || controller.middleware.CanAccess(ctx) == nil

	result, err := controller.service.VersionList(ctx.UserContext(), packageName, ctx.BaseURL(), publicOnly)

	if err != nil {
		return controller.handleControllerError(ctx, "api/packages/"+packageName, err)
	}

	return controller.responseService.SendSuccessDetailResponse(ctx, 200, result)
}

func (controller *pubController) handleVersionDetail(ctx *fiber.Ctx) error {
	packageName := ctx.Params("package")
	version := ctx.Params("version")

	publicOnly := !utils.HasJwt(ctx) || controller.middleware.CanAccess(ctx) == nil

	result, err := controller.service.VersionDetail(ctx.UserContext(), packageName, version, ctx.BaseURL(), publicOnly)

	if err != nil {
		return controller.handleControllerError(ctx, "api/packages/"+packageName+"/versions/"+version, err)
	}

	return controller.responseService.SendSuccessDetailResponse(ctx, 200, result)
}

func (controller *pubController) handleGetUploadUrl(ctx *fiber.Ctx) error {
	packageName := ctx.Params("package")
	version := ctx.Params("version")

	result, err := controller.service.VersionDetail(ctx.UserContext(), packageName, version, ctx.BaseURL(), false)

	if err != nil {
		return controller.handleControllerError(ctx, "api/packages/"+packageName+"/versions/"+version, err)
	}

	return controller.responseService.SendSuccessDetailResponse(ctx, 200, result)
}

func (controller *pubController) handleDoUpload(ctx *fiber.Ctx) error {
	packageName := ctx.Params("package")
	version := ctx.Params("version")

	result, err := controller.service.VersionDetail(ctx.UserContext(), packageName, version, ctx.BaseURL(), false)

	if err != nil {
		return controller.handleControllerError(ctx, "api/packages/"+packageName+"/versions/"+version, err)
	}

	return controller.responseService.SendSuccessDetailResponse(ctx, 200, result)
}

func (controller *pubController) handleFinishUpload(ctx *fiber.Ctx) error {
	packageName := ctx.Params("package")
	version := ctx.Params("version")

	result, err := controller.service.VersionDetail(ctx.UserContext(), packageName, version, ctx.BaseURL(), false)

	if err != nil {
		return controller.handleControllerError(ctx, "api/packages/"+packageName+"/versions/"+version, err)
	}

	return controller.responseService.SendSuccessDetailResponse(ctx, 200, result)
}

// handlers end

func (controller *pubController) handleControllerError(ctx *fiber.Ctx, currentPath string, err error) error {
	if err == fiber.ErrNotFound {
		if url := controller.service.GetUpstreamUrl(currentPath); url != nil {
			ctx.Redirect(*url, 302)
		}
		return ctx.Status(404).JSON(map[string]interface{}{
			"error": map[string]interface{}{
				"code":    "404",
				"message": "Not Found",
			},
		})
	}

	if err == fiber.ErrForbidden {
		return ctx.Status(403).JSON(map[string]interface{}{
			"error": map[string]interface{}{
				"code":    "403",
				"message": "Forbidden",
			},
		})
	}

	return ctx.Status(500).JSON(map[string]interface{}{
		"error": map[string]interface{}{
			"code":    "500",
			"message": err.Error(),
		},
	})
}
