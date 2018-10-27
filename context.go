package main

import (
	"io/ioutil"
	"path"

	"github.com/ghodss/yaml"
	homedir "github.com/mitchellh/go-homedir"
)

// RetrieveAllContexts reads the k8s configuration file and return all the
// available contexts
func RetrieveAllContexts() ([]string, error) {
	type kubeConfig struct {
		Contexts []struct {
			Name string
		}
	}

	configPath, err := getDefaultConfigPath()
	if err != nil {
		return nil, err
	}

	rawCfgFile, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var cfg kubeConfig
	err = yaml.Unmarshal(rawCfgFile, &cfg)
	if err != nil {
		return nil, err
	}

	contextNames := make([]string, len(cfg.Contexts))

	for idx, ctx := range cfg.Contexts {
		contextNames[idx] = ctx.Name
	}

	return contextNames, nil
}

func getDefaultConfigPath() (string, error) {
	home, err := homedir.Dir()
	if err != nil {
		return "", err
	}

	configPath := path.Join(home, ".kube", "config")

	return configPath, nil
}
