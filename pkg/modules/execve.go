package modules

import (
    "fmt"
    bpf "github.com/iovisor/gobpf/bcc"
    "time"

    "github.com/phoenixxc/elf-load-analyser/pkg/bcc"
    "github.com/phoenixxc/elf-load-analyser/pkg/factory"
)

//goland:noinspection ALL
const source = `
#include <linux/fs.h>

struct exec_event {
	char filename[256];
};
BPF_PERF_OUTPUT(events);

int kprobe__do_execveat_common(struct pt_regs* ctx, int fd, struct filename* filename) {
	struct exec_event event = {};
	bpf_probe_read_kernel(&(event.filename), sizeof(event.filename), (void*)filename->name);
	events.perf_submit((void*)ctx, (void*)&event, sizeof(event));
	return 0;
}
`

func init() {
    m := bcc.NewMonitor(source, []string{}, func(m *bpf.Module) {
        time.Sleep(20 * time.Second)
        fmt.Println("Do something, resolve data and send result....")
    })
    e := bcc.NewKprobeEvent("kprobe__do_execveat_common", "do_execveat_common", -1)
    m.AddEvent(e)
    factory.Register("execve_handle", m)
}
