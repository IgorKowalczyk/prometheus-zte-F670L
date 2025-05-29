package ont

import (
	"encoding/xml"
	"errors"
	"io"
	"strconv"
	"time"
)

type DeviceInfo struct {
	Manufacturer            string
	ManufacturerOui         string
	VersionDate             string
	BootVersion             string
	SofwareVersion          string
	SoftwareVersionExtended string
	SerialNumber            string
	Model                   string
	HardwareVersion         string

	CPUUsage1   int
	CPUUsage2   int
	CPUUsage3   int
	CPUUsage4   int
	MemoryUsage int

	Uptime int
}

type InformationResponse struct {
	XMLName      xml.Name `xml:"ajax_response_xml_root"`
	Text         string   `xml:",chardata"`
	IFERRORPARAM string   `xml:"IF_ERRORPARAM"`
	IFERRORTYPE  string   `xml:"IF_ERRORTYPE"`
	IFERRORSTR   string   `xml:"IF_ERRORSTR"`
	IFERRORID    string   `xml:"IF_ERRORID"`
	OBJDEVINFOID struct {
		Text     string `xml:",chardata"`
		Instance struct {
			Text      string   `xml:",chardata"`
			ParaName  []string `xml:"ParaName"`
			ParaValue []string `xml:"ParaValue"`
		} `xml:"Instance"`
	} `xml:"OBJ_DEVINFO_ID"`
	OBJCPUMEMUSAGEID struct {
		Text     string `xml:",chardata"`
		Instance struct {
			Text      string   `xml:",chardata"`
			ParaName  []string `xml:"ParaName"`
			ParaValue []string `xml:"ParaValue"`
		} `xml:"Instance"`
	} `xml:"OBJ_CPUMEMUSAGE_ID"`
	OBJPOWERONTIMEID struct {
		Text     string `xml:",chardata"`
		Instance struct {
			Text      string   `xml:",chardata"`
			ParaName  []string `xml:"ParaName"`
			ParaValue []string `xml:"ParaValue"`
		} `xml:"Instance"`
	} `xml:"OBJ_POWERONTIME_ID"`
}

func (s *Session) LoadDeviceInfo() (*DeviceInfo, error) {
	_, _ = s.Get(s.Endpoint + "/?_type=menuView&_tag=statusMgr&Menu3Location=0&_=" + strconv.FormatInt(time.Now().Unix(), 10))
	resp, err := s.Get(s.Endpoint + "/?_type=menuData&_tag=devmgr_statusmgr_lua.lua&_=" + strconv.FormatInt(time.Now().Unix(), 10))

	if err != nil {
		return nil, err
	}

	defer func() {
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}()

	var result InformationResponse
	if err := xml.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	if result.IFERRORSTR == "SessionTimeout" {
		return nil, errors.New("session timeout")
	}

	return result.Convert(), nil
}

func (result InformationResponse) Convert() *DeviceInfo {
	var deviceInfo DeviceInfo
	for i, name := range result.OBJDEVINFOID.Instance.ParaName {
		if i < len(result.OBJDEVINFOID.Instance.ParaValue) {
			switch name {
			case "ManuFacturer":
				deviceInfo.Manufacturer = result.OBJDEVINFOID.Instance.ParaValue[i]
			case "ManuFacturerOui":
				deviceInfo.ManufacturerOui = result.OBJDEVINFOID.Instance.ParaValue[i]
			case "VerDate":
				deviceInfo.VersionDate = result.OBJDEVINFOID.Instance.ParaValue[i]
			case "BootVer":
				deviceInfo.BootVersion = result.OBJDEVINFOID.Instance.ParaValue[i]
			case "SoftwareVer":
				deviceInfo.SofwareVersion = result.OBJDEVINFOID.Instance.ParaValue[i]
			case "SoftwareVerExtent":
				deviceInfo.SoftwareVersionExtended = result.OBJDEVINFOID.Instance.ParaValue[i]
			case "SerialNumber":
				deviceInfo.SerialNumber = result.OBJDEVINFOID.Instance.ParaValue[i]
			case "ModelName":
				deviceInfo.Model = result.OBJDEVINFOID.Instance.ParaValue[i]
			case "HardwareVer":
				deviceInfo.HardwareVersion = result.OBJDEVINFOID.Instance.ParaValue[i]
			}
		}
	}

	for i, name := range result.OBJCPUMEMUSAGEID.Instance.ParaName {
		if i < len(result.OBJCPUMEMUSAGEID.Instance.ParaValue) {

			switch name {
			case "CpuUsage1":
				deviceInfo.CPUUsage1, _ = strconv.Atoi(result.OBJCPUMEMUSAGEID.Instance.ParaValue[i])
			case "CpuUsage2":
				deviceInfo.CPUUsage2, _ = strconv.Atoi(result.OBJCPUMEMUSAGEID.Instance.ParaValue[i])
			case "CpuUsage3":
				deviceInfo.CPUUsage3, _ = strconv.Atoi(result.OBJCPUMEMUSAGEID.Instance.ParaValue[i])
			case "CpuUsage4":
				deviceInfo.CPUUsage4, _ = strconv.Atoi(result.OBJCPUMEMUSAGEID.Instance.ParaValue[i])
			case "MemUsage":
				deviceInfo.MemoryUsage, _ = strconv.Atoi(result.OBJCPUMEMUSAGEID.Instance.ParaValue[i])
			}
		}
	}

	for i, name := range result.OBJPOWERONTIMEID.Instance.ParaName {
		if i < len(result.OBJPOWERONTIMEID.Instance.ParaValue) {
			switch name {
			case "PowerOnTime":
				deviceInfo.Uptime, _ = strconv.Atoi(result.OBJPOWERONTIMEID.Instance.ParaValue[i])
			}
		}
	}

	return &deviceInfo
}
