package modules

import (
    "fmt"
    "github.com/phoenixxc/elf-load-analyser/pkg/bcc"
    "github.com/phoenixxc/elf-load-analyser/pkg/data"
)

//goland:noinspection ALL
const doExecveatCommonSource = `
#include <linux/fs.h>
#include <linux/sched.h>

struct exec_event {
    int fd;
    int flags;
    char filename[256];
};
BPF_PERF_OUTPUT(events);

int kprobe__do_execveat_common(struct pt_regs* ctx, int fd, struct filename* filename, int flags) {
    // Returns the process ID in the lower 32 bits (kernel's view of the PID, 
    // which in user space is usually presented as the thread ID), 
    // and the thread group ID in the upper 32 bits (what user space often thinks of as the PID).
    // By directly setting this to a u32, we discard the upper 32 bits.
    if ((bpf_get_current_pid_tgid() >> 32) != _PID_) {
        return 0;
    }
    struct exec_event event = {};
    bpf_probe_read_kernel(&(event.fd), sizeof(int), (void*)&fd);
    bpf_probe_read_kernel(&(event.flags), sizeof(int), (void*)&flags);
    bpf_probe_read_kernel_str(&(event.filename), sizeof(event.filename), (void*)filename->name);
    events.perf_submit((void*)ctx, (void*)&event, sizeof(event));
    return 0;
}
`

const (
    monitorName = "hook_execveat"
)

type execEvent struct {
    Fd       int32
    Flags    int32
    Filename [256]byte
}

func (e *execEvent) Render() *data.AnalyseData {
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
    m.RegisterTable("events", false, func(data []byte) (*data.AnalyseData, error) {
        return m.Render(data, &execEvent{})
    })
    ModuleInit(m, false)
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