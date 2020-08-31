package guest

import (
	"github.com/giantswarm/aws-collector/flag/service/guest/ignition"
	"github.com/giantswarm/aws-collector/flag/service/guest/ssh"
)

type Guest struct {
	Ignition ignition.Ignition
	SSH      ssh.SSH
}
