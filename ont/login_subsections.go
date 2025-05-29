package ont

import (
	"encoding/json"
	"encoding/xml"
	"io"
)

type SessionTokenResponse struct {
	LockingTime  int    `json:"lockingTime"`
	LoginErrMsg  string `json:"loginErrMsg"`
	PromptMsg    string `json:"promptMsg"`
	SessionToken string `json:"sess_token"`
}

func (s *Session) GetSessionToken() (string, error) {
	resp, err := s.Get(s.Endpoint + "/?_type=loginData&_tag=login_entry")

	if err != nil {
		return "", err
	}

	var result SessionTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}
	return result.SessionToken, nil
}

type LoginToken struct {
	XMLName xml.Name `xml:"ajax_response_xml_root"`
	Value   string   `xml:",chardata"`
}

func (s *Session) GetLoginToken() (string, error) {
	resp, err := s.Get(s.Endpoint + "/?_type=loginData&_tag=login_token")
	if err != nil {
		return "", err
	}
	defer func() {
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}()

	var result LoginToken
	if err := xml.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	return result.Value, nil
}
