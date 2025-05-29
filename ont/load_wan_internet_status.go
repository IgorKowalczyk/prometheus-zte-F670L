package ont

import (
	"encoding/xml"
	"errors"
	"io"
	"strconv"
	"time"
)

type WanInternetStatus struct {
	ConnTrigger       string
	UpTime            int
	IsNAT             int
	UserName          string
	ConnError         string
	XdslMode          string
	WanType           string
	WANCName          string
	IpMode            string
	TransType         string
	PPPoeServiceName  string
	Mode              string
	Uplink            int
	PageType          int
	VlanEnable        int
	StrServList       string
	ConnStatus6       string
	InstID            string
	Enable            int
	DSCP              int
	Priority          int
	VLANID            int
	SubnetMask        string
	AuthType          string
	MTU               int
	DNS1              string
	DNS3              string
	GateWay           string
	WorkIFMac         string
	ServList          string
	LinkMode          string
	IsDefGW           int
	Password          string
	IPAddress         string
	DNS2              string
	EnablePassThrough int
	ConnStatus        string
}

type wanInternetStatusResponse struct {
	XMLName      xml.Name `xml:"ajax_response_xml_root"`
	IFERRORPARAM string   `xml:"IF_ERRORPARAM"`
	IFERRORTYPE  string   `xml:"IF_ERRORTYPE"`
	IFERRORSTR   string   `xml:"IF_ERRORSTR"`
	IFERRORID    string   `xml:"IF_ERRORID"`
	IDWANCONFIG  struct {
		Instance wanInternetStatusInstance `xml:"Instance"`
	} `xml:"ID_WAN_COMFIG"`
}

type wanInternetStatusInstance struct {
	ParaName  []string `xml:"ParaName"`
	ParaValue []string `xml:"ParaValue"`
}

func (s *Session) LoadWanInternetStatus() (*WanInternetStatus, error) {
	_, _ = s.Get(s.Endpoint + "/?_type=menuView&_tag=ethWanStatus&Menu3Location=0&_" + strconv.FormatInt(time.Now().Unix(), 10))
	url := s.Endpoint + "/?_type=menuData&_tag=wan_internetstatus_lua.lua&TypeUplink=2&pageType=1&_=" + strconv.FormatInt(time.Now().Unix(), 10)
	resp, err := s.Get(url)

	if err != nil {
		return nil, err
	}

	defer func() {
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}()

	var result wanInternetStatusResponse
	if err := xml.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	if result.IFERRORSTR != "SUCC" {
		return nil, errors.New(result.IFERRORSTR)
	}
	return result.Convert(), nil
}

func (r wanInternetStatusResponse) Convert() *WanInternetStatus {
	s := WanInternetStatus{}
	for i, name := range r.IDWANCONFIG.Instance.ParaName {
		if i >= len(r.IDWANCONFIG.Instance.ParaValue) {
			continue
		}
		val := r.IDWANCONFIG.Instance.ParaValue[i]
		switch name {
		case "ConnTrigger":
			s.ConnTrigger = val
		case "UpTime":
			s.UpTime, _ = strconv.Atoi(val)
		case "IsNAT":
			s.IsNAT, _ = strconv.Atoi(val)
		case "UserName":
			s.UserName = val
		case "ConnError":
			s.ConnError = val
		case "xdslMode":
			s.XdslMode = val
		case "wantype":
			s.WanType = val
		case "WANCName":
			s.WANCName = val
		case "IpMode":
			s.IpMode = val
		case "TransType":
			s.TransType = val
		case "PPPoeServiceName":
			s.PPPoeServiceName = val
		case "mode":
			s.Mode = val
		case "uplink":
			s.Uplink, _ = strconv.Atoi(val)
		case "pageType":
			s.PageType, _ = strconv.Atoi(val)
		case "VlanEnable":
			s.VlanEnable, _ = strconv.Atoi(val)
		case "StrServList":
			s.StrServList = val
		case "ConnStatus6":
			s.ConnStatus6 = val
		case "_InstID":
			s.InstID = val
		case "Enable":
			s.Enable, _ = strconv.Atoi(val)
		case "DSCP":
			s.DSCP, _ = strconv.Atoi(val)
		case "Priority":
			s.Priority, _ = strconv.Atoi(val)
		case "VLANID":
			s.VLANID, _ = strconv.Atoi(val)
		case "SubnetMask":
			s.SubnetMask = val
		case "AuthType":
			s.AuthType = val
		case "MTU":
			s.MTU, _ = strconv.Atoi(val)
		case "DNS1":
			s.DNS1 = val
		case "DNS3":
			s.DNS3 = val
		case "GateWay":
			s.GateWay = val
		case "WorkIFMac":
			s.WorkIFMac = val
		case "ServList":
			s.ServList = val
		case "linkMode":
			s.LinkMode = val
		case "IsDefGW":
			s.IsDefGW, _ = strconv.Atoi(val)
		case "Password":
			s.Password = val
		case "IPAddress":
			s.IPAddress = val
		case "DNS2":
			s.DNS2 = val
		case "EnablePassThrough":
			s.EnablePassThrough, _ = strconv.Atoi(val)
		case "ConnStatus":
			s.ConnStatus = val
		}
	}
	return &s
}
