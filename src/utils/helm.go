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

func InstallHelmPackage(name string, url string, version string) (helmPackage types.Package) {
	// Set up Helm action configuration
	helmPackage = types.Package{
		Name:    name,
		Url:     url,
		Version: version,
		Type:    "helm",
	}

	actionConfig := new(action.Configuration)
	if err := actionConfig.Init(cli.New().RESTClientGetter(), "default", os.Getenv("HELM_DRIVER"), func(format string, v ...interface{}) {
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

	// Install the Operator
	install := action.NewInstall(actionConfig)
	install.ReleaseName = name
	install.Namespace = "default"

	chartPath, err := install.LocateChart(name+"/"+name, settings)
	if err != nil {
		panic(err.Error())
	}

	chart, err := loader.Load(chartPath)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(len(chart.CRDObjects()))
	for _, crd := range chart.CRDObjects() {
		helmPackage.Dependencies = append(helmPackage.Dependencies, &types.Package{
			Name:    crd.Name,
			Url:     crd.Filename,
			Type:    "k8s.crd",
			Version: crd.File.Name,
		})
		fmt.Println(crd.Filename)
		fmt.Println(crd.Name)
	}
	delete := action.NewUninstall(actionConfig)
	deletedrelease, _ := delete.Run(name)
	if deletedrelease != nil {
		fmt.Printf("removed previous installed release %s \n", deletedrelease.Release.Name)
	}

	_, err = install.Run(chart, map[string]interface{}{})
	if err != nil {
		panic(err.Error())
	}
	return
}
