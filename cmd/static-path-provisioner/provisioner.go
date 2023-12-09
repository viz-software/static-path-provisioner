package main

import (
	"context"
	"path/filepath"

	v1 "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog/v2"
	"sigs.k8s.io/sig-storage-lib-external-provisioner/v8/controller"
)

// validate that `StaticPathProvisioner` has correctly implemented the `controller.Provisioner` interface
var _ controller.Provisioner = (*StaticPathProvisioner)(nil)

// StaticPathProvisioner stores config and states for a static path provisioner
type StaticPathProvisioner struct{}

// NewProvisioner creates a static path provisioner that implements "Provision" and "Delete" actions
func NewProvisioner() (*StaticPathProvisioner, error) {
	provisioner := &StaticPathProvisioner{}

	return provisioner, nil
}

// Provision persistent volume implementation
func (*StaticPathProvisioner) Provision(context context.Context, opts controller.ProvisionOptions) (*v1.PersistentVolume, controller.ProvisioningState, error) {
	// collect options
	storageClass := opts.StorageClass
	pvc := opts.PVC
	pvName := opts.PVName

	// build the host path for the directory to store volume content, use `DEFAULT_STORAGE_PATH` if the `storagePath` parameter is not provided for the StorageClass
	// note unlike local path provisioner, the host path is static or deterministic (`{storagePath}/{namespaceName}/{pvcName}`, no random UUID in `pvName`)
	storagePath := storageClass.Parameters["storagePath"]
	if storagePath == "" {
		storagePath = DEFAULT_STORAGE_PATH
	}
	hostPath := filepath.Join(storagePath, pvc.Namespace, pvc.Name)

	klog.Infof("Creating volume %v for claim %v/%v at %v...", pvName, pvc.Namespace, pvc.Name, hostPath)

	// return the persistent volume to be created by the provisioner controller
	filesystemVolumeMode := v1.PersistentVolumeFilesystem
	hostPathType := v1.HostPathDirectoryOrCreate
	storageKey := v1.ResourceName(v1.ResourceStorage)

	return &v1.PersistentVolume{
		ObjectMeta: meta.ObjectMeta{
			Name: pvName,
		},
		Spec: v1.PersistentVolumeSpec{
			PersistentVolumeReclaimPolicy: *storageClass.ReclaimPolicy,
			AccessModes:                   pvc.Spec.AccessModes,
			VolumeMode:                    &filesystemVolumeMode,
			Capacity: v1.ResourceList{
				storageKey: pvc.Spec.Resources.Requests[storageKey],
			},
			PersistentVolumeSource: v1.PersistentVolumeSource{
				HostPath: &v1.HostPathVolumeSource{
					Path: hostPath,
					Type: &hostPathType,
				},
			},
		},
	}, controller.ProvisioningFinished, nil
}

// Delete persistent volume implementation
func (*StaticPathProvisioner) Delete(context context.Context, pv *v1.PersistentVolume) error {
	klog.Infof("Deleting volume %v...", pv.Name)

	// we are not going to delete content currently
	klog.Warningf("Not going to delete the volume directory or content in it for volume %v...", pv.Name)

	return nil
}
