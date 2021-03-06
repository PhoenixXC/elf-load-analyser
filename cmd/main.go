package main

import (
    "flag"
    "github.com/phoenixxc/elf-load-analyser/pkg/bcc"
    "github.com/phoenixxc/elf-load-analyser/pkg/env"
    "github.com/phoenixxc/elf-load-analyser/pkg/factory"
    "github.com/phoenixxc/elf-load-analyser/pkg/log"
    _ "github.com/phoenixxc/elf-load-analyser/pkg/modules/module"
    "github.com/phoenixxc/elf-load-analyser/pkg/render"
    "github.com/phoenixxc/elf-load-analyser/pkg/web"
    "golang.org/x/sys/unix"
    "os"
    "os/signal"
    "os/user"
    "path/filepath"
    "strconv"
    "strings"
    "syscall"
)

type cmdArgs struct {
    path          string // exec file path
    args          string // exec args
    user          string // exec user
    level         string
    uid, gid      int
    in, out, eOut string // child process input and output
    iFd, oFd, eFd uintptr
    port          uint
}

var (
    cmd          = &cmdArgs{}
    closeHandler []func()
)

func init() {
    flag.StringVar(&cmd.user, "user", "", "run user")
    flag.StringVar(&cmd.path, "exec", "", "program path")
    flag.StringVar(&cmd.in, "in", "", "(optional) target program input")
    flag.StringVar(&cmd.out, "out", "", "(optional) target program output")
    flag.StringVar(&cmd.eOut, "err", "", "(optional) target program error output")
    flag.UintVar(&cmd.port, "port", 0, "(optional) web server port, default use random port")
    flag.StringVar(&cmd.level, "log", "", "(optional) log level(info debug warn error), default: info")
    flag.StringVar(&cmd.args, "arg", "", "(optional) transform program parameter, split by space, default: ''")

    flag.Parse()
}

func main() {
    preProcessing(cmd)
    render.PreAnalyse(render.Content{Filepath: cmd.path})

    childPID := buildProcess(cmd)
    pool, _ := factory.LoadMonitors(bcc.Context{Pid: childPID})
    wakeChild(childPID)

    web.VisualAnalyseData(pool, cmd.port)

    log.Info(log.Emphasize("Press [CTRL+C] to exit"))
    exit := make(chan os.Signal, 1)
    signal.Notify(exit, os.Interrupt, syscall.SIGTERM)
    <-exit

    defer closeHandle()
}

func closeHandle() {
    if len(closeHandler) == 0 {
        return
    }
    for _, f := range closeHandler {
        f()
    }
}

func preProcessing(c *cmdArgs) {
    if e := log.SetConfigLevel(c.level); e != nil {
        flag.Usage()
        log.Error(e)
    }
    // child
    transExecPath, isChild := os.LookupEnv(ChildFlagEnv)
    if isChild {
        childProcess(transExecPath)
        return
    }
    // args
    treatingArgs(c)
    // banner
    env.EchoBanner()
    // env
    env.CheckEnv()
}

func treatingArgs(c *cmdArgs) {
    // path check
    if len(c.path) == 0 {
        flag.Usage()
        os.Exit(1)
    }
    absPath, err := filepath.Abs(c.path)
    if err != nil {
        log.Errorf("Get absolute path error, %v", err)
    }
    c.path = absPath

    if err := unix.Access(c.path, unix.X_OK); err != nil {
        log.Errorf("Check %q error, %v", c.path, err)
    }

    // user check
    c.uid, c.gid = getUIDGid(c.user)
    // input output
    c.iFd, c.oFd, c.eFd = getIOFd(c.in, c.out, c.eOut)
}

func getUIDGid(username string) (uid, gid int) {
    s := strings.TrimSpace(username)
    if len(s) != 0 {
        u, err := user.Lookup(s)
        if err == nil {
            uid, _ := strconv.Atoi(u.Uid)
            gid, _ := strconv.Atoi(u.Gid)
            return uid, gid
        }
        log.Error("Invalid user")
    }
    flag.Usage()
    os.Exit(1)
    return
}

func getIOFd(in, out, eOut string) (iFd, outFd, eFd uintptr) {
    iFd, outFd, eFd = 0, 1, 2
    if len(in) > 0 {
        iFd = getFd(in, true)
    }
    if len(out) > 0 {
        outFd = getFd(out, false)
    }
    if len(eOut) > 0 {
        eFd = getFd(eOut, false)
    }
    return
}

func getFd(file string, read bool) uintptr {
    var f *os.File
    var e error
    if read {
        f, e = os.Open(file)
    } else {
        f, e = os.OpenFile(file, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
    }
    if e != nil {
        closeHandle()
        log.Errorf("Get file fd: %v", e)
    }
    closeHandler = append(closeHandler, func() { f.Close() })
    return f.Fd()
}
