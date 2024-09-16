package user

import (
	"private-pub-repo/modules/app"
	"private-pub-repo/modules/app/appmodel"
	"private-pub-repo/modules/jwt"
	"private-pub-repo/modules/user/userdto"
	"private-pub-repo/utils"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

const (
	validationError = "Validation Error"
)

type userController struct {
	service         UserService
	responseService app.ResponseService
	validator       *validator.Validate
}

func newUserController(service UserService, responseService app.ResponseService, validator *validator.Validate) *userController {
	return &userController{
		service:         service,
		responseService: responseService,
		validator:       validator,
	}
}

// handlers start

func (controller *userController) handleCreate(ctx *fiber.Ctx) error {
	request := userdto.RegisterDTO{}
	ctx.BodyParser(&request)
	err := controller.validator.Struct(request)

	if err != nil {
		return controller.responseService.SendValidationErrorResponse(ctx, 400, validationError, err.(validator.ValidationErrors))
	}

	model, err := controller.service.Insert(ctx.UserContext(), request.ToModel())

	if err != nil {
		return fiber.NewError(400, err.Error())
	}

	return controller.responseService.SendSuccessResponse(ctx, 201, model)
}

func (controller *userController) handleLogin(ctx *fiber.Ctx) error {
	request := userdto.LoginDTO{}
	ctx.BodyParser(&request)
	err := controller.validator.Struct(request)

	if err != nil {
		return controller.responseService.SendValidationErrorResponse(ctx, 400, validationError, err.(validator.ValidationErrors))
	}

	response, err := controller.service.Login(ctx.UserContext(), &request)

	if err != nil {
		return fiber.NewError(400, err.Error())
	}
	return controller.responseService.SendSuccessResponse(ctx, 200, response)
}

func (controller *userController) handleProfile(ctx *fiber.Ctx) error {
	var user *userdto.UserDTO
	var err error

	id, err := utils.GetFiberJwtUserId(ctx)

	if err == nil {
		user, err = controller.service.Detail(ctx.UserContext(), id)
	}

	if err != nil {
		return fiber.NewError(400, err.Error())
	}

	user.CreatedAt = nil

	return controller.responseService.SendSuccessResponse(ctx, 200, user)
}

func (controller *userController) handleRefresh(ctx *fiber.Ctx) error {
	var user *jwt.JWTTokenModel
	var err error

	user, err = controller.service.RefreshToken(ctx.UserContext(), utils.GetFiberJwtClaims(ctx))

	if err != nil {
		return fiber.NewError(400, err.Error())
	}

	return controller.responseService.SendSuccessResponse(ctx, 200, user)
}

func (controller *userController) handleList(ctx *fiber.Ctx) error {
	request := appmodel.NewGetListRequest(ctx.Query("page"), ctx.Query("limit"), ctx.Query("search"))
	err := controller.validator.Struct(request)

	if err != nil {
		return controller.responseService.SendValidationErrorResponse(ctx, 400, validationError, err.(validator.ValidationErrors))
	}

	list, err := controller.service.List(ctx.UserContext(), request)

	if err != nil {
		return fiber.NewError(400, err.Error())
	}

	return controller.responseService.SendSuccessResponse(ctx, 200, appmodel.PaginationResponse{
		List: list,
	})
}

func (controller *userController) handleDetail(ctx *fiber.Ctx) error {
	userId, err := uuid.Parse(ctx.Params("id"))

	if err != nil {
		return fiber.NewError(400, err.Error())
	}

	user, err := controller.service.Detail(ctx.UserContext(), userId)

	if err != nil {
		return fiber.NewError(400, err.Error())
	}
	return controller.responseService.SendSuccessResponse(ctx, 200, user)
}

func (controller *userController) handleUpdate(ctx *fiber.Ctx) error {
	request := userdto.UpdateUserDTO{}
	ctx.BodyParser(&request)
	err := controller.validator.Struct(request)

	if err != nil {
		return controller.responseService.SendValidationErrorResponse(ctx, 400, validationError, err.(validator.ValidationErrors))
	}

	userId, err := uuid.Parse(ctx.Params("id"))

	if err != nil {
		return fiber.NewError(400, err.Error())
	}

	user, err := controller.service.Update(ctx.UserContext(), userId, &request)

	if err != nil {
		return fiber.NewError(400, err.Error())
	}
	return controller.responseService.SendSuccessResponse(ctx, 200, user)
}

func (controller *userController) handleDelete(ctx *fiber.Ctx) error {
	userId, err := uuid.Parse(ctx.Params("id"))

	if err != nil {
		return fiber.NewError(400, err.Error())
	}

	err = controller.service.Delete(ctx.UserContext(), userId)

	if err != nil {
		return fiber.NewError(400, err.Error())
	}
	return controller.responseService.SendSuccessResponse(ctx, 200, nil)
}

// handlers end
