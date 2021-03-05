package helper

// IfElse simple if else
func IfElse(condition bool, established interface{}, elVal interface{}) interface{} {
    if condition {
        return established
    }
    return elVal
}
