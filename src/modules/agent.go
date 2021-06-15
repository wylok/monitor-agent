package modules

import (
	"bytes"
	_ "bytes"
	"encoding/json"
	_ "github.com/CodyGuo/godaemon"
	"github.com/asmcos/requests"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
	"github.com/shirou/gopsutil/process"
	"monitor-agent/config"
	"monitor-agent/kits"
	"os"
	"strconv"
	"strings"
	"time"
)

func PsIo() map[string]uint64 {
	// 计算间隔监控数据
	bytesSent := make([]uint64, 1)
	bytesRecv := make([]uint64, 1)
	packetsSent := make([]uint64, 1)
	packetsRecv := make([]uint64, 1)
	diskIo := make(map[string]uint64)
	p, _ := disk.Partitions(false)
	dm := strings.Split(p[0].Device, "/")
	sd, _ := disk.IOCounters()
	diskIo["before_ReadBytes"] = sd[dm[len(dm)-1]].ReadBytes
	diskIo["before_WriteBytes"] = sd[dm[len(dm)-1]].WriteBytes
	n, _ := net.IOCounters(true)
	for _, v := range n {
		if !strings.HasPrefix(v.Name, "lo") && !strings.HasPrefix(v.Name, "docker") {
			bytesSent = append(bytesSent, v.BytesSent)
			bytesRecv = append(bytesRecv, v.BytesRecv)
			packetsSent = append(packetsSent, v.PacketsSent)
			packetsRecv = append(packetsRecv, v.PacketsRecv)
		}
	}
	diskIo["before_lan_intraffic"] = kits.Max(bytesRecv)
	diskIo["before_lan_outtraffic"] = kits.Max(bytesSent)
	diskIo["before_lan_inpkg"] = kits.Max(packetsRecv)
	diskIo["before_lan_outpkg"] = kits.Max(packetsSent)
	//间隔时间重新获取数据
	time.Sleep(time.Duration(config.Interval) * time.Second)
	sd, _ = disk.IOCounters()
	diskIo["after_ReadBytes"] = sd[dm[len(dm)-1]].ReadBytes
	diskIo["after_WriteBytes"] = sd[dm[len(dm)-1]].WriteBytes
	n, _ = net.IOCounters(true)
	for _, v := range n {
		if !strings.HasPrefix(v.Name, "lo") && !strings.HasPrefix(v.Name, "docker") {
			bytesSent = append(bytesSent, v.BytesSent)
			bytesRecv = append(bytesRecv, v.BytesRecv)
			packetsSent = append(packetsSent, v.PacketsSent)
			packetsRecv = append(packetsRecv, v.PacketsRecv)
		}
	}
	diskIo["after_lan_intraffic"] = kits.Max(bytesRecv)
	diskIo["after_lan_outtraffic"] = kits.Max(bytesSent)
	diskIo["after_lan_inpkg"] = kits.Max(packetsRecv)
	diskIo["after_lan_outpkg"] = kits.Max(packetsSent)
	return diskIo
}

