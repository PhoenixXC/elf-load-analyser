package modules

import (
    "bytes"
    _ "embed"
    "encoding/binary"
    bpf "github.com/iovisor/gobpf/bcc"
    "github.com/phoenixxc/elf-load-analyser/pkg/bcc"
    "github.com/phoenixxc/elf-load-analyser/pkg/data"
)

//go:embed src/execve.cpp
var execveSource string

type sysExecve struct {
    *BaseMonitorModule
}

func init() {
    m := NewPerfResolveMonitorModule(&sysExecve{})
    m.RegisterOnceTable("call_events", func(_ []byte) (*data.AnalyseData, error) {
        return data.NewAnalyseData("syscall:execve", data.NewData(data.MarkdownType, "Start call syscall execve...")), nil
    })
    m.RegisterOnceTable("ret_events", func(d []byte) (*data.AnalyseData, error) {
        r := new(int8)
        err := binary.Read(bytes.NewBuffer(d), bpf.GetHostByteOrder(), r)
        if err != nil {
            return nil, err
        }
        desc := "success"
        if *r != 0 {
            desc = "failed"
        }
        return data.NewAnalyseData("syscall:execve", data.NewData(data.MarkdownType, "Start call syscall execve "+desc)), nil
    })
    ModuleInit(m, true)
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
