# Static Path Provisioner

[![Publish Container Image](https://github.com/viz-software/static-path-provisioner/actions/workflows/publish-container-image.yaml/badge.svg?branch=main)](https://github.com/viz-software/static-path-provisioner/actions/workflows/publish-container-image.yaml)

## Overview

Static Path Provisioner simply creates `hostPath`-based persistent volumes based on persistent volumes claims.

## Motivation

I want to dynamically provision `hostPath`-based persistent volumes to deterministic directories for a single-node cluster.

For example, if I want a PV created from a PVC that stores content in `/srv/data/{namespaceName}/{pvcName}`, I can't use any existing provisioner to achieve this.

[Local Path Provisioner](https://github.com/rancher/local-path-provisioner) is able to dynamically provision `hostPath`-based persistent volume, however the path is not deterministic so I created the Static Path Provisioner.

## Features

- Simple: < 200 lines of code
- Works "out of the box": No environment variables, command line arguments, or config maps needed
- Support multiple storage class with different `storagePath`

## Usage

Simply apply the kustomization manifest to create a default storage class provisioned by Static Path Provisioner:

```bash
kubectl apply -k deploy/kustomize/default-storage-class
```

You may also create a non-default storage class, and optionally change the `storagePath` parameter in `deploy/kustomize/custom-storage-path/storage-class.yaml` before applying.

```bash
kubectl apply -k deploy/kustomize/custom-storage-path
```

A directory will be created at `{storagePath}/{namespaceName}/{pvcName}` when pod starts using the PV created by PVC.

## Note

Directory or content for a persistent volume will not be deleted when its persistent volume gets deleted, as intended. I might add a parameter to storage class to delete the content in the future (will not be the default behaviour), as I don't need it for now.

Static Path Provisioner is designed for single-node clusters for now, using NFS-mounted path as `storagePath` may allows you to dynamically provision persistent volumes in multi-node clusters, but it's not tested.

Do **not** use it in production.
