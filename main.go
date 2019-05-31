package main

import "os"

func main() {
	kc := selectKubeconfig()
	startNewShell(kc)

	// Avoid stacked shell sessions, when exit/ctrl+D caller shell is killed
	process, _ := os.FindProcess(os.Getppid())
	process.Kill()
}
