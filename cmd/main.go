package main

import (
    "flag"
    "github.com/phoenixxc/elf-load-analyser/pkg/bcc"
    "github.com/phoenixxc/elf-load-analyser/pkg/env"
    "github.com/phoenixxc/elf-load-analyser/pkg/factory"
    "github.com/phoenixxc/elf-load-analyser/pkg/log"
    _ "github.com/phoenixxc/elf-load-analyser/pkg/modules/module"
    "github.com/phoenixxc/elf-load-analyser/pkg/render"
    "golang.org/x/sys/unix"
    "os"
    "os/user"
    "path/filepath"
    "strconv"
    "strings"
)

type cmdArgs struct {
    path     string // exec file path
    args     string // exec args
    user     string // exec user
    level    string
    uid, gid int
}

var cmd = &cmdArgs{}

func init() {
    flag.StringVar(&cmd.user, "u", "", "run user")
    flag.StringVar(&cmd.path, "e", "", "program path")
    flag.StringVar(&cmd.level, "l", "", "log level (info debug warn error)")
    flag.StringVar(&cmd.args, "p", "", "transform program parameter, split by space")

    flag.Parse()
}

func main() {
    preProcessing(cmd)
    render.PreAnalyse(render.Content{Filepath: cmd.path})

    childPID := buildProcess(cmd)
    pool, _ := factory.LoadMonitors(bcc.Context{Pid: childPID})
    wakeChild(childPID)

    render.VisualAnalyseData(pool, true)
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
    c.uid, c.gid = getUidGid(c.user)
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
        log.Error("Invalid user")
    }
    flag.Usage()
    os.Exit(1)
    return
}
