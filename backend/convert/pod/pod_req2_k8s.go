package pod

import (
	coreV1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	pod_req "kubejiangnan/model/pod/request"
	"strconv"
	"strings"
)

const (
	probe_http = "http"
	probe_tcp  = "tcp"
	probe_exec = "exec"
)

const (
	volume_emptyDir = "emptyDir"
)

type Req2K8sConvert struct {
}

// turn pod request format into k8s struct format
func (pc *Req2K8sConvert) PodReq2K8s(podReq pod_req.Pod) *coreV1.Pod {
	return &coreV1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      podReq.Base.Name,
			Namespace: podReq.Base.NameSpace,
			Labels:    pc.getK8sLabels(podReq.Base.Labels),
		},
		Spec: coreV1.PodSpec{
			InitContainers: pc.getK8sContainers(podReq.InitContainers),
			Containers:     pc.getK8sContainers(podReq.Containers),
			Volumes:        pc.getK8sVolumes(podReq.Volumes),
			DNSConfig: &coreV1.PodDNSConfig{
				Nameservers: podReq.NetWorking.DnsConfig.NameServers,
			},
			DNSPolicy:     coreV1.DNSPolicy(podReq.NetWorking.DnsPolicy),
			HostAliases:   pc.getK8sHostAliases(podReq.NetWorking.HostAliases),
			Hostname:      podReq.NetWorking.HostName,
			RestartPolicy: coreV1.RestartPolicy(podReq.Base.RestartPolicy),
		},
	}
}

func (pc *Req2K8sConvert) getK8sLabels(podReqLabels []pod_req.ListMapItem) map[string]string {
	podK8sLabels := make(map[string]string)
	for _, label := range podReqLabels {
		podK8sLabels[label.Key] = label.Value
	}
	return podK8sLabels
}

func (pc *Req2K8sConvert) getK8sContainers(podReqContainers []pod_req.Container) []coreV1.Container {
	podK8sContainers := make([]coreV1.Container, 0)
	for _, container := range podReqContainers {
		podK8sContainers = append(podK8sContainers, pc.getK8sContainer(container))
	}
	return podK8sContainers
}

func (pc *Req2K8sConvert) getK8sContainer(podReqContainer pod_req.Container) coreV1.Container {
	return coreV1.Container{
		Name:            podReqContainer.Name,
		Image:           podReqContainer.Image,
		ImagePullPolicy: coreV1.PullPolicy(podReqContainer.ImagePullPolicy),
		TTY:             podReqContainer.Tty,
		Command:         podReqContainer.Command,
		Args:            podReqContainer.Args,
		WorkingDir:      podReqContainer.WorkingDir,
		SecurityContext: &coreV1.SecurityContext{
			Privileged: &podReqContainer.Privileged,
		},
		Ports:          pc.getK8sPorts(podReqContainer.Ports),
		Env:            pc.getK8sEnv(podReqContainer.Envs),
		VolumeMounts:   pc.getK8sVolumeMount(podReqContainer.VolumeMounts),
		StartupProbe:   pc.getK8sProbe(podReqContainer.StartupProbe),
		LivenessProbe:  pc.getK8sProbe(podReqContainer.LivenessProbe),
		ReadinessProbe: pc.getK8sProbe(podReqContainer.ReadinessProbe),
		Resources:      pc.getK8sResources(podReqContainer.Resources),
	}
}

func (pc *Req2K8sConvert) getK8sEnv(podReqEnvs []pod_req.ListMapItem) []coreV1.EnvVar {
	podK8sEnvs := make([]coreV1.EnvVar, 0)
	for _, env := range podReqEnvs {
		podK8sEnvs = append(podK8sEnvs, coreV1.EnvVar{
			Name:  env.Key,
			Value: env.Value,
		})
	}
	return podK8sEnvs
}

func (pc *Req2K8sConvert) getK8sVolumeMount(podReqVolumeMounts []pod_req.VolumeMount) []coreV1.VolumeMount {
	podK8sVolumeMounts := make([]coreV1.VolumeMount, 0)
	for _, volumeMount := range podReqVolumeMounts {
		podK8sVolumeMounts = append(podK8sVolumeMounts, coreV1.VolumeMount{
			Name:      volumeMount.MountName,
			MountPath: volumeMount.MountPath,
			ReadOnly:  volumeMount.ReadOnly,
		})
	}
	return podK8sVolumeMounts
}

