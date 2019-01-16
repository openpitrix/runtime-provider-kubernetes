// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package runtime_provider

import (
	"bytes"
	"context"
	"fmt"
	"strings"

	"google.golang.org/grpc/transport"
	"k8s.io/helm/pkg/chartutil"
	"k8s.io/helm/pkg/proto/hapi/chart"
	"k8s.io/helm/pkg/proto/hapi/release"

	appclient "openpitrix.io/openpitrix/pkg/client/app"
	runtimeclient "openpitrix.io/openpitrix/pkg/client/runtime"
	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/sender"
	"openpitrix.io/openpitrix/pkg/util/funcutil"
	"openpitrix.io/openpitrix/pkg/util/pbutil"
)

func getChart(ctx context.Context, versionId string) (*chart.Chart, error) {
	appClient, err := appclient.NewAppManagerClient()
	if err != nil {
		return nil, err
	}

	req := pb.GetAppVersionPackageRequest{
		VersionId: pbutil.ToProtoString(versionId),
	}

	resp, err := appClient.GetAppVersionPackage(ctx, &req)
	if err != nil {
		return nil, err
	}

	pkg := resp.GetPackage()
	r := bytes.NewReader(pkg)

	c, err := chartutil.LoadArchive(r)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (p *Server) ParseClusterConf(ctx context.Context, req *pb.ParseClusterConfRequest) (*pb.ParseClusterConfResponse, error) {
	versionId := req.GetVersionId().GetValue()
	runtimeId := req.GetRuntimeId().GetValue()
	conf := req.GetConf().GetValue()
	cluster := models.PbToClusterWrapper(req.GetCluster())

	c, err := getChart(ctx, versionId)
	if err != nil {
		logger.Error(ctx, "Load helm chart from app version [%s] failed: %+v", versionId, err)
		return nil, err
	}

	runtime, err := runtimeclient.NewRuntime(ctx, runtimeId)
	if err != nil {
		return nil, err
	}
	namespace := runtime.Zone

	parser := Parser{
		ctx:       ctx,
		Chart:     c,
		Conf:      conf,
		VersionId: versionId,
		RuntimeId: runtimeId,
		Namespace: namespace,
	}
	err = parser.Parse(cluster)
	if err != nil {
		logger.Error(ctx, "Parse app version [%s] failed: %+v", versionId, err)
		return nil, err
	}

	return &pb.ParseClusterConfResponse{
		Cluster: models.ClusterWrapperToPb(cluster),
	}, err
}

func (p *Server) SplitJobIntoTasks(ctx context.Context, req *pb.SplitJobIntoTasksRequest) (*pb.SplitJobIntoTasksResponse, error) {
	job := models.PbToJob(req.GetJob())
	jobDirective, err := decodeJobDirective(ctx, job.Directive)
	if err != nil {
		return nil, err
	}

	tl := new(models.TaskLayer)

	switch job.JobAction {
	case constants.ActionCreateCluster:
		td := TaskDirective{
			VersionId:         job.VersionId,
			Namespace:         jobDirective.Namespace,
			RuntimeId:         jobDirective.RuntimeId,
			Values:            jobDirective.Values,
			ClusterName:       jobDirective.ClusterName,
			RawClusterWrapper: job.Directive,
		}
		tdj := encodeTaskDirective(td)

		task := models.NewTask(constants.PlaceHolder, job.JobId, "", jobDirective.RuntimeId, constants.ActionCreateCluster, tdj, sender.OwnerPath(job.OwnerPath), false)
		tl = &models.TaskLayer{
			Tasks: []*models.Task{task},
			Child: nil,
		}
	case constants.ActionUpgradeCluster:
		td := TaskDirective{
			VersionId:         job.VersionId,
			Namespace:         jobDirective.Namespace,
			RuntimeId:         jobDirective.RuntimeId,
			Values:            jobDirective.Values,
			ClusterName:       jobDirective.ClusterName,
			RawClusterWrapper: job.Directive,
		}
		tdj := encodeTaskDirective(td)

		task := models.NewTask(constants.PlaceHolder, job.JobId, "", jobDirective.RuntimeId, constants.ActionUpgradeCluster, tdj, sender.OwnerPath(job.OwnerPath), false)
		tl = &models.TaskLayer{
			Tasks: []*models.Task{task},
			Child: nil,
		}
	case constants.ActionUpdateClusterEnv:
		td := TaskDirective{
			VersionId:         job.VersionId,
			Namespace:         jobDirective.Namespace,
			RuntimeId:         jobDirective.RuntimeId,
			Values:            jobDirective.Values,
			ClusterName:       jobDirective.ClusterName,
			RawClusterWrapper: job.Directive,
		}
		tdj := encodeTaskDirective(td)

		task := models.NewTask(constants.PlaceHolder, job.JobId, "", jobDirective.RuntimeId, constants.ActionUpgradeCluster, tdj, sender.OwnerPath(job.OwnerPath), false)
		tl = &models.TaskLayer{
			Tasks: []*models.Task{task},
			Child: nil,
		}
	case constants.ActionRollbackCluster:
		td := TaskDirective{
			Namespace:         jobDirective.Namespace,
			RuntimeId:         jobDirective.RuntimeId,
			ClusterName:       jobDirective.ClusterName,
			RawClusterWrapper: job.Directive,
		}
		tdj := encodeTaskDirective(td)

		task := models.NewTask(constants.PlaceHolder, job.JobId, "", jobDirective.RuntimeId, constants.ActionRollbackCluster, tdj, sender.OwnerPath(job.OwnerPath), false)
		tl = &models.TaskLayer{
			Tasks: []*models.Task{task},
			Child: nil,
		}
	case constants.ActionDeleteClusters:
		td := TaskDirective{
			RuntimeId:   jobDirective.RuntimeId,
			ClusterName: jobDirective.ClusterName,
		}
		tdj := encodeTaskDirective(td)

		task := models.NewTask(constants.PlaceHolder, job.JobId, "", jobDirective.RuntimeId, constants.ActionDeleteClusters, tdj, sender.OwnerPath(job.OwnerPath), false)
		tl = &models.TaskLayer{
			Tasks: []*models.Task{task},
			Child: nil,
		}
	case constants.ActionCeaseClusters:
		td := TaskDirective{
			RuntimeId:   jobDirective.RuntimeId,
			ClusterName: jobDirective.ClusterName,
		}
		tdj := encodeTaskDirective(td)

		task := models.NewTask(constants.PlaceHolder, job.JobId, "", jobDirective.RuntimeId, constants.ActionCeaseClusters, tdj, sender.OwnerPath(job.OwnerPath), false)
		tl = &models.TaskLayer{
			Tasks: []*models.Task{task},
			Child: nil,
		}
	default:
		return nil, fmt.Errorf("the job action [%s] is not supported", job.JobAction)
	}
	return &pb.SplitJobIntoTasksResponse{
		TaskLayer: models.TaskLayerToPb(tl),
	}, nil
}

func (p *Server) HandleSubtask(ctx context.Context, req *pb.HandleSubtaskRequest) (*pb.HandleSubtaskResponse, error) {
	task := models.PbToTask(req.GetTask())
	taskDirective, err := decodeTaskDirective(task.Directive)
	if err != nil {
		return nil, err
	}

	helmHandler := GetHelmHandler(ctx, taskDirective.RuntimeId)

	switch task.TaskAction {
	case constants.ActionCreateCluster:
		c, err := getChart(ctx, taskDirective.VersionId)
		if err != nil {
			return nil, err
		}

		rawVals, err := ConvertJsonToYaml([]byte(taskDirective.Values))
		if err != nil {
			return nil, err
		}

		logger.Debug(ctx, "Install helm release with name [%+v], namespace [%+v], values [%s]", taskDirective.ClusterName, taskDirective.Namespace, rawVals)

		err = helmHandler.InstallReleaseFromChart(c, taskDirective.Namespace, rawVals, taskDirective.ClusterName)
		if err != nil {
			return nil, err
		}
	case constants.ActionUpgradeCluster:
		c, err := getChart(ctx, taskDirective.VersionId)
		if err != nil {
			return nil, err
		}

		rawVals, err := ConvertJsonToYaml([]byte(taskDirective.Values))
		if err != nil {
			return nil, err
		}

		logger.Debug(ctx, "Update helm release [%+v] with values [%s]", taskDirective.ClusterName, rawVals)

		err = helmHandler.UpdateReleaseFromChart(taskDirective.ClusterName, c, rawVals)
		if err != nil {
			return nil, err
		}
	case constants.ActionRollbackCluster:
		err = helmHandler.RollbackRelease(taskDirective.ClusterName)
		if err != nil {
			return nil, err
		}
	case constants.ActionDeleteClusters:
		err = helmHandler.DeleteRelease(taskDirective.ClusterName, false)
		if err != nil {
			return nil, err
		}
	case constants.ActionCeaseClusters:
		err = helmHandler.DeleteRelease(taskDirective.ClusterName, true)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("the task action [%s] is not supported", task.TaskAction)
	}

	return &pb.HandleSubtaskResponse{
		Task: models.TaskToPb(task),
	}, nil
}

func (p *Server) WaitSubtask(ctx context.Context, req *pb.WaitSubtaskRequest) (*pb.WaitSubtaskResponse, error) {
	task := models.PbToTask(req.GetTask())
	taskDirective, err := decodeTaskDirective(task.Directive)
	if err != nil {
		return nil, err
	}

	helmHandler := GetHelmHandler(ctx, taskDirective.RuntimeId)

	err = funcutil.WaitForSpecificOrError(func() (bool, error) {
		switch task.TaskAction {
		case constants.ActionCreateCluster:
			fallthrough
		case constants.ActionUpgradeCluster:
			fallthrough
		case constants.ActionRollbackCluster:
			resp, err := helmHandler.ReleaseStatus(taskDirective.ClusterName)
			if err != nil {
				if _, ok := err.(transport.ConnectionError); ok {
					return false, nil
				}
				return true, err
			}

			switch resp.Info.Status.Code {
			case release.Status_FAILED:
				logger.Debug(ctx, "Helm release gone to failed")
				return true, fmt.Errorf("release failed")
			case release.Status_DEPLOYED:
				clusterWrapper, err := models.NewClusterWrapper(ctx, taskDirective.RawClusterWrapper)
				if err != nil {
					return true, err
				}

				kubeHandler := GetKubeHandler(ctx, taskDirective.RuntimeId)
				err = kubeHandler.WaitWorkloadReady(
					taskDirective.RuntimeId,
					taskDirective.Namespace,
					clusterWrapper.ClusterRoles,
					task.GetTimeout(constants.WaitTaskTimeout),
					constants.WaitTaskInterval,
				)
				if err != nil {
					return true, err
				}

				return true, nil
			}
		case constants.ActionDeleteClusters:
			resp, err := helmHandler.ReleaseStatus(taskDirective.ClusterName)
			if err != nil {
				if _, ok := err.(transport.ConnectionError); ok {
					return false, nil
				}
				if strings.Contains(err.Error(), "not found") {
					logger.Warn(nil, "Waiting on a helm release not existed, %+v", err)
					return true, nil
				}
				return true, err
			}

			if resp.Info.Status.Code == release.Status_DELETED {
				return true, nil
			}
		case constants.ActionCeaseClusters:
			_, err := helmHandler.ReleaseStatus(taskDirective.ClusterName)
			if err != nil {
				if _, ok := err.(transport.ConnectionError); ok {
					return false, nil
				}
				if strings.Contains(err.Error(), "not found") {
					logger.Warn(nil, "Waiting on a helm release not existed, %+v", err)
					return true, nil
				}
				return true, nil
			}
		}
		return false, nil
	}, task.GetTimeout(constants.WaitHelmTaskTimeout), constants.WaitTaskInterval)

	if err != nil {
		return nil, err
	} else {
		return &pb.WaitSubtaskResponse{
			Task: models.TaskToPb(task),
		}, nil
	}
}

func (p *Server) DescribeSubnets(ctx context.Context, req *pb.DescribeSubnetsRequest) (*pb.DescribeSubnetsResponse, error) {
	return nil, fmt.Errorf("the action DescribeSubnets is not supported")
}

func (p *Server) CheckResource(ctx context.Context, req *pb.CheckResourceRequest) (*pb.CheckResourceResponse, error) {
	cluster := models.PbToClusterWrapper(req.GetCluster())
	helmHandler := GetHelmHandler(ctx, cluster.Cluster.RuntimeId)

	err := helmHandler.CheckClusterNameIsUnique(cluster.Cluster.Name)
	if err != nil {
		logger.Error(ctx, "Cluster name [%s] already existed in runtime [%s]: %+v",
			cluster.Cluster.Name, cluster.Cluster.RuntimeId, err)
		return &pb.CheckResourceResponse{
			Ok: pbutil.ToProtoBool(false),
		}, err
	} else {
		return &pb.CheckResourceResponse{
			Ok: pbutil.ToProtoBool(true),
		}, nil
	}
}

func (p *Server) DescribeVpc(ctx context.Context, req *pb.DescribeVpcRequest) (*pb.DescribeVpcResponse, error) {
	return nil, fmt.Errorf("the action DescribeSubnets is not supported")
}

func (p *Server) DescribeClusterDetails(ctx context.Context, req *pb.DescribeClusterDetailsRequest) (*pb.DescribeClusterDetailsResponse, error) {
	cluster := models.PbToClusterWrapper(req.GetCluster())
	kubeHandler := GetKubeHandler(ctx, cluster.Cluster.RuntimeId)
	err := kubeHandler.DescribeClusterDetails(cluster)
	return &pb.DescribeClusterDetailsResponse{
		Cluster: models.ClusterWrapperToPb(cluster),
	}, err
}

func (p *Server) ValidateRuntime(ctx context.Context, req *pb.ValidateRuntimeRequest) (*pb.ValidateRuntimeResponse, error) {
	runtimeId := req.GetRuntimeId().GetValue()
	zone := req.GetZone().GetValue()
	needCreate := req.GetNeedCreate().GetValue()
	runtimeCredential := models.PbToRuntimeCredential(req.GetRuntimeCredential())
	kubeHandler := GetKubeHandler(ctx, runtimeId)
	err := kubeHandler.ValidateRuntime(zone, runtimeCredential, needCreate)
	if err != nil {
		return &pb.ValidateRuntimeResponse{
			Ok: pbutil.ToProtoBool(false),
		}, err
	} else {
		return &pb.ValidateRuntimeResponse{
			Ok: pbutil.ToProtoBool(true),
		}, nil
	}
}

func (p *Server) DescribeZones(ctx context.Context, req *pb.DescribeZonesRequest) (*pb.DescribeZonesResponse, error) {
	runtimeCredential := models.PbToRuntimeCredential(req.GetRuntimeCredential())
	kubeHandler := GetKubeHandler(ctx, "")
	zones, err := kubeHandler.DescribeRuntimeProviderZones(runtimeCredential)
	return &pb.DescribeZonesResponse{
		Zones: zones,
	}, err
}
