package uploader

import (
	"github.com/aquasecurity/trivy-plugin-aqua/pkg/buildClient"
	"github.com/aquasecurity/trivy-plugin-aqua/pkg/log"
	"github.com/aquasecurity/trivy-plugin-aqua/pkg/proto/buildsecurity"
	"github.com/argonsecurity/go-environments/models"
)

// Upload forwards the results to the configured client
func Upload(client buildClient.Client,
	results []*buildsecurity.Result,
	tags map[string]string,
	avdUrlMap buildClient.ResultIdToUrlMap,
	pipelines []*buildsecurity.Pipeline,
	dependencies map[string]*buildsecurity.PackageDependencies,
	envConfig *models.Configuration) error {
	log.Logger.Debugf("Uploading scan with tags. %v", tags)
	return client.Upload(results, tags, avdUrlMap, pipelines, dependencies, envConfig)
}
