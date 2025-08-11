package helpers

import (
	"reflect"

	"go.mongodb.org/mongo-driver/bson"
)

func BuildUpdateBson(input interface{}) bson.M {
	update := bson.M{}
	val := reflect.ValueOf(input)
	typ := reflect.TypeOf(input)

	for i := 0; i < val.NumField(); i++ {
		fieldVal := val.Field(i)
		fieldType := typ.Field(i)

		// Ambil tag `bson` atau nama field
		bsonTag := fieldType.Tag.Get("bson")
		if bsonTag == "" || bsonTag == "-" {
			continue
		}

		// Skip zero values
		if !fieldVal.IsZero() {
			update[bsonTag] = fieldVal.Interface()
		}
	}
	return update
}
