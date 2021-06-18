package collector

import (
	"context"
	"encoding/json"
	"time"

	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/sync/errgroup"

	clientaws "github.com/giantswarm/aws-collector/client/aws"
	"github.com/giantswarm/aws-collector/service/controller/key"
	"github.com/giantswarm/aws-collector/service/internal/cache"
)

const (
	// __SubnetCache__ is used as temporal cache key to save Subnet response.
	prefixSubnetcacheKey = "__SubnetCache__"
)

const (
	subsystemSubnet = "subnet"
)

var (
	subnetsDesc *prometheus.Desc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystemSubnet, "available_ips"),
		"Subnet information.",
		[]string{
			labelAccountID,
			labelCIDR,
			labelCluster,
			labelID,
			labelInstallation,
			labelOrganization,
			labelStack,
			labelState,
			labelAvailabilityZone,
			labelAccount,
			labelVPC,
		},
		nil,
	)
)

type SubnetConfig struct {
	Helper *helper
	Logger micrologger.Logger

	InstallationName string
}

type Subnet struct {
	cache  *subnetCache
	helper *helper
	logger micrologger.Logger

	installationName string
}

type subnetCache struct {
	cache *cache.StringCache
}

type subnetInfoResponse struct {
	Subnets []subnetInfo
}

type subnetInfo struct {
	Name         string
	AvailableIPs int64
	Tags         map[string]string
}

func NewSubnet(config SubnetConfig) (*Subnet, error) {
	if config.Helper == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Helper must not be empty", config)
	}
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Logger must not be empty", config)
	}

	if config.InstallationName == "" {
		return nil, microerror.Maskf(invalidConfigError, "%T.InstallationName must not be empty", config)
	}

	e := &Subnet{
		cache:  newSubnetCache(time.Minute * 5),
		helper: config.Helper,
		logger: config.Logger,

		installationName: config.InstallationName,
	}

	return e, nil
}

func newSubnetCache(expiration time.Duration) *subnetCache {
	cache := &subnetCache{
		cache: cache.NewStringCache(expiration),
	}

	return cache
}

func (n *subnetCache) Get(key string) (*subnetInfoResponse, error) {
	var c subnetInfoResponse
	raw, exists := n.cache.Get(getSubnetCacheKey(key))
	if exists {
		err := json.Unmarshal(raw, &c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	return &c, nil
}

func (n *subnetCache) Set(key string, content subnetInfoResponse) error {
	contentSerialized, err := json.Marshal(content)
	if err != nil {
		return microerror.Mask(err)
	}

	n.cache.Set(getSubnetCacheKey(key), contentSerialized)

	return nil
}

func getSubnetCacheKey(key string) string {
	return prefixSubnetcacheKey + key
}

func (e *Subnet) Collect(ch chan<- prometheus.Metric) error {
	reconciledClusters, err := e.helper.ListReconciledClusters()
	if err != nil {
		return microerror.Mask(err)
	}

	awsClientsList, err := e.helper.GetAWSClients(context.Background(), reconciledClusters)
	if err != nil {
		return microerror.Mask(err)
	}

	var g errgroup.Group

	for _, item := range awsClientsList {
		awsClients := item

		g.Go(func() error {
			err := e.collectForAccount(context.Background(), ch, awsClients)
			if err != nil {
				return microerror.Mask(err)
			}

			return nil
		})
	}

	err = g.Wait()
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}

func (e *Subnet) Describe(ch chan<- *prometheus.Desc) error {
	ch <- subnetsDesc
	return nil
}

func (e *Subnet) collectForAccount(ctx context.Context, ch chan<- prometheus.Metric, awsClients clientaws.Clients) error {
	account, err := e.helper.AWSAccountID(awsClients)
	if err != nil {
		return microerror.Mask(err)
	}

	var subnetInfo *subnetInfoResponse
	// Check if response is cached
	subnetInfo, err = e.cache.Get(account)
	if err != nil {
		return microerror.Mask(err)
	}

	//Cache empty, getting from API
	if subnetInfo == nil || subnetInfo.Subnets == nil {
		subnetInfo, err = e.getSubnetInfoFromAPI(ctx, awsClients)
		if err != nil {
			return microerror.Mask(err)
		}

		if subnetInfo != nil {
			err = e.cache.Set(account, *subnetInfo)
			if err != nil {
				return microerror.Mask(err)
			}
		}
	}

	if subnetInfo != nil {
		for _, subnet := range subnetInfo.Subnets {
			ch <- prometheus.MustNewConstMetric(
				subnetsDesc,
				prometheus.GaugeValue,
				float64(subnet.AvailableIPs),
				account,
				subnet.Tags["CidrBlock"],
				subnet.Tags[key.TagCluster],
				subnet.Name,
				e.installationName,
				subnet.Tags[key.TagOrganization],
				subnet.Tags[key.TagStack],
				subnet.Tags["State"],
				subnet.Tags["AvailabilityZone"],
				subnet.Tags["OwnerId"],
				subnet.Tags["VpcId"],
			)
		}
	}

	return nil
}

// getSubnetInfoFromAPI collects Subnet Info from AWS API
func (e *Subnet) getSubnetInfoFromAPI(ctx context.Context, awsClients clientaws.Clients) (*subnetInfoResponse, error) {
	var res subnetInfoResponse
	o, err := awsClients.EC2.DescribeSubnets(&ec2.DescribeSubnetsInput{})
	if err != nil {
		return nil, microerror.Mask(err)
	}

	var subnets []subnetInfo
	for _, sn := range o.Subnets {
		subnet := subnetInfo{
			Name:         *sn.SubnetId,
			AvailableIPs: *sn.AvailableIpAddressCount,
			Tags: map[string]string{
				"CidrBlock":        *sn.CidrBlock,
				"AvailabilityZone": *sn.AvailabilityZone,
				"OwnerId":          *sn.OwnerId,
				"VpcId":            *sn.VpcId,
				"State":            *sn.State,
			},
		}
		for _, t := range sn.Tags {
			subnet.Tags[*t.Key] = *t.Value
		}

		if subnet.Tags[key.TagInstallation] != e.installationName {
			continue
		}
		subnets = append(subnets, subnet)
	}
	res.Subnets = subnets
	return &res, nil
}
