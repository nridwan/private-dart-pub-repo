package storage

import (
	"private-pub-repo/base"
	"private-pub-repo/modules/config"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"go.uber.org/fx"
)

type StorageModule struct {
	s3            *s3.S3
	s3Public      *s3.S3
	uploader      *s3manager.Uploader
	bucket        string
	enablePresign bool
	presignTime   int
}

func NewModule(config config.ConfigService) *StorageModule {
	endpoint := aws.String(config.Getenv("S3_ENDPOINT", ""))
	publicEndpoint := aws.String(config.Getenv("S3_PUBLIC_ENDPOINT", *endpoint))
	region := aws.String(config.Getenv("S3_REGION", ""))
	credentials := credentials.NewStaticCredentials(
		config.Getenv("S3_KEY_ID", ""),
		config.Getenv("S3_ACCESS_KEY", ""),
		"",
	)
	usePathStyle := aws.Bool(config.Getenv("S3_USE_PATH_STYLE", "") == "true")

	s3Session, err := session.NewSession(&aws.Config{
		Endpoint:         endpoint,
		Region:           region,
		Credentials:      credentials,
		S3ForcePathStyle: usePathStyle,
	})

	if err != nil {
		panic(err)
	}

	s3PublicSession, err := session.NewSession(&aws.Config{
		Endpoint:         publicEndpoint,
		Region:           region,
		Credentials:      credentials,
		S3ForcePathStyle: usePathStyle,
	})

	if err != nil {
		panic(err)
	}

	presignTime, err := strconv.Atoi(config.Getenv("S3_PRESIGN_TIME", "15"))

	if err != nil {
		presignTime = 15
	}

	return &StorageModule{
		s3:            s3.New(s3Session),
		s3Public:      s3.New(s3PublicSession),
		uploader:      s3manager.NewUploader(s3Session),
		bucket:        config.Getenv("S3_BUCKET", ""),
		enablePresign: config.Getenv("S3_ENABLE_PRESIGN", "false") == "true",
		presignTime:   presignTime,
	}
}

func ProvideService(module *StorageModule) StorageService {
	return module
}

func fxRegister(lifeCycle fx.Lifecycle, module *StorageModule) {
	base.FxRegister(module, lifeCycle)
}

func SetupModule(config *config.ConfigModule) *StorageModule {
	return NewModule(config)
}

var FxModule = fx.Module("Storage", fx.Provide(NewModule), fx.Provide(ProvideService), fx.Invoke(fxRegister))

// implements `BaseModule` of `base/module.go` start

func (module *StorageModule) OnStart() error {
	return nil
}

func (module *StorageModule) OnStop() error {
	return nil
}

// implements `BaseModule` of `base/module.go` end
