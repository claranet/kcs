package main

import (
	"io/ioutil"
	"log"
	"os"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

func getConfig() string {
	home := os.Getenv("HOME")
	kcp, _ := ioutil.ReadFile(home + "/.kcs")
	kcs := strings.Split(string(kcp), "=")
	if kcs[0] == "KUBECONFIG_FILES" && len(kcs[1]) != 0 {
		return strings.TrimSuffix(kcs[1], "\n")
	} else {
		log.Fatal("wrong config")
	}
	return ""
}

type kubeconfigyaml struct {
	Contexts []context `yaml:"contexts"`
}

type context struct {
	Name string `yaml:"name"`
}

func getKubeconfigContexts(kcf string) string {
	var c kubeconfigyaml

	yamlFile, err := ioutil.ReadFile(kcf)
	if err != nil {
		return ""
	}
	err = yaml.Unmarshal(yamlFile, &c)
	if err != nil {
		return ""
	}
	var v string
	for _, cx := range c.Contexts {
		v = v + cx.Name + " | "
	}
	return v[:len(v)-3]
}

func getCurrentKubeConfig() string {
	if os.Getenv("KUBECONFIG") == "" {
		return "None"
	}
	k := strings.Split(os.Getenv("KUBECONFIG"), "/")
	return k[len(k)-1]
}
