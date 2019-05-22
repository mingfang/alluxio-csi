package alluxio

import (
    "fmt"
    "log"
    "os"
    "os/exec"
    "strings"

    "github.com/container-storage-interface/spec/lib/go/csi"
    "github.com/kubernetes-csi/drivers/pkg/csi-common"
    "golang.org/x/net/context"
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
    "k8s.io/kubernetes/pkg/util/mount"
)

type nodeServer struct {
    *csicommon.DefaultNodeServer
}

func (ns *nodeServer) NodePublishVolume(ctx context.Context, req *csi.NodePublishVolumeRequest) (*csi.NodePublishVolumeResponse, error) {
    targetPath := req.GetTargetPath()

    notMnt, err := mount.New("").IsLikelyNotMountPoint(targetPath)
    log.Println("err:", err)
    if err != nil {
        if os.IsNotExist(err) {
            if err := os.MkdirAll(targetPath, 0750); err != nil {
                log.Println("MkdirAll", "err:", err)
                return nil, status.Error(codes.Internal, err.Error())
            }
            notMnt = true
        } else {
            return nil, status.Error(codes.Internal, err.Error())
        }
    }

    if !notMnt {
        return &csi.NodePublishVolumeResponse{}, nil
    }

    mo := req.GetVolumeCapability().GetMount().GetMountFlags()
    if req.GetReadonly() {
        mo = append(mo, "ro")
    }

    masterHost := req.GetVolumeContext()["alluxio.master.hostname"]
    masterPort := req.GetVolumeContext()["alluxio.master.port"]

    command := exec.Command("/opt/alluxio/integration/fuse/bin/alluxio-fuse",
        "mount", "-o", "allow_other", targetPath, "/",
    )
    alluxioJavaOpts := fmt.Sprintf("ALLUXIO_JAVA_OPTS=-Dalluxio.master.hostname=%s -Dalluxio.master.port=%s", masterHost, masterPort)
    command.Env = append(os.Environ(), alluxioJavaOpts)
    err = command.Run()
    if err != nil {
        if os.IsPermission(err) {
            return nil, status.Error(codes.PermissionDenied, err.Error())
        }
        if strings.Contains(err.Error(), "invalid argument") {
            return nil, status.Error(codes.InvalidArgument, err.Error())
        }
        return nil, status.Error(codes.Internal, err.Error())
    }

    return &csi.NodePublishVolumeResponse{}, nil
}

func (ns *nodeServer) NodeUnpublishVolume(ctx context.Context, req *csi.NodeUnpublishVolumeRequest) (*csi.NodeUnpublishVolumeResponse, error) {
    targetPath := req.GetTargetPath()
    notMnt, err := mount.New("").IsLikelyNotMountPoint(targetPath)

    if err != nil {
        if os.IsNotExist(err) {
            return nil, status.Error(codes.NotFound, "Targetpath not found")
        } else {
            return nil, status.Error(codes.Internal, err.Error())
        }
    }
    if notMnt {
        return nil, status.Error(codes.NotFound, "Volume not mounted")
    }

    command := exec.Command("/opt/alluxio/integration/fuse/bin/alluxio-fuse",
        "unmount", targetPath,
    )
    err = command.Run()

    err = mount.CleanupMountPoint(req.GetTargetPath(), mount.New(""), false)
    if err != nil {
        return nil, status.Error(codes.Internal, err.Error())
    }

    return &csi.NodeUnpublishVolumeResponse{}, nil
}

func (ns *nodeServer) NodeUnstageVolume(ctx context.Context, req *csi.NodeUnstageVolumeRequest) (*csi.NodeUnstageVolumeResponse, error) {
    return &csi.NodeUnstageVolumeResponse{}, nil
}

func (ns *nodeServer) NodeStageVolume(ctx context.Context, req *csi.NodeStageVolumeRequest) (*csi.NodeStageVolumeResponse, error) {
    return &csi.NodeStageVolumeResponse{}, nil
}

func (ns *nodeServer) NodeExpandVolume(ctx context.Context, req *csi.NodeExpandVolumeRequest) (*csi.NodeExpandVolumeResponse, error) {
    return nil, status.Error(codes.Unimplemented, "")
}
