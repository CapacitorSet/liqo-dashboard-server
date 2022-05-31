package api

import (
	"context"
	"github.com/ghodss/yaml"
	"github.com/labstack/echo/v4"
	discoveryv1alpha1 "github.com/liqotech/liqo/apis/discovery/v1alpha1"
	sharingv1alpha1 "github.com/liqotech/liqo/apis/sharing/v1alpha1"
	liqoconst "github.com/liqotech/liqo/pkg/consts"
	"net/http"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type PeeringEntry struct {
	Incoming bool
	Outgoing bool
}

func (s APIServer) GetClusters(c echo.Context) error {
	var liqoClusters discoveryv1alpha1.ForeignClusterList
	var roList sharingv1alpha1.ResourceOfferList
	incomingPeerings := make(map[string]bool)
	outgoingPeerings := make(map[string]bool)
	jsonClusters := []ForeignCluster{} // do not convert to var, or it will be serialized to `null` if there are no clusters
	ctx := context.TODO()
	err := s.List(ctx, &liqoClusters, &client.ListOptions{})
	if err != nil {
		return err
	}
	err = s.List(ctx, &roList, &client.ListOptions{})
	if err != nil {
		return err
	}
	for _, ro := range roList.Items {
		if ro.ObjectMeta.Labels[liqoconst.ReplicationRequestedLabel] == "true" {
			outgoingPeerings[ro.ObjectMeta.Labels[liqoconst.ReplicationDestinationLabel]] = true
		} else {
			incomingPeerings[ro.ObjectMeta.Labels[liqoconst.ReplicationOriginLabel]] = true
		}
	}
	for _, cluster := range liqoClusters.Items {
		clusterID := cluster.Spec.ClusterIdentity.ClusterID
		rawYamlBytes, err := yaml.Marshal(cluster.Spec)
		if err != nil {
			return err
		}
		rawYaml := string(rawYamlBytes)
		jsonClusters = append(jsonClusters, ForeignCluster{
			Id:              clusterID,
			IncomingPeering: incomingPeerings[clusterID],
			OutgoingPeering: outgoingPeerings[clusterID],
			Ip:              cluster.Spec.ForeignAuthURL,
			Name:            cluster.ObjectMeta.Name,
			RawYaml:         rawYaml,
		})
	}
	return c.JSON(http.StatusOK, jsonClusters)
}
