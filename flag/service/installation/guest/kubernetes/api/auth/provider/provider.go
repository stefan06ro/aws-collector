package provider

import (
	"github.com/giantswarm/aws-collector/flag/service/installation/guest/kubernetes/api/auth/provider/oidc"
)

type Provider struct {
	OIDC oidc.OIDC
}
