package configuration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"net"
	"reflect"
	"runtime/debug"

	"github.com/spf13/viper"
)

const (
	EnvLocal = "local"
)

func NewConfiguration() (conf *Configuration, err error) {
	conf = new(Configuration)

	method := reflect.ValueOf(conf).MethodByName("SetDefaults")
	if method.IsValid() {
		method.Call(nil)
	}

	viper.AutomaticEnv()
	viper.SetConfigType(`json`)

	jsonConf, err := json.Marshal(conf)
	if err != nil {
		return conf, fmt.Errorf("json marshal: %w", err)
	}

	err = viper.MergeConfig(bytes.NewBuffer(jsonConf))
	if err != nil {
		return conf, fmt.Errorf("merge config: %w", err)
	}

	err = viper.Unmarshal(conf)
	if err != nil {
		return conf, fmt.Errorf("viper unmarshal: %w", err)
	}

	return conf, nil
}

type Service struct {
	HostName    string `json:"service_host" mapstructure:"service_host"`
	Port        string `json:"service_port" mapstructure:"service_port"`
	Name        string `json:"service_name" mapstructure:"service_name"`
	Environment string `json:"service_environment" mapstructure:"service_environment"`
	CommitHash  string `json:"service_commit_hash"`
	GoVersion   string `json:"service_go_version"`
}

func (s *Service) SetDefaults() {
	s.Port = "8080"
	s.Environment = EnvLocal
	buildInfo, ok := debug.ReadBuildInfo()
	if ok {
		for _, v := range buildInfo.Settings {
			if v.Key == "vcs.revision" && len(v.Value) >= 7 {
				s.CommitHash = v.Value[:7]
			}
		}
		s.GoVersion = buildInfo.GoVersion
	}
}

func (s Service) ToLogFields() []any {
	return []any{slog.String("name", s.Name), slog.String("env", s.Environment), slog.String("commit", s.CommitHash), slog.String("go_version", s.GoVersion)}
}

func (s Service) ToHost() string {
	return net.JoinHostPort(s.HostName, s.Port)
}

type Configuration struct {
	Ethereum `mapstructure:"ethereum"`
	Service  `mapstructure:",squash"`
}

func (c *Configuration) SetDefaults() {
	c.Service.SetDefaults()
	c.Ethereum.SetDefaults()
}

type Ethereum struct {
	RPCEndpoint   string   `json:"rpc_endpoint" mapstructure:"rpc_endpoint"`
	InitialHeight string   `json:"initial_height" mapstructure:"initial_height"`
	SeedAddresses []string `json:"seed_addresses" mapstructure:"seed_addresses"`
}

func (e *Ethereum) SetDefaults() {
	e.RPCEndpoint = "https://ethereum-rpc.publicnode.com/"
	e.InitialHeight = "0x14b66a0"
	e.SeedAddresses = []string{"0x2527d2ed1dd0e7de193cf121f1630caefc23ac70", "0xf70da97812cb96acdf810712aa562db8dfa3dbef"}
}
