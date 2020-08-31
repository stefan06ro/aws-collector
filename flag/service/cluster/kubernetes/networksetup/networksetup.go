package networksetup

import (
	"github.com/giantswarm/aws-collector/flag/service/cluster/kubernetes/networksetup/docker"
)

type NetworkSetup struct {
	Docker docker.Docker
}
