package db_test_tools

import (
	"context"
	"fmt"
	"reflect"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
)

type (
	// TableEntity -
	TableEntity interface {
		TableName() string
	}
)

var pgQb = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

// Select -
func Select[T TableEntity, S ~[]T](ctx context.Context, entities *S) error {
	if entities == nil {
		initializedEntities := make(S, 0)
		entities = &initializedEntities
	}

	qb := pgQb.Select(columns(*entities)...).
		From(table(entities))

	query, args, err := qb.ToSql()
	if err != nil {
		return err
	}

	rows, err := PgxPool.Query(ctx, query, args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	values, err := pgx.CollectRows(rows, pgx.RowToStructByName[T])
	if err != nil {
		return err
	}

	*entities = (*entities)[:0]
	*entities = append(*entities, values...)

	return nil
}

// Insert -
func Insert[T TableEntity, S ~[]T](ctx context.Context, entities S) error {
	if len(entities) == 0 {
		return nil
	}

	tableName := entities[0].TableName()
	qb := pgQb.Insert(tableName).Columns(columns(entities)...)
	for _, entity := range entities {
		qb = qb.Values(values(entity)...)
	}

	query, args, err := qb.ToSql()
	if err != nil {
		return err
	}

	_, err = PgxPool.Exec(ctx, query, args...)
	return err
}

// TruncateAll -
func TruncateAll(ctx context.Context) error {
	const (
		tablesQuery = `
            SELECT tablename
            FROM pg_catalog.pg_tables
            WHERE schemaname = 'public' AND NOT tablename LIKE '%goose%';
        `
		truncateQuery = `TRUNCATE TABLE %s;`
	)

	rows, err := PgxPool.Query(ctx, tablesQuery)
	if err != nil {
		return err
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var table string
		if err := rows.Scan(&table); err != nil {
			return err
		}

		tables = append(tables, table)
	}

	for _, table := range tables {
		query := fmt.Sprintf(truncateQuery, pgx.Identifier{table}.Sanitize())
		if _, err := PgxPool.Exec(ctx, query); err != nil {
			return err
		}
	}

	return nil
}

func table[T TableEntity, S ~[]T](entities *S) string {
	typeOf := reflect.TypeOf(entities).Elem().Elem()
	value := reflect.New(typeOf)
	tableEntity := value.Interface().(TableEntity)
	return tableEntity.TableName()
}

func columns[T any, S ~[]T](S) []string {
	var value T
	t := reflect.TypeOf(value)
	columnNames := make([]string, 0)

	for i := range t.NumField() {
		dbTag, ok := t.Field(i).Tag.Lookup("db")
		if !ok || dbTag == "-" {
			continue
		}

		columnNames = append(columnNames, dbTag)
	}

	return columnNames
}

func values[T any](value T) []any {
	values := make([]any, 0)
	t := reflect.TypeOf(value)
	v := reflect.ValueOf(value)

	for i := range t.NumField() {
		dbTag, ok := t.Field(i).Tag.Lookup("db")
		if !ok || dbTag == "-" {
			continue
		}

		values = append(values, v.Field(i).Interface())
	}
	return values
}
