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
	"weblogic-operator/pkg/server"
	"weblogic-operator/pkg/types"
	"weblogic-operator/pkg/domain"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"io"
)

// Operator operates things!
type Operator struct {
	Controllers []controllers.Controller
}

// NewWeblogicOperator instantiates a Weblogic Operator.
func NewWeblogicOperator(restConfig *rest.Config) (*Operator, error) {
	managedServerRESTClient, err := newManagedServerRESTClient(restConfig)
	domainRESTClient, err := newDomainRESTClient(restConfig)
	if err != nil {
		return nil, err
	}

	var clientSet kubernetes.Interface
	clientSet, err = kubernetes.NewForConfig(restConfig)
	if err != nil {
		return nil, err
	}

	serverController, err := server.NewController(clientSet, managedServerRESTClient, 30*time.Second, v1.NamespaceAll)
	domainController, err := domain.NewController(clientSet, domainRESTClient, 30*time.Second, v1.NamespaceAll)
	if err != nil {
		return nil, err
	}

	err = copyScripts()
	if err != nil {
		return nil, err
	}

	return NewWithControllers([]controllers.Controller{serverController, domainController}), nil
}

func NewPersistentVolume() v1.PersistentVolume {
	storageSize, err := resource.ParseQuantity("12Gi")

	persistentVolume := v1.PersistentVolume{
		ObjectMeta: metav1.ObjectMeta{
			Name: "k8s-weblogic-volume",
		},
		Spec: v1.PersistentVolumeSpec{
			AccessModes: []v1.PersistentVolumeAccessMode{
				v1.PersistentVolumeAccessMode("ReadWriteMany"),
			},
			Capacity: v1.ResourceList{
				v1.ResourceStorage: storageSize,
			},
			PersistentVolumeReclaimPolicy: v1.PersistentVolumeReclaimRetain,
		},
	}
	if err != nil {
		panic(err)
	}
	return persistentVolume
}

func NewPersistentVolumeClaim() v1.PersistentVolumeClaim {
	requestedStorageSize, err := resource.ParseQuantity("10Gi")
	if err != nil {
		panic(err)
	}
	persistenetVolumeClaim := v1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name: "weblogic-claim",
		},
		Spec: v1.PersistentVolumeClaimSpec{
			AccessModes: []v1.PersistentVolumeAccessMode{
				v1.PersistentVolumeAccessMode("ReadWriteMany"),
			},
			Resources: v1.ResourceRequirements{
				Requests: v1.ResourceList{
					v1.ResourceStorage: requestedStorageSize,
				},
			},
		},
	}

	return persistenetVolumeClaim
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

func newManagedServerRESTClient(config *rest.Config) (*rest.RESTClient, error) {
	//if err := types.AddToScheme(scheme.Scheme); err != nil {
	//	return nil, err
	//}
	config.GroupVersion = &types.WeblogicManagedServerSchemeGroupVersion
	config.APIPath = "/apis"
	config.ContentType = runtime.ContentTypeJSON
	config.NegotiatedSerializer = serializer.DirectCodecFactory{CodecFactory: serializer.NewCodecFactory(scheme.Scheme)}

	return rest.RESTClientFor(config)
}

func newDomainRESTClient(config *rest.Config) (*rest.RESTClient, error) {
	//if err := types.AddToScheme(scheme.Scheme); err != nil {
	//	return nil, err
	//}
	config.GroupVersion = &types.WebLogicDomainSchemeGroupVersion
	config.APIPath = "/apis"
	config.ContentType = runtime.ContentTypeJSON
	config.NegotiatedSerializer = serializer.DirectCodecFactory{CodecFactory: serializer.NewCodecFactory(scheme.Scheme)}

	return rest.RESTClientFor(config)
}

func copyScripts() error {
	glog.Infof("Copying scripts to %s...", "/u01/oracle/user_projects")

	source := "/scripts"
	dest := "/u01/oracle/user_projects"

	directory, _ := os.Open(source)
	objects, _ := directory.Readdir(-1)

	for _, obj := range objects {
		sourcefilepointer := source + "/" + obj.Name()
		destinationfilepointer := dest + "/" + obj.Name()

		sourcefile, err := os.Open(sourcefilepointer)
		if err != nil {
			return err
		}

		defer sourcefile.Close()

		destfile, err := os.Create(destinationfilepointer)
		if err != nil {
			return err
		}

		defer destfile.Close()

		_, err = io.Copy(destfile, sourcefile)
		if err == nil {
			_, err := os.Stat(source)
			if err != nil {
				err = os.Chmod(dest, os.FileMode(777))
			}

		}
	}
	return nil
}
