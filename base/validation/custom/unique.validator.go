package custom

import (
    "fmt"
    "github.com/backend-timedoor/gtimekeeper-framework/base/database"
    "github.com/backend-timedoor/gtimekeeper-framework/container"
    "github.com/go-playground/validator/v10"
    "strings"
    "time"
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
        excludeID, excludeCol = params[2], params[3]
    }
    var excludeVal int64
    if excludeID != "" {
        excludeVal = fl.Parent().FieldByName(excludeID).Int()
    }

    var record struct {
        ID        int32
        DeletedAt time.Time
    }
    db := container.Get(database.ContainerName).(*database.Database)
    query := db.DB.Table(tableName).
        Order("id desc").
        Where(fmt.Sprintf("%s = ?", fieldName), fieldValue)
    if excludeCol != "" && excludeVal != 0 {
        query = query.Where(fmt.Sprintf("%s != ?", excludeCol), excludeVal)
    }
    query.First(&record)

    if record.ID == 0 {
        return true
    }

    if !record.DeletedAt.IsZero() {
        return true
    }

    return false
}
