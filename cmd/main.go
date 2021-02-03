package main

import (
    "flag"
    "github.com/phoenixxc/elf-load-analyser/pkg/bcc"
    "github.com/phoenixxc/elf-load-analyser/pkg/factory"
    _ "github.com/phoenixxc/elf-load-analyser/pkg/modules"
    "github.com/phoenixxc/elf-load-analyser/pkg/system"

    "log"
    "os"
    "path/filepath"
    "time"

    "golang.org/x/sys/unix"
)

var (
    execPath   string // exec file path
    execArgStr string // exec args
)

func init() {
    flag.StringVar(&execPath, "e", "", "the analyse program path")
    flag.StringVar(&execArgStr, "p", "", "the analyse program parameter, split by space")

    flag.Parse()
}

func main() {
    // child
    transExecPath, isChild := os.LookupEnv(ChildFlagEnv)
    if isChild {
        childProcess(transExecPath)
    }

    // handle flag
    checkFlag()

    // system, kernel version, kernel config and depend software check
    system.CheckEnv()

    // fork, get pid, block until receive signal
    childPID := buildProcess(execArgStr)

    // bcc handler update, hook pid, load modules, begin hook
    factory.LoadMonitors(bcc.Context{
        Pid: childPID,
    })

    // wake up chile to exec binary
    wakeChild(childPID)

    // cache load detail data, render use html(use graphviz build images, if no graphviz, show code use <code> tag)

    // when load ok, save html to disk

    // optional: start web server show message

    // optional: start to monitor dynamic link at real time, use websocket
    time.Sleep(1 * time.Hour)
}

func checkFlag() {
    // -e
    if len(execPath) == 0 {
        flag.Usage()
        os.Exit(1)
    }
    absPath, err := filepath.Abs(execPath)
    if err != nil {
        log.Fatalf("Get absolute path error, %v", err)
    }
    execPath = absPath

    if err := unix.Access(execPath, unix.X_OK); err != nil {
        log.Fatalf("Check %q error, %v", execPath, err)
    }
}
