package operator

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/golang/glog"

	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"weblogic-operator/pkg/controllers"
	"weblogic-operator/pkg/types"
)

// Operator operates things!
type Operator struct {
	Controllers []controllers.Controller
}

// NewWeblogicOperator instantiates a Weblogic Operator.
func NewWeblogicOperator(restConfig *rest.Config) (*Operator, error) {
	restClient, err := newRESTClient(restConfig)
	if err != nil {
		return nil, err
	}

	var clientSet kubernetes.Interface
	clientSet, err = kubernetes.NewForConfig(restConfig)
	if err != nil {
		return nil, err
	}

	serverController, err := controllers.NewController(clientSet, restClient, 30*time.Second, v1.NamespaceAll)
	if err != nil {
		return nil, err
	}

	return NewWithControllers([]controllers.Controller{serverController}), nil
}

// NewWithControllers creates an new operator for the given controllers.
func NewWithControllers(controllers []controllers.Controller) *Operator {
	return &Operator{Controllers: controllers}
}

// Run runs the operator until SIGINT or SIGTERM signal is received.
func (o *Operator) Run() {
	// Multiple signals will get dropped
	signalChan := make(chan os.Signal, 1)
	stopChan := make(chan struct{})
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	for _, controller := range o.Controllers {
		go controller.Run(stopChan)
	}
	select {
	case signal := <-signalChan:
		glog.Infof("Received %s, shutting down...", signal.String())
		close(stopChan)
	}
}

func newRESTClient(config *rest.Config) (*rest.RESTClient, error) {
	config.GroupVersion = &types.SchemeGroupVersion
	config.APIPath = "/apis"
	config.ContentType = runtime.ContentTypeJSON
	config.NegotiatedSerializer = serializer.DirectCodecFactory{CodecFactory: serializer.NewCodecFactory(scheme.Scheme)}

	return rest.RESTClientFor(config)
}
