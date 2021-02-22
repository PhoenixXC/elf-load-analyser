package main

import (
    "flag"
    "fmt"
    "github.com/phoenixxc/elf-load-analyser/pkg/bcc"
    "github.com/phoenixxc/elf-load-analyser/pkg/data"
    "github.com/phoenixxc/elf-load-analyser/pkg/factory"
    _ "github.com/phoenixxc/elf-load-analyser/pkg/modules/module"
    "github.com/phoenixxc/elf-load-analyser/pkg/system"
    "os/user"
    "strconv"
    "strings"

    "golang.org/x/sys/unix"
    "log"
    "os"
    "path/filepath"
)

var (
    execPath   string // exec file path
    execArgStr string // exec args
    execUser   string // exec user
)

func init() {
    flag.StringVar(&execPath, "e", "", "the analyse program path")
    flag.StringVar(&execArgStr, "p", "", "the analyse program parameter, split by space")
    flag.StringVar(&execUser, "u", "", "the analyse program run user")

    flag.Parse()
}

func main() {
    // child
    transExecPath, isChild := os.LookupEnv(ChildFlagEnv)
    if isChild {
        childProcess(transExecPath)
    }

    checkFlag()
    u, g := getUidGid(execUser)
    system.CheckEnv()

    // fork, get pid, block until receive signal
    childPID := buildProcess(execCtx{
        execArgs: execArgStr,
        uid:      u,
        gid:      g,
    })

    // bcc handler update, hook pid, load modules, begin hook
    ok := make(chan struct{})
    pool := factory.LoadMonitors(bcc.Context{Pid: childPID}, ok)

    // wake up chile to exec binary
    wakeChild(childPID)

    // wait until data collection ok
    <-ok
    d := pool.Data()
    for _, analyseData := range d {
        // Just for debug
        fmt.Println(analyseData.Timestamp(), analyseData.Name(),
            func() interface{} {
                if analyseData.Status() == data.Success {
                    return analyseData.Data().Data
                }
                return analyseData.Desc()
            }(),
        )
    }

    // cache load detail data, render use html(use graphviz build images, if no graphviz, show code use <code> tag)
    // save html to disk
    // render data result

    // optional: start web server show message

    // optional: start to monitor dynamic link at real time, use websocket
    // if start web server, wait server exit, if not, save html and exit
    //time.Sleep(1 * time.Hour)
}

func checkFlag() {
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

func getUidGid(username string) (uid, gid int) {
    s := strings.TrimSpace(username)
    if len(s) != 0 {
        u, err := user.Lookup(s)
        if err == nil {
            uid, _ := strconv.Atoi(u.Uid)
            gid, _ := strconv.Atoi(u.Gid)
            return uid, gid
        }
        fmt.Println("Invalid user")
    }
    flag.Usage()
    os.Exit(1)
    return
}
