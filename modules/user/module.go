package user

import (
	"gofiber-boilerplate/base"
	"gofiber-boilerplate/modules/app"
	"gofiber-boilerplate/modules/db"
	"gofiber-boilerplate/modules/jwt"
	"gofiber-boilerplate/modules/monitor"
	"gofiber-boilerplate/modules/user/usermodel"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
)

type UserModule struct {
	Service    UserService
	Middleware UserJwtMiddleware
	controller *userController
	jwtService jwt.JwtService
	db         db.DbService
	app        *fiber.App
}

func NewModule(service UserService, middleware UserJwtMiddleware, controller *userController, jwtService jwt.JwtService, db db.DbService, app *fiber.App) *UserModule {
	return &UserModule{Service: service, Middleware: middleware, jwtService: jwtService, controller: controller, db: db, app: app}
}

func fxRegister(lifeCycle fx.Lifecycle, module *UserModule) {
	base.FxRegister(module, lifeCycle)
}

func SetupModule(app *app.AppModule, db *db.DbModule, jwt *jwt.JwtModule, monitor *monitor.MonitorModule) *UserModule {
	service := NewUserService(jwt, monitor.Service)
	middleware := NewUserJwtMiddleware(jwt, monitor.Service)
	controller := newUserController(service, app.ResponseService, app.Validator)
	return NewModule(service, middleware, controller, jwt, db, app.App)
}

var FxModule = fx.Module("User", fx.Provide(NewUserService), fx.Provide(NewUserJwtMiddleware), fx.Provide(newUserController), fx.Provide(NewModule), fx.Invoke(fxRegister))

// implements `BaseModule` of `base/module.go` start

func (module *UserModule) OnStart() error {
	if module.db.AutoMigrate() {
		module.db.Default().AutoMigrate(&usermodel.UserModel{})
	}

	//run seeder
	module.Service.Init(module.db)
	module.registerRoutes()
	return nil
}

func (module *UserModule) OnStop() error {
	return nil
}

// implements `BaseModule` of `base/module.go` end
