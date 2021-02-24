package log

import (
    "fmt"
    "io"
    "log"
    "os"
)

type Level int8

const (
    debugLevel Level = iota
    infoLevel
    warnLevel
    errorLevel
)

var (
    currentLevel = infoLevel
    defaultFlags = log.Lmicroseconds | log.Ltime

    debugLogger logger = newBaseLogger(debugLevel, os.Stdout, "DEBUG ", defaultFlags)
    infoLogger  logger = newBaseLogger(infoLevel, os.Stdout, "INFO ", defaultFlags)
    warnLogger  logger = newBaseLogger(warnLevel, os.Stdout, "WARN ", defaultFlags).SetHandle(warn)
    errorLogger logger = newBaseLogger(errorLevel, os.Stderr, "ERROR ", defaultFlags).SetHandle(err)
)

type logger interface {
    Level() Level
    Log(a interface{})
    Logf(format string, a ...interface{})
}

func ConfigLevel() Level {
    return currentLevel
}

func SetConfigLevel(l string) {
    if len(l) == 0 {
        return
    }
    switch l {
    case "debug":
        currentLevel = debugLevel
    case "info":
        currentLevel = infoLevel
    case "warn":
        currentLevel = warnLevel
    case "error":
        currentLevel = errorLevel
    default:
        log.Fatalf("Invalid log level %q", l)
    }
}

type baseLogger struct {
    log    *log.Logger
    level  Level
    handle func(s string) string
}

func newBaseLogger(level Level, w io.Writer, prefix string, flag int) *baseLogger {
    return &baseLogger{level: level, log: log.New(w, prefix, flag)}
}

func (b *baseLogger) SetHandle(handle func(s string) string) *baseLogger {
    b.handle = handle
    return b
}

func (b baseLogger) Log(a interface{}) {
    var beforeMsg, afterMsg string
    afterMsg = fmt.Sprint(a)
    if b.handle != nil {
        beforeMsg = b.handle(afterMsg)
    }
    if len(beforeMsg) > 0 {
        b.log.Println(beforeMsg)
    } else {
        b.log.Println(afterMsg)
    }
}

func (b baseLogger) Logf(format string, a ...interface{}) {
    b.Log(fmt.Sprintf(format, a...))
}

func (b baseLogger) Level() Level {
    return b.level
}

func Info(a interface{}) {
    innerLog(infoLogger, a)
}

func Infof(format string, a ...interface{}) {
    innerLogf(infoLogger, format, a...)
}

func Warn(a interface{}) {
    innerLog(warnLogger, a)
}

func Warnf(format string, a ...interface{}) {
    innerLogf(warnLogger, format, a...)
}

func Debug(a interface{}) {
    innerLog(debugLogger, a)
}

func Debugf(format string, a ...interface{}) {
    innerLogf(debugLogger, format, a...)
}

func Error(a interface{}) {
    innerLog(errorLogger, a)
    os.Exit(1)
}

func Errorf(format string, a ...interface{}) {
    innerLogf(errorLogger, format, a...)
    os.Exit(1)
}

func innerLog(log logger, a interface{}) {
    if log.Level() >= ConfigLevel() {
        log.Log(a)
    }
}

func innerLogf(log logger, format string, a ...interface{}) {
    if log.Level() >= ConfigLevel() {
        log.Logf(format, a...)
    }
}
