package fabricconfig

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config holds the configuration for a Fabric client, including
// the user's identity and network connection details.
type FabricConfig struct {
	Identity struct {
		CertPath string `yaml:"cert_path"`
		KeyPath  string `yaml:"key_path"`
		MspID    string `yaml:"msp_id"`
	} `yaml:"identity"`

	Network struct {
		PeerEndpoint  string `yaml:"peer_endpoint"`
		TLSCertPath   string `yaml:"tls_cert_path"`
		TLSHostname   string `yaml:"tls_hostname"`
		ChannelName   string `yaml:"channel_name"`
		ChaincodeName string `yaml:"chaincode_name"`
	} `yaml:"network"`
}

// LoadConfig reads a YAML configuration file from the given path
// and unmarshals it into a Config struct. Returns an error if the
// file cannot be read or if the YAML is invalid.
func LoadConfig(path string) (*FabricConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	var cfg FabricConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &cfg, nil
}
