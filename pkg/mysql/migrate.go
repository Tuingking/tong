package mysql

import (
	"bytes"
	"database/sql"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/iancoleman/strcase"
)

// tag: 'orm'
// value:
// - autoIncrement
// - primaryKey
// - unique
// - index

func Migrate(schema interface{}) {
	val := reflect.ValueOf(schema)

	structName := val.Type().Name()
	tableName := strcase.ToSnake(structName)

	buf := new(bytes.Buffer)
	bufIdx := new(bytes.Buffer)
	buf.WriteString(fmt.Sprintf("CREATE TABLE IF NOT EXISTS `%s` (\n", tableName))
	for i := 0; i < val.NumField(); i++ {
		fieldName := strcase.ToSnake(val.Type().Field(i).Name)
		fieldVal := val.Field(i).Interface()
		ormTag := val.Type().Field(i).Tag.Get("orm") // tag: "orm"

		autoIncrement := ""
		if ormTag != "" {
			tags := strings.Split(ormTag, ";")
			for _, tag := range tags {
				if bufIdx.Len() > 0 {
					bufIdx.WriteString(",\n")
				}

				switch tag {
				case "autoIncrement":
					autoIncrement = " AUTO_INCREMENT"
				case "primaryKey":
					bufIdx.WriteString(fmt.Sprintf("\tPRIMARY KEY (`%s`)", fieldName))
				case "unique":
					bufIdx.WriteString(fmt.Sprintf("\tCONSTRAINT `%s_%s_uq` UNIQUE (%s)", tableName, fieldName, fieldName))
				case "index":
					bufIdx.WriteString(fmt.Sprintf("\tKEY `%s_%s_ix` (%s)", tableName, fieldName, fieldName))
				}
			}
		}

		var fieldType string
		switch fieldVal.(type) {
		case string:
			fieldType = "varchar(255) NOT NULL"
		case sql.NullString:
			fieldType = "varchar(255) NULL"
		case int, int32, int64:
			if autoIncrement != "" {
				fieldType = "int(11) NOT NULL" + autoIncrement
			} else {
				fieldType = "int(11) NOT NULL DEFAULT '0'" + autoIncrement
			}
		case sql.NullInt32, sql.NullInt64:
			fieldType = "int(11) NULL"
		case float32, float64:
			fieldType = "decimal(12,2) NOT NULL DEFAULT '0.00'"
		case sql.NullFloat64:
			fieldType = "decimal(12,2) NULL"
		case time.Time:
			fieldType = "timestamp(6) NOT NULL"
		case sql.NullTime:
			fieldType = "timestamp(6) NULL"
		}

		buf.WriteString(fmt.Sprintf("\t%-15s%s", fieldName, fieldType))

		if i != val.NumField()-1 {
			buf.WriteString(",\n")
		}
	}

	// primary key + index
	if bufIdx.Len() != 0 {
		buf.WriteString(",\n" + bufIdx.String())
	}

	buf.WriteString("\n) ENGINE=InnoDB;")

	fmt.Printf("[DEBUG] BUFFER RESULT: %+v\n", buf.String())
}
