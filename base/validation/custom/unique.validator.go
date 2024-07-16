package custom

import (
	"fmt"
	"github.com/backend-timedoor/gtimekeeper-framework/base/database"
	"github.com/backend-timedoor/gtimekeeper-framework/container"
	"github.com/go-playground/validator/v10"
	"strings"
)

type UniqueValidator struct{}

func (u *UniqueValidator) Signature() string {
	return "unique"
}

func (u *UniqueValidator) Handle(fl validator.FieldLevel) bool {
	params := strings.Split(fl.Param(), ".")
	fieldValue := fl.Field().String()
	if len(params) < 2 {
		return true
	}

	tableName, fieldName := params[0], params[1]

	var excludeID, excludeCol string
	if len(params) >= 4 {
		excludeCol, excludeID = params[2], params[3]
	}

	db := container.Get(database.ContainerName).(*database.Database)

	var count int64

	query := db.DB.Table(tableName).
		Order("id desc").
		Where(fmt.Sprintf("%s = ?", fieldName), fieldValue)

	if excludeCol != "" && excludeID != "" {
		query = query.Where(fmt.Sprintf("%s != ?", excludeCol), excludeID).
			Where("deleted_at IS NULL")
	}
	query.Count(&count)

	return count <= 0
}
