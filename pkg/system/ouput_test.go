package system

import (
	"log"
	"testing"
)

func TestWarn(t *testing.T) {
	log.Printf(Warn("This is warn message: %q\n"), "Miss file")
	log.Printf(Success("This is success message: %q\n"), "build success")
	log.Printf(Error("This is err message: %q\n"), "Run err")
	log.Printf(Check("TEST PASSED: %q\n", true), "OK")
	log.Printf(Check("TEST_FAILED: %q\n", false), "equal err")
}
