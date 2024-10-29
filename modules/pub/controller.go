package pub

import (
	"net/url"
	"private-pub-repo/modules/app"
	"private-pub-repo/modules/pubtoken"
	"private-pub-repo/utils"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

const jsonResponseType = "application/vnd.pub.v2+json"

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

	publicOnly := !utils.HasJwt(ctx) || controller.middleware.HasAccess(ctx) != nil

	result, err := controller.service.VersionList(ctx.UserContext(), packageName, ctx.BaseURL(), publicOnly)

	if err != nil {
		return controller.handleControllerError(ctx, "api/packages/"+packageName, err)
	}

	return ctx.Status(200).JSON(result, jsonResponseType)
}

func (controller *pubController) handleVersionDetail(ctx *fiber.Ctx) error {
	packageName := ctx.Params("package")
	version := ctx.Params("version")

	publicOnly := !utils.HasJwt(ctx) || controller.middleware.HasAccess(ctx) != nil

	result, err := controller.service.VersionDetail(ctx.UserContext(), packageName, version, ctx.BaseURL(), publicOnly)

	if err != nil {
		return controller.handleControllerError(ctx, "api/packages/"+packageName+"/versions/"+version, err)
	}

	return ctx.Status(200).JSON(result, jsonResponseType)
}

func (controller *pubController) handleDownloadPath(ctx *fiber.Ctx) error {
	packageName := ctx.Params("package")
	version := ctx.Params("version")

	publicOnly := !utils.HasJwt(ctx) || controller.middleware.HasAccess(ctx) != nil

	downloadUrl, err := controller.service.GetDownloadUrl(ctx.UserContext(), packageName, version, ctx.BaseURL(), publicOnly)

	if err != nil {
		return controller.handleControllerError(ctx, "api/packages/"+packageName+"/versions/"+version, err)
	}

	return ctx.Redirect(*downloadUrl, fiber.StatusFound)
}

func (controller *pubController) handleGetUploadUrl(ctx *fiber.Ctx) error {
	return ctx.Status(200).JSON(map[string]interface{}{
		"url":    ctx.BaseURL() + "/" + uploadUrlPath,
		"fields": map[string]interface{}{},
	}, jsonResponseType)
}

func (controller *pubController) handleDoUpload(ctx *fiber.Ctx) error {
	// Get the uploaded file from the "file" field
	file, err := ctx.FormFile("file")
	if err != nil {
		return controller.processError(ctx, fiber.StatusBadRequest, err.Error())
	}

	err = controller.service.UploadVersion(ctx.UserContext(), file, controller.middleware.GetPubUserId(ctx))

	if err != nil {
		return ctx.Redirect(ctx.BaseURL()+"/"+finishUploadUrlPath+"?error="+url.QueryEscape(err.Error()), fiber.StatusNoContent)
	}

	return ctx.Redirect(ctx.BaseURL()+"/"+finishUploadUrlPath, fiber.StatusNoContent)
}

func (controller *pubController) handleFinishUpload(ctx *fiber.Ctx) error {
	errorMsg := ctx.Query("error")

	if errorMsg != "" {
		return controller.processError(ctx, fiber.StatusBadRequest, errorMsg)
	}

	return ctx.JSON(map[string]interface{}{
		"success": map[string]interface{}{
			"message": "Successfully uploaded package.",
		},
	}, jsonResponseType)
}

// handlers end

func (controller *pubController) handleControllerError(ctx *fiber.Ctx, currentPath string, err error) error {
	if err == fiber.ErrNotFound {
		if url := controller.service.GetUpstreamUrl(ctx.UserContext(), currentPath); url != nil {
			ctx.Redirect(*url, fiber.StatusFound)
		}
		return controller.processError(ctx, fiber.StatusNotFound, "Not Found")
	}

	if err == fiber.ErrForbidden {
		return controller.processError(ctx, fiber.StatusForbidden, "Forbidden")
	}

	return controller.processError(ctx, fiber.StatusBadRequest, err.Error())
}

func (controller *pubController) processError(ctx *fiber.Ctx, status int, message string) error {
	return ctx.Status(status).JSON(map[string]interface{}{
		"error": map[string]interface{}{
			"code":    strconv.Itoa(status),
			"message": message,
		},
	}, jsonResponseType)
}
