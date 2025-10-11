package ipfsconfig

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// IpfsConfig holds the configuration necessary for connecting to an IPFS node.
type IpfsConfig struct {
	Ipfs struct {
		NodePath string `yaml:"node_path"`
	} `yaml:"ipfs"`
}

// LoadConfig reads a YAML configuration file from the given path
// and unmarshals it into an IpfsConfig struct. Returns an error if the
// file cannot be read or if the YAML is invalid.
func LoadConfig(path string) (*IpfsConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	var cfg IpfsConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &cfg, nil
}