func Collection() config.CollectionData {
	//获取系统监控信息
	cod := config.CollectionData{}
	cod.Platform = "agent"
	cod.Resource = "server"
	cod.Item = "system"
	// 当前时间
	tt := time.Now()
	//如果设置的间隔时间大于5分钟均以5分钟间隔进行执行
	if config.Interval > 300 {
		config.Interval = 300
	}
	//分钟级抓取数据需要进行日期格式化
	if config.Interval%60 == 0 {
		cod.NowTime = tt.Format("2006-01-02T15:04:00Z")
	} else {
		cod.NowTime = tt.Format("2006-01-02T15:04:05Z")
	}
	//获取HostId
	hostId := kits.GetHostId(config.HostIdFile)
	if hostId != "" {
		cod.HostId = hostId
	}
	//获取系统load值
	info, _ := load.Avg()
	cod.CpuLoadavg = info.Load5
	//获取cpu平均使用率
	c2, _ := cpu.Percent(0, false)
	cod.CpuUsage = c2[0]
	//获取内存信息
	m, _ := mem.VirtualMemory()
	cod.MemUsed = m.Used
	cod.MemAvailable = m.Free
	cod.MemPavailable = 100 - m.UsedPercent
	cod.MemPused = m.UsedPercent
	//获取磁盘信息
	d, _ := disk.Usage("/")
	cod.DiskUsage = d.UsedPercent
	do := PsIo()
	cod.DiskReadTraffic = (do["after_ReadBytes"] - do["before_ReadBytes"]) / config.Interval
	cod.DiskWriteTraffic = (do["after_WriteBytes"] - do["before_WriteBytes"]) / config.Interval
	//获取流量信息
	cod.LanIntraffic = (do["after_lan_intraffic"] - do["before_lan_intraffic"]) / config.Interval
	cod.LanOuttraffic = (do["after_lan_outtraffic"] - do["before_lan_outtraffic"]) / config.Interval
	cod.LanInpkg = (do["after_lan_inpkg"] - do["before_lan_inpkg"]) / config.Interval
	cod.LanOutpkg = (do["after_lan_outpkg"] - do["before_lan_outpkg"]) / config.Interval
	//获取tcp连接数
	Estab := 0
	e, err := net.Connections("inet4")
	if err != nil {
		kits.Log(err.Error(), "error", "Collection")
	} else {
		for _, k := range e {
			if k.Status == "ESTABLISHED" {
				Estab++
			}
		}
	}
	cod.TcpEstab = Estab
	//获取进程信息
	if len(config.Process) > 0 {
		Pro := make(map[string]config.ProcessInfo)
		if err != nil {
			kits.Log(err.Error(), "error", "Collection")
		} else {
			p, err := process.Processes()
			if err != nil {
				kits.Log(err.Error(), "error", "Collection")
			} else {
				for _, name := range config.Process {
					Info := config.ProcessInfo{}
					for _, pid := range p {
						n, _ := pid.Name()
						if n == name {
							Cpu, _ := pid.CPUPercent()
							Info.Cpu = Info.Cpu + Cpu
							Mem, _ := pid.MemoryPercent()
							Info.Mem = Info.Mem + Mem
							Info.Num += 1
							Pro[name] = Info
							cod.Process = Pro
						}
					}
				}
			}
		}
	}
	return cod
}

func Agent(cnf config.Config) {
	//req.Debug = 1
	JsonData := config.ApiJson{}
	Co := config.AgentConf{}
	req := requests.Requests()
	Encrypt := kits.NewEncrypt([]byte(config.CryptKey), 16)
	co := Collection()
	if co.HostId != "" {
		data := Encrypt.EncryptString(co.CollectionDataString())
		m := map[string]string{"agent_id": cnf.AgentId, "version": config.Version, "data": data}
		jsonData, _ := json.Marshal(m)
		var stringBuilder bytes.Buffer
		stringBuilder.WriteString(kits.ExportUrl())
		stringBuilder.WriteString(config.UrlPath)
		resp, err := req.PostJson(stringBuilder.String(), string(jsonData))
		if err != nil {
			kits.Log(err.Error(), "error", "Agent")
		} else {
			err = resp.Json(&JsonData)
			if err != nil {
				kits.Log(err.Error(), "error", "Agent")
			} else {
				if JsonData.Success {
					v, _ := Encrypt.DecryptString(JsonData.Data.(string))
					err = json.Unmarshal(v, &Co)
					if err != nil {
						kits.Log(err.Error(), "error", "Agent")
					} else {
						config.Process = Co.Process
						// 版本升级
						Version, _ := strconv.Atoi(config.Version)
						NewVersion, _ := strconv.Atoi(Co.Version)
						if Co.Action == "upgrade" && Version < NewVersion {
							if kits.CheckFile("/usr/bin/monitor-agent") {
								_ = os.Remove("/usr/bin/monitor-agent")
							}
							os.Exit(0)
						}
					}
				}
			}
		}
	}
}
