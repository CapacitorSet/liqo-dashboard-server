package api

import (
	"github.com/labstack/echo/v4"
	discoveryv1alpha1 "github.com/liqotech/liqo/apis/discovery/v1alpha1"
	sharingv1alpha1 "github.com/liqotech/liqo/apis/sharing/v1alpha1"
	oapi_client "github.com/CapacitorSet/liqo-dashboard-server/client"
	"k8s.io/client-go/kubernetes/scheme"
	"net/http"
	"net/http/httptest"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

var (
	_         = discoveryv1alpha1.AddToScheme(scheme.Scheme)
	_         = sharingv1alpha1.AddToScheme(scheme.Scheme)
	k8sClient = fake.NewClientBuilder().WithScheme(scheme.Scheme).Build()

	echoSrv = echo.New()
	srv     = APIServer{
		Client:        k8sClient,
		EchoServer:    echoSrv,
		ListenAddress: "",
	}
)

// LocalHTTPClient sends HTTP requests directly to the echo.Server.
type LocalHTTPClient struct {
	APIMethod func(echo.Context) error
}

func (l LocalHTTPClient) Do(req *http.Request) (*http.Response, error) {
	rec := httptest.NewRecorder()
	c := echoSrv.NewContext(req, rec)
	err := l.APIMethod(c)
	if err != nil {
		return nil, err
	}
	return rec.Result(), nil
}

// NewClientForMethod returns an OpenAPI client that mocks a specific method
func NewClientForMethod(method func(echo.Context) error) *oapi_client.ClientWithResponses {
	return &oapi_client.ClientWithResponses{ClientInterface: &oapi_client.Client{
		Server:         "http://unused/",
		Client:         LocalHTTPClient{APIMethod: method},
		RequestEditors: nil,
	}}
}
