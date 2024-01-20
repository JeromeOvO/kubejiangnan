package pod

import (
	"fmt"
	coreV1 "k8s.io/api/core/v1"
	pod_req "kubejiangnan/model/pod/request"
	pod_res "kubejiangnan/model/pod/response"
	"strings"
)

const volume_type_emptydir = "emptyDir"

type K8s2ReqConvert struct {
	volumeMap map[string]string
}

func (pc *K8s2ReqConvert) PodK8s2ItemRes(pod coreV1.Pod) pod_res.PodListItem {

	var totalContainers, readyContainers, restartContainer int32
	for _, containerStatus := range pod.Status.ContainerStatuses {
		if containerStatus.Ready {
			readyContainers++
		}
		restartContainer += containerStatus.RestartCount
		totalContainers++
	}

	var podStatus string
	if pod.Status.Phase != "Running" {
		podStatus = "Error"
	} else {
		podStatus = "Running"
	}

	return pod_res.PodListItem{
		Name:    pod.Name,
		Ready:   fmt.Sprintf("%d/%d", readyContainers, totalContainers),
		Status:  podStatus,
		Restart: restartContainer,
		Age:     pod.CreationTimestamp.Unix(),
		IP:      pod.Status.PodIP,
		Node:    pod.Spec.NodeName,
	}
}

func (pc *K8s2ReqConvert) PodK8s2Req(podK8s coreV1.Pod) pod_req.Pod {
	return pod_req.Pod{
		Base:           pc.getReqBase(podK8s),
		NetWorking:     pc.getReqNetworking(podK8s),
		Volumes:        pc.getReqVolumes(podK8s.Spec.Volumes),
		Containers:     pc.getReqContainers(podK8s.Spec.Containers),
		InitContainers: pc.getReqContainers(podK8s.Spec.InitContainers),
	}
}

func (pc *K8s2ReqConvert) getReqBase(pod coreV1.Pod) pod_req.Base {
	return pod_req.Base{
		Name:          pod.Name,
		NameSpace:     pod.Namespace,
		Labels:        pc.getReqLabels(pod.Labels),
		RestartPolicy: string(pod.Spec.RestartPolicy),
	}
}

func (pc *K8s2ReqConvert) getReqLabels(data map[string]string) []pod_req.ListMapItem {
	reqLabels := make([]pod_req.ListMapItem, 0)
	for k, v := range data {
		reqLabels = append(reqLabels, pod_req.ListMapItem{
			Key:   k,
			Value: v,
		})
	}
	return reqLabels
}

func (pc *K8s2ReqConvert) getReqNetworking(pod coreV1.Pod) pod_req.NetWorking {
	return pod_req.NetWorking{
		HostNetWork: pod.Spec.HostNetwork,
		HostName:    pod.Spec.Hostname,
		DnsPolicy:   string(pod.Spec.DNSPolicy),
		DnsConfig:   pc.getReqDnsConfig(pod.Spec.DNSConfig),
		HostAliases: pc.getReqHostAliases(pod.Spec.HostAliases),
	}
}

func (pc *K8s2ReqConvert) getReqDnsConfig(dnsConfig *coreV1.PodDNSConfig) pod_req.DnsConfig {
	var reqDnsConfig pod_req.DnsConfig
	if dnsConfig != nil {
		reqDnsConfig.NameServers = dnsConfig.Nameservers
	}
	return reqDnsConfig
}

func (pc *K8s2ReqConvert) getReqHostAliases(hostAlias []coreV1.HostAlias) []pod_req.ListMapItem {
	reqHostAliases := make([]pod_req.ListMapItem, 0)
	for _, alias := range hostAlias {
		reqHostAliases = append(reqHostAliases, pod_req.ListMapItem{
			Key:   alias.IP,
			Value: strings.Join(alias.Hostnames, ","),
		})
	}
	return reqHostAliases
}

func (pc *K8s2ReqConvert) getReqVolumes(volumes []coreV1.Volume) []pod_req.Volume {
	reqVolumes := make([]pod_req.Volume, 0)
	for _, volume := range volumes {
		if volume.EmptyDir == nil {
			continue
		}

		if pc.volumeMap == nil {
			pc.volumeMap = make(map[string]string)
		}

		pc.volumeMap[volume.Name] = ""

		reqVolumes = append(reqVolumes, pod_req.Volume{
			Type: volume_type_emptydir,
			Name: volume.Name,
		})
	}

	return reqVolumes
}

func (pc *K8s2ReqConvert) getReqContainers(containers []coreV1.Container) []pod_req.Container {
	reqContainers := make([]pod_req.Container, 0)

	for _, container := range containers {
		reqContainers = append(reqContainers, pc.getReqContainer(container))
	}
	return reqContainers
}

