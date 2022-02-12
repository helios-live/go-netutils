package netutils

import (
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/alex-eftimie/utils"
)

// ApiCreateServer is the response given by HTTP-Proxy on api servercreation
type ApiCreateServer struct {
	AuthToken  string
	Email      string
	Username   string
	Password   string
	Host       string
	Port       int
	ID         string
	Bandwidth  *string           `json:",omitempty"`
	MaxThreads *int              `json:",omitempty"`
	MaxIPs     *int              `json:",omitempty"`
	Time       *utils.Duration   `json:",omitempty"`
	ExpireAt   *utils.Time       `json:",omitempty"`
	Devices    map[string]string `json:",omitempty"`
}

// GetAuth decodes the Proxy-Authorization header
// @param *http.Request
// @returns *UserInfo
func GetAuth(r *http.Request) *UserInfo {
	s := r.Header.Get("Proxy-Authorization")
	if s == "" {
		return nil
	}
	ss := strings.Split(s, " ")
	if ss[0] != "Basic" {
		return nil
	}
	b, err := base64.StdEncoding.DecodeString(ss[1])
	if err != nil {
		return nil
	}
	ss = strings.Split(string(b), ":")
	return &UserInfo{
		User: ss[0],
		Pass: ss[1],
	}
}
