package main

import (
	"k8s.io/client-go/tools/clientcmd"

	"github.com/golang/glog"
	"github.com/spf13/pflag"

	"github.com/sczachariah/weblogic-operator/pkg/operator"
	"github.com/sczachariah/weblogic-operator/pkg/util/flags"
	"github.com/sczachariah/weblogic-operator/pkg/util/logs"
)

func main() {
	var kubeConfigFile = pflag.String("kubeconfig", "", "Path to kubeconfig file with authorization and master location information.")

	flags.InitFlags()
	logs.InitLogs()
	defer logs.FlushLogs()

	glog.V(2).Info("Starting Weblogic operator")

	cfg, err := clientcmd.BuildConfigFromFlags("", *kubeConfigFile)
	if err != nil {
		glog.Errorf("Failed to build REST config: %s", err)
		panic(err.Error())
	}

	operator, err := operator.NewWeblogicOperator(cfg)
	if err != nil {
		glog.Errorf("Failed to initialize the operator: %s", err)
		panic(err.Error())
	}

	operator.Run()
}
