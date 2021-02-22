// system env about, eg kernel version, system type
package system

import (
    "bufio"
    "compress/gzip"
    "github.com/phoenixxc/elf-load-analyser/pkg/helper"
    "io/ioutil"
    "log"
    "os"
    "runtime"
    "strings"
    "sync"
)

const (
	kernelReleaseFile  = "/proc/sys/kernel/osrelease"
	kernelConfigGzFile = "/proc/config.gz"
)

// cache version result
var once sync.Once
var kernelVersion string

// GetSysOS get os type and arch
func GetSysOS() string {
	return runtime.GOOS
}

// GetKernelVersion get linux version
func GetKernelVersion() string {
	once.Do(extraKernelVersion)
	return kernelVersion
}

// GetKernelConfigs get kernel configs from kernelConfigGzFile
func GetKernelConfigs() map[string]bool {
	file, err := os.Open(kernelConfigGzFile)
	if err != nil {
		log.Fatalf("Open config file %q failed, %v", kernelConfigGzFile, err)
	}
	defer file.Close()

	reader, err := gzip.NewReader(file)
	if err != nil {
		log.Fatalf("Reset file %q err, %v", kernelConfigGzFile, err)
	}
	//goland:noinspection GoUnhandledErrorResult
	defer reader.Close()

	configs := make(map[string]bool)

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		item := strings.TrimSpace(scanner.Text())
		if len(item) > 0 && !strings.HasPrefix(strings.TrimSpace(item), "#") {
			kv := strings.SplitN(item, "=", 2)
			if len(kv) > 1 {
				configs[kv[0]] = strings.ToLower(kv[1]) == "y"
			}
		}
	}

	return configs
}

func extraKernelVersion() {
	// check system type
	helper.EqualWithTip("linux", GetSysOS(), "Unsupported system, the toolkit just for linux")

	file, err := os.Open(kernelReleaseFile)
	if err != nil {
		log.Fatalf("Open release file %q failed, %v", kernelReleaseFile, err)
	}
	defer file.Close()
	release, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatalf("Read %q failed, %v", kernelReleaseFile, err)
	}
	kernelVersion = strings.TrimSpace(string(release))
}
