package collector

import (
	"context"
	"math"
	"strconv"
	"time"

	infrastructurev1alpha2 "github.com/giantswarm/apiextensions/v2/pkg/apis/infrastructure/v1alpha2"
	"github.com/giantswarm/k8smetadata/pkg/annotation"
	"github.com/giantswarm/k8smetadata/pkg/label"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/senseyeio/duration"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	// subsystemUpdate will become the second part of the metric name, right after
	// namespace.
	subsystemUpdate  = "update"
	DefaultBatchSize = "0.3"
	DefaultPauseTime = "PT15M"
	labelNodepool    = "node_pool_id"
)

var (
	updateBatchPercentage *prometheus.Desc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystemUpdate, "max_batch_percentage"),
		"Max percentage of worker nodes that can be rolled at once during an upgrade for a given node pool.",
		[]string{
			labelCluster,
			labelNodepool,
		},
		nil,
	)

	updateBatchNumber *prometheus.Desc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystemUpdate, "max_batch_number"),
		"Max number of worker nodes that can be rolled at once during an upgrade for a given node pool.",
		[]string{
			labelCluster,
			labelNodepool,
		},
		nil,
	)

	updatePauseTime *prometheus.Desc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystemUpdate, "pause_time_seconds"),
		"The pause time in seconds between rolling batches of worker nodes during an upgrade for a given node pool.",
		[]string{
			labelCluster,
			labelNodepool,
		},
		nil,
	)
)

type UpdateConfig struct {
	Helper *helper
	Logger micrologger.Logger
}

type Update struct {
	helper *helper
	logger micrologger.Logger
}

type updateInfo struct {
	nodePoolID      string
	clusterID       string
	pauseTime       float64
	batchPercentage float64
	batchNumber     float64
}

func NewUpdate(config UpdateConfig) (*Update, error) {
	if config.Helper == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Helper must not be empty", config)
	}
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Logger must not be empty", config)
	}

	np := &Update{
		helper: config.Helper,
		logger: config.Logger,
	}

	return np, nil
}

func (np *Update) Collect(ch chan<- prometheus.Metric) error {
	ctx := context.Background()

	var list infrastructurev1alpha2.AWSMachineDeploymentList
	{
		err := np.helper.clients.CtrlClient().List(
			ctx,
			&list,
		)
		if err != nil {
			return microerror.Mask(err)
		}
	}
	var nodePools []updateInfo
	{
		for _, md := range list.Items {
			batch, pause, err := np.getUpdateAnnotations(md)
			if err != nil {
				return microerror.Mask(err)
			}
			pauseTime, batchNumber, batchPercentage, err := calculateUpdateMetrics(batch, pause, md.Spec.NodePool.Scaling.Min, md.Spec.NodePool.Scaling.Max)
			if err != nil {
				return microerror.Mask(err)
			}
			nodePool := updateInfo{
				nodePoolID:      md.Labels[label.MachineDeployment],
				clusterID:       md.Labels[label.Cluster],
				pauseTime:       math.Ceil(pauseTime),
				batchPercentage: batchPercentage,
				batchNumber:     batchNumber,
			}
			nodePools = append(nodePools, nodePool)
		}
	}

	for _, np := range nodePools {
		ch <- prometheus.MustNewConstMetric(
			updateBatchPercentage,
			prometheus.GaugeValue,
			np.batchPercentage,
			np.clusterID,
			np.nodePoolID,
		)

		ch <- prometheus.MustNewConstMetric(
			updateBatchNumber,
			prometheus.GaugeValue,
			np.batchNumber,
			np.clusterID,
			np.nodePoolID,
		)

		ch <- prometheus.MustNewConstMetric(
			updatePauseTime,
			prometheus.GaugeValue,
			float64(np.pauseTime),
			np.clusterID,
			np.nodePoolID,
		)
	}

	return nil
}

func (np *Update) Describe(ch chan<- *prometheus.Desc) error {
	ch <- updateBatchPercentage
	ch <- updateBatchNumber
	ch <- updatePauseTime

	return nil
}

func (np *Update) getUpdateAnnotations(md infrastructurev1alpha2.AWSMachineDeployment) (string, string, error) {
	var batch, time string

	// get info from AWSMachineDeployment
	if md.Annotations != nil {
		batch = md.Annotations[annotation.AWSUpdateMaxBatchSize]
		time = md.Annotations[annotation.AWSUpdatePauseTime]
	}
	if batch != "" && time != "" {
		return batch, time, nil
	}

	// get info from AWSCluster
	var cl infrastructurev1alpha2.AWSCluster
	{
		var clusterList infrastructurev1alpha2.AWSClusterList
		err := np.helper.clients.CtrlClient().List(
			context.Background(),
			&clusterList,
			client.MatchingLabels{label.Cluster: md.Labels[label.Cluster]},
		)
		if err != nil {
			return "", "", microerror.Mask(err)
		}
		if len(clusterList.Items) != 1 {
			return "", "", microerror.Maskf(notFoundError, "Tried to find one AWSCluster CR with ID %s, found %v.", md.Labels[label.Cluster], len(clusterList.Items))
		}
		cl = clusterList.Items[0]
	}

	if cl.Annotations != nil {
		if batch == "" {
			batch = cl.Annotations[annotation.AWSUpdateMaxBatchSize]
		}
		if time == "" {
			time = cl.Annotations[annotation.AWSUpdatePauseTime]
		}
	}
	if batch != "" && time != "" {
		return batch, time, nil
	}
	// get default values
	if batch == "" {
		batch = DefaultBatchSize
	}
	if time == "" {
		time = DefaultPauseTime
	}
	return batch, time, nil
}

func calculateUpdateMetrics(batch string, pause string, min int, max int) (float64, float64, float64, error) {
	var batchNumber, batchPercentage, pauseTime float64

	//check whether max batch size is given as percentage or integer and calculate the other accordingly
	if n, err := strconv.Atoi(batch); err == nil {
		// In case the batch is given as an integer, we calculate the max percentage of rolling nodes based on the minimum scaled setting.
		// This can be more than 100 percent! The lowest min scaling we consider is 1 to avoid division by 0.
		batchNumber = float64(n)
		batchPercentage = float64(n) / math.Max(float64(min), 1)
	} else if p, err := strconv.ParseFloat(batch, 64); err == nil {
		// In case the batch is given as a percentage, we calculate the actual max number based on the maximum scaled setting.
		batchNumber = float64(max) * p
		batchPercentage = p
	} else if err != nil {
		return 0, 0, 0, microerror.Mask(err)
	}

	// Calculate the pause time in seconds
	// The pause time is given as ISO8601 string and we transform it to time in seconds
	duration, err := duration.ParseISO8601(pause)
	if err != nil {
		return 0, 0, 0, microerror.Mask(err)
	}
	shifted := duration.Shift(time.Now())
	pauseTime = time.Until(shifted).Seconds()

	return pauseTime, batchNumber, batchPercentage, nil
}
