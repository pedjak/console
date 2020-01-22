package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/coreos/pkg/capnslog"

	"github.com/openshift/console/pkg/helm/actions"
)

var (
	plog = capnslog.NewPackageLogger("github.com/openshift/console", "helm")
)

// HelmHandlers provides handlers to handle helm related requests
type HelmHandlers struct {
	ApiServerHost string
	Transport http.RoundTripper
	UserToken string
}

func(h *HelmHandlers) HandleHelmRenderManifests(w http.ResponseWriter, r *http.Request) {
	var req HelmRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	conf := actions.GetActionConfigurations(h.ApiServerHost, req.Namespace, h.UserToken, &h.Transport)
	resp, err := actions.RenderManifests(req.Name, req.ChartUrl, req.Values, conf)
	if err != nil {
		sendResponse(w, http.StatusBadGateway, apiError{fmt.Sprintf("Failed to render manifests: %v", err)})
	}

	w.Header().Set("Content-Type", "text/yaml")
	w.Write([]byte(resp))
}

func (h *HelmHandlers) HandleHelmInstall(w http.ResponseWriter, r *http.Request) {
	var req HelmRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	conf := actions.GetActionConfigurations(h.ApiServerHost, req.Namespace, h.UserToken, &h.Transport)
	resp, err := actions.InstallChart(req.Namespace, req.Name, req.ChartUrl, req.Values, conf)
	if err != nil {
		sendResponse(w, http.StatusBadGateway, apiError{fmt.Sprintf("Failed to install helm chart: %v", err)})
	}

	w.Header().Set("Content-Type", "application/json")
	res, _ := json.Marshal(resp)
	w.Write(res)
}

func (h *HelmHandlers) HandleHelmList(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	ns := params.Get("ns")

	conf := actions.GetActionConfigurations(h.ApiServerHost, ns, h.UserToken, &h.Transport)
	resp, err := actions.ListReleases(conf)
	if err != nil {
		sendResponse(w, http.StatusBadGateway, apiError{fmt.Sprintf("Failed to list helm releases: %v", err)})
	}

	w.Header().Set("Content-Type", "application/json")

	res, _ := json.Marshal(resp)
	w.Write(res)
}

// Copied from Server package to maintain error response consistency
func sendResponse(rw http.ResponseWriter, code int, resp interface{}) {
	enc, err := json.Marshal(resp)
	if err != nil {
		plog.Printf("Failed JSON-encoding HTTP response: %v", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(code)

	_, err = rw.Write(enc)
	if err != nil {
		plog.Errorf("Failed sending HTTP response body: %v", err)
	}
}

type apiError struct {
	Err string `json:"error"`
}
