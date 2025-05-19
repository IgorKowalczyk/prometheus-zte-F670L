package ont

import (
	"encoding/xml"
	"errors"
	"strconv"
	"time"
)

type LanDHCPSettings struct {
	InstID          string
	SubMask         string
	DNSServer1      string
	DNSServer2      string
	LeaseTime       int
	MaxAddress      string
	SubnetMask      string
	DnsServerSource string
	IPAddr          string
	ServerEnable    int
	MinAddress      string

	// LANDNS
	Ipv4DnsOrigin   string
	IPv4AssignLANIP string
	Ipv6DnsOrigin   string
	IPv6AssignLANIP string
}

func (s *Session) LoadLanDHCPSettings() (*LanDHCPSettings, error) {
	url := s.Endpoint + "/?_type=menuData&_tag=Localnet_LanMgrIpv4_DHCPBasicCfg_lua.lua&_=" + strconv.FormatInt(time.Now().Unix(), 10)
	resp, err := s.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result lanDHCPSettingsResponse
	if err := xml.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	if result.IFERRORSTR != "SUCC" {
		return nil, errors.New(result.IFERRORSTR)
	}
	return result.Convert(), nil
}

type lanDHCPSettingsResponse struct {
	XMLName                xml.Name `xml:"ajax_response_xml_root"`
	IFERRORPARAM           string   `xml:"IF_ERRORPARAM"`
	IFERRORTYPE            string   `xml:"IF_ERRORTYPE"`
	IFERRORSTR             string   `xml:"IF_ERRORSTR"`
	IFERRORID              string   `xml:"IF_ERRORID"`
	OBJBr0AndDhcpsHosCfgID struct {
		Instance dhcpSettingsInstance `xml:"Instance"`
	} `xml:"OBJ_Br0AndDhcpsHosCfg_ID"`
	OBJLANDNSID struct {
		Instance dhcpSettingsInstance `xml:"Instance"`
	} `xml:"OBJ_LANDNS_ID"`
}

type dhcpSettingsInstance struct {
	Params []struct {
		XMLName xml.Name
		Value   string `xml:",chardata"`
	} `xml:",any"`
}

func (inst *dhcpSettingsInstance) ToMap() map[string]string {
	m := make(map[string]string)
	var lastKey string
	for _, p := range inst.Params {
		if p.XMLName.Local == "ParaName" {
			lastKey = p.Value
		} else if p.XMLName.Local == "ParaValue" && lastKey != "" {
			m[lastKey] = p.Value
			lastKey = ""
		}
	}
	return m
}

func (r lanDHCPSettingsResponse) Convert() *LanDHCPSettings {
	m := r.OBJBr0AndDhcpsHosCfgID.Instance.ToMap()
	m2 := r.OBJLANDNSID.Instance.ToMap()
	leaseTime, _ := strconv.Atoi(m["LeaseTime"])
	serverEnable, _ := strconv.Atoi(m["ServerEnable"])
	return &LanDHCPSettings{
		InstID:          m["_InstID"],
		SubMask:         m["SubMask"],
		DNSServer1:      m["DNSServer1"],
		DNSServer2:      m["DNSServer2"],
		LeaseTime:       leaseTime,
		MaxAddress:      m["MaxAddress"],
		SubnetMask:      m["SubnetMask"],
		DnsServerSource: m["DnsServerSource"],
		IPAddr:          m["IPAddr"],
		ServerEnable:    serverEnable,
		MinAddress:      m["MinAddress"],
		Ipv4DnsOrigin:   m2["Ipv4DnsOrigin"],
		IPv4AssignLANIP: m2["IPv4AssignLANIP"],
		Ipv6DnsOrigin:   m2["Ipv6DnsOrigin"],
		IPv6AssignLANIP: m2["IPv6AssignLANIP"],
	}
}
