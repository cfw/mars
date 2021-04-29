package discovk8s

import (
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"path/filepath"
	"time"
)

type Client struct {
	ClientSet       *kubernetes.Clientset
	InformerFactory informers.SharedInformerFactory
	stop            chan struct{}
}

func NewClient() (*Client, error) {

	//config, err := getInClusterConfig()
	config, err := getConfig(getKubeConfig())
	if err != nil {
		return nil, err
	}

	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	factory := informers.NewSharedInformerFactory(clientSet, time.Hour*24)

	stopCh := make(chan struct{})

	return &Client{
		ClientSet:       clientSet,
		InformerFactory: factory,
		stop:            stopCh,
	}, nil
}

func getConfig(kubeConfig string) (*rest.Config, error) {
	var (
		config *rest.Config
		err    error
	)

	if kubeConfig != "" {
		// Use the current context in kubeConfig
		config, err = clientcmd.BuildConfigFromFlags("", kubeConfig)
		if err != nil {
			return nil, err
		}
	} else {
		// Creates the in-cluster config
		config, err = rest.InClusterConfig()
		if err != nil {
			return nil, err
		}
	}

	config.Timeout = 10 * time.Second
	return config, nil
}
func getKubeConfig() string {

	c := filepath.Join(os.Getenv("HOME"), ".kube", "config")

	_, err := os.Lstat(c)

	if err == nil || os.IsExist(err) {
		return c
	}
	return ""
}
