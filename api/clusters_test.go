package api

import (
	"context"
	discoveryv1alpha1 "github.com/liqotech/liqo/apis/discovery/v1alpha1"
	sharingv1alpha1 "github.com/liqotech/liqo/apis/sharing/v1alpha1"
	liqoconst "github.com/liqotech/liqo/pkg/consts"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
)

func TestGetClusters(t *testing.T) {
	fc1 := &discoveryv1alpha1.ForeignCluster{
		ObjectMeta: metav1.ObjectMeta{Name: "fc1"},
		Spec: discoveryv1alpha1.ForeignClusterSpec{
			ClusterIdentity: discoveryv1alpha1.ClusterIdentity{ClusterID: "cid-fc1"},
		},
	}
	err := k8sClient.Create(context.Background(), fc1)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	ro1 := &sharingv1alpha1.ResourceOffer{
		ObjectMeta: metav1.ObjectMeta{
			Name: "ro1", Namespace: "test",
			Labels: map[string]string{
				liqoconst.ReplicationRequestedLabel: "false",
				liqoconst.ReplicationOriginLabel:    "cid-fc1",
			},
		},
	}
	err = k8sClient.Create(context.Background(), ro1)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	resp, err := NewClientForMethod(srv.GetClusters).
		GetClustersWithResponse(context.Background())
	if err != nil {
		panic(err)
	}
	clusters := *resp.JSON200
	if len(clusters) != 1 {
		t.Fail()
	}
	if clusters[0].Id != "cid-fc1" {
		t.Fail()
	}
	if clusters[0].IncomingPeering != true {
		t.Fail()
	}
	if clusters[0].OutgoingPeering != false {
		t.Fail()
	}
}
