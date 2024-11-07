package pubtoken

import (
	"private-pub-repo/modules/app"
	"private-pub-repo/modules/app/appmodel"
	"private-pub-repo/modules/pubtoken/pubtokendto"
	"private-pub-repo/utils"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

const (
	validationError = "Validation Error"
)

type pubTokenController struct {
	service         PubTokenService
	responseService app.ResponseService
	validator       *validator.Validate
}

func newPubTokenController(service PubTokenService, responseService app.ResponseService, validator *validator.Validate) *pubTokenController {
	return &pubTokenController{
		service:         service,
		responseService: responseService,
		validator:       validator,
	}
}

// handlers start

func (controller *pubTokenController) handleCreate(ctx *fiber.Ctx) error {
	var err error

	id, err := utils.GetFiberJwtUserId(ctx)

	if err != nil {
		return fiber.NewError(400, err.Error())
	}

	request := pubtokendto.CreateTokenDTO{}
	ctx.BodyParser(&request)
	err = controller.validator.Struct(request)

	if err != nil {
		return controller.responseService.SendValidationErrorResponse(ctx, 400, validationError, err.(validator.ValidationErrors))
	}

	token, err := controller.service.Insert(ctx.UserContext(), request.ToModel(id, utils.IsFiberJwtCanWrite(ctx)))

	if err != nil {
		return fiber.NewError(400, err.Error())
	}

	return controller.responseService.SendSuccessDetailResponse(ctx, 201, token)
}

func (controller *pubTokenController) handleList(ctx *fiber.Ctx) error {
	request := appmodel.NewGetListRequest(ctx.Query("page"), ctx.Query("limit"), ctx.Query("search"))
	err := controller.validator.Struct(request)

	if err != nil {
		return controller.responseService.SendValidationErrorResponse(ctx, 400, validationError, err.(validator.ValidationErrors))
	}

	userId, err := utils.GetFiberJwtUserId(ctx)

	if err != nil {
		return fiber.NewError(400, err.Error())
	}

	list, err := controller.service.List(ctx.UserContext(), request, &userId)

	if err != nil {
		return fiber.NewError(400, err.Error())
	}

	return controller.responseService.SendSuccessResponse(ctx, 200, appmodel.PaginationResponse{
		List: list,
	})
}

func (controller *pubTokenController) handleDetail(ctx *fiber.Ctx) error {
	tokenId, err := uuid.Parse(ctx.Params("id"))

	if err != nil {
		return fiber.NewError(400, err.Error())
	}

	userId, err := utils.GetFiberJwtUserId(ctx)

	if err != nil {
		return fiber.NewError(400, err.Error())
	}

	pubToken, err := controller.service.Detail(ctx.UserContext(), tokenId, &userId)

	if err != nil {
		return fiber.NewError(400, err.Error())
	}
	return controller.responseService.SendSuccessDetailResponse(ctx, 200, pubToken)
}

func (controller *pubTokenController) handleUpdate(ctx *fiber.Ctx) error {
	tokenId, err := uuid.Parse(ctx.Params("id"))

	if err != nil {
		return fiber.NewError(400, err.Error())
	}

	userId, err := utils.GetFiberJwtUserId(ctx)

	if err != nil {
		return fiber.NewError(400, err.Error())
	}

	request := pubtokendto.UpdateTokenDTO{}
	ctx.BodyParser(&request)
	err = controller.validator.Struct(request)

	if err != nil {
		return controller.responseService.SendValidationErrorResponse(ctx, 400, validationError, err.(validator.ValidationErrors))
	}

	canWrite := utils.IsFiberJwtCanWrite(ctx) && *request.Write

	request.Write = &canWrite

	pubToken, err := controller.service.Update(ctx.UserContext(), tokenId, &userId, &request)

	if err != nil {
		return fiber.NewError(400, err.Error())
	}
	return controller.responseService.SendSuccessDetailResponse(ctx, 200, pubToken)
}

func (controller *pubTokenController) handleDelete(ctx *fiber.Ctx) error {
	tokenId, err := uuid.Parse(ctx.Params("id"))

	if err != nil {
		return fiber.NewError(400, err.Error())
	}

	userId, err := utils.GetFiberJwtUserId(ctx)

	if err != nil {
		return fiber.NewError(400, err.Error())
	}

	err = controller.service.Delete(ctx.UserContext(), tokenId, &userId)

	if err != nil {
		return fiber.NewError(400, err.Error())
	}
	return controller.responseService.SendSuccessDetailResponse(ctx, 200, nil)
}

// handlers end
