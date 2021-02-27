package helper

// emm莫得lambda表达式和泛型..

// IfElse simple if else
func IfElse(condition bool, established interface{}, elVal interface{}) interface{} {
    if condition {
        return established
    }
    return elVal
}
