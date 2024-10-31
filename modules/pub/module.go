package pub

import (
	"private-pub-repo/base"
	"private-pub-repo/modules/app"
	"private-pub-repo/modules/config"
	"private-pub-repo/modules/db"
	"private-pub-repo/modules/jwt"
	"private-pub-repo/modules/monitor"
	"private-pub-repo/modules/pub/pubmodel"
	"private-pub-repo/modules/pubtoken"
	"private-pub-repo/modules/storage"
	"private-pub-repo/modules/user"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
)

type PubModule struct {
	Service        PubService
	middleware     pubtoken.PubTokenJwtMiddleware
	userMiddleware user.UserJwtMiddleware
	controller     *pubController
	jwtService     jwt.JwtService
	db             db.DbService
	app            *fiber.App
}

func NewModule(service PubService, middleware pubtoken.PubTokenJwtMiddleware, userMiddleware user.UserJwtMiddleware, controller *pubController, jwtService jwt.JwtService, db db.DbService, app *fiber.App) *PubModule {
	return &PubModule{Service: service, middleware: middleware, userMiddleware: userMiddleware, jwtService: jwtService, controller: controller, db: db, app: app}
}

func fxRegister(lifeCycle fx.Lifecycle, module *PubModule) {
	base.FxRegister(module, lifeCycle)
}

func SetupModule(
	app *app.AppModule, db *db.DbModule, jwt *jwt.JwtModule, pubToken *pubtoken.PubTokenModule,
	user *user.UserModule, monitor *monitor.MonitorModule, config *config.ConfigModule,
	storage *storage.StorageModule,
) *PubModule {
	service := NewPubService(jwt, monitor.Service, config, storage)
	controller := newPubController(service, app.ResponseService, app.Validator, pubToken.Middleware, user.Middleware)
	return NewModule(service, pubToken.Middleware, user.Middleware, controller, jwt, db, app.App)
}

var FxModule = fx.Module("Pub", fx.Provide(NewPubService), fx.Provide(newPubController), fx.Provide(NewModule), fx.Invoke(fxRegister))

// implements `BaseModule` of `base/module.go` start

func (module *PubModule) OnStart() error {
	if module.db.AutoMigrate() {
		module.db.Default().AutoMigrate(&pubmodel.PubPackageModel{}, &pubmodel.PubVersionModel{})
	}

	//run seeder
	module.Service.Init(module.db)
	module.registerRoutes()
	return nil
}

func (module *PubModule) OnStop() error {
	return nil
}

// implements `BaseModule` of `base/module.go` end
