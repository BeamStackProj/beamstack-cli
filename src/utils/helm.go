package utils

import (
	"fmt"
	"os"

	"github.com/BeamStackProj/beamstack-cli/src/types"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/getter"
	"helm.sh/helm/v3/pkg/repo"
)

func InstallHelmPackage(name string, index string, url string, version string, namespace string, values *map[string]interface{}) (helmPackage types.Package) {
	// Set up Helm action configuration
	helmPackage = types.Package{
		Name:    name,
		Url:     url,
		Version: version,
		Type:    "helm",
	}

	if namespace == "" {
		namespace = "default"
	}

	actionConfig := new(action.Configuration)
	if err := actionConfig.Init(cli.New().RESTClientGetter(), namespace, os.Getenv("HELM_DRIVER"), func(format string, v ...interface{}) {
		fmt.Printf(format, v...)
	}); err != nil {
		panic(err.Error())
	}

	// Add the Flink Operator Helm repository
	repoEntry := repo.Entry{
		Name: name,
		URL:  url,
	}
	chartRepo, err := repo.NewChartRepository(&repoEntry, getter.All(&cli.EnvSettings{}))
	if err != nil {
		panic(err.Error())
	}
	if _, err := chartRepo.DownloadIndexFile(); err != nil {
		panic(err.Error())
	}

	// Update the repositories
	settings := cli.New()
	chartRepoList := repo.File{Repositories: []*repo.Entry{&repoEntry}}
	if err := chartRepoList.WriteFile(settings.RepositoryConfig, 0644); err != nil {
		panic(err.Error())
	}

	err = updateHelmRepositories()
	if err != nil {
		fmt.Println("error updating helm repo")
		fmt.Println(err)
	}
	install := action.NewInstall(actionConfig)
	install.ReleaseName = name
	install.Namespace = namespace
	if index == "" {
		index = name
	}
	chartPath, err := install.LocateChart(name+"/"+index, settings)
	if err != nil {
		panic(err.Error())
	}

	chart, err := loader.Load(chartPath)
	if err != nil {
		panic(err.Error())
	}

	for _, crd := range chart.CRDObjects() {
		helmPackage.Dependencies = append(helmPackage.Dependencies, &types.Package{
			Name:    crd.Name,
			Url:     crd.Filename,
			Type:    "k8s.crd",
			Version: crd.File.Name,
		})
	}
	delete := action.NewUninstall(actionConfig)
	_, _ = delete.Run(name)
	_, err = install.Run(chart, *values)
	if err != nil {
		panic(err.Error())
	}
	return
}

func updateHelmRepositories() error {
	settings := cli.New()

	// Load the existing repositories file
	repoFile := settings.RepositoryConfig
	file, err := repo.LoadFile(repoFile)
	if err != nil {
		return fmt.Errorf("failed to load repository file: %w", err)
	}

	// If no repositories are found, return early
	if len(file.Repositories) == 0 {
		fmt.Println("No repositories found in the Helm configuration.")
		return nil
	}
	// Iterate over each repository entry and update the index
	for _, repoEntry := range file.Repositories {
		chartRepo, err := repo.NewChartRepository(repoEntry, getter.All(settings))
		if err != nil {
			return fmt.Errorf("failed to create chart repository for %s: %w", repoEntry.Name, err)
		}

		// Download the latest index file
		_, err = chartRepo.DownloadIndexFile()
		if err != nil {
			return fmt.Errorf("failed to update repository %s: %w", repoEntry.Name, err)
		}
	}

	// Write the updated repositories back to the file
	if err := file.WriteFile(repoFile, 0644); err != nil {
		return fmt.Errorf("failed to write updated repository file: %w", err)
	}

	return nil
}
