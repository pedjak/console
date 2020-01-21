package actions

import (
	"github.com/stretchr/testify/assert"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chartutil"
	kubefake "helm.sh/helm/v3/pkg/kube/fake"
	"helm.sh/helm/v3/pkg/release"
	"helm.sh/helm/v3/pkg/storage"
	"helm.sh/helm/v3/pkg/storage/driver"
	"helm.sh/helm/v3/pkg/time"
	"io/ioutil"
	"testing"
)

func TestListReleases(t *testing.T) {
	store := storage.Init(driver.NewMemory())
	err := store.Create(&release.Release{
		Name:      "test",
		Namespace: "test-namespace",
		Info: &release.Info{
			FirstDeployed: time.Time{},
			Status:        "deployed",
		},
		Chart: &chart.Chart{
			Metadata: &chart.Metadata{
				Name:    "influxdb",
				Version: "3.0.2",
			},
		},
	})
	actionConfig := &action.Configuration{
		Releases:     store,
		KubeClient:   &kubefake.PrintingKubeClient{Out: ioutil.Discard},
		Capabilities: chartutil.DefaultCapabilities,
		Log:          func(format string, v ...interface{}) {},
	}
	rels, err := ListReleases(actionConfig)
	if err != nil {
		t.Error(err.Error())
	}
	assert.Len(t, rels, 1, "Release list should return 1 release")
	assert.Equal(t, rels[0].Name, "test", "Release name isn't matching")
	assert.Equal(t, rels[0].Namespace, "test-namespace", "Namespace name isn't matching")
	assert.Equal(t, rels[0].Info.Status, release.StatusDeployed, "Chart status is not deployed")
	assert.Equal(t, "influxdb", rels[0].Chart.Metadata.Name, "Chart name is not matching")
	assert.Equal(t, "3.0.2", rels[0].Chart.Metadata.Version, "Chart version is not matching")
}
