package modules

import (
    "fmt"
    bpf "github.com/iovisor/gobpf/bcc"
    "github.com/phoenixxc/elf-load-analyser/pkg/bcc"
    "github.com/phoenixxc/elf-load-analyser/pkg/data"
    "github.com/phoenixxc/elf-load-analyser/pkg/system"
    "log"
)

// about stack 512byte limit, see: https://stackoverflow.com/questions/53627094/ebpf-track-values-longer-than-stack-size
//goland:noinspection SpellCheckingInspection
const allowBprmSource = `
#include <linux/fs.h>
#include <linux/sched.h>
#include <linux/binfmts.h>

struct alloc_bprm_event {
    char filename[256];
    char fdpath[256];
    char interp[256];
};
BPF_PERF_OUTPUT(events);
BPF_PERCPU_ARRAY(alloc_bprm_array, struct alloc_bprm_event, 1);

int kretprobe__alloc_bprm(struct pt_regs *ctx) {
    if ((bpf_get_current_pid_tgid() >> 32) != _PID_) {
        return 0;
    }
    int zero = 0;
    struct linux_binprm* bprm = (struct linux_binprm*)PT_REGS_RC(ctx);
    struct alloc_bprm_event* e = alloc_bprm_array.lookup(&zero);
    if (!e)
        return 0;
    bpf_probe_read_kernel(&e->filename, sizeof(e->filename), (void*)bprm->filename);
    bpf_probe_read_kernel(&e->fdpath, sizeof(e->fdpath), (void*)bprm->fdpath);
    bpf_probe_read_kernel(&e->interp, sizeof(e->interp), (void*)bprm->interp);
    events.perf_submit(ctx, e, sizeof(*e));
    return 0;
}
`

type allocBprmEvent struct {
    Filename [256]byte
    Fdpath   [256]byte
    Interp   [256]byte
}

func (a *allocBprmEvent) Render() *data.AnalyseData {
    s := fmt.Sprintf("after `%v`, filename: %q, fdpath: %q, interp: %q\n",
        "alloc_bprm", bytes2Str(a.Filename[:]), bytes2Str(a.Fdpath[:]), bytes2Str(a.Interp[:]))
    return data.NewAnalyseData("alloc_bprm", data.NewData(data.MarkdownType, s))
}

type allocBprm struct {
    *BaseMonitorModule
}

func init() {
    ModuleInit(allocBprm{}, true)
}

func (a allocBprm) Monitor() string {
    return "alloc_bprm"
}

func (a allocBprm) Source() string {
    return allowBprmSource
}

func (a allocBprm) Events() []*bcc.Event {
    ke := bcc.NewKretprobeEvent("kretprobe__alloc_bprm", "alloc_bprm", -1)
    return []*bcc.Event{ke}
}

func (a allocBprm) Resolve(m *bpf.Module, ch chan<- *data.AnalyseData, ready chan<- struct{}, stop <-chan struct{}) {
    table := bpf.NewTable(m.TableId("events"), m)

    channel := make(chan []byte)
    perMap, err := bpf.InitPerfMap(table, channel, nil)
    if err != nil {
        log.Fatalf(system.Error("(%s, %s) Failed to init perf map: %v\n"), a.Monitor(), "events", err)
    }

    ok := make(chan []struct{})
    go func() {
        defer func() { close(ok) }()
        for {
            select {
            case d := <-channel:
                analyseData, err := a.Render(d, &allocBprmEvent{})
                if err != nil {
                    fmt.Println(err)
                } else {
                    ch <- analyseData
                }
                return
            case ready <- struct{}{}:
                ready = make(chan struct{})
            }
        }
    }()

    perMap.Start()
    <-ok
    perMap.Stop()
}
