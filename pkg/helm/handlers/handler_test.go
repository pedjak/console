package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/openshift/console/pkg/auth"
	"github.com/openshift/console/pkg/helm/fake"
	"helm.sh/helm/v3/pkg/release"
)

var fakeReleaseList = []*release.Release{
	{
		Name: "Test",
	},
}

var fakeRelease = release.Release{
	Name: "Test",
}

var fakeReleaseManifest = "manifest-data"

func fakeHelmHandler() HelmHandlers {
	return HelmHandlers{
		getActionConfigurations: fake.GetFakeActionConfigurations,
	}
}

func TestHelmHandlers_HandleHelmList(t *testing.T) {
	tests := []struct {
		test        string
		expectedMsg string
		releaseList []*release.Release
		err         error
	}{
		{test: "invalid test", expectedMsg: "{\"error\":\"Failed to list helm releases: unknown error occurred\"}"},
		{test: "valid test", expectedMsg: "[{\"name\":\"Test\"}]"},
	}
	for _, tt := range tests {
		t.Run("List helm releases", func(t *testing.T) {
			handlers := fakeHelmHandler()
			if tt.test == "valid test" {
				handlers.listReleases = fake.FakeListReleases(fakeReleaseList, nil)
			} else if tt.test == "invalid test" {
				handlers.listReleases = fake.FakeListReleases(nil, errors.New("unknown error occurred"))
			}

			request := httptest.NewRequest("", "/foo", strings.NewReader("{}"))
			response := httptest.NewRecorder()

			handlers.HandleHelmList(&auth.User{}, response, request)
			switch tt.test {
			case "valid test":
				if response.Code != http.StatusOK {
					t.Error("Failed to install release")
				}
				if bytes.Compare(response.Body.Bytes(), []byte(tt.expectedMsg)) != 0 {
					t.Errorf("response body not matching expected is %s and received is %s", tt.expectedMsg, string(response.Body.Bytes()))
				}
			case "invalid test":
				if response.Code != http.StatusBadGateway {
					t.Error("response code should be 400")
				}
				if bytes.Compare(response.Body.Bytes(), []byte(tt.expectedMsg)) != 0 {
					t.Errorf("response body not matching expected is %s and received is %s", tt.expectedMsg, string(response.Body.Bytes()))
				}
			}
		})
	}
}

func TestHelmHandlers_HandleHelmInstall(t *testing.T) {
	tests := []struct {
		test        string
		expectedMsg string
	}{
		{test: "invalid test", expectedMsg: "{\"error\":\"Failed to install helm chart: Chart path is invalid\"}"},
		{test: "valid test", expectedMsg: "{\"name\":\"Test\"}"},
	}
	for _, tt := range tests {
		t.Run("Install Helm release", func(t *testing.T) {
			handlers := fakeHelmHandler()
			if tt.test == "valid test" {
				handlers.installChart = fake.FakeInstallChart(&fakeRelease, nil)
			} else if tt.test == "invalid test" {
				handlers.installChart = fake.FakeInstallChart(nil, errors.New("Chart path is invalid"))
			}

			req := HelmRequest{
				Name:      "test",
				Namespace: "test",
				ChartUrl:  "../testdata/influxdb-3.0.2.tgz",
				Values:    nil,
			}
			body, err := json.Marshal(req)
			if err != nil {
				t.Error("Failed to marshal request")
			}
			request := httptest.NewRequest("", "/foo", bytes.NewReader(body))
			response := httptest.NewRecorder()

			handlers.HandleHelmInstall(&auth.User{}, response, request)
			switch tt.test {
			case "valid test":
				if response.Code != http.StatusOK {
					t.Error("Failed to install release")
				}
				if bytes.Compare(response.Body.Bytes(), []byte(tt.expectedMsg)) != 0 {
					t.Errorf("response body not matching expected is %s and received is %s", tt.expectedMsg, string(response.Body.Bytes()))
				}
			case "invalid test":
				if response.Code != http.StatusBadGateway {
					t.Error("response code should be 400")
				}
				if bytes.Compare(response.Body.Bytes(), []byte(tt.expectedMsg)) != 0 {
					t.Errorf("response body not matching expected is %s and received is %s", tt.expectedMsg, string(response.Body.Bytes()))
				}
			}
		})
	}
}

func TestHelmHandlers_HandleHelmRenderManifest(t *testing.T) {
	tests := []struct {
		test        string
		expectedMsg string
	}{
		{test: "invalid test", expectedMsg: "{\"error\":\"Failed to render manifests: Chart path is invalid\"}"},
		{test: "valid test", expectedMsg: fakeReleaseManifest},
	}
	for _, tt := range tests {
		t.Run("Render Manifests", func(t *testing.T) {
			handlers := fakeHelmHandler()
			if tt.test == "valid test" {
				handlers.renderManifests = fake.FakeGetManifest(fakeReleaseManifest, nil)
			} else if tt.test == "invalid test" {
				handlers.renderManifests = fake.FakeGetManifest("", errors.New("Chart path is invalid"))
			}

			req := HelmRequest{
				Name:      "test",
				Namespace: "test",
				ChartUrl:  "../testdata/influxdb-3.0.2.tgz",
				Values:    nil,
			}
			body, err := json.Marshal(req)
			if err != nil {
				t.Error("Failed to marshal request")
			}
			request := httptest.NewRequest("", "/foo", bytes.NewReader(body))
			response := httptest.NewRecorder()

			handlers.HandleHelmRenderManifests(&auth.User{}, response, request)
			switch tt.test {
			case "valid test":
				if response.Code != http.StatusOK {
					t.Error("Failed to install release")
				}
				if bytes.Compare(response.Body.Bytes(), []byte(tt.expectedMsg)) != 0 {
					t.Errorf("response body not matching expected is %s and received is %s", tt.expectedMsg, string(response.Body.Bytes()))
				}
			case "invalid test":
				if response.Code != http.StatusBadGateway {
					t.Error("response code should be 502")
				}
				if bytes.Compare(response.Body.Bytes(), []byte(tt.expectedMsg)) != 0 {
					t.Errorf("response body not matching expected is %s and received is %s", tt.expectedMsg, string(response.Body.Bytes()))
				}
			}
		})
	}
}
