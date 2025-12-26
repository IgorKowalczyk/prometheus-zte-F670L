package ont

import (
	"encoding/xml"
	"errors"
	"io"
	"strconv"
	"time"
)

type WlanAP struct {
	InstID             string
	Alias              string
	ESSID              string
	BSSID              string
	Band               string
	Encryption         string
	Enable             bool
	Channel            int
	TotalBytesSent     uint64
	TotalBytesReceived uint64
}

type wlanAPInstance struct {
	Params []struct {
		XMLName xml.Name
		Value   string `xml:",chardata"`
	} `xml:",any"`
}

type wlanAPsResponse struct {
	XMLName      xml.Name `xml:"ajax_response_xml_root"`
	IFERRORPARAM string   `xml:"IF_ERRORPARAM"`
	IFERRORTYPE  string   `xml:"IF_ERRORTYPE"`
	IFERRORSTR   string   `xml:"IF_ERRORSTR"`
	IFERRORID    string   `xml:"IF_ERRORID"`
	OBJWLANAPID  struct {
		Instances []wlanAPInstance `xml:"Instance"`
	} `xml:"OBJ_WLANAP_ID"`
	OBJWLANCONFIGDRVID struct {
		Instances []wlanAPInstance `xml:"Instance"`
	} `xml:"OBJ_WLANCONFIGDRV_ID"`
	OBJWLANSETTINGID struct {
		Instances []wlanAPInstance `xml:"Instance"`
	} `xml:"OBJ_WLANSETTING_ID"`
}

func (inst *wlanAPInstance) ToMap() map[string]string {
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

func (s *Session) LoadWlanInfo() ([]WlanAP, error) {
	respMenu, _ := s.Get(s.Endpoint + "/?_type=menuView&_tag=localNetStatus&_=" + strconv.FormatInt(time.Now().Unix(), 10))
	if respMenu != nil {
		io.Copy(io.Discard, respMenu.Body)
		respMenu.Body.Close()
	}

	url := s.Endpoint + "/?_type=menuData&_tag=wlan_wlanstatus_lua.lua&_=" + strconv.FormatInt(time.Now().Unix(), 10)
	resp, err := s.Get(url)
	if err != nil {
		return nil, err
	}
	defer func() {
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}()

	var result wlanAPsResponse
	if err := xml.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	if result.IFERRORSTR != "SUCC" {
		return nil, errors.New(result.IFERRORSTR)
	}

	return result.Convert(), nil
}

func (r *wlanAPsResponse) Convert() []WlanAP {
	bandMap := make(map[string]string)

	for _, inst := range r.OBJWLANSETTINGID.Instances {
		m := inst.ToMap()
		if id := m["_InstID"]; id != "" {
			bandMap[id] = m["Band"]
		}
	}

	var aps []WlanAP

	for _, inst := range r.OBJWLANAPID.Instances {
		m := inst.ToMap()

		ap := WlanAP{
			InstID:     m["_InstID"],
			Alias:      m["Alias"],
			ESSID:      m["ESSID"],
			Encryption: m["WPAEncryptType"],
			Enable:     m["Enable"] == "1",
		}

		if enc := m["11iEncryptType"]; enc != "" {
			ap.Encryption = enc
		}

		wlanViewName := m["WLANViewName"]

		for _, drv := range r.OBJWLANCONFIGDRVID.Instances {
			drvMap := drv.ToMap()
			if drvMap["_InstID"] == ap.InstID {
				ap.BSSID = drvMap["Bssid"]

				ap.Channel, _ = strconv.Atoi(drvMap["ChannelInUsed"])
				ap.TotalBytesSent, _ = strconv.ParseUint(drvMap["TotalBytesSent"], 10, 64)
				ap.TotalBytesReceived, _ = strconv.ParseUint(drvMap["TotalBytesReceived"], 10, 64)
				break
			}
		}

		ap.Band = bandMap[wlanViewName]
		aps = append(aps, ap)
	}

	return aps
}
