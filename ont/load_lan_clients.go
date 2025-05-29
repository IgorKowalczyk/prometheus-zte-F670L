package ont

import (
	"encoding/xml"
	"errors"
	"io"
	"strconv"
	"time"
)

type LanClient struct {
	HostName    string
	IPAddress   string
	IPV6Address string
	MACAddress  string
	AliasName   string
}

type LanClientsResponse struct {
	XMLName        xml.Name `xml:"ajax_response_xml_root"`
	IFERRORPARAM   string   `xml:"IF_ERRORPARAM"`
	IFERRORTYPE    string   `xml:"IF_ERRORTYPE"`
	IFERRORSTR     string   `xml:"IF_ERRORSTR"`
	IFERRORID      string   `xml:"IF_ERRORID"`
	OBJACCESSDEVID struct {
		Instances []LanClientInstance `xml:"Instance"`
	} `xml:"OBJ_ACCESSDEV_ID"`
}

type LanClientInstance struct {
	ParaName  []string `xml:"ParaName"`
	ParaValue []string `xml:"ParaValue"`
}

func (s *Session) LoadLanClients() ([]LanClient, error) {
	// Load the LAN clients
	url := s.Endpoint + "/?_type=menuData&_tag=accessdev_landevs_lua.lua&_=" + strconv.FormatInt(time.Now().Unix(), 10)
	resp, err := s.Get(url)
	if err != nil {
		return nil, err
	}
	defer func() {
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}()

	var result LanClientsResponse
	if err := xml.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	if result.IFERRORSTR != "SUCC" {
		return nil, errors.New(result.IFERRORSTR)
	}
	return result.Convert(), nil
}

func (r LanClientsResponse) Convert() []LanClient {
	var clients []LanClient
	for _, inst := range r.OBJACCESSDEVID.Instances {
		client := LanClient{}
		for i, name := range inst.ParaName {
			if i >= len(inst.ParaValue) {
				continue
			}
			val := inst.ParaValue[i]
			switch name {
			case "HostName":
				client.HostName = val
			case "IPAddress":
				client.IPAddress = val
			case "IPV6Address":
				client.IPV6Address = val
			case "MACAddress":
				client.MACAddress = val
			case "AliasName":
				client.AliasName = val
			}
		}
		clients = append(clients, client)
	}
	return clients
}
