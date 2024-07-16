package custom

import (
	"fmt"
	"github.com/backend-timedoor/gtimekeeper-framework/base/database"
	"github.com/backend-timedoor/gtimekeeper-framework/container"
	"github.com/go-playground/validator/v10"
	"strings"
)

type ExistsValidator struct{}

func (u *ExistsValidator) Signature() string {
	return "exists"
}

func (u *ExistsValidator) Handle(fl validator.FieldLevel) bool {
	params := strings.Split(fl.Param(), ".")
	fieldValue := fl.Field().Int()

	if len(params) < 2 {
		return true
	}

	tableName, fieldName := params[0], params[1]

	db := container.Get(database.ContainerName).(*database.Database)

	var count int64

	query := db.DB.Table(tableName).
		Order("id desc").
		Where(fmt.Sprintf("%s = ?", fieldName), fieldValue).
		Where("deleted_at IS NULL")

	query.Count(&count)

	return count > 0
}
