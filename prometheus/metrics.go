package prometheus

import "github.com/prometheus/client_golang/prometheus"

var (
	// Device Info metrics
	deviceInfoDesc = prometheus.NewDesc(
		"ont_device_info",
		"Device information",
		[]string{"manufacturer", "manufacturer_oui", "version_date", "boot_version",
			"software_version", "software_version_extended", "serial_number",
			"model", "hardware_version"},
		nil,
	)
	cpuUsageDesc = prometheus.NewDesc(
		"ont_usage_cpu",
		"CPU usage percentage per core",
		[]string{"core"},
		nil,
	)
	memoryUsageDesc = prometheus.NewDesc(
		"ont_usage_memory",
		"Memory usage percentage",
		nil,
		nil,
	)
	uptimeDesc = prometheus.NewDesc(
		"ont_device_uptime",
		"Device uptime in seconds",
		nil,
		nil,
	)

	// LAN Info metrics
	bytesDesc = prometheus.NewDesc(
		"ont_octets_total",
		"Number of bytes transmitted/received",
		[]string{"direction"},
		nil,
	)
	packetsDesc = prometheus.NewDesc(
		"ont_packets_total",
		"Number of packets transmitted/received",
		[]string{"direction", "type"},
		nil,
	)
	errorsDesc = prometheus.NewDesc(
		"ont_packets_errors_total",
		"Number of network errors",
		[]string{"direction"},
		nil,
	)
	discardsDesc = prometheus.NewDesc(
		"ont_packets_discards_total",
		"Number of discarded packets",
		[]string{"direction"},
		nil,
	)
	networkStatusDesc = prometheus.NewDesc(
		"ont_ethernet_status",
		"Ethernet interface status. 10/100/1000 Mbps, full/half duplex",
		[]string{"speed", "duplex"},
		nil,
	)

	// WLAN Info metrics
	wlanClientStatusDesc = prometheus.NewDesc(
		"ont_wlan_client_status",
		"WLAN Client Status",
		[]string{"hostname", "ip", "ipv6", "mac", "alias", "essid", "rssi", "tx_rate", "rx_rate", "snr", "noise", "link_time", "mode", "mcs", "band"},
		nil,
	)

	// LAN Client metrics
	lanClientStatusDesc = prometheus.NewDesc(
		"ont_lan_client_status",
		"LAN Client Status",
		[]string{"hostname", "ip", "ipv6", "mac", "alias"}, nil,
	)

	// WAN Internet Status metrics
	wanInternetStatusDesc = prometheus.NewDesc(
		"ont_wan_internet_status",
		"WAN Internet status info (all fields as labels, value is 1 if present)",
		[]string{
			"conn_trigger", "uptime", "is_nat", "conn_error", "xdsl_mode", "wan_type", "wan_cname", "ip_mode", "trans_type", "pppoe_service_name", "mode", "uplink", "page_type", "vlan_enable", "str_serv_list", "conn_status6", "inst_id", "enable", "dscp", "priority", "vlanid", "subnet_mask", "auth_type", "mtu", "dns1", "dns3", "gateway", "work_if_mac", "serv_list", "link_mode", "is_def_gw", "ip_address", "dns2", "enable_pass_through", "conn_status",
		},
		nil,
	)

	wlanAPStatusDesc = prometheus.NewDesc(
		"ont_wlan_ap_status",
		"WLAN AP status and statistics",
		[]string{
			"inst_id", "alias", "essid", "bssid", "band", "enable", "channel", "encryption", "bytes_sent", "bytes_received",
		},
		nil,
	)

	lanDHCPHostDesc = prometheus.NewDesc(
		"ont_lan_dhcp_host",
		"DHCP host info from ONT",
		[]string{"inst_id", "phy_port_name", "ip_addr", "expired_time", "mac_addr", "host_name"},
		nil,
	)

	lanDHCPSettingsDesc = prometheus.NewDesc(
		"ont_lan_dhcp_settings",
		"DHCP server settings from ONT",
		[]string{
			"inst_id", "sub_mask", "dns_server1", "dns_server2", "lease_time", "max_address", "subnet_mask",
			"dns_server_source", "ip_addr", "server_enable", "min_address",
			"ipv4_dns_origin", "ipv4_assign_lan_ip", "ipv6_dns_origin", "ipv6_assign_lan_ip",
		},
		nil,
	)
)
