package request

type Pod struct {
	Base           Base        `json:"base"` // basic definition information
	Volumes        []Volume    `json:"volumes"`
	NetWorking     NetWorking  `json:"netWorking"` // network configuration
	InitContainers []Container `json:"initContainers"`
	Containers     []Container `json:"containers"`
}

type Base struct {
	Name          string        `json:"name"`
	Labels        []ListMapItem `json:"labels"`
	NameSpace     string        `json:"namespace"`
	RestartPolicy string        `json:"restartPolicy"` // Always | Never | On-Failure
}

type ListMapItem struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Volume struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type NetWorking struct {
	HostNetWork bool          `json:"hostNetWork"`
	HostName    string        `json:"hostName"`
	DnsPolicy   string        `json:"dnsPolicy"`
	DnsConfig   DnsConfig     `json:"dnsConfig"`
	HostAliases []ListMapItem `json:"hostAliases"`
}

type DnsConfig struct {
	NameServers []string `json:"nameServers"`
}

type Container struct {
	Name            string          `json:"name"`
	Image           string          `json:"image"`
	ImagePullPolicy string          `json:"ImagePullPolicy"`
	Tty             bool            `json:"tty"` // Turn on Fake Terminal or not
	Ports           []ContainerPort `json:"ports"`
	WorkingDir      string          `json:"workingDir"`
	Command         []string        `json:"command"`
	Args            []string        `json:"args"`
	Envs            []ListMapItem   `json:"envs"`
	Privileged      bool            `json:"privileged"` // Turn on Privileged Model
	Resources       Resources       `json:"resources"`
	VolumeMounts    []VolumeMount   `json:"volumeMounts"`
	StartupProbe    ContainerProbe  `json:"startupProbe"`
	LivenessProbe   ContainerProbe  `json:"livenessProbe"`
	ReadinessProbe  ContainerProbe  `json:"readinessProbe"`
}

type Resources struct {
	Enable        bool  `json:"enable"`
	MemoryRequest int32 `json:"memoryRequest"`
	MemoryLimit   int32 `json:"memoryLimit"`
	CpuRequest    int32 `json:"cpuRequest"`
	CpuLimit      int32 `json:"cpuLimit"`
}

type VolumeMount struct {
	MountName string `json:"mountName"`
	MountPath string `json:"mountPath"` // mount path in container
	ReadOnly  bool   `json:"readOnly"`
}

type ContainerProbe struct {
	Enable    bool           `json:"enable"` // Turn on probe or not
	Type      string         `json:"type"`
	HttpGet   ProbeHttpGet   `json:"httpGet"` // Kind of probe: http | command | tcp
	Exec      ProbeCommand   `json:"exec"`
	TcpSocket ProbeTcpSocket `json:"tcpSocket"`
	ProbeTime ProbeTime      `json:"probeTime"`
}

type ProbeHttpGet struct {
	Scheme      string        `json:"scheme"` // http | https
	Host        string        `json:"host"`   // if nil, request in pod
	Path        string        `json:"path"`
	Port        int32         `json:"port"`
	HttpHeaders []ListMapItem `json:"httpHeaders"`
}

type ProbeCommand struct {
	Command []string `json:"command"`
}

type ProbeTcpSocket struct {
	Host string `json:"host"`
	Port int32  `json:"port"`
}

type ProbeTime struct {
	InitialDelaySeconds int32 `json:"initialDelaySeconds"` // After a few seconds to start probe
	PeriodSeconds       int32 `json:"periodSeconds"`       // Probe in a period time
	TimeOutSeconds      int32 `json:"timeOutSeconds"`      // Probe Waiting time before fail
	SuccessThreshold    int32 `json:"successThreshold"`    // Probe for a few times successfully before it successes
	FailureThreshold    int32 `json:"failureThreshold"`    // Probe for a few times fail before it fail
}

type ContainerPort struct {
	Name          string `json:"name"`
	ContainerPort int32  `json:"containerPort"`
	HostPort      int32  `json:"hostPort"`
}
