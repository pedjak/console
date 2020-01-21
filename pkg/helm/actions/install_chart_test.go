package actions

import (
	"github.com/stretchr/testify/assert"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chartutil"
	kubefake "helm.sh/helm/v3/pkg/kube/fake"
	"helm.sh/helm/v3/pkg/release"
	"helm.sh/helm/v3/pkg/storage"
	"helm.sh/helm/v3/pkg/storage/driver"
	"io/ioutil"
	"testing"
)

func TestInstallChart(t *testing.T) {
	chart := "../testdata/influxdb-3.0.2.tgz"
	store := storage.Init(driver.NewMemory())
	actionConfig := &action.Configuration{
		Releases:     store,
		KubeClient:   &kubefake.PrintingKubeClient{Out: ioutil.Discard},
		Capabilities: chartutil.DefaultCapabilities,
		Log:          func(format string, v ...interface{}) {},
	}
	rel, err := InstallChart("test-namespace", "test", chart, nil, actionConfig)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, rel.Name, "test", "Release name isn't matching")
	assert.Equal(t, rel.Namespace, "test-namespace", "Namespace name isn't matching")
	assert.Equal(t, rel.Info.Status, release.StatusDeployed, "Chart status is not deployed")
	assert.Equal(t, "influxdb", rel.Chart.Metadata.Name, "Chart name is not matching")
	assert.Equal(t, "3.0.2", rel.Chart.Metadata.Version, "Chart version is not matching")
}
