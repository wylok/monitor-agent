package config

const (
	Version    = "2021042306"
	UrlPath    = "/xxx"
	CryptKey   = "4096779a2529ca11b8508805ahf88a2d"
	FilePath   = "/xxx"
	PidPath    = FilePath + "/"
	PidFile    = PidPath + "monitor-agent.pid"
	LogFile    = PidPath + "monitor-agent.log"
	HostIdFile = FilePath + "instance-id"
	ProxyUrl   = "http://xxx:12345"
)

var Interval uint64 = 60
var Process []string

type Config struct {
	AgentId string `yaml:"agent_id"`
}

type ApiJson struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type ProcessInfo struct {
	Cpu float64 `json:"cpu_usage"`
	Mem float32 `json:"mem_pused"`
	Num int32   `json:"numbers"`
}

type AgentConf struct {
	Action  string   `json:"action"`
	Process []string `json:"process"`
	Version string   `json:"version"`
}

type CollectionData struct {
	Platform         string                 `json:"platform"`
	Resource         string                 `json:"resource"`
	Item             string                 `json:"item"`
	NowTime          string                 `json:"now_time"`
	HostId           string                 `json:"host_id"`
	CpuLoadavg       float64                `json:"cpu_loadavg"`
	CpuUsage         float64                `json:"cpu_usage"`
	MemUsed          uint64                 `json:"mem_used"`
	MemAvailable     uint64                 `json:"mem_available"`
	MemPavailable    float64                `json:"mem_pavailable"`
	MemPused         float64                `json:"mem_pused"`
	DiskUsage        float64                `json:"disk_usage"`
	DiskReadTraffic  uint64                 `json:"disk_read_traffic"`
	DiskWriteTraffic uint64                 `json:"disk_write_traffic"`
	LanOuttraffic    uint64                 `json:"lan_outtraffic"`
	LanIntraffic     uint64                 `json:"lan_intraffic"`
	LanOutpkg        uint64                 `json:"lan_outpkg"`
	LanInpkg         uint64                 `json:"lan_inpkg"`
	TcpEstab         int                    `json:"tcp_estab"`
	Process          map[string]ProcessInfo `json:"process"`
}
