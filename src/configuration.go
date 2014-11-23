package goule

import (
	"./exec"
	"./server"
	"bytes"
	"encoding/json"
	"io/ioutil"
)

type Service struct {
	Name         string          `json:"name"`
	ForwardRules []ForwardRule   `json:"forward_rules"`
	Executables  []exec.Settings `json:"executables"`
}

type ServerSettings struct {
	Enabled bool `json:"enabled"`
	Port    int  `json:"port"`
}

type AdminSettings struct {
	Rules          []SourceURL `json:"rules"`
	PasswordHash   string      `json:"password_hash"`
	SessionTimeout int         `json:"session_timeout"`
}

type Configuration struct {
	LoadedPath    string         `json:"-"`
	Services      []Service      `json:"services"`
	TLS           server.TLSInfo `json:"tls"`
	HTTPSettings  ServerSettings `json:"http"`
	HTTPSSettings ServerSettings `json:"https"`
	Admin         AdminSettings  `json:"admin"`
}

func NewConfiguration() *Configuration {
	return &Configuration{}
}

func (self *Configuration) Read(path string) error {
	configData, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(configData, self); err != nil {
		return err
	}
	self.LoadedPath = path
	return nil
}

func (self *Configuration) Save() error {
	if data, err := json.Marshal(self); err == nil {
		var out bytes.Buffer
		json.Indent(&out, data, "", "  ")
		return ioutil.WriteFile(self.LoadedPath, out.Bytes(), 0700)
	} else {
		return err
	}
}
