package pub

import (
	"private-pub-repo/modules/app"
	"private-pub-repo/modules/pubtoken"
	"private-pub-repo/utils"
	"strconv"

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

	return ctx.Status(200).JSON(result, "application/vnd.pub.v2+json")
}

func (controller *pubController) handleVersionDetail(ctx *fiber.Ctx) error {
	packageName := ctx.Params("package")
	version := ctx.Params("version")

	publicOnly := !utils.HasJwt(ctx) || controller.middleware.CanAccess(ctx) == nil

	result, err := controller.service.VersionDetail(ctx.UserContext(), packageName, version, ctx.BaseURL(), publicOnly)

	if err != nil {
		return controller.handleControllerError(ctx, "api/packages/"+packageName+"/versions/"+version, err)
	}

	return ctx.Status(200).JSON(result, "application/vnd.pub.v2+json")
}

func (controller *pubController) handleGetUploadUrl(ctx *fiber.Ctx) error {
	return ctx.Status(200).JSON(map[string]interface{}{
		"url":    ctx.BaseURL() + "/" + uploadUrlPath,
		"fields": map[string]interface{}{},
	}, "application/vnd.pub.v2+json")
}

func (controller *pubController) handleDoUpload(ctx *fiber.Ctx) error {
	//open "file" multipart
	//extract tar.gz, read "pubspec.yaml", "changelog.md", "readme.md"
	//upload tar.gz
	//save package (if not exist) and version to db
	return ctx.Redirect(ctx.BaseURL()+"/"+finishUploadUrlPath, fiber.StatusFound)
}

func (controller *pubController) handleFinishUpload(ctx *fiber.Ctx) error {
	errorMsg := ctx.Query("error")

	if errorMsg != "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(map[string]interface{}{
			"error": map[string]interface{}{
				"code":    strconv.Itoa(fiber.StatusBadRequest),
				"message": errorMsg,
			},
		}, "application/vnd.pub.v2+json")
	}

	return ctx.JSON(map[string]interface{}{
		"success": map[string]interface{}{
			"message": "Successfully uploaded package.",
		},
	}, "application/vnd.pub.v2+json")
}

// handlers end

func (controller *pubController) handleControllerError(ctx *fiber.Ctx, currentPath string, err error) error {
	if err == fiber.ErrNotFound {
		if url := controller.service.GetUpstreamUrl(currentPath); url != nil {
			ctx.Redirect(*url, fiber.StatusFound)
		}
		return ctx.Status(fiber.StatusNotFound).JSON(map[string]interface{}{
			"error": map[string]interface{}{
				"code":    strconv.Itoa(fiber.StatusNotFound),
				"message": "Not Found",
			},
		}, "application/vnd.pub.v2+json")
	}

	if err == fiber.ErrForbidden {
		return ctx.Status(fiber.StatusForbidden).JSON(map[string]interface{}{
			"error": map[string]interface{}{
				"code":    strconv.Itoa(fiber.StatusForbidden),
				"message": "Forbidden",
			},
		}, "application/vnd.pub.v2+json")
	}

	return ctx.Status(fiber.StatusBadRequest).JSON(map[string]interface{}{
		"error": map[string]interface{}{
			"code":    strconv.Itoa(fiber.StatusBadRequest),
			"message": err.Error(),
		},
	}, "application/vnd.pub.v2+json")
}
