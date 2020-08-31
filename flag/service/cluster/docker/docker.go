package docker

import (
	"github.com/giantswarm/aws-collector/flag/service/cluster/docker/daemon"
)

type Docker struct {
	Daemon daemon.Daemon
}
