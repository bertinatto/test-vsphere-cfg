package main

import (
	"bytes"
	"fmt"
	"io"
	"testing"

	"gopkg.in/gcfg.v1"
)

// e2eTestConfig contains vSphere connection detail and kubernetes cluster-id
type e2eTestConfig struct {
	Global struct {
		// Kubernetes Cluster-ID
		ClusterID string `gcfg:"cluster-id"`
		// Kubernetes Cluster-Distribution
		ClusterDistribution string `gcfg:"cluster-distribution"`
		// vCenter username.
		User string `gcfg:"user"`
		// vCenter password in clear text.
		Password string `gcfg:"password"`
		// vmc vCentre cloudadmin username
		VmcCloudUser string `gcfg:"vmc-cloudadminuser"`
		// vmc vCentre cloudadmin password
		VmcCloudPassword string `gcfg:"cloudadminpassword"`
		// vmc vCentre Devops username
		VmcDevopsUser string `gcfg:"vmc-devopsuser"`
		// vmc vCentre Devops password
		VmcDevopsPassword string `gcfg:"vmc-devopspassword"`
		// vCenter Hostname.
		VCenterHostname string `gcfg:"hostname"`
		// vCenter port.
		VCenterPort string `gcfg:"port"`
		// True if vCenter uses self-signed cert.
		InsecureFlag bool `gcfg:"insecure-flag"`
		// Datacenter in which VMs are located.
		Datacenters string `gcfg:"datacenters"`
		// CnsRegisterVolumesCleanupIntervalInMin specifies the interval after which
		// successful CnsRegisterVolumes will be cleaned up.
		CnsRegisterVolumesCleanupIntervalInMin int `gcfg:"cnsregistervolumes-cleanup-intervalinmin"`
		// preferential topology
		CSIFetchPreferredDatastoresIntervalInMin int `gcfg:"csi-fetch-preferred-datastores-intervalinmin"`
		// QueryLimit specifies the number of volumes that can be fetched by CNS QueryAll API at a time
		QueryLimit int `gcfg:"query-limit"`
		// ListVolumeThreshold specifies the maximum number of differences in volume that can exist between CNS
		// and kubernetes
		ListVolumeThreshold int `gcfg:"list-volume-threshold"`
	}
}

// readConfig parses e2e tests config file into Config struct.
func readConfig(config io.Reader) (e2eTestConfig, error) {
	if config == nil {
		err := fmt.Errorf("no config file given")
		return e2eTestConfig{}, err
	}
	var cfg e2eTestConfig
	err := gcfg.ReadInto(&cfg, config)
	return cfg, err
}

func TestReadConfig(t *testing.T) {
	tests := []struct {
		name          string
		config        io.Reader
		expectedError bool
	}{
		{
			name: "escaped password",
			config: bytes.NewBufferString(`
[Global]
password=\"&)<;^}.
`),
			expectedError: false,
		},
		{
			name: "non-escaped password should fail",
			config: bytes.NewBufferString(`
[Global]
password="&)<;^}.
`),
			expectedError: true,
		},
		{
			name: "blah",
			config: bytes.NewBufferString(`
[Global]
password="&)<;^}."
`),
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := readConfig(tt.config)
			if tt.expectedError {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Fatalf("did not expect error, got %v", err)
				}
			}
		})
	}
}
