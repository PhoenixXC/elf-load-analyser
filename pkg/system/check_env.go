package system

import (
    "log"
    "os/exec"
    "runtime"

    "github.com/phoenixxc/elf-load-analyser/pkg/validate"
)

// requiredConfigs kernel required configuration, see https://github.com/iovisor/bcc/blob/master/INSTALL.md#kernel-configuration
var requiredConfigs = []string{
    "CONFIG_BPF",
    "CONFIG_BPF_SYSCALL",
    "CONFIG_BPF_JIT",
    // 不支持三目运算符也太...
    func() string {
        if GetKernelVersion() >= "4.7" {
            return "CONFIG_HAVE_EBPF_JIT"
        }
        return "CONFIG_HAVE_BPF_JIT"
    }(),
    "CONFIG_BPF_EVENTS",
}

func CheckEnv() {
    os, arch := GetSysOS(), runtime.GOARCH
    log.Printf("OS: %s\tARCH: %s\n", os, arch)

    // Check os
    validate.EqualWithTip("linux", os, "Unsupported platform, just work on linux")

    // Check kernel version
    validate.WithTip("4.1", GetKernelVersion(),
        func(expected, actual interface{}) bool {
            return actual.(string) >= expected.(string)
        },
        "Kernel version too slow, linux kernel version 4.1 or newer is required\n"+
            "You can see \"https://github.com/iovisor/bcc/blob/master/INSTALL.md\"",
    )

    // Check kernel config
    validate.WithTip(requiredConfigs, GetKernelConfigs(),
        func(expected, actual interface{}) bool {
            kernelConfigs := actual.(map[string]bool)
            for _, entry := range expected.([]string) {
                if !kernelConfigs[entry] {
                    return false
                }
            }
            return true
        },
        `The kernel should have been compiled with the following flags set:
    CONFIG_BPF=y
    CONFIG_BPF_SYSCALL=y
    CONFIG_BPF_JIT=y
    # [for Linux kernel versions 4.1 through 4.6]
    CONFIG_HAVE_BPF_JIT=y
    # [for Linux kernel versions 4.7 and later]
    CONFIG_HAVE_EBPF_JIT=y
    # [optional, for kprobes]
    CONFIG_BPF_EVENTS=y

you can check your kernel config use cmd "zcat /proc/config.gz | zgrep '^[^#].*BPF.*'"
for more information, see "https://github.com/iovisor/bcc/blob/master/INSTALL.md#kernel-configuration"
`)

    // Check if bcc is installed

    log.Printf(Emphasize("The program depend on bcc, please make sure you had install bcc, "+
        "for more information, see %q\n"), "https://github.com/iovisor/bcc/blob/master/INSTALL.md")
    binary, lookErr := exec.LookPath("bcc")
    if lookErr != nil {
        log.Fatalf("Bcc cannot find, %v", lookErr)
    }
    log.Printf("Found bcc: %q\n", binary)

}
