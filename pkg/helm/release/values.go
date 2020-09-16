// Copyright 2020-present Open Networking Foundation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package release

import (
	"encoding/csv"
	"github.com/iancoleman/strcase"
	"reflect"
	"strings"
)

// mergeValues recursively merges values 'b' into values 'a', returning a new merged values map
func mergeValues(a, b map[string]interface{}) map[string]interface{} {
	out := make(map[string]interface{}, len(a))
	for k, v := range a {
		out[k] = v
	}
	for k, v := range b {
		if v, ok := v.(map[string]interface{}); ok {
			if bv, ok := out[k]; ok {
				if bv, ok := bv.(map[string]interface{}); ok {
					out[k] = mergeValues(bv, v)
					continue
				}
			}
		}
		out[k] = v
	}
	return out
}

// setValue sets a value in a map by its path
func setValue(config map[string]interface{}, path string, value interface{}) {
	keys := getValueKeys(path)
	parentKeys, childKey := keys[:len(keys)-1], keys[len(keys)-1]
	parent := getValuesRef(config, parentKeys)
	parent[childKey] = value
}

// getValuesRef gets a reference to the given values path
func getValuesRef(parent map[string]interface{}, path []string) map[string]interface{} {
	if len(path) == 0 {
		return parent
	}
	child, ok := parent[path[0]]
	if !ok {
		child = make(map[string]interface{})
		parent[path[0]] = child
	}
	return getValuesRef(child.(map[string]interface{}), path[1:])
}

// getValueKeys splits the given value path into a list of keys
func getValueKeys(path string) []string {
	r := csv.NewReader(strings.NewReader(path))
	r.Comma = '.'
	names, err := r.Read()
	if err != nil {
		panic(err)
	}
	return names
}

// normalizeValues normalizes the given values map, converting structs into maps
func normalizeValues(values map[string]interface{}) map[string]interface{} {
	return normalizeValue(values).(map[string]interface{})
}

// normalizeValue recursively normalizes the given value, converting structs into maps
func normalizeValue(value interface{}) interface{} {
	kind := reflect.ValueOf(value).Kind()
	if kind == reflect.Struct {
		return normalizeStruct(value.(struct{}))
	} else if kind == reflect.Map {
		return normalizeMap(value.(map[string]interface{}))
	} else if kind == reflect.Slice {
		return normalizeSlice(value.([]interface{}))
	}
	return value
}

// normalizeStruct recursively normalizes the given struct, converting it to a map
func normalizeStruct(value struct{}) interface{} {
	elem := reflect.ValueOf(value).Elem()
	elemType := elem.Type()
	normalized := make(map[string]interface{})
	for i := 0; i < elem.NumField(); i++ {
		key := getFieldKey(elemType.Field(i))
		value := normalizeValue(elem.Field(i).Interface())
		normalized[key] = value
	}
	return normalized
}

// normalizeMap recursively normalizes the given map, converting structs into maps
func normalizeMap(values map[string]interface{}) interface{} {
	normalized := make(map[string]interface{})
	for key, value := range values {
		normalized[key] = normalizeValue(value)
	}
	return normalized
}

// normalizeValue recursively normalizes the given slice, converting structs into maps
func normalizeSlice(values []interface{}) interface{} {
	normalized := make([]interface{}, len(values))
	for i, value := range values {
		normalized[i] = normalizeValue(value)
	}
	return normalized
}

// getFieldKey returns the map key name for the given struct field
func getFieldKey(field reflect.StructField) string {
	tag := field.Tag.Get("yaml")
	if tag != "" {
		return strings.Split(tag, ",")[0]
	}
	return strcase.ToLowerCamel(field.Name)
}