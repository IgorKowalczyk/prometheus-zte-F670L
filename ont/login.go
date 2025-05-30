package ont

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
)

type LoginResponse struct {
	SessionToken     string `json:"sess_token"`
	LoginNeedRefresh bool   `json:"login_need_refresh"`
}

func Login(endpoint, username, password string) (*Session, error) {
	jar, _ := cookiejar.New(nil)
	session := &Session{
		Client: &http.Client{
			Jar: jar,
		},
		Endpoint: endpoint,
	}

	sessionToken, err := session.GetSessionToken()
	if err != nil {
		panic(err)
	}

	loginToken, err := session.GetLoginToken()
	if err != nil {
		panic(err)
	}

	preparedHash := sha256.New()
	preparedHash.Write([]byte(password + loginToken))

	var payload url.Values = map[string][]string{
		"action":        {"login"},
		"Username":      {username},
		"Password":      {hex.EncodeToString(preparedHash.Sum(nil))},
		"_sessionTOKEN": {sessionToken},
	}

	resp, err := session.Post(session.Endpoint+"/?_type=loginData&_tag=login_entry", "application/x-www-form-urlencoded; charset=UTF-8", strings.NewReader(payload.Encode()))

	if err != nil {
		panic(err)
	}

	defer func() {
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}()

	var result LoginResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	if result.LoginNeedRefresh {
		resp2, _ := session.Get(session.Endpoint)
		if resp2 != nil {
			io.Copy(io.Discard, resp2.Body)
			resp2.Body.Close()
		}
		return session, nil
	}

	return nil, errors.New("failed to login")
}
