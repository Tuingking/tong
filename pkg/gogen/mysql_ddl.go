package gogen

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"github.com/Tuingking/tong/pkg/gofparser"
	"github.com/iancoleman/strcase"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

var GoToMysqlTypeMap = map[string]string{
	"string":         "varchar(255) NOT NULL",
	"sql.NullString": "varchar(255) NULL",

	"int":           "int(11) NOT NULL",
	"int32":         "int(11) NOT NULL",
	"int64":         "int(11) NOT NULL",
	"sql.NullInt32": "int(11) NULL",
	"sql.NullInt64": "int(11) NULL",

	"float32":         "decimal(12,2) NOT NULL",
	"float64":         "decimal(12,2) NOT NULL",
	"sql.NullFloat64": "decimal(12,2) NULL",

	"time.Time":    "timestamp(6) NOT NULL",
	"sql.NullTime": "timestamp(6) NULL",
}

func GenerateDDL(filename, output string) (string, error) {
	// get current dirrectory
	wd, err := os.Getwd()
	if err != nil {
		logrus.Fatal(errors.Wrap(err, "get working dir"))
	}
	filepath := path.Join(wd, filename)
	logrus.Infof("filePath: %+v\n", filepath)

	// parse go file and get list of struct
	structList, err := gofparser.GetStructDecl(filepath)
	if err != nil {
		logrus.Fatal(errors.Wrap(err, "get struct declaration"))
	}

	// interate to available struct and construct DDL
	buf := new(bytes.Buffer)
	for sIdx, s := range structList {
		tableName := strcase.ToSnake(s.Name)
		if strings.HasSuffix(tableName, "_param") || strings.HasPrefix(tableName, "param_") {
			continue
		}

		if sIdx > 0 {
			buf.WriteString("\n\n")
		}

		var (
			pk string
			uq string
			ix string
		)

		buf.WriteString(fmt.Sprintf("CREATE TABLE IF NOT EXISTS `%s`(\n", tableName))
		for fIdx, field := range s.Fields {
			// add indentation
			buf.WriteRune('\t')

			// field
			fName := strcase.ToSnake(field.Name)
			buf.WriteString("`" + fName + "`")

			// type
			if fTyp, found := GoToMysqlTypeMap[field.Type]; found {
				buf.WriteString(" " + fTyp)
			} else {
				buf.WriteString(" " + field.Type)
			}

			// auto increment
			if field.Name == "ID" && field.Type == "int64" {
				buf.WriteString(" AUTO_INCREMENT")
			}

			// parse struct tag (primaryKey, indexKey, uniqueKey)
			tag := field.ParseTagAndGet("orm")
			for _, v := range strings.Split(tag, ";") {
				switch v {
				case "primaryKey":
					pk = fName
				case "uniqueKey":
					uq = fName
				case "indexKey":
					ix = fName
				}
			}

			// if last
			isLast := fIdx == len(s.Fields)-1
			if !isLast {
				buf.WriteString(",\n")
			}
		}

		// primary key
		if pk != "" {
			buf.WriteString(",\n\tPRIMARY KEY (`" + pk + "`)")
		}

		// unique key
		if uq != "" {
			buf.WriteString(fmt.Sprintf(",\n\tCONSTRAINT `%s_%s_uq` UNIQUE (`%s`)", tableName, uq, uq))
		}

		// index
		if ix != "" {
			buf.WriteString(fmt.Sprintf(",\n\tKEY `%s_%s_ix` (`%s`)", tableName, ix, ix))
		}

		buf.WriteString("\n)ENGINE=InnoDB;")
	}

	result := buf.String()

	// write to output
	if output != "" {
		outName := fmt.Sprintf("migration__%d.sql", time.Now().Unix())
		out, err := os.Create(outName)
		if err != nil {
			logrus.Error(errors.Wrap(err, "create file"))
		}
		defer out.Close()
		buf.WriteTo(out)
		logrus.Infof("check file %s in current dir ", outName)
	}

	return result, nil
}
