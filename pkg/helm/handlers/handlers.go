package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/coreos/pkg/capnslog"

	"github.com/openshift/console/pkg/auth"
	"github.com/openshift/console/pkg/helm/actions"
)

var (
	plog = capnslog.NewPackageLogger("github.com/openshift/console", "helm")
)

type request struct {
	Name      string                 `json:"name"`
	Namespace string                 `json:"namespace"`
	ChartUrl  string                 `json:"chart_url"`
	Values    map[string]interface{} `json:"values"`
}

// HelmHandlers provides handlers to handle helm related requests
type HelmHandlers struct {
	kubeApiURL string
	Transport  *http.RoundTripper
}

func Handlers(kubeApiURL string, transport *http.RoundTripper) *[]auth.UserAuthEndPointHandler {
	h := HelmHandlers{
		kubeApiURL: kubeApiURL,
		Transport:  transport,
	}
	return &[]auth.UserAuthEndPointHandler{
		{
			Path: "/api/helm/template",
			Handler: func(user *auth.User, w http.ResponseWriter, r *http.Request) {
				renderManifests(&h, user, w, r)
			},
		},
		{
			Path: "/api/helm/release",
			Handler: func(user *auth.User, w http.ResponseWriter, r *http.Request) {
				installChart(&h, user, w, r)
			},
		},
		{
			Path: "/api/helm/releases",
			Handler: func(user *auth.User, w http.ResponseWriter, r *http.Request) {
				listReleases(&h, user, w, r)
			},
		},
	}
}

func renderManifests(h *HelmHandlers, user *auth.User, w http.ResponseWriter, r *http.Request) {
	var req request

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	conf := actions.GetActionConfigurations(h.kubeApiURL, req.Namespace, user.Token, h.Transport)
	resp, err := actions.RenderManifests(req.Name, req.ChartUrl, req.Values, conf)
	if err != nil {
		sendResponse(w, http.StatusBadGateway, apiError{fmt.Sprintf("Failed to render manifests: %v", err)})
	}

	w.Header().Set("Content-Type", "text/yaml")
	w.Write([]byte(resp))
}

func installChart(h *HelmHandlers, user *auth.User, w http.ResponseWriter, r *http.Request) {
	var req request

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	conf := actions.GetActionConfigurations(h.kubeApiURL, req.Namespace, user.Token, h.Transport)
	resp, err := actions.InstallChart(req.Namespace, req.Name, req.ChartUrl, req.Values, conf)
	if err != nil {
		sendResponse(w, http.StatusBadGateway, apiError{fmt.Sprintf("Failed to install helm chart: %v", err)})
	}

	w.Header().Set("Content-Type", "application/json")
	res, _ := json.Marshal(resp)
	w.Write(res)
}

func listReleases(h *HelmHandlers, user *auth.User, w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	ns := params.Get("ns")

	conf := actions.GetActionConfigurations(h.kubeApiURL, ns, user.Token, h.Transport)
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
