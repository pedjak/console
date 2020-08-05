package chartproxy

import (
	"crypto/tls"
	"helm.sh/helm/v3/pkg/repo"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"sigs.k8s.io/yaml"
)

type helmRepo struct {
	Name    string
	Url     *url.URL
	TlsClientConfig *tls.Config
}

type TLSConfigGetter interface {
	Get() (*tls.Config, error)
}

func (repo helmRepo) Get() (*tls.Config, error) {
	return repo.TlsClientConfig, nil
}

func (repo helmRepo) HttpClient() (*http.Client, error) {
	tlsConfig, err := repo.Get()
	if err != nil {
		return nil, err
	}
	tr := &http.Transport{
		TLSClientConfig: tlsConfig,
	}

	client := &http.Client{Transport: tr}
	return client, nil
}

func (hr helmRepo) IndexFile() (*repo.IndexFile, error) {
	var indexFile repo.IndexFile
	httpClient, err := hr.HttpClient()
	if err != nil {
		return nil, err
	}
	indexUrl := hr.Url.String()
	if !strings.HasSuffix(indexUrl, "/index.yaml") {
		indexUrl += "/index.yaml"
	}
	resp, err := httpClient.Get(indexUrl)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(body, &indexFile)
	if err != nil {
		return nil, err
	}
	return &indexFile, nil
}

