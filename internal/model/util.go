package model

import (
	"reflect"
	"strings"

	"github.com/pkg/errors"
)

type Table struct {
	Name       string
	Columns    []string
	Insertable []string
	Updatable  []string
}

type Column struct {
	Name      string
	Generated bool
	Immutable bool
}

func StructToTable(target any, tableName string) (*Table, error) {
	var structValue reflect.Value

	value := reflect.ValueOf(target)
	table := Table{Name: tableName}
	columns := make([]Column, 0)

	if !isStructOrStructPointer(value) {
		return nil, errors.New("repository: target must be a struct or a struct pointer")
	}

	if value.Kind() == reflect.Ptr {
		structValue = value.Elem()
	} else {
		structValue = value
	}

	for i := range structValue.NumField() {
		field := structValue.Type().Field(i)
		tag := field.Tag.Get("db")

		if tag == "" {
			return nil, errors.Errorf("repository: missing db tag on field %s", field.Name)
		}

		column, err := parseTag(tag, field.Name)
		if err != nil {
			return nil, err
		}

		columns = append(columns, column)
	}

	table = columnsToTable(columns)

	return &table, nil
}

func parseTag(tag string, field string) (Column, error) {
	parts := strings.Split(tag, ",")
	column := Column{Name: parts[0]}

	for _, part := range parts[1:] {
		switch part {
		case "generated":
			column.Generated = true
		case "immutable":
			column.Immutable = true
		default:
			return Column{}, errors.Errorf("repository: invalid tag option %s for field %s", part, field)
		}
	}

	return column, nil
}

func columnsToTable(columns []Column) Table {
	var table Table

	for _, column := range columns {
		table.Columns = append(table.Columns, column.Name)

		if !column.Generated {
			table.Insertable = append(table.Insertable, column.Name)
		}

		if !column.Immutable {
			table.Updatable = append(table.Updatable, column.Name)
		}
	}

	return table
}

func isStructOrStructPointer(value reflect.Value) bool {
	return value.Kind() == reflect.Struct ||
		(value.Kind() == reflect.Ptr && value.Elem().Kind() == reflect.Struct)
}
