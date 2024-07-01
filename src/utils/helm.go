package utils

import (
	"fmt"
	"os"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/getter"
	"helm.sh/helm/v3/pkg/repo"
)

func InstallHelmPackage(name string, url string) {

	// Set up Helm action configuration
	actionConfig := new(action.Configuration)
	if err := actionConfig.Init(cli.New().RESTClientGetter(), "default", os.Getenv("HELM_DRIVER"), func(format string, v ...interface{}) {
		fmt.Sprintf(format, v...)
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

	delete := action.NewUninstall(actionConfig)
	deletedrelease, _ := delete.Run(name)
	if deletedrelease != nil {
		fmt.Printf("removed previous installed release %s \n", deletedrelease.Release.Name)
	}

	_, err = install.Run(chart, map[string]interface{}{})
	if err != nil {
		panic(err.Error())
	}

	fmt.Printf("%s Operator installed successfully\n", name)
}
