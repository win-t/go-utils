package structural

import (
	"reflect"
)

// Assign every field from source to target.
//
// Only field that exported, have same name, and same type is got assigned.
//
// target must be pointer to a struct, and source must be struct or pointer to a struct.
func Assign(target interface{}, sources ...interface{}) {
	var targetVal reflect.Value
	if v := reflect.ValueOf(target); v.Kind() == reflect.Ptr {
		targetVal = v.Elem()
	}
	if targetVal.Kind() != reflect.Struct {
		panic("structural: target must be non-nil pointer to struct")
	}
	targetTyp := targetVal.Type()

	for _, s := range sources {
		sourceVal := reflect.ValueOf(s)
		if sourceVal.Kind() == reflect.Ptr {
			sourceVal = sourceVal.Elem()
		}
		if sourceVal.Kind() != reflect.Struct {
			panic("structural: source must be struct or non-nil pointer to struct")
		}

		for i := 0; i < targetTyp.NumField(); i++ {
			f := targetTyp.Field(i)
			if !f.IsExported() {
				continue
			}

			s := sourceVal.FieldByName(f.Name)
			if !s.IsValid() {
				continue
			}

			if s.Type() != f.Type {
				continue
			}

			targetVal.Field(i).Set(s)
		}
	}
}
