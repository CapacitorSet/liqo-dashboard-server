package api

import (
	"context"
	"github.com/labstack/echo/v4"
	"errors"
	sharingv1alpha1 "github.com/liqotech/liqo/apis/sharing/v1alpha1"
	liqoconst "github.com/liqotech/liqo/pkg/consts"
	corev1 "k8s.io/api/core/v1"
	"net/http"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func (s APIServer) GetResources(c echo.Context, which string) error {
	var listOptions client.MatchingLabels
	if which == "local" {
		listOptions = client.MatchingLabels{liqoconst.ReplicationRequestedLabel: "true"}
	} else if which == "remote" {
		listOptions = client.MatchingLabels{liqoconst.ReplicationRequestedLabel: "false"}
	} else {
		return errors.New("GetResources: invalid parameter")
	}

	var roList sharingv1alpha1.ResourceOfferList
	var jsonRoList []corev1.ResourceList
	ctx := context.TODO()
	err := s.List(ctx, &roList, &listOptions)
	if err != nil {
		return err
	}
	for _, ro := range roList.Items {
		jsonRoList = append(jsonRoList, (ro.Spec.ResourceQuota.Hard))
	}
	return c.JSON(http.StatusOK, jsonRoList)
}

func (s APIServer) GetResourcesLocal(c echo.Context) error {
	return s.GetResources(c, "local")
}

func (s APIServer) GetResourcesRemote(c echo.Context) error {
	return s.GetResources(c, "remote")
}