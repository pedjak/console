package actions

import (
	"net/http"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/cli"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/rest"
	"k8s.io/klog"
)

var settings = cli.New()
var k8sInClusterCA = "/var/run/secrets/kubernetes.io/serviceaccount/ca.crt"

type configFlagsWithTransport struct {
	*genericclioptions.ConfigFlags
	Transport *http.RoundTripper
}

func (c configFlagsWithTransport) ToRESTConfig() (*rest.Config, error) {
	return &rest.Config{
		Host:        *c.APIServer,
		BearerToken: *c.BearerToken,
		Transport:   *c.Transport,
	}, nil
}

func GetActionConfigurations(host, ns, token string, inCluster bool, transport *http.RoundTripper) *action.Configuration {
	confFlags := &configFlagsWithTransport{
		ConfigFlags: &genericclioptions.ConfigFlags{
			APIServer:   &host,
			BearerToken: &token,
			Namespace:   &ns,
		},
		Transport: transport,
	}
	if inCluster {
		confFlags.CAFile = &k8sInClusterCA
	} else {
		truePtr := true
		confFlags.Insecure = &truePtr
	}
	conf := new(action.Configuration)
	conf.Init(confFlags, ns, "secrets", klog.Infof)

	return conf
}
