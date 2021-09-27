//go:build !linux && !darwin
// +build !linux,!darwin

package graceful

import (
	"os"
)

func getInterruptSigs() []os.Signal {
	return []os.Signal{os.Interrupt}
}
