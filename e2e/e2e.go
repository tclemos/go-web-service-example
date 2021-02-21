package e2e

import (
	"context"
	"fmt"
)

type ctxKey string

const (
	valuesKey = ctxKey("values")
)

// AddValue allows a value to be stored during the TestMain to be used within tests
func AddValue(key string, value interface{}) {
	values := Ctx.Value(valuesKey).(map[string]interface{})
	values[key] = value
	Ctx = context.WithValue(Ctx, valuesKey, values)
	fmt.Printf("[e2e]added value: %s %v\n", key, value)
}

// GetValue allows a value to be retrieved by its key
func GetValue(key string) interface{} {
	values := Ctx.Value(valuesKey).(map[string]interface{})
	return values[key]
}

// GetValues gets all stored values
func GetValues() map[string]interface{} {
	values := Ctx.Value(valuesKey).(map[string]interface{})
	return values
}
