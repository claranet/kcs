package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/manifoldco/promptui"
)

func startNewShell(kc string) {
	// Get the current working directory.
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	// Set envvar with path to kubeconfig file
	os.Setenv("KUBECONFIG", kc)

	// Transfer stdin, stdout, and stderr to the new process
	// and also set target directory for the shell to start in.
	pa := os.ProcAttr{
		Files: []*os.File{os.Stdin, os.Stdout, os.Stderr},
		Dir:   cwd,
	}

	// Start up a new shell.
	proc, err := os.StartProcess(os.Getenv("SHELL"), []string{os.Getenv("SHELL")}, &pa)
	if err != nil {
		panic(err)
	}

	// Wait until user exits the shell
	_, err = proc.Wait()
	if err != nil {
		panic(err)
	}
}

type kubeconfig struct {
	Name    string
	Path    string
	Context string
}

func selectKubeconfig() string {
	kcpath := getConfig()

	files, err := ioutil.ReadDir(kcpath)
	if err != nil {
		log.Fatal(err)
	}

	kubeconfigList := []kubeconfig{}

	for _, f := range files {
		kubeconfigList = append(kubeconfigList, kubeconfig{
			f.Name(),
			kcpath,
			getKubeconfigContexts(kcpath + "/" + f.Name()),
		})
	}

	templates := &promptui.SelectTemplates{
		Label:    "{{ . | bold }}",
		Active:   "> {{ .Name | cyan | bold }}",
		Inactive: "  {{ .Name }}",
		Details: `
{{ "Contexts: " }}{{ .Context }}`,
	}

	searcher := func(input string, index int) bool {
		kc := kubeconfigList[index]
		name := strings.Replace(strings.ToLower(kc.Name), " ", "", -1)
		input = strings.Replace(strings.ToLower(input), " ", "", -1)

		return strings.Contains(name, input)
	}

	prompt := promptui.Select{
		Label:             "Current Config : " + getCurrentKubeConfig(),
		Items:             kubeconfigList,
		Templates:         templates,
		Size:              10,
		Searcher:          searcher,
		StartInSearchMode: true,
		HideSelected:      true,
		// HideHelp:          true,
	}

	i, _, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return ""
	}

	return kubeconfigList[i].Path + "/" + kubeconfigList[i].Name
}
