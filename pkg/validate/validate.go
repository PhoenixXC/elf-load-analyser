package validate

import (
    "fmt"
    "log"
    "reflect"
)

func EqualWithTip(expected, actual interface{}, errorMsg string) {
    Equal(expected, actual, func(e, a interface{}) {
        log.Fatal(errorMsg)
    })
}

func Equal(expected, actual interface{}, handler func(e, a interface{})) {
    Validate(expected, actual, func(expected, actual interface{}) bool {
        return expected == actual
    }, handler)
}

func WithTip(expected, actual interface{}, predicate func(expected, actual interface{}) bool, errorMsg string) {
    Validate(expected, actual, predicate, func(e, a interface{}) {
        log.Fatal(errorMsg)
    })
}

func Validate(expected, actual interface{},
    predicate func(expected, actual interface{}) bool,
    handler func(e, a interface{})) {
    expectedVal, actualVal := getValue(expected), getValue(actual)
    if !predicate(expectedVal, actualVal) {
        handler(expectedVal, actualVal)
    }
}

func isNotFunc(val interface{}) bool {
    if val == nil {
        return false
    }
    return reflect.TypeOf(val).Kind() != reflect.Func
}

func getValue(val interface{}) interface{} {
    if isNotFunc(val) {
        return val
    }
    switch val.(type) {
    case func() interface{}:
        return val.(func() interface{})()
    default:
        panic(fmt.Sprintf("argument %T should be normal value or func() interface{} type", val))
    }
}
