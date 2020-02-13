package fake

import (
	"net/http"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/release"
)

func FakeInstallChart(mockedRelease *release.Release, err error) func(ns string, name string, url string, values map[string]interface{}, conf *action.Configuration) (*release.Release, error) {
	return func(ns string, name string, url string, values map[string]interface{}, conf *action.Configuration) (r *release.Release, er error) {
		return mockedRelease, err
	}
}

func FakeListReleases(mockedReleases []*release.Release, err error) func(conf *action.Configuration) ([]*release.Release, error) {
	return func(conf *action.Configuration) (releases []*release.Release, er error) {
		return mockedReleases, err
	}
}

func FakeGetManifest(mockedManifest string, err error) func(name string, url string, values map[string]interface{}, conf *action.Configuration) (string, error) {
	return func(name string, url string, values map[string]interface{}, conf *action.Configuration) (r string, er error) {
		return mockedManifest, err
	}
}

func GetFakeActionConfigurations(string, string, string, *http.RoundTripper) *action.Configuration {
	return &action.Configuration{}
}
