#include <linux/binfmts.h>
#include <linux/fs.h>
#include <linux/sched.h>

#include "_dev.h"

struct alloc_bprm_event {
    uint64_t ts;
    char filename[256];
    char fdpath[256];
    char interp[256];
    uint64_t cur_top_of_mem;
    uint64_t rlim_cur;
    uint64_t rlim_max;
};
BPF_PERF_OUTPUT(events);
BPF_PERCPU_ARRAY(alloc_bprm_array, struct alloc_bprm_event, 1);

int kretprobe__alloc_bprm(struct pt_regs* ctx) {
    if ((bpf_get_current_pid_tgid() >> 32) != _PID_) {
        return 0;
    }
    int zero = 0;
    __kernel_ulong_t rlim_val = 0;
    unsigned long tmp_val = 0;
    struct linux_binprm* bprm = (struct linux_binprm*)PT_REGS_RC(ctx);
    struct alloc_bprm_event* e = alloc_bprm_array.lookup(&zero);
    if (!e) return 0;
    bpf_probe_read_kernel(&e->filename, sizeof(e->filename),
                          (void*)bprm->filename);
    bpf_probe_read_kernel(&e->fdpath, sizeof(e->fdpath), (void*)bprm->fdpath);
    bpf_probe_read_kernel(&e->interp, sizeof(e->interp), (void*)bprm->interp);
    bpf_probe_read_kernel(&rlim_val, sizeof(rlim_val),
                          (void*)&bprm->rlim_stack.rlim_cur);
    e->rlim_cur = (uint64_t)rlim_val;
    bpf_probe_read_kernel(&rlim_val, sizeof(rlim_val),
                          (void*)&bprm->rlim_stack.rlim_max);
    e->rlim_max = (uint64_t)rlim_val;
    bpf_probe_read_kernel(&tmp_val, sizeof(bprm->p), (void*)&bprm->p);
    e->cur_top_of_mem = tmp_val;
    uint64_t ns = bpf_ktime_get_ns();
    bpf_probe_read_kernel(&e->ts, sizeof(e->ts), (void*)&ns);
    events.perf_submit(ctx, e, sizeof(*e));
    return 0;
}