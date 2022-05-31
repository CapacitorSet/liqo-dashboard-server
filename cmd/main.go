package main

import (
	"flag"
	"os"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/liqotech/liqo/pkg/utils/apiserver"
	// "github.com/liqotech/liqo/pkg/utils/args"
	discoveryv1alpha1 "github.com/liqotech/liqo/apis/discovery/v1alpha1"
	sharingv1alpha1 "github.com/liqotech/liqo/apis/sharing/v1alpha1"
	"github.com/liqotech/liqo/pkg/utils/restcfg"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"

	"github.com/CapacitorSet/liqo-dashboard-server/api"
)

var (
	scheme = runtime.NewScheme()
)

func init() {
	_ = clientgoscheme.AddToScheme(scheme)

	_ = discoveryv1alpha1.AddToScheme(scheme)
	_ = sharingv1alpha1.AddToScheme(scheme)
	// +kubebuilder:scaffold:scheme
}

func main() {
	// klog.Info("Starting")

	address := flag.String("address", ":8000", "The address the service binds to")

	// Configure the flags concerning the exposed API server connection parameters.
	apiserver.InitFlags(nil)

	restcfg.InitFlags(nil)
	flag.Parse()

	// klog.Info("Namespace: ", *namespace)

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:                 scheme,
		MetricsBindAddress:     "",
		HealthProbeBindAddress: "",
	})
	if err != nil {
		// klog.Error(err)
		os.Exit(1)
	}

	k8sClient := mgr.GetClient()
	_ = k8sClient

	ctx := ctrl.SetupSignalHandler()

	apiServer := api.NewAPIServer(k8sClient, *address)
	mgr.Add(apiServer)

	if err := mgr.Start(ctx); err != nil {
		panic(err)
	}
}
