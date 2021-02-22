#include <linux/fs.h>
#include <linux/sched.h>
#include <linux/timekeeping.h>
#include <uapi/linux/ptrace.h>

#include "_dev.h"

struct call_sys_execve {
    uint64_t ts;
};

struct ret_sys_execve {
    uint64_t ts;
    int8_t ret;
};

BPF_PERF_OUTPUT(call_events);
BPF_PERF_OUTPUT(ret_events);

int syscall__execve(struct pt_regs *ctx, const char __user *filename,
                    const char __user *const __user *__argv,
                    const char __user *const __user *__envp) {
    if ((bpf_get_current_pid_tgid() >> 32) != _PID_) {
        return 0;
    }
    struct call_sys_execve e = {};
    uint64_t ns = bpf_ktime_get_ns();
    bpf_probe_read_kernel(&e.ts, sizeof(e.ts), (void *)(&ns));
    call_events.perf_submit(ctx, &e, sizeof(struct call_sys_execve));
    return 0;
}
int do_ret_sys_execve(struct pt_regs *ctx) {
    if ((bpf_get_current_pid_tgid() >> 32) != _PID_) {
        return 0;
    }
    int retval = PT_REGS_RC(ctx);
    struct ret_sys_execve e = {};
    bpf_probe_read_kernel(&e.ret, sizeof(e.ret), (void *)&retval);
    uint64_t ns = bpf_ktime_get_ns();
    bpf_probe_read_kernel(&e.ts, sizeof(e.ts), (void *)(&ns));
    ret_events.perf_submit(ctx, &e, sizeof(e));
    return 0;
}