#include <linux/fs.h>
#include <linux/sched.h>
#include <linux/timekeeping.h>
#include <uapi/linux/ptrace.h>

#include "_dev.h"
#include "common.h"

TDATA(ret_sys_execve, 
    int8_t ret; 
);

TDATA(call_sys_execve, TEMPTY);

BPF_PERF_OUTPUT(call_event);
BPF_PERF_OUTPUT(ret_event);

int syscall__execve(struct pt_regs *ctx, const char __user *filename,
                    const char __user *const __user *__argv,
                    const char __user *const __user *__envp) {
    if ((bpf_get_current_pid_tgid() >> 32) != _PID_) {
        return 0;
    }
    struct call_sys_execve e = {};
    init_tdata(&e);
    bpf_trace_printk("call_syscall_execve => ns: %llu\n", e.ts); 
    call_event.perf_submit(ctx, &e, sizeof(struct call_sys_execve));
    return 0;
}
int do_ret_sys_execve(struct pt_regs *ctx) {
    if ((bpf_get_current_pid_tgid() >> 32) != _PID_) {
        return 0;
    }
    int retval = PT_REGS_RC(ctx);
    struct ret_sys_execve e = {};
    init_tdata(&e);
    bpf_probe_read_kernel(&e.ret, sizeof(e.ret), (void *)&retval);
    ret_event.perf_submit(ctx, &e, sizeof(struct ret_sys_execve));

    return 0;
}