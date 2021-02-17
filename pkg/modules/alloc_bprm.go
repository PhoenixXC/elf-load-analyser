package modules

import (
    "fmt"
    "github.com/phoenixxc/elf-load-analyser/pkg/bcc"
    "github.com/phoenixxc/elf-load-analyser/pkg/data"
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
    uint64_t cur_top_of_mem;
    uint64_t rlim_cur;
    uint64_t rlim_max;
};
BPF_PERF_OUTPUT(events);
BPF_PERCPU_ARRAY(alloc_bprm_array, struct alloc_bprm_event, 1);

int kretprobe__alloc_bprm(struct pt_regs *ctx) {
    if ((bpf_get_current_pid_tgid() >> 32) != _PID_) {
        return 0;
    }
    int zero = 0;
    __kernel_ulong_t rlim_val = 0;
    unsigned long tmp_val = 0;
    struct linux_binprm* bprm = (struct linux_binprm*)PT_REGS_RC(ctx);
    struct alloc_bprm_event* e = alloc_bprm_array.lookup(&zero);
    if (!e)
        return 0;
    bpf_probe_read_kernel(&e->filename, sizeof(e->filename), (void*)bprm->filename);
    bpf_probe_read_kernel(&e->fdpath, sizeof(e->fdpath), (void*)bprm->fdpath);
    bpf_probe_read_kernel(&e->interp, sizeof(e->interp), (void*)bprm->interp);
    bpf_probe_read_kernel(&rlim_val, sizeof(rlim_val), (void*)&bprm->rlim_stack.rlim_cur);
    e->rlim_cur = (uint64_t)rlim_val;
    bpf_probe_read_kernel(&rlim_val, sizeof(rlim_val), (void*)&bprm->rlim_stack.rlim_max);
    e->rlim_max = (uint64_t)rlim_val;
    bpf_probe_read_kernel(&tmp_val, sizeof(bprm->p), (void*)&bprm->p);
    e->cur_top_of_mem = tmp_val;
    events.perf_submit(ctx, e, sizeof(*e));
    return 0;
}
`

type allocBprmEvent struct {
    Filename    [256]byte
    Fdpath      [256]byte
    Interp      [256]byte
    CurTopOfMem uint64
    RlimCur     uint64
    RlimMax     uint64
}

func (a *allocBprmEvent) Render() *data.AnalyseData {
    s := fmt.Sprintf("after `%v`, filename: %q, fdpath: %q, interp: %q, rlimit stack cur: 0x%X,"+
        " rlimit stack max: 0x%X, current of top mem: 0x%X\n",
        "alloc_bprm", bytes2Str(a.Filename[:]), bytes2Str(a.Fdpath[:]), bytes2Str(a.Interp[:]),
        a.RlimCur, a.RlimMax, a.CurTopOfMem)
    return data.NewAnalyseData("alloc_bprm", data.NewData(data.MarkdownType, s))
}

type allocBprm struct {
    *BaseMonitorModule
}

func init() {
    m := NewPerfResolveMonitorModule(&allocBprm{})
    m.RegisterTable("events", false, func(data []byte) (*data.AnalyseData, error) {
        log.Println(m.Monitor(), "Call alloc bprm finished")
        return m.Render(data, &allocBprmEvent{})
    })
    ModuleInit(m, true)
}

func (a *allocBprm) Monitor() string {
    return "alloc_bprm"
}

func (a *allocBprm) Source() string {
    return allowBprmSource
}

func (a *allocBprm) Events() []*bcc.Event {
    ke := bcc.NewKretprobeEvent("kretprobe__alloc_bprm", "alloc_bprm", -1)
    return []*bcc.Event{ke}
}
