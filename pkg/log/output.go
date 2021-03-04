package log

import "bytes"

// Emphasize warp string with bolder in terminal
func Emphasize(message string) string {
    return wrap("\u001B[1m", message, "\u001B[0m")
}

func warn(message string) string {
    return wrap("\033[1;33m", message, "\033[0m")
}

func minor(message string) string {
    return wrap("\033[0;37m", message, "\033[0m")
}

func success(message string) string {
    return wrap("\033[1;32m", message, "\033[0m")
}

func err(message string) string {
    return wrap("\033[1;31m", message, "\033[0m")
}

func wrap(before string, message string, after string) string {
    var buff bytes.Buffer
    buff.Grow(len(message) + len(before) + len(after))
    buff.WriteString(before)
    buff.WriteString(message)
    buff.WriteString(after)
    return buff.String()
}
