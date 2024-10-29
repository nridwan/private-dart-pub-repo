package pubtoken

import (
	"private-pub-repo/base"
	"private-pub-repo/modules/app"
	"private-pub-repo/modules/db"
	"private-pub-repo/modules/jwt"
	"private-pub-repo/modules/monitor"
	"private-pub-repo/modules/pubtoken/pubtokenmodel"
	"private-pub-repo/modules/user"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
)

type PubTokenModule struct {
	Service        PubTokenService
	Middleware     PubTokenJwtMiddleware
	userMiddleware user.UserJwtMiddleware
	controller     *pubTokenController
	jwtService     jwt.JwtService
	db             db.DbService
	app            *fiber.App
}

func NewModule(service PubTokenService, middleware PubTokenJwtMiddleware, controller *pubTokenController, jwtService jwt.JwtService, db db.DbService, userMiddleware user.UserJwtMiddleware, app *fiber.App) *PubTokenModule {
	return &PubTokenModule{Service: service, Middleware: middleware, userMiddleware: userMiddleware, jwtService: jwtService, controller: controller, db: db, app: app}
}

func fxRegister(lifeCycle fx.Lifecycle, module *PubTokenModule) {
	base.FxRegister(module, lifeCycle)
}

func SetupModule(app *app.AppModule, db *db.DbModule, user *user.UserModule, jwt *jwt.JwtModule, monitor *monitor.MonitorModule) *PubTokenModule {
	service := NewPubTokenService(jwt, monitor.Service)
	middleware := NewPubTokenJwtMiddleware(jwt, service, monitor.Service)
	controller := newPubTokenController(service, app.ResponseService, app.Validator)
	return NewModule(service, middleware, controller, jwt, db, user.Middleware, app.App)
}

var FxModule = fx.Module("PubToken", fx.Provide(NewPubTokenService), fx.Provide(NewPubTokenJwtMiddleware), fx.Provide(newPubTokenController), fx.Provide(NewModule), fx.Invoke(fxRegister))

// implements `BaseModule` of `base/module.go` start

func (module *PubTokenModule) OnStart() error {
	if module.db.AutoMigrate() {
		module.db.Default().AutoMigrate(&pubtokenmodel.PubTokenModel{})
	}

	//run seeder
	module.Service.Init(module.db)
	module.registerRoutes()
	return nil
}

func (module *PubTokenModule) OnStop() error {
	return nil
}

// implements `BaseModule` of `base/module.go` end
