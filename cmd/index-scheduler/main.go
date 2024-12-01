package main

import (
	"os"

	"k8s.io/component-base/cli"
	"k8s.io/kubernetes/cmd/kube-scheduler/app"

	"github.com/9506hqwy/k8s-scheduler-extension/pkg/indexscheduling"
)

func main() {
	command := app.NewSchedulerCommand(
		app.WithPlugin(indexscheduling.Name, indexscheduling.New),
	)

	code := cli.Run(command)
	os.Exit(code)
}
