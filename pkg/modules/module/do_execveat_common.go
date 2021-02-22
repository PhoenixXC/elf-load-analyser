package module

import (
    _ "embed"
    "fmt"
    "github.com/phoenixxc/elf-load-analyser/pkg/bcc"
    "github.com/phoenixxc/elf-load-analyser/pkg/data"
    "github.com/phoenixxc/elf-load-analyser/pkg/modules"
    "github.com/phoenixxc/elf-load-analyser/pkg/modules/enhance"
)

//go:embed src/do_execveat_common.cpp
var doExecveatCommonSource string

const (
    monitorName = "hook_execveat"
)

type execveatComEvent struct {
    enhance.TimeEventResult
    Fd       int32
    Flags    int32
    Filename [256]byte
}

func (e execveatComEvent) Render() *data.AnalyseData {
    s := data.TrimBytes2Str(e.Filename[:])
    msg := fmt.Sprintf("Do `%s` function, with fd = %d, flags = %d, filename = %s\n",
        "do_execveat_common", e.Fd, e.Flags, s)
    return data.NewAnalyseData(monitorName, data.NewData(data.MarkdownType, msg))
}

type doExecveatCommon struct {
    modules.MonitorModule
}

func init() {
    m := modules.NewPerfResolveMonitorModule(&doExecveatCommon{})
    m.RegisterOnceTable("events", func(data []byte) (*data.AnalyseData, error) {
        return modules.Render(data, &execveatComEvent{}, true)
    })
    modules.ModuleDefaultInit(m)
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