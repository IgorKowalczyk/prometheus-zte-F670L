package prometheus

import (
	"cmp"
	"log"
	"os"
	"prometheus_F670L/ont"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

// ONTCollector implements the prometheus.Collector interface
type ONTCollector struct {
	session *ont.Session
}

// NewONTCollector creates a new ONT metrics collector
func NewONTCollector(session *ont.Session) *ONTCollector {
	return &ONTCollector{
		session: session,
	}
}

func mapDuplex(val int) string {
	switch val {
	case 1:
		return "half"
	case 2:
		return "full"
	default:
		return "unknown"
	}
}

func mapSpeed(val int) string {
	switch val {
	case 1:
		return "10"
	case 2:
		return "100"
	case 3:
		return "1000"
	default:
		return "0"
	}
}

// Describe implements prometheus.Collector
func (c *ONTCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- deviceInfoDesc
	ch <- cpuUsageDesc
	ch <- memoryUsageDesc
	ch <- uptimeDesc
	ch <- bytesDesc
	ch <- packetsDesc
	ch <- errorsDesc
	ch <- discardsDesc
	ch <- networkStatusDesc
	ch <- wlanClientStatusDesc
	ch <- lanClientStatusDesc
	ch <- lanDHCPHostDesc
	ch <- lanDHCPSettingsDesc

}

func sleepQuit(reaason string) {
	sleepTimeString := cmp.Or(os.Getenv("ONT_SLEEP_QUIT"), "60")
	sleepTime, _ := strconv.Atoi(sleepTimeString)

	log.Printf("[SleepQuit] %s, sleeping for %d seconds before quitting...\n", reaason, sleepTime)
	time.Sleep(time.Duration(sleepTime) * time.Second)
	log.Println("[SleepQuit] Sleep time is over, exiting...")
	os.Exit(1)
}

// Collect implements prometheus.Collector
func (c *ONTCollector) Collect(ch chan<- prometheus.Metric) {
	// Collect Device Info
	deviceInfo, err := c.session.LoadDeviceInfo()
	if err != nil {
		log.Printf("Error loading device info: %v", err)
		sleepQuit(err.Error())
		return
	}

	ch <- prometheus.MustNewConstMetric(
		deviceInfoDesc,
		prometheus.GaugeValue,
		1,
		deviceInfo.Manufacturer,
		deviceInfo.ManufacturerOui,
		deviceInfo.VersionDate,
		deviceInfo.BootVersion,
		deviceInfo.SofwareVersion,
		deviceInfo.SoftwareVersionExtended,
		deviceInfo.SerialNumber,
		deviceInfo.Model,
		deviceInfo.HardwareVersion,
	)

	// CPU Usage metrics (loop for each core)
	cpuUsages := []int{deviceInfo.CPUUsage1, deviceInfo.CPUUsage2, deviceInfo.CPUUsage3, deviceInfo.CPUUsage4}
	for i, usage := range cpuUsages {
		ch <- prometheus.MustNewConstMetric(
			cpuUsageDesc,
			prometheus.GaugeValue,
			float64(usage),
			strconv.Itoa(i+1),
		)
	}

	// Memory Usage metric
	ch <- prometheus.MustNewConstMetric(
		memoryUsageDesc,
		prometheus.GaugeValue,
		float64(deviceInfo.MemoryUsage),
	)

	// Uptime metric
	ch <- prometheus.MustNewConstMetric(
		uptimeDesc,
		prometheus.CounterValue,
		float64(deviceInfo.Uptime),
	)

	// Collect LAN Info
	lanInfo, err := c.session.LoadLanInfo()
	if err != nil {
		log.Printf("Error loading LAN info: %v", err)
		sleepQuit(err.Error())
		return
	}

	// Network traffic metrics (correct direction)
	ch <- prometheus.MustNewConstMetric(
		bytesDesc,
		prometheus.CounterValue,
		float64(lanInfo.BytesIn),
		"in",
	)
	ch <- prometheus.MustNewConstMetric(
		bytesDesc,
		prometheus.CounterValue,
		float64(lanInfo.BytesOut),
		"out",
	)

	// Packet metrics (loop for unicast/multicast, in/out)
	packetMetrics := []struct {
		desc  *prometheus.Desc
		value int
		dir   string
		ptype string
	}{
		{packetsDesc, lanInfo.PacketsUnicastIn, "in", "unicast"},
		{packetsDesc, lanInfo.PacketsUnicastOut, "out", "unicast"},
		{packetsDesc, lanInfo.PacketsMulticastIn, "in", "multicast"},
		{packetsDesc, lanInfo.PacketsMulticastOut, "out", "multicast"},
	}
	for _, m := range packetMetrics {
		ch <- prometheus.MustNewConstMetric(
			m.desc,
			prometheus.CounterValue,
			float64(m.value),
			m.dir, m.ptype,
		)
	}

	// Error and discard metrics
	ch <- prometheus.MustNewConstMetric(errorsDesc, prometheus.CounterValue, float64(lanInfo.PacketsErrorIn), "in")
	ch <- prometheus.MustNewConstMetric(errorsDesc, prometheus.CounterValue, float64(lanInfo.PacketsErrorOut), "out")
	ch <- prometheus.MustNewConstMetric(discardsDesc, prometheus.CounterValue, float64(lanInfo.PacketsDiscardedIn), "in")
	ch <- prometheus.MustNewConstMetric(discardsDesc, prometheus.CounterValue, float64(lanInfo.PacketsDiscardedOut), "out")

	// Status metric
	duplexInt, err := strconv.Atoi(lanInfo.Duplex)
	if err != nil {
		duplexInt = 0
	}
	ch <- prometheus.MustNewConstMetric(
		networkStatusDesc,
		prometheus.GaugeValue,
		float64(lanInfo.Status),
		mapSpeed(lanInfo.Speed),
		mapDuplex(duplexInt),
	)

	// WLAN Info
	wlanInfo, err := c.session.LoadWlanClientsInfo()
	if err != nil {
		log.Printf("Error loading WLAN info: %v", err)
	} else {
		apMap := make(map[string]string)
		for _, client := range wlanInfo.Clients {
			essid := apMap[client.AliasName]
			ch <- prometheus.MustNewConstMetric(
				wlanClientStatusDesc,
				prometheus.GaugeValue,
				1,
				client.HostName,
				client.IPAddress,
				client.IPV6Address,
				client.MACAddress,
				client.AliasName,
				essid,
				strconv.Itoa(client.RSSI),
				strconv.Itoa(client.TxRate),
				strconv.Itoa(client.RxRate),
				strconv.Itoa(client.SNR),
				strconv.Itoa(client.NOISE),
				strconv.Itoa(client.LinkTime),
				client.CurrentMode,
				strconv.Itoa(client.MCS),
				client.BAND,
			)
		}
	}

	// Collect LAN Clients
	lanClients, err := c.session.LoadLanClients()
	if err != nil {
		log.Printf("Error loading LAN clients: %v", err)
	} else {
		for _, client := range lanClients {
			ch <- prometheus.MustNewConstMetric(
				lanClientStatusDesc,
				prometheus.GaugeValue,
				1,
				client.HostName, client.IPAddress, client.IPV6Address, client.MACAddress, client.AliasName,
			)
		}
	}

	// WAN Internet Status
	wanStatus, err := c.session.LoadWanInternetStatus()
	if err != nil {
		log.Printf("Error loading WAN Internet status: %v", err)
	} else {
		ch <- prometheus.MustNewConstMetric(
			wanInternetStatusDesc,
			prometheus.GaugeValue,
			1,
			wanStatus.ConnTrigger,
			strconv.Itoa(wanStatus.UpTime),
			strconv.Itoa(wanStatus.IsNAT),
			wanStatus.ConnError,
			wanStatus.XdslMode,
			wanStatus.WanType,
			wanStatus.WANCName,
			wanStatus.IpMode,
			wanStatus.TransType,
			wanStatus.PPPoeServiceName,
			wanStatus.Mode,
			strconv.Itoa(wanStatus.Uplink),
			strconv.Itoa(wanStatus.PageType),
			strconv.Itoa(wanStatus.VlanEnable),
			wanStatus.StrServList,
			wanStatus.ConnStatus6,
			wanStatus.InstID,
			strconv.Itoa(wanStatus.Enable),
			strconv.Itoa(wanStatus.DSCP),
			strconv.Itoa(wanStatus.Priority),
			strconv.Itoa(wanStatus.VLANID),
			wanStatus.SubnetMask,
			wanStatus.AuthType,
			strconv.Itoa(wanStatus.MTU),
			wanStatus.DNS1,
			wanStatus.DNS3,
			wanStatus.GateWay,
			wanStatus.WorkIFMac,
			wanStatus.ServList,
			wanStatus.LinkMode,
			strconv.Itoa(wanStatus.IsDefGW),
			wanStatus.IPAddress,
			wanStatus.DNS2,
			strconv.Itoa(wanStatus.EnablePassThrough),
			wanStatus.ConnStatus,
		)
	}

	wlanAPs, err := c.session.LoadWlanInfo()
	if err != nil {
		log.Printf("Error loading WLAN info: %v", err)
	} else {
		for _, ap := range wlanAPs {
			ch <- prometheus.MustNewConstMetric(
				wlanAPStatusDesc,
				prometheus.GaugeValue,
				1,
				ap.InstID,
				ap.Alias,
				ap.ESSID,
				ap.BSSID,
				ap.Band,
				ap.Enable,
				ap.Channel,
				ap.Encryption,
				ap.TotalBytesSent,
				ap.TotalBytesReceived,
			)
		}
	}

	lanDHCPHosts, err := c.session.LoadLanDHCPInfo()
	if err != nil {
		log.Printf("Error loading LAN DHCP hosts: %v", err)
	} else {
		for _, host := range lanDHCPHosts {
			ch <- prometheus.MustNewConstMetric(
				lanDHCPHostDesc,
				prometheus.GaugeValue,
				1,
				host.InstID,
				host.PhyPortName,
				host.IPAddr,
				strconv.Itoa(host.ExpiredTime),
				host.MACAddr,
				host.HostName,
			)
		}
	}

	lanDHCPSettings, err := c.session.LoadLanDHCPSettings()
	if err != nil {
		log.Printf("Error loading LAN DHCP settings: %v", err)
	} else {
		ch <- prometheus.MustNewConstMetric(
			lanDHCPSettingsDesc,
			prometheus.GaugeValue,
			1,
			lanDHCPSettings.InstID,
			lanDHCPSettings.SubMask,
			lanDHCPSettings.DNSServer1,
			lanDHCPSettings.DNSServer2,
			strconv.Itoa(lanDHCPSettings.LeaseTime),
			lanDHCPSettings.MaxAddress,
			lanDHCPSettings.SubnetMask,
			lanDHCPSettings.DnsServerSource,
			lanDHCPSettings.IPAddr,
			strconv.Itoa(lanDHCPSettings.ServerEnable),
			lanDHCPSettings.MinAddress,
			lanDHCPSettings.Ipv4DnsOrigin,
			lanDHCPSettings.IPv4AssignLANIP,
			lanDHCPSettings.Ipv6DnsOrigin,
			lanDHCPSettings.IPv6AssignLANIP,
		)
	}
}
