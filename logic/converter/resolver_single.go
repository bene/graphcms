package converter

import (
	"database/sql"
	"fmt"
	"github.com/graphql-go/graphql"
	_ "github.com/lib/pq"
	"log"
	"strings"
)

func MaskToSingleField(object *graphql.Object) (graphql.Field, error) {

	db, err := sql.Open("postgres", "postgres://postgres:mysecretpassword@localhost:5431/evobend?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	return graphql.Field{
		Type: object,
		Resolve: func(p graphql.ResolveParams) (i interface{}, err error) {

			tableName := p.Info.FieldName
			tableFields := []string{fmt.Sprintf("%s.id", tableName)}
			tableJoins := []string{}

			selections := getSelectedFields([]string{p.Info.FieldName}, p)
			for _, s := range selections {
				if !s.HasSubSelect {
					tableFields = append(tableFields, fmt.Sprintf("%s.%s", tableName, s.Name))
				} else {

					// TODO: Make recursive for deeper levels
					subSelections := getSelectedFields([]string{tableName, s.Name}, p)
					for _, ss := range subSelections {
						tableFields = append(tableFields, fmt.Sprintf("%s.%s AS %s___%s", s.Name, ss.Name, s.Name, ss.Name))
					}
					tableJoins = append(tableJoins, fmt.Sprintf(" LEFT JOIN %s_%s %s ON %s.id = %s.%s_id", tableName, s.Name, s.Name, tableName, s.Name, tableName))
				}
			}

			query := fmt.Sprintf("SELECT %s FROM %s%s", strings.Join(tableFields, ", "), tableName, strings.Join(tableJoins, " "))

			rows, err := db.QueryContext(p.Context, query)
			if err != nil {
				return nil, err
			}
			defer rows.Close()

			cols, err := rows.Columns()
			if err != nil {
				return nil, err
			}

			results := map[string][]map[string]interface{}{}
			for rows.Next() {

				columns := make([]interface{}, len(cols))
				columnPointers := make([]interface{}, len(cols))
				for i, _ := range columns {
					columnPointers[i] = &columns[i]
				}

				// Scan the result into the column pointers...
				if err := rows.Scan(columnPointers...); err != nil {
					return nil, err
				}

				m := make(map[string]interface{})
				for i, colName := range cols {
					val := columnPointers[i].(*interface{})
					m[colName] = *val
				}

				id := fmt.Sprint(m["id"])

				if _, ok := results[id]; ok {
					results[id] = append(results[id], m)
				} else {
					results[id] = []map[string]interface{}{m}
				}
			}

			reduceds := []map[string]interface{}{}

			for _, one := range results {
				reduceds = append(reduceds, reduceRowsOfOne(one))
			}

			// TODO: Change to == 1 when args support built in
			if len(reduceds) >= 1 {
				return reduceds[0], nil
			}

			return nil, nil
		},
	}, nil
}
