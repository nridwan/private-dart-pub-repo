package pubtoken

import (
	"private-pub-repo/base"
	"private-pub-repo/modules/app"
	"private-pub-repo/modules/db"
	"private-pub-repo/modules/jwt"
	"private-pub-repo/modules/monitor"
	"private-pub-repo/modules/user"
	"private-pub-repo/modules/user/usermodel"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
)

type UserModule struct {
	Service        PubTokenService
	Middleware     PubTokenJwtMiddleware
	userMiddleware user.UserJwtMiddleware
	controller     *pubTokenController
	jwtService     jwt.JwtService
	db             db.DbService
	app            *fiber.App
}

func NewModule(service PubTokenService, middleware PubTokenJwtMiddleware, controller *pubTokenController, jwtService jwt.JwtService, db db.DbService, userMiddleware user.UserJwtMiddleware, app *fiber.App) *UserModule {
	return &UserModule{Service: service, Middleware: middleware, userMiddleware: userMiddleware, jwtService: jwtService, controller: controller, db: db, app: app}
}

func fxRegister(lifeCycle fx.Lifecycle, module *UserModule) {
	base.FxRegister(module, lifeCycle)
}

func SetupModule(app *app.AppModule, db *db.DbModule, user *user.UserModule, jwt *jwt.JwtModule, monitor *monitor.MonitorModule) *UserModule {
	service := NewPubTokenService(jwt, monitor.Service)
	middleware := NewUserJwtMiddleware(jwt, service, monitor.Service)
	controller := newPubTokenController(service, app.ResponseService, app.Validator)
	return NewModule(service, middleware, controller, jwt, db, user.Middleware, app.App)
}

var FxModule = fx.Module("User", fx.Provide(NewPubTokenService), fx.Provide(NewUserJwtMiddleware), fx.Provide(newPubTokenController), fx.Provide(NewModule), fx.Invoke(fxRegister))

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
