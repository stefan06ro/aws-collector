package aws

import (
	"github.com/giantswarm/aws-collector/flag/service/aws/hostaccesskey"
	"github.com/giantswarm/aws-collector/flag/service/aws/trustedadvisor"
)

type AWS struct {
	HostAccessKey  hostaccesskey.HostAccessKey
	Region         string
	TrustedAdvisor trustedadvisor.TrustedAdvisor
}
