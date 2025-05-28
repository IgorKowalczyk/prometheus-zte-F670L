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
	Enable             string
	Channel            string
	Encryption         string
	TotalBytesSent     string
	TotalBytesReceived string
}

// Helper struct to parse alternating ParaName/ParaValue pairs
type wlanAPInstance struct {
	Params []struct {
		XMLName xml.Name
		Value   string `xml:",chardata"`
	} `xml:",any"`
}

// Convert Params slice to a map
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
	_, _ = s.Get(s.Endpoint + "/?_type=menuView&_tag=localNetStatus&Menu3Location=0&_" + strconv.FormatInt(time.Now().Unix(), 10))
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
	// Build band map from OBJ_WLANSETTING_ID
	bandMap := make(map[string]string)
	for _, inst := range r.OBJWLANSETTINGID.Instances {
		m := inst.ToMap()
		instID := m["_InstID"]
		band := m["Band"]
		if instID != "" && band != "" {
			bandMap[instID] = band
		}
	}

	var aps []WlanAP
	for _, inst := range r.OBJWLANAPID.Instances {
		m := inst.ToMap()
		ap := WlanAP{
			InstID:     m["_InstID"],
			Alias:      m["Alias"],
			ESSID:      m["ESSID"],
			Enable:     m["Enable"],
			Encryption: m["WPAEncryptType"],
		}

		if enc := m["11iEncryptType"]; enc != "" {
			ap.Encryption = enc
		}
		wlanViewName := m["WLANViewName"]

		// Find BSSID, Channel, TotalBytesSent, TotalBytesReceived from OBJ_WLANCONFIGDRV_ID
		for _, drv := range r.OBJWLANCONFIGDRVID.Instances {
			drvMap := drv.ToMap()
			if drvMap["_InstID"] == ap.InstID {
				ap.BSSID = drvMap["Bssid"]
				ap.Channel = drvMap["ChannelInUsed"]
				ap.TotalBytesSent = drvMap["TotalBytesSent"]
				ap.TotalBytesReceived = drvMap["TotalBytesReceived"]
				break
			}
		}
		ap.Band = bandMap[wlanViewName]
		aps = append(aps, ap)
	}
	return aps
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
