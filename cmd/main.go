package main

import (
    "flag"
    "github.com/phoenixxc/elf-load-analyser/pkg/bcc"
    "github.com/phoenixxc/elf-load-analyser/pkg/factory"
    "github.com/phoenixxc/elf-load-analyser/pkg/log"
    _ "github.com/phoenixxc/elf-load-analyser/pkg/modules/module"
    "github.com/phoenixxc/elf-load-analyser/pkg/render"
    "golang.org/x/sys/unix"
    "os"
    "path/filepath"
)

var (
    execPath   string // exec file path
    execArgStr string // exec args
    execUser   string // exec user
    logLevel   string
)

func init() {
    flag.StringVar(&execUser, "u", "", "run user")
    flag.StringVar(&execPath, "e", "", "program path")
    flag.StringVar(&logLevel, "l", "", "log level (info debug warn error)")
    flag.StringVar(&execArgStr, "p", "", "transform program parameter, split by space")

    flag.Parse()
}

func main() {
    log.SetConfigLevel(logLevel)

    // child
    transExecPath, isChild := os.LookupEnv(ChildFlagEnv)
    if isChild {
        childProcess(transExecPath)
    }

    checkArgs()
    render.PreAnalyse(render.Content{Filepath: execPath})

    // fork, get pid, block until receive signal
    childPID := buildProcess(execCtx{args: execArgStr, user: execUser})
    // bcc handler update, hook pid, load modules, begin hook
    pool, _ := factory.LoadMonitors(bcc.Context{Pid: childPID})
    // wake up chile to exec binary
    wakeChild(childPID)

    render.VisualAnalyseData(pool, true)
}

func checkArgs() {
    if len(execPath) == 0 {
        flag.Usage()
        os.Exit(1)
    }
    absPath, err := filepath.Abs(execPath)
    if err != nil {
        log.Errorf("Get absolute path error, %v", err)
    }
    execPath = absPath

    if err := unix.Access(execPath, unix.X_OK); err != nil {
        log.Errorf("Check %q error, %v", execPath, err)
    }
}
