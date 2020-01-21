package actions

import (
	"encoding/json"
	"fmt"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chartutil"
	kubefake "helm.sh/helm/v3/pkg/kube/fake"
	"helm.sh/helm/v3/pkg/storage"
	"helm.sh/helm/v3/pkg/storage/driver"
	"io/ioutil"
	"os"

	"testing"
)

func Test(t *testing.T) {
	store := storage.Init(driver.NewMemory())
	actionConfig := &action.Configuration{
		Releases:     store,
		KubeClient:   &kubefake.PrintingKubeClient{Out: ioutil.Discard},
		Capabilities: chartutil.DefaultCapabilities,
		Log:          func(format string, v ...interface{}) {},
	}

	chart, err := GetChart("https://kubernetes-charts.storage.googleapis.com/mariadb-7.3.5.tgz", actionConfig)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	val, _ := json.MarshalIndent(chart, "", "  ")
	fmt.Fprintln(os.Stdout, string(val))
}
