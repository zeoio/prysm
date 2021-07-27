
package migration

import (
	v1Alpha1 "github.com/prysmaticlabs/prysm/proto/prysm/v1alpha1"
	v2 "github.com/prysmaticlabs/prysm/proto/prysm/v2"
)

func V1Alpha1ToV2VoluntaryExit(src *v1Alpha1.VoluntaryExit) *v2.VoluntaryExit {
	return &v2.VoluntaryExit{}
}

func V2ToV1Alpha1VoluntaryExit(src *v2.VoluntaryExit) *v1Alpha1.VoluntaryExit {
	return &v1Alpha1.VoluntaryExit{}
}
