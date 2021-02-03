package system

import (
    "fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetOSArch(t *testing.T) {
	os := GetSysOS()
	assert.Equal(t, "linux", os)
}

func TestGetKernelVersion(t *testing.T) {
	assert.Equal(t, "5.10.7-3-MANJARO", GetKernelVersion())
}

func TestGetKernelConfigs(t *testing.T) {
	fmt.Println(GetKernelConfigs())
}
