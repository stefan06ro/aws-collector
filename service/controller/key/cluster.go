package key

import (
	infrastructurev1alpha3 "github.com/giantswarm/apiextensions/v3/pkg/apis/infrastructure/v1alpha3"
)

const (
	// CloudConfigVersion defines the version of k8scloudconfig in use. It is used
	// in the main stack output and S3 object paths.
	CloudConfigVersion = "v_6_1_0"
	CloudProvider      = "aws"
)

const (
	EC2RoleK8s   = "EC2-K8S-Role"
	EC2PolicyK8s = "EC2-K8S-Policy"
)

const (
	EtcdPort                     = 2379
	EtcdPrefix                   = "giantswarm.io"
	KubernetesSecurePort         = 443
	KubernetesApiHealthCheckPort = 8089
)

const (
	HAMasterSnapshotIDValue = "ha-master-migration"
)

const (
	// KubernetesAPIHealthzVersion is a tag representing the version of
	// https://github.com/giantswarm/k8s-api-healthz/ used.
	KubernetesAPIHealthzVersion = "0.1.1"
	// K8sSetupNetworkEnvironment is a tag representing the version of
	// https://github.com/giantswarm/k8s-setup-network-environment used.
	K8sSetupNetworkEnvironment = "0.2.0"
)

// AWS Tags used for cost analysis and general resource tagging.
const (
	TagAvailabilityZone        = "giantswarm.io/availability-zone"
	TagCluster                 = "giantswarm.io/cluster"
	TagClusterType             = "giantswarm.io/cluster-type"
	TagClusterTypeControlPlane = "control-plane"
	TagControlPlane            = "giantswarm.io/control-plane"
	TagInstallation            = "giantswarm.io/installation"
	TagMachineDeployment       = "giantswarm.io/machine-deployment"
	TagOrganization            = "giantswarm.io/organization"
	TagRouteTableType          = "giantswarm.io/route-table-type"
	TagStack                   = "giantswarm.io/stack"
	TagSnapshot                = "giantswarm.io/snapshot"
	TagSubnetType              = "giantswarm.io/subnet-type"
)

const (
	StackTCCP  = "tccp"
	StackTCCPF = "tccpf"
	StackTCCPI = "tccpi"
	StackTCCPN = "tccpn"
	StackTCNP  = "tcnp"
	StackTCNPF = "tcnpf"
)

const (
	LifeCycleHookControlPlane = "ControlPlane"
	LifeCycleHookNodePool     = "NodePool"
)

const (
	RefWorkerASG = "workerAutoScalingGroup"
)

const (
	// ComponentOS is the name of the component specified in a Release CR which
	// determines the version of the OS to be used for tenant cluster nodes and
	// is ultimately transformed into an AMI based on TC region.
	ComponentOS = "containerlinux"
)

func CredentialName(cluster infrastructurev1alpha3.AWSCluster) string {
	return cluster.Spec.Provider.CredentialSecret.Name
}

func CredentialNamespace(cluster infrastructurev1alpha3.AWSCluster) string {
	return cluster.Spec.Provider.CredentialSecret.Namespace
}
