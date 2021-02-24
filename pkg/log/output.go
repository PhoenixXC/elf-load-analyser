package log

// Emphasize warp string with bolder in terminal
func Emphasize(message string) string {
	return "\u001B[1m" + message + "\u001B[0m"
}

func warn(message string) string {
	return "\033[1;33m" + message + "\033[0m"
}

func success(message string) string {
	return "\033[1;32m" + message + "\033[0m"
}

func err(message string) string {
	return "\033[1;31m" + message + "\033[0m"
}