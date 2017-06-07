package util

import (
	"fmt"
	"reflect"
	"strings"
)

//Json TagName
const (
	TagName = "json" //json tag
)

//BuildCommandURL creates the rest url depending the given command options (interface)
func BuildCommandURL(f interface{}) string {

	val := reflect.ValueOf(f).Elem()
	var url string

	for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i)
		typeField := val.Type().Field(i)

		switch valueField.Kind() {
		case reflect.Slice:
			if valueField.Len() == 0 {
				continue
			}
			url += fmt.Sprintf("%s=", typeField.Tag.Get(TagName))
			sliceVal := make([]interface{}, valueField.Len())
			for j := range sliceVal {
				sliceVal[j] = valueField.Index(j).Interface()
				url += fmt.Sprintf("%v,", sliceVal[j])
			}
			//remove last comma
			url = strings.TrimRight(url, ",")

		default:
			url += fmt.Sprintf("%s=%v", typeField.Tag.Get(TagName), valueField.Interface())

		}

		url += "&"
	}

	//remove last &
	url = strings.TrimRight(url, "&")

	return url
}