func (pc *K8s2ReqConvert) getReqContainer(container coreV1.Container) pod_req.Container {
	return pod_req.Container{
		Name:            container.Name,
		Image:           container.Image,
		ImagePullPolicy: string(container.ImagePullPolicy),
		Tty:             container.TTY,
		Ports:           pc.getReqContainerPorts(container.Ports),
		WorkingDir:      container.WorkingDir,
		Command:         container.Command,
		Args:            container.Args,
		Envs:            pc.getReqEnvs(container.Env),
		Privileged:      pc.getReqPrivileged(container.SecurityContext),
		Resources:       pc.getReqResources(container.Resources),
		VolumeMounts:    pc.getReqContainerVolumeMounts(container.VolumeMounts),
		StartupProbe:    pc.getReqContainerProbe(container.StartupProbe),
		LivenessProbe:   pc.getReqContainerProbe(container.LivenessProbe),
		ReadinessProbe:  pc.getReqContainerProbe(container.ReadinessProbe),
	}
}

func (pc *K8s2ReqConvert) getReqContainerPorts(ports []coreV1.ContainerPort) []pod_req.ContainerPort {
	reqContainerPorts := make([]pod_req.ContainerPort, 0)
	for _, port := range ports {
		reqContainerPorts = append(reqContainerPorts, pod_req.ContainerPort{
			Name:          port.Name,
			HostPort:      port.HostPort,
			ContainerPort: port.ContainerPort,
		})
	}
	return reqContainerPorts
}

func (pc *K8s2ReqConvert) getReqEnvs(envs []coreV1.EnvVar) []pod_req.ListMapItem {
	reqEnvs := make([]pod_req.ListMapItem, 0)
	for _, env := range envs {
		reqEnvs = append(reqEnvs, pod_req.ListMapItem{
			Key:   env.Name,
			Value: env.Value,
		})
	}
	return reqEnvs
}

func (pc *K8s2ReqConvert) getReqPrivileged(ctx *coreV1.SecurityContext) (privileged bool) {
	if ctx != nil {
		privileged = *ctx.Privileged
	}
	return
}

func (pc *K8s2ReqConvert) getReqResources(requirements coreV1.ResourceRequirements) pod_req.Resources {
	reqResources := pod_req.Resources{
		Enable: false,
	}
	requests := requirements.Requests
	limits := requirements.Limits

	if requests != nil {
		reqResources.Enable = true
		reqResources.CpuRequest = int32(requests.Cpu().MilliValue())
		reqResources.MemoryRequest = int32(requests.Memory().MilliValue())
	}

	if limits != nil {
		reqResources.Enable = true
		reqResources.CpuLimit = int32(limits.Cpu().MilliValue())
		reqResources.MemoryLimit = int32(limits.Memory().MilliValue())
	}

	return reqResources
}

func (pc *K8s2ReqConvert) getReqContainerVolumeMounts(volumeMounts []coreV1.VolumeMount) []pod_req.VolumeMount {
	reqVolumeMounts := make([]pod_req.VolumeMount, 0)
	for _, volumeMount := range volumeMounts {
		// filter none empty dir
		if _, ok := pc.volumeMap[volumeMount.Name]; ok {
			reqVolumeMounts = append(reqVolumeMounts, pod_req.VolumeMount{
				MountName: volumeMount.Name,
				MountPath: volumeMount.MountPath,
				ReadOnly:  volumeMount.ReadOnly,
			})
		}
	}
	return reqVolumeMounts
}

func (pc *K8s2ReqConvert) getReqContainerProbe(probe *coreV1.Probe) pod_req.ContainerProbe {
	reqContainerProbe := pod_req.ContainerProbe{
		Enable: false,
	}

	if probe != nil {
		reqContainerProbe.Enable = true
		if probe.Exec != nil {
			reqContainerProbe.Type = probe_exec
			reqContainerProbe.Exec.Command = probe.Exec.Command
		} else if probe.HTTPGet != nil {
			reqContainerProbe.Type = probe_http
			httpGet := probe.HTTPGet
			reqHeaders := make([]pod_req.ListMapItem, 0)
			for _, httpHeader := range httpGet.HTTPHeaders {
				reqHeaders = append(reqHeaders, pod_req.ListMapItem{
					Key:   httpHeader.Name,
					Value: httpHeader.Value,
				})
			}
			reqContainerProbe.HttpGet = pod_req.ProbeHttpGet{
				Host:        httpGet.Host,
				Port:        httpGet.Port.IntVal,
				Scheme:      string(httpGet.Scheme),
				Path:        httpGet.Path,
				HttpHeaders: reqHeaders,
			}
		} else if probe.TCPSocket != nil {
			reqContainerProbe.Type = probe_tcp
			reqContainerProbe.TcpSocket = pod_req.ProbeTcpSocket{
				Host: probe.TCPSocket.Host,
				Port: probe.TCPSocket.Port.IntVal,
			}
		} else {
			reqContainerProbe.Type = probe_http
			return reqContainerProbe
		}

		reqContainerProbe.ProbeTime.InitialDelaySeconds = probe.InitialDelaySeconds
		reqContainerProbe.ProbeTime.PeriodSeconds = probe.PeriodSeconds
		reqContainerProbe.ProbeTime.TimeOutSeconds = probe.TimeoutSeconds
		reqContainerProbe.ProbeTime.SuccessThreshold = probe.SuccessThreshold
		reqContainerProbe.ProbeTime.FailureThreshold = probe.FailureThreshold
	}

	return reqContainerProbe
}
