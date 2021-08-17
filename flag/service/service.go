package service

import (
	"github.com/giantswarm/operatorkit/v4/pkg/flag/service/kubernetes"

	"github.com/giantswarm/aws-collector/flag/service/aws"
	"github.com/giantswarm/aws-collector/flag/service/installation"
)

type Service struct {
	AWS          aws.AWS
	Installation installation.Installation
	Kubernetes   kubernetes.Kubernetes
}
