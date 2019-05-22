# alluxio-csi
Container Storage Interface For Alluxio

Currently Alluxio can be accessed as a POSIX filesystem using FUSE.
The alluxio-csi (this repo) provides the following advantages.

1. Currently the FUSE mount must be provisioned ahead of time, before the application.
This operational dependency is not ideal when using Kubernetes.  
The alluxio-csi enables the creation of StorageClass objects to enable dynamic provisioning the Alluxio FUSE mount just-in-time of application pod creation.

2. Currently the FUSE mount is global. All access to the FUSE mount shares the Alluxio mount path and permissions.
This is not ideal when using Kubernetes where different containers may be using different Alluxio mount points with different permissions.
The alluxio-csi enables the creation of different StorageClass objects that specifies different Alluxio mount points and permissions.
Pods can be configured to use the different StorageClass objects by creating their pod specific PersistentVolumeClaims objects.


[Example provisioning of Alluxio Cluster with CSI](https://github.com/mingfang/terraform-provider-k8s/blob/81a9ceb02625bf06da873b206b7612ad56cf62cf/examples/alluxio/main.tf#L112)
[Example usage of alluxio-csi using PersistentVolumeClaim](https://github.com/mingfang/terraform-provider-k8s/blob/81a9ceb02625bf06da873b206b7612ad56cf62cf/examples/dremio/main.tf#L103)

  