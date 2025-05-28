package ont

import (
	"encoding/xml"
	"errors"
	"io"
	"strconv"
	"time"
)

type WlanClient struct {
	InstID      string
	AliasName   string
	RxRate      int
	HostName    string
	RSSI        int
	LinkTime    int
	TxRate      int
	NOISE       int
	IPAddress   string
	IPV6Address string
	SNR         int
	MACAddress  string
	CurrentMode string
	MCS         int
	BAND        string
}

type WlanInfo struct {
	Clients []WlanClient
}

func (s *Session) LoadWlanClientsInfo() (*WlanInfo, error) {
	url := s.Endpoint + "/?_type=menuData&_tag=wlan_client_stat_lua.lua&_=" + strconv.FormatInt(time.Now().Unix(), 10)
	resp, err := s.Get(url)

	if err != nil {
		return nil, err
	}

	defer func() {
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}()

	var result wlanInfoResponse
	if err := xml.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	if result.IFERRORSTR != "SUCC" {
		return nil, errors.New(result.IFERRORSTR)
	}
	return result.Convert(), nil
}

type wlanInfoResponse struct {
	XMLName      xml.Name `xml:"ajax_response_xml_root"`
	IFERRORPARAM string   `xml:"IF_ERRORPARAM"`
	IFERRORTYPE  string   `xml:"IF_ERRORTYPE"`
	IFERRORSTR   string   `xml:"IF_ERRORSTR"`
	IFERRORID    string   `xml:"IF_ERRORID"`
	OBJWLANADID  struct {
		Instances []wlanClientInstance `xml:"Instance"`
	} `xml:"OBJ_WLAN_AD_ID"`
	OBJWLANAPID struct {
		Instances []wlanAPInstance `xml:"Instance"`
	} `xml:"OBJ_WLANAP_ID"`
}

type wlanClientInstance struct {
	ParaName  []string `xml:"ParaName"`
	ParaValue []string `xml:"ParaValue"`
}

func (r wlanInfoResponse) Convert() *WlanInfo {
	info := &WlanInfo{}

	for _, inst := range r.OBJWLANADID.Instances {
		client := WlanClient{}
		for i, name := range inst.ParaName {
			if i >= len(inst.ParaValue) {
				continue
			}
			val := inst.ParaValue[i]
			switch name {
			case "_InstID":
				client.InstID = val
			case "AliasName":
				client.AliasName = val
			case "RxRate":
				client.RxRate, _ = strconv.Atoi(val)
			case "HostName":
				client.HostName = val
			case "RSSI":
				client.RSSI, _ = strconv.Atoi(val)
			case "LinkTime":
				client.LinkTime, _ = strconv.Atoi(val)
			case "TxRate":
				client.TxRate, _ = strconv.Atoi(val)
			case "NOISE":
				client.NOISE, _ = strconv.Atoi(val)
			case "IPAddress":
				client.IPAddress = val
			case "IPV6Address":
				client.IPV6Address = val
			case "SNR":
				client.SNR, _ = strconv.Atoi(val)
			case "MACAddress":
				client.MACAddress = val
			case "CurrentMode":
				client.CurrentMode = val
			case "MCS":
				client.MCS, _ = strconv.Atoi(val)
			case "BAND":
				client.BAND = val
			}
		}
		info.Clients = append(info.Clients, client)
	}
	return info
}
