package user

import (
	"private-pub-repo/base"
	"private-pub-repo/modules/user/usermodel"

	"github.com/google/uuid"
)

func (module *UserModule) RunSeeder() {
	println("Insert User seeders...")
	service := module.Service
	db := module.db.Default()

	defaultPassword, err := service.GenerateHashPassword("Password")

	if err != nil {
		return
	}

	db.FirstOrCreate(&usermodel.UserModel{
		BaseModel: base.BaseModel{
			ID: uuid.MustParse("6d9b7354-b127-46dc-bcae-ff289c2bdcac"),
		},
		Name:     "Admin",
		Email:    "admin@m8zn.work",
		Password: defaultPassword,
		IsAdmin:  true,
		CanWrite: true,
	})
}
