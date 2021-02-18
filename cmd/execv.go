package main

import (
    "fmt"
    "log"
    "os"
    "os/signal"
    "strings"
    "syscall"
)

const (
    ChildFlagEnv  = "_CHILD_46cf632ae5fb1e8cbb556aa49964ac7d"
    ChildArgsFlag = "_ARGS_51cee58d2009745b72acb79005a881e4"
)

type execCtx struct {
    execArgs string
    uid int
    gid int
}

func childProcess(execPath string) {
    // wait until parent load bcc modules ok
    startSignals := make(chan os.Signal, 1)
    signal.Notify(startSignals, syscall.SIGUSR1)
    <-startSignals

    argsEnv, _ := os.LookupEnv(ChildArgsFlag)
    execArgs := []string{execPath}
    execArgs = append(execArgs, strings.Fields(argsEnv)...)
    log.Printf("Boot binary %q with \"%v\" to analyse load data...\n", execPath, execArgs)
    if err := syscall.Exec(execPath, execArgs, os.Environ()); err != nil {
        log.Fatalf("Call binary failed, %v", err)
    }

    // exit, do not touch parent process exec stream
    os.Exit(0)
}

func buildProcess(ctx execCtx) int {
    execArgs := ctx.execArgs
    uid, gid := ctx.uid, ctx.gid

    args := os.Args
    pwd, err := os.Getwd()
    if err != nil {
        log.Fatalf("Get pwd error, %v", err)
    }
    childEnvItem := fmt.Sprintf("%s=%s", ChildFlagEnv, execPath)
    childArgItem := fmt.Sprintf("%s=%s", ChildArgsFlag, strings.TrimSpace(execArgs))
    childPID, err := syscall.ForkExec(args[0], args, &syscall.ProcAttr{
        Dir: pwd,
        Env: append(os.Environ(), childEnvItem, childArgItem),
        Sys: &syscall.SysProcAttr{
            Setsid: true,
            Credential: &syscall.Credential{
                Uid: uint32(uid),
                Gid: uint32(gid),
            },
        },
        // TODO 重定向输入输出
        Files: []uintptr{0, 1, 2},
    })
    if err != nil {
        log.Fatalf("Create process failed, %v\n", err)
    }

    log.Printf("Create child process %d, parent: %d success\n", childPID, os.Getppid())
    return childPID
}

func wakeChild(childPID int) {
    err := syscall.Kill(childPID, syscall.SIGUSR1)
    if err != nil {
        log.Fatalf("Wake child process error, %v\n", err)
    }
}