package path

import (
	"encoding/json"
	"reflect"

	"github.com/gopatchy/jsrest"
)

func Merge(to, from any) {
	MergeValue(reflect.ValueOf(to), reflect.ValueOf(from))
}

func MergeValue(to, from reflect.Value) {
	to = reflect.Indirect(to)
	from = reflect.Indirect(from)

	for i := 0; i < to.NumField(); i++ {
		toField := to.Field(i)
		fromField := from.Field(i)

		if fromField.IsZero() {
			continue
		}

		if reflect.Indirect(fromField).Kind() == reflect.Struct {
			MergeValue(toField, fromField)
			continue
		}

		toField.Set(fromField)
	}
}

func MergeMap(to any, from map[string]any) error {
	m, err := ToMap(to)
	if err != nil {
		return jsrest.Errorf(jsrest.ErrInternalServerError, "converting to map failed (%w)", err)
	}

	MergeMaps(m, from)

	return FromMap(to, m)
}

func MergeMaps(to map[string]any, from map[string]any) {
	for k, v := range from {
		if vMap, isMap := v.(map[string]any); isMap {
			if _, ok := to[k].(map[string]any); !ok {
				// Either key doesn't exist or it's a different type
				// If different type, error will happen during json decode
				to[k] = map[string]any{}
			}

			MergeMaps(to[k].(map[string]any), vMap)
		} else {
			to[k] = v
		}
	}
}

func ToMap(from any) (map[string]any, error) {
	js, err := json.Marshal(from)
	if err != nil {
		return nil, jsrest.Errorf(jsrest.ErrInternalServerError, "json marshal failed (%w)", err)
	}

	ret := map[string]any{}

	err = json.Unmarshal(js, &ret)
	if err != nil {
		return nil, jsrest.Errorf(jsrest.ErrInternalServerError, "json unmarshal failed (%w)", err)
	}

	return ret, nil
}

func FromMap(to any, from map[string]any) error {
	js, err := json.Marshal(from)
	if err != nil {
		return jsrest.Errorf(jsrest.ErrInternalServerError, "json marshal failed (%w)", err)
	}

	err = json.Unmarshal(js, to)
	if err != nil {
		return jsrest.Errorf(jsrest.ErrInternalServerError, "json unmarshal failed (%w)", err)
	}

	return nil
}
