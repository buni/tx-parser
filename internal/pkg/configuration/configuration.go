package configuration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"net"
	"runtime/debug"

	"github.com/spf13/viper"
)

const (
	EnvLocal = "local"
)

func NewConfiguration() (conf *Configuration, err error) {
	conf = new(Configuration)

	conf.SetDefaults()

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
	s.Port = "8181"
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
	Ethereum `mapstructure:",squash"`
	Database `mapstructure:",squash"`
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
	e.InitialHeight = "21718689"
	e.SeedAddresses = []string{"0x2527d2ed1dd0e7de193cf121f1630caefc23ac70", "0xf70da97812cb96acdf810712aa562db8dfa3dbef"}
}

type Database struct {
	Host     string `json:"database_host" mapstructure:"database_host"`
	Port     string `json:"database_port" mapstructure:"database_port"`
	User     string `json:"database_user" mapstructure:"database_user"`
	Password string `json:"database_password" mapstructure:"database_password"`
	Name     string `json:"database_name" mapstructure:"database_name"`
	URL      string `json:"database_url" mapstructure:"database_url"`
}

func (db Database) ToURL() string {
	if db.URL != "" {
		return db.URL
	}
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		db.User,
		db.Password,
		db.Host,
		db.Port,
		db.Name,
	)
}
