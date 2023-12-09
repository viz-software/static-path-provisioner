package main

import (
	"github.com/pkg/errors"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// NewClientset creates clientset to interact with Kubernetes APIs
func NewClientset() (*kubernetes.Clientset, error) {
	var err error

	// We only support in-cluster config, make sure you have set a service account with sufficient permission for the pod running the provisioner
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, errors.Wrap(err, "unable to get in-cluster config")
	}

	// construct clientset with the in-cluster config
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return clientset, nil
}
