package modules

import (
    _ "embed"
    bpf "github.com/iovisor/gobpf/bcc"
    "github.com/phoenixxc/elf-load-analyser/pkg/bcc"
)

//go:embed src/execve.cpp
var execveSource string

type sysExecve struct {
    *BaseMonitorModule
}

func init() {
    //m := NewPerfResolveMonitorModule(&sysExecve{})
    //ModuleInit(m, true)
}

func (e *sysExecve) Monitor() string {
    return "syscall:execve"
}

func (e *sysExecve) Source() string {
    return execveSource
}

func (e *sysExecve) Events() []*bcc.Event {
    fnName := bpf.GetSyscallFnName("execve")
    k := bcc.NewKprobeEvent("syscall__execve", fnName, -1)
    ke := bcc.NewKretprobeEvent("do_ret_sys_execve", fnName, -1)
    return []*bcc.Event{ke, k}
}
