package ont

import (
	"encoding/xml"
	"errors"
	"io"
	"strconv"
	"time"
)

type LanInfo struct {
	PacketsDiscardedIn  int
	PacketsDiscardedOut int

	PacketsErrorIn  int
	PacketsErrorOut int

	PacketsMulticastIn  int
	PacketsMulticastOut int

	PacketsUnicastIn  int
	PacketsUnicastOut int

	BytesIn  int
	BytesOut int

	PacketsIn  int
	PacketsOut int

	Status int
	Duplex string
	Speed  int
}

type LanInfoResponse struct {
	XMLName                 xml.Name `xml:"ajax_response_xml_root"`
	Text                    string   `xml:",chardata"`
	IFERRORPARAM            string   `xml:"IF_ERRORPARAM"`
	IFERRORTYPE             string   `xml:"IF_ERRORTYPE"`
	IFERRORSTR              string   `xml:"IF_ERRORSTR"`
	IFERRORID               string   `xml:"IF_ERRORID"`
	OBJPONPORTBASICSTATUSID struct {
		Text     string `xml:",chardata"`
		Instance struct {
			Text      string   `xml:",chardata"`
			ParaName  []string `xml:"ParaName"`
			ParaValue []string `xml:"ParaValue"`
		} `xml:"Instance"`
	} `xml:"OBJ_PON_PORT_BASIC_STATUS_ID"`
}

func (s *Session) LoadLanInfo() (*LanInfo, error) {
	// Trigger the menu to load the lan info
	respMenu, _ := s.Get(s.Endpoint + "/?_type=menuView&_tag=localNetStatus&Menu3Location=0&_" + strconv.FormatInt(time.Now().Unix(), 10))
	if respMenu != nil {
		io.Copy(io.Discard, respMenu.Body)
		respMenu.Body.Close()
	}

	// Load the lan info
	resp, err := s.Get(s.Endpoint + "/?_type=menuData&_tag=status_lan_info_lua.lua&_=" + strconv.FormatInt(time.Now().Unix(), 10))

	if err != nil {
		return nil, err
	}

	defer func() {
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}()

	var result LanInfoResponse
	if err := xml.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	if result.IFERRORSTR == "SessionTimeout" {
		return nil, errors.New("session timeout")
	}

	return result.Convert(), nil
}

func (result LanInfoResponse) Convert() *LanInfo {
	var lanInfos []LanInfo
	names := result.OBJPONPORTBASICSTATUSID.Instance.ParaName
	values := result.OBJPONPORTBASICSTATUSID.Instance.ParaValue

	const fieldsPerIface = 16
	for i := 0; i+fieldsPerIface <= len(values); i += fieldsPerIface {
		var lanInfo LanInfo
		for j := 0; j < fieldsPerIface; j++ {
			name := names[j]
			val := values[i+j]
			switch name {
			case "InDiscard":
				lanInfo.PacketsDiscardedIn, _ = strconv.Atoi(val)
			case "OutDiscard":
				lanInfo.PacketsDiscardedOut, _ = strconv.Atoi(val)
			case "InError":
				lanInfo.PacketsErrorIn, _ = strconv.Atoi(val)
			case "OutError":
				lanInfo.PacketsErrorOut, _ = strconv.Atoi(val)
			case "InMulticast":
				lanInfo.PacketsMulticastIn, _ = strconv.Atoi(val)
			case "OutMulticast":
				lanInfo.PacketsMulticastOut, _ = strconv.Atoi(val)
			case "InUnicast":
				lanInfo.PacketsUnicastIn, _ = strconv.Atoi(val)
			case "OutUnicast":
				lanInfo.PacketsUnicastOut, _ = strconv.Atoi(val)
			case "InBytes":
				lanInfo.BytesIn, _ = strconv.Atoi(val)
			case "OutBytes":
				lanInfo.BytesOut, _ = strconv.Atoi(val)
			case "InPkts":
				lanInfo.PacketsIn, _ = strconv.Atoi(val)
			case "OutPkts":
				lanInfo.PacketsOut, _ = strconv.Atoi(val)
			case "Status":
				lanInfo.Status, _ = strconv.Atoi(val)
			case "Duplex":
				lanInfo.Duplex = val
			case "Speed":
				lanInfo.Speed, _ = strconv.Atoi(val)
			}
		}
		lanInfos = append(lanInfos, lanInfo)
	}

	if len(lanInfos) > 0 {
		return &lanInfos[0]
	}
	return &LanInfo{}
}
