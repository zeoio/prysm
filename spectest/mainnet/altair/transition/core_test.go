package transition

import (
	"testing"

	"github.com/prysmaticlabs/prysm/spectest/shared/altair/transition"
)

func TestMainnet_Altair_Transition(t *testing.T) {
	transition.RunCoreTests(t, "mainnet")
}
