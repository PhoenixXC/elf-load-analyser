#include <linux/fs.h>
#include <linux/sched.h>
#include <uapi/linux/ptrace.h>

#include "_dev.h"

struct call_sys_execve {};

BPF_PERF_OUTPUT(call_events);
BPF_PERF_OUTPUT(ret_events);

int syscall__execve(struct pt_regs *ctx, const char __user *filename,
                    const char __user *const __user *__argv,
                    const char __user *const __user *__envp) {
    if ((bpf_get_current_pid_tgid() >> 32) != _PID_) {
        return 0;
    }
    struct call_sys_execve e = {};
    call_events.perf_submit(ctx, &e, sizeof(struct call_sys_execve));
    return 0;
}
int do_ret_sys_execve(struct pt_regs *ctx) {
    if ((bpf_get_current_pid_tgid() >> 32) != _PID_) {
        return 0;
    }
    int retval = PT_REGS_RC(ctx);
    int8_t ret;
    bpf_probe_read_kernel(&ret,sizeof(ret), (void*)&retval);
    ret_events.perf_submit(ctx, &ret, sizeof(ret));
    return 0;
}