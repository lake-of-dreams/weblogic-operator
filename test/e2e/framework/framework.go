package framework

import (
	"flag"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	"weblogic-operator/pkg/clients/weblogic-operator"
)

// Global framework.
var Global *Framework

// Framework handles communication with the kube cluster in e2e tests.
type Framework struct {
	KubeClient    kubernetes.Interface
	MySQLOpClient weblogic_operator.Interface
	Namespace     string
}

// Setup sets up a test framework and initialises framework.Global.
func Setup() error {
	kubeconfig := flag.String("kubeconfig", "", "Path to kubeconfig file with authorization and master location information.")
	namespace := flag.String("namespace", "default", "e2e test namespace")
	flag.Parse()

	cfg, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		return err
	}

	kubeClient, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return err
	}

	mysqlopClient, err := weblogic_operator.NewForConfig(cfg)
	if err != nil {
		return err
	}
	Global = &Framework{
		KubeClient:    kubeClient,
		MySQLOpClient: mysqlopClient,
		Namespace:     *namespace,
	}

	return nil
}

// Teardown shuts down the test framework and cleans up.
func Teardown() error {
	// TODO: wait for all resources deleted.
	Global = nil
	return nil
}
