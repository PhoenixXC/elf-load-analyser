package system

const (
	successFlag = "[√]"
	failFlag    = "[×]"
)

// Emphasize warp string with bolder in terminal
func Emphasize(message string) string {
	return "\u001B[1m" + message + "\u001B[0m"
}

// Warn warn message
func Warn(message string) string {
	return warn("WARN " + message)
}

func Success(message string) string {
	return success("SUCCESS " + message)
}

func Error(message string) string {
	return error("ERROR " + message)
}

func Check(message string, ok bool) string {
	var prefix string
	if ok {
		prefix = success(successFlag)
	} else {
		prefix = error(failFlag)
	}
	return prefix + " " + message
}

func warn(message string) string {
	return "\033[1;33m" + message + "\033[0m"
}

func success(message string) string {
	return "\033[1;32m" + message + "\033[0m"
}

func error(message string) string {
	return "\033[1;31m" + message + "\033[0m"
}