package ont

import (
	"encoding/xml"
	"errors"
	"io"
	"strconv"
	"time"
)

type LanDHCPHost struct {
	InstID      string
	PhyPortName string
	IPAddr      string
	ExpiredTime int
	MACAddr     string
	HostName    string
}

type LanDHCPHostsResponse struct {
	XMLName           xml.Name `xml:"ajax_response_xml_root"`
	IFERRORPARAM      string   `xml:"IF_ERRORPARAM"`
	IFERRORTYPE       string   `xml:"IF_ERRORTYPE"`
	IFERRORSTR        string   `xml:"IF_ERRORSTR"`
	IFERRORID         string   `xml:"IF_ERRORID"`
	OBJDHCPHOSTINFOID struct {
		Instances []dhcpHostInstance `xml:"Instance"`
	} `xml:"OBJ_DHCPHOSTINFO_ID"`
}

func (s *Session) LoadLanDHCPInfo() ([]LanDHCPHost, error) {
	_, _ = s.Get(s.Endpoint + "/?_type=menuView&_tag=lanMgrIpv4&Menu3Location=0&_" + strconv.FormatInt(time.Now().Unix(), 10))
	url := s.Endpoint + "/?_type=menuData&_tag=Localnet_LanMgrIpv4_DHCPHostInfo_lua.lua&_=" + strconv.FormatInt(time.Now().Unix(), 10)
	resp, err := s.Get(url)
	if err != nil {
		return nil, err
	}
	defer func() {
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}()

	var result LanDHCPHostsResponse
	if err := xml.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	if result.IFERRORSTR != "SUCC" {
		return nil, errors.New(result.IFERRORSTR)
	}
	return result.Convert(), nil
}

type dhcpHostInstance struct {
	Params []struct {
		XMLName xml.Name
		Value   string `xml:",chardata"`
	} `xml:",any"`
}

func (inst *dhcpHostInstance) ToMap() map[string]string {
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

func (r LanDHCPHostsResponse) Convert() []LanDHCPHost {
	var hosts []LanDHCPHost
	for _, inst := range r.OBJDHCPHOSTINFOID.Instances {
		m := inst.ToMap()
		expiredTime, _ := strconv.Atoi(m["ExpiredTime"])
		host := LanDHCPHost{
			InstID:      m["_InstID"],
			PhyPortName: m["PhyPortName"],
			IPAddr:      m["IPAddr"],
			ExpiredTime: expiredTime,
			MACAddr:     m["MACAddr"],
			HostName:    m["HostName"],
		}
		hosts = append(hosts, host)
	}
	return hosts
}
