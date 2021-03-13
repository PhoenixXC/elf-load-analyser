package main

import (
    "fmt"
    "github.com/phoenixxc/elf-load-analyser/pkg/log"
    "os"
    "os/signal"
    "strings"
    "syscall"
)

const (
    ChildFlagEnv  = "_CHILD_46cf632ae5fb1e8cbb556aa49964ac7d"
    ChildArgsFlag = "_ARGS_51cee58d2009745b72acb79005a881e4"
)

// TODO
//  - 子进程挂掉结束程序，心跳否
//  - 进程启动后监控子进程堆栈、动态链接记录

func childProcess(execPath string) {
    // wait until parent load bcc modules ok
    startSignals := make(chan os.Signal, 1)
    signal.Notify(startSignals, syscall.SIGUSR1)
    <-startSignals

    argsEnv, _ := os.LookupEnv(ChildArgsFlag)
    execArgs := []string{execPath}
    execArgs = append(execArgs, strings.Fields(argsEnv)...)
    if err := syscall.Exec(execPath, execArgs, os.Environ()); err != nil {
        log.Errorf("Call binary failed, %v", err)
    }

    // exit, do not touch parent process exec stream
    os.Exit(0)
}

func buildProcess(ctx *cmdArgs) int {
    execArgs := ctx.args

    args := os.Args
    pwd, err := os.Getwd()
    if err != nil {
        log.Errorf("Get pwd error, %v", err)
    }
    childEnvItem := fmt.Sprintf("%s=%s", ChildFlagEnv, ctx.path)
    childArgItem := fmt.Sprintf("%s=%s", ChildArgsFlag, strings.TrimSpace(execArgs))
    childPID, err := syscall.ForkExec(args[0], args, &syscall.ProcAttr{ //nolint:gosec
        Dir: pwd,
        Env: append(os.Environ(), childEnvItem, childArgItem),
        Sys: &syscall.SysProcAttr{
            Setsid: true,
            Credential: &syscall.Credential{
                Uid: uint32(ctx.uid),
                Gid: uint32(ctx.gid),
            },
        },
        Files: []uintptr{cmd.iFd, cmd.oFd, cmd.eFd},
    })
    if err != nil {
        log.Errorf("Create process failed, %v", err)
    }

    log.Infof("Create child process %d(parent: %d) success", childPID, os.Getppid())
    return childPID
}

func wakeChild(childPID int) {
    err := syscall.Kill(childPID, syscall.SIGUSR1)
    if err != nil {
        log.Errorf("Wake child process error, %v", err)
    }
}
