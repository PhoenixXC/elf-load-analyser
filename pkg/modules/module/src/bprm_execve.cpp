// https://elixir.bootlin.com/linux/v5.10.7/source/fs/exec.c#L1775

#include <linux/fs.h>
#include <linux/sched.h>

#include "_dev.h"
#include "common.h"

TDATA(bprm_execve_event, 
); 
BPF_PERF_OUTPUT(call_event);

int kprobe__bprm_execve(struct pt_regs* ctx) {
    if ((bpf_get_current_pid_tgid() >> 32) != _PID_) {
        return 0;
    }
    struct bprm_execve_event event = {};
    init_tdata(&event);
    init_tdata(&event);
    call_event.perf_submit((void*)ctx, &event, sizeof(event));
    return 0;
}