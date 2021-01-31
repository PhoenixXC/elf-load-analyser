package main

import (
    "flag"
    "log"
    "os"
    "path/filepath"
    "syscall"
    
    "github.com/phoenixxc/elf-load-analyser/pkg/system/env"
)

var (
    execPath string // exec file path
)

func init() {
    flag.StringVar(&execPath, "e", "", "the analyse program")
    
    flag.Parse()
}

func main() {
    // check flag
    checkFlag()
    
    // system, kernel version, kernel config and depend software check
    env.CheckEnv()
    
    // fork, get pid, block until receive signal
    // inject env for child progress
    
    // bcc handler update, hook pid, load modules, begin hook
    
    // exec binary
    
    // cache load detail data, render use html(use graphviz build images, if no graphviz, show code use <code> tag)
    
    // save html to disk
    
    // optional: start web server show message
    
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
    if err := syscall.Access(execPath, syscall.O_EXCL & syscall.F_OK); err != nil {
        log.Fatalf("Exec %q not exist or no executable perm, %v", execPath, err)
    }
}