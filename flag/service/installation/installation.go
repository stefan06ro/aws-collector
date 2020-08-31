package installation

import (
	"github.com/giantswarm/aws-collector/flag/service/installation/guest"
)

type Installation struct {
	Name  string
	Guest guest.Guest
}
