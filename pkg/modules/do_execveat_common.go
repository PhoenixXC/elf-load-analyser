package modules

import (
    _ "embed"
    "fmt"
    "github.com/phoenixxc/elf-load-analyser/pkg/bcc"
    "github.com/phoenixxc/elf-load-analyser/pkg/data"
)

//go:embed src/do_execveat_common.cpp
var doExecveatCommonSource string

const (
    monitorName = "hook_execveat"
)

type execveatComEvent struct {
    Fd       int32
    Flags    int32
    Filename [256]byte
}

func (e *execveatComEvent) Render() *data.AnalyseData {
    s := bytes2Str(e.Filename[:])
    msg := fmt.Sprintf("Do `%s` function, with fd = %d, flags = %d, filename = %s\n",
        "do_execveat_common", e.Fd, e.Flags, s)
    return data.NewAnalyseData(monitorName, data.NewData(data.MarkdownType, msg))
}

type doExecveatCommon struct {
    *BaseMonitorModule
}

func init() {
    m := NewPerfResolveMonitorModule(&doExecveatCommon{})
    m.RegisterOnceTable("events", func(data []byte) (*data.AnalyseData, error) {
        return m.Render(data, &execveatComEvent{})
    })
    ModuleDefaultInit(m)
}

func (c *doExecveatCommon) Monitor() string {
    return "hook_execveat"
}

func (c *doExecveatCommon) Source() string {
    return doExecveatCommonSource
}

func (c *doExecveatCommon) Events() []*bcc.Event {
    ke := bcc.NewKprobeEvent("kprobe__do_execveat_common", "do_execveat_common", -1)
    return []*bcc.Event{ke}
}