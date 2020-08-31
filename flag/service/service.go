package service

import (
	"github.com/giantswarm/operatorkit/v2/pkg/flag/service/kubernetes"

	"github.com/giantswarm/aws-collector/flag/service/aws"
	"github.com/giantswarm/aws-collector/flag/service/cluster"
	"github.com/giantswarm/aws-collector/flag/service/guest"
	"github.com/giantswarm/aws-collector/flag/service/installation"
	"github.com/giantswarm/aws-collector/flag/service/registry"
)

type Service struct {
	AWS          aws.AWS
	Cluster      cluster.Cluster
	Guest        guest.Guest
	Installation installation.Installation
	Kubernetes   kubernetes.Kubernetes
	Registry     registry.Registry
}
