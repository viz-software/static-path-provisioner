package main

import (
	"github.com/pkg/errors"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/sig-storage-lib-external-provisioner/v8/controller"
)

// NewController creates a provisioner controller
func NewController(clientset *kubernetes.Clientset) (*controller.ProvisionController, error) {
	var err error

	// construct the static path provisioner
	provisioner, err := NewProvisioner()
	if err != nil {
		return nil, errors.Wrap(err, "unable to create provisioner")
	}

	// create provisioner controller with name "viz.software/static-path-provisioner", and use the static path provisioner to handle actions
	opts := []func(*controller.ProvisionController) error{
		controller.LeaderElection(false),
		controller.Threadiness(1),
	}
	provisionController := controller.NewProvisionController(clientset, "viz.software/static-path-provisioner", provisioner, opts...)

	return provisionController, nil
}
