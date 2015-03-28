package models

import (
	"database/sql/driver"
	"fmt"
	"github.com/revel/revel"
	"log"
	"reflect"
	"regexp"
	"strings"
	"time"
)

type Logger struct {
	*log.Logger
}

// Format log
var logger Logger
var sqlRegexp = regexp.MustCompile(`(\$\d+)|\?`)

func (logger Logger) Print(values ...interface{}) {

	if len(values) > 1 {
		level := values[0]
		sourceNodes := strings.Split(fmt.Sprintf("%v", values[1]), "/")
		source := sourceNodes[len(sourceNodes)-1:]
		formattedSource := fmt.Sprintf("\033[34m%v\033[0m", source)
		messages := []interface{}{formattedSource}

		if level == "sql" {
			messages = append(messages, fmt.Sprintf("\033[36;1m(%.2fms)\033[0m", float64(values[2].(time.Duration).Nanoseconds()/1e4)/100.0))
			// sql
			var formatedValues []interface{}
			for _, value := range values[4].([]interface{}) {
				indirectValue := reflect.Indirect(reflect.ValueOf(value))
				if indirectValue.IsValid() {
					value = indirectValue.Interface()
					if t, ok := value.(time.Time); ok {
						formatedValues = append(formatedValues, fmt.Sprintf("'%v'", t.Format(time.RFC3339)))
					} else if b, ok := value.([]byte); ok {
						formatedValues = append(formatedValues, fmt.Sprintf("'%v'", string(b)))
					} else if r, ok := value.(driver.Valuer); ok {
						if value, err := r.Value(); err == nil && value != nil {
							formatedValues = append(formatedValues, fmt.Sprintf("'%v'", value))
						} else {
							formatedValues = append(formatedValues, "NULL")
						}
					} else {
						formatedValues = append(formatedValues, fmt.Sprintf("'%v'", value))
					}
				} else {
					formatedValues = append(formatedValues, fmt.Sprintf("'%v'", value))
				}
			}

			messages = append(messages, fmt.Sprintf(sqlRegexp.ReplaceAllString(values[3].(string), "%v"), formatedValues...))

		}
		revel.INFO.Println(messages...)
	} else {
		revel.INFO.Println(values)
	}

}
