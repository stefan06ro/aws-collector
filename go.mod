module github.com/giantswarm/aws-collector

go 1.14

require (
	github.com/aws/aws-sdk-go v1.39.4
	github.com/giantswarm/apiextensions/v2 v2.6.2
	github.com/giantswarm/exporterkit v0.2.1
	github.com/giantswarm/k8sclient/v4 v4.1.0
	github.com/giantswarm/k8smetadata v0.3.0
	github.com/giantswarm/microendpoint v0.2.0
	github.com/giantswarm/microerror v0.3.0
	github.com/giantswarm/microkit v0.2.2
	github.com/giantswarm/micrologger v0.5.0
	github.com/giantswarm/operatorkit/v2 v2.0.2
	github.com/patrickmn/go-cache v2.1.0+incompatible
	github.com/prometheus/client_golang v1.11.0
	github.com/senseyeio/duration v0.0.0-20180430131211-7c2a214ada46
	github.com/spf13/viper v1.8.1
	golang.org/x/crypto v0.0.0-20201002170205-7f63de1d35b0 // indirect
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c
	k8s.io/api v0.18.19
	k8s.io/apimachinery v0.18.19
	k8s.io/client-go v0.18.19
	sigs.k8s.io/cluster-api v0.4.0
	sigs.k8s.io/controller-runtime v0.6.3
)

replace (
	github.com/coreos/etcd v3.3.10+incompatible => github.com/coreos/etcd v3.3.25+incompatible
	github.com/coreos/etcd v3.3.13+incompatible => github.com/coreos/etcd v3.3.25+incompatible
	github.com/dgrijalva/jwt-go => github.com/dgrijalva/jwt-go/v4 v4.0.0-preview1
	github.com/gogo/protobuf v1.3.1 => github.com/gogo/protobuf v1.3.2
	sigs.k8s.io/cluster-api => github.com/giantswarm/cluster-api v0.3.13-gs
)
