package findfield

import (
	"fmt"

	"github.com/Tuingking/tong/pkg/logger"
	"github.com/Tuingking/tong/pkg/mysql"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

const howToUse = `tong sql find-field -f [FIELD_NAME]`

var (
	// flag
	fieldName string
)

var Cmd = &cobra.Command{
	Use:   "find-field",
	Short: "Find field in database table.",
	Long: fmt.Sprintf(`Find field in database table.
	example: %s`, howToUse),
	RunE: runE,
}

func init() {
	Cmd.PersistentFlags().StringVarP(&fieldName, "field", "f", "", "field name to search")
}

func runE(cmd *cobra.Command, args []string) error {
	return findField(fieldName)
}

func findField(name string) error {
	db := mysql.New(mysql.DefaultOption())

	query := `select 
		TABLE_SCHEMA,
		TABLE_NAME,
		COLUMN_NAME,
		COLUMN_TYPE,
		IS_NULLABLE 
	from INFORMATION_SCHEMA.COLUMNS 
	where COLUMN_NAME like ? and TABLE_SCHEMA not in ('information_schema', 'performance_schema', 'sys', 'mysql')
	order by TABLE_NAME`

	rows, err := db.Query(query, "%"+name+"%")
	if err != nil {
		logger.Logger.Error("failed exec query", zap.Error(err))
		return err
	}

	result := InformationSchema{}
	for rows.Next() {
		var res InformationSchemaColumn
		rows.Scan(
			&res.TableSchema,
			&res.TableName,
			&res.ColumnName,
			&res.ColumnType,
			&res.IsNullable,
		)

		result.Columns = append(result.Columns, res)
	}

	result.Print()

	return nil
}
