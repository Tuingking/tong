package findfield

import (
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
)

type InformationSchema struct {
	Columns []InformationSchemaColumn
}

func (v *InformationSchema) Print() {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{
		"#",
		"Database",
		"Table",
		"Field",
		"Type",
		"Nullable",
	})
	t.AppendSeparator()
	for i, col := range v.Columns {
		t.AppendRow([]interface{}{
			i + 1,
			col.TableSchema,
			col.TableName,
			col.ColumnName,
			col.ColumnType,
			col.IsNullable,
		})
	}
	t.Render()
}

type InformationSchemaColumn struct {
	TableSchema string `db:"TABLE_SCHEMA"`
	TableName   string `db:"TABLE_NAME"`
	ColumnName  string `db:"COLUMN_NAME"`
	ColumnType  string `db:"COLUMN_TYPE"`
	IsNullable  string `db:"IS_NULLABLE"`
}
