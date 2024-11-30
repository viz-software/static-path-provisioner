package main

import (
	"context"

	"github.com/pkg/errors"
	"k8s.io/klog/v2"
)

func main() {
	var err error

	ctx := context.Background()
	logger := klog.FromContext(ctx)

	klog.Infof("%s %s", APPNAME, VERSION)

	// construct the clientset to interact with Kubernetes APIs
	clientset, err := NewClientset()
	if err != nil {
		klog.Fatalf("Critical error: %v", errors.Wrap(err, "unable to create clientset"))
	}

	// construct the provisioner controller
	controller, err := NewController(logger, clientset)
	if err != nil {
		klog.Fatalf("Critical error: %v", errors.Wrap(err, "unable to create controller"))
	}

	// start the provisioner controller
	klog.Infof("%s is ready.", APPNAME)
	controller.Run(context.Background())
}
