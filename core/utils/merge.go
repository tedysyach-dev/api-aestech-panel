package utils

import (
	"reflect"
)

func MergeRequest(dst interface{}, src interface{}) {
	dstVal := reflect.ValueOf(dst).Elem()
	srcVal := reflect.ValueOf(src).Elem()
	srcType := reflect.TypeOf(src).Elem()

	for i := 0; i < srcVal.NumField(); i++ {
		srcField := srcVal.Field(i)
		srcFieldType := srcType.Field(i)

		dstField := dstVal.FieldByName(srcFieldType.Name)
		if !dstField.IsValid() || !dstField.CanSet() {
			continue
		}

		// Handle pointer source fields
		if srcField.Kind() == reflect.Ptr {
			if !srcField.IsNil() {
				// Jika dst field juga pointer
				if dstField.Kind() == reflect.Ptr {
					dstField.Set(srcField)
				} else {
					dstField.Set(srcField.Elem())
				}
			}
			continue
		}

		// Handle non-pointer fields langsung
		dstField.Set(srcField)
	}
}
