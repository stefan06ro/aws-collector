package auth

import (
	"github.com/giantswarm/aws-collector/flag/service/installation/guest/kubernetes/api/auth/provider"
)

type Auth struct {
	Provider provider.Provider
}