func (pc *Req2K8sConvert) getK8sProbe(podReqProbe pod_req.ContainerProbe) *coreV1.Probe {
	if !podReqProbe.Enable {
		return nil
	}
	var podK8sProbe coreV1.Probe
	switch podReqProbe.Type {
	case probe_http:
		httpGet := podReqProbe.HttpGet
		k8sHttpHeaders := make([]coreV1.HTTPHeader, 0)

		for _, httpHeader := range httpGet.HttpHeaders {
			k8sHttpHeaders = append(k8sHttpHeaders, coreV1.HTTPHeader{
				Name:  httpHeader.Key,
				Value: httpHeader.Value,
			})
		}

		podK8sProbe.HTTPGet = &coreV1.HTTPGetAction{
			Scheme:      coreV1.URIScheme(httpGet.Scheme),
			Host:        httpGet.Host,
			Port:        intstr.FromInt(int(httpGet.Port)),
			Path:        httpGet.Path,
			HTTPHeaders: k8sHttpHeaders,
		}
	case probe_tcp:
		tcpSocket := podReqProbe.TcpSocket
		podK8sProbe.TCPSocket = &coreV1.TCPSocketAction{
			Host: tcpSocket.Host,
			Port: intstr.FromInt(int(tcpSocket.Port)),
		}
	case probe_exec:
		exec := podReqProbe.Exec
		podK8sProbe.Exec = &coreV1.ExecAction{
			Command: exec.Command,
		}
	}

	return &podK8sProbe
}

func (pc *Req2K8sConvert) getK8sResources(podReqSources pod_req.Resources) coreV1.ResourceRequirements {
	var podK8sResources coreV1.ResourceRequirements
	if !podReqSources.Enable {
		return podK8sResources
	}
	podK8sResources.Requests = coreV1.ResourceList{
		coreV1.ResourceCPU:    resource.MustParse(strconv.Itoa(int(podReqSources.CpuRequest)) + "m"),
		coreV1.ResourceMemory: resource.MustParse(strconv.Itoa(int(podReqSources.MemoryRequest)) + "Mi"),
	}
	podK8sResources.Limits = coreV1.ResourceList{
		coreV1.ResourceCPU:    resource.MustParse(strconv.Itoa(int(podReqSources.CpuLimit)) + "m"),
		coreV1.ResourceMemory: resource.MustParse(strconv.Itoa(int(podReqSources.MemoryLimit)) + "Mi"),
	}
	return podK8sResources
}

func (pc *Req2K8sConvert) getK8sPorts(podReqPorts []pod_req.ContainerPort) []coreV1.ContainerPort {
	podK8sContainedPorts := make([]coreV1.ContainerPort, 0)
	for _, port := range podReqPorts {
		podK8sContainedPorts = append(podK8sContainedPorts, coreV1.ContainerPort{
			Name:          port.Name,
			HostPort:      port.HostPort,
			ContainerPort: port.ContainerPort,
		})
	}
	return podK8sContainedPorts
}

func (pc *Req2K8sConvert) getK8sVolumes(podReqVolumes []pod_req.Volume) []coreV1.Volume {
	podK8sVolumes := make([]coreV1.Volume, 0)

	for _, volume := range podReqVolumes {
		if volume.Type != volume_emptyDir {
			continue
		}
		source := coreV1.VolumeSource{
			EmptyDir: &coreV1.EmptyDirVolumeSource{},
		}
		podK8sVolumes = append(podK8sVolumes, coreV1.Volume{
			VolumeSource: source,
			Name:         volume.Name,
		})
	}

	return podK8sVolumes
}

func (pc *Req2K8sConvert) getK8sHostAliases(podReqHostAliases []pod_req.ListMapItem) []coreV1.HostAlias {
	podK8sHostAliases := make([]coreV1.HostAlias, 0)

	for _, hostAlias := range podReqHostAliases {
		podK8sHostAliases = append(podK8sHostAliases, coreV1.HostAlias{
			IP:        hostAlias.Key,
			Hostnames: strings.Split(hostAlias.Value, ","),
		})
	}

	return podK8sHostAliases
}
