package mail

import (
	"private-pub-repo/base"
	"private-pub-repo/modules/config"
	"strconv"

	"gopkg.in/gomail.v2"

	"go.uber.org/fx"
)

type MailModule struct {
	dialer    *gomail.Dialer
	fromName  string
	fromEmail string
}

func NewModule(config config.ConfigService) *MailModule {
	smtpHost := config.Getenv("SMTP_HOST", "")
	smtpPort, err := strconv.Atoi(config.Getenv("SMTP_PORT", "587"))
	if err != nil {
		smtpPort = 587
	}
	smtpUsername := config.Getenv("SMTP_USERNAME", "")
	smtpPassword := config.Getenv("SMTP_PASSWORD", "")
	dialer := gomail.NewDialer(smtpHost, smtpPort, smtpUsername, smtpPassword)

	return &MailModule{
		dialer:    dialer,
		fromName:  config.Getenv("SMTP_FROM_NAME", ""),
		fromEmail: config.Getenv("SMTP_FROM_EMAIL", ""),
	}
}

func ProvideService(module *MailModule) MailService {
	return module
}

func fxRegister(lifeCycle fx.Lifecycle, module *MailModule) {
	base.FxRegister(module, lifeCycle)
}

func SetupModule(config *config.ConfigModule) *MailModule {
	return NewModule(config)
}

var FxModule = fx.Module("Mail", fx.Provide(NewModule), fx.Provide(ProvideService), fx.Invoke(fxRegister))

// implements `BaseModule` of `base/module.go` start

func (module *MailModule) OnStart() error {
	return nil
}

func (module *MailModule) OnStop() error {
	return nil
}

// implements `BaseModule` of `base/module.go` end
