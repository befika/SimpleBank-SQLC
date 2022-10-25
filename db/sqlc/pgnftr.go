package db

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/jackc/pgx/v4/pgxpool"
)

func GetFilter(f *FilterParams) func(db *pgxpool.Conn) *pgxpool.Conn {
	log.Printf("FILTER: %#v", f)
	return func(db *pgxpool.Conn) *pgxpool.Conn {
		if !(f == nil) {

			where := ""
			q := ""
			f.LinkOperator = strings.ToUpper(f.LinkOperator)
			for _, filter := range f.Filter {
				if filter.ColumnField != "q" {
					switch filter.OperatorValue {
					case "=", "!=", ">", ">=", "<", "<=":
						where += fmt.Sprintf("%v %v '%v'", filter.ColumnField, filter.OperatorValue, filter.Value)
					case "is empty":
						where += fmt.Sprintf("%v = 'NULL'", filter.ColumnField)
					case "is not empty":
						where += fmt.Sprintf("%v = 'NOT NULL'", filter.ColumnField)
					case "contains":
						where += fmt.Sprintf("%v ILIKE '%v'", filter.ColumnField, "%"+filter.Value+"%")
					case "equals":
						where += fmt.Sprintf("%v = '%v'", filter.ColumnField, filter.Value)
					case "starts with":
						where += fmt.Sprintf("%v ILIKE '%v'", filter.ColumnField, filter.Value+"%")
					case "ends with":
						where += fmt.Sprintf("%v ILIKE '%v'", filter.ColumnField, "%"+filter.Value)
					case "is":
						where += fmt.Sprintf("%v = '%v'", filter.ColumnField, filter.Value)
					case "is not":
						where += fmt.Sprintf("%v != '%v'", filter.ColumnField, filter.Value)
					case "is after":
						where += fmt.Sprintf("%v > TIMESTAMPTZ '%v'", filter.ColumnField, filter.Value)
					case "is on or after":
						where += fmt.Sprintf("%v >= TIMESTAMPTZ '%v'", filter.ColumnField, filter.Value)
					case "is before":
						where += fmt.Sprintf("%v < TIMESTAMPTZ '%v'", filter.ColumnField, filter.Value)
					case "is on or before":
						where += fmt.Sprintf("%v <= TIMESTAMPTZ '%v'", filter.ColumnField, filter.Value)
					default:
						continue
					}
					where += " " + f.LinkOperator + " "
				} else {
					var searchQ Filter2Search
					err := json.Unmarshal([]byte(filter.Value), &searchQ)
					if err != nil {
						log.Println("err unmarshaling q: ", err)
						continue
					}
					for _, column := range searchQ.Columns {
						q += fmt.Sprintf("%v ILIKE '%v' OR ", column, "%"+searchQ.Value+"%")
					}
				}
			}
			if f == nil {
				f = &FilterParams{}
			}
			if f.Sort == nil {
				if !f.NoSort {
					f.Sort = append(f.Sort, "created_at")
				}
			}
			return db
		}
		return db
	}
}
