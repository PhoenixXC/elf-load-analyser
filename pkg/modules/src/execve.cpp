#include <uapi/linux/ptrace.h>
#include <linux/sched.h>
#include <linux/fs.h>

#include "_dev.h"

BPF_PERF_OUTPUT(events);

int syscall__execve(struct pt_regs *ctx,
    const char __user *filename,
    const char __user *const __user *__argv,
    const char __user *const __user *__envp)
{
    
    return 0;
}
int do_ret_sys_execve(struct pt_regs *ctx)
{
    int retval = PT_REGS_RC(ctx);
    return 0;
}