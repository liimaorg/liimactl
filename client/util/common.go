package util

import (
	"fmt"
	"reflect"
	"strconv"
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

//TagExists checks if a given "json" tag exists in the given interface
func TagExists(f interface{}, actTag string) bool {
	val := reflect.ValueOf(f).Elem()

	for i := 0; i < val.NumField(); i++ {
		// valueField := val.Field(i)
		typeField := val.Type().Field(i)
		tag := typeField.Tag
		if tag.Get(TagName) == actTag {
			return true
		}

		// fmt.Printf("Field Name: %s,\t Field Value: %v,\t Tag Value: %s\n", typeField.Name, valueField.Interface(), tag.Get("tag_name"))
	}

	return false
}

//SetValueIfTagExists set a value if a given "json" tag exists in the given interface
func SetValueIfTagExists(f interface{}, actTag string, setVal string) {
	v := reflect.ValueOf(f).Elem()
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		vField := v.Field(i)
		tField := t.Field(i)
		tag := tField.Tag.Get(TagName)

		if tag == actTag {

			switch vField.Kind() {
			case reflect.Int:
				iVal, _ := strconv.ParseInt(setVal, 10, 64)
				vField.SetInt(iVal)
			case reflect.Float64:
				fVal, _ := strconv.ParseFloat(setVal, 64)
				vField.SetFloat(fVal)
			case reflect.String:
				vField.SetString(setVal)
			case reflect.Bool:
				bVal, _ := strconv.ParseBool(setVal)
				vField.SetBool(bVal)
			default:
				fmt.Printf("Unsupported kind: '%s' Value: %s", vField.Kind(), vField)
			}

		}
	}
}
