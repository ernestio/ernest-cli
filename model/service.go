/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package model

// Service : Model representing ernest.io service json responses
type Service struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Datacenter int    `json:"datacenter_id"`
	Version    string `json:"version"`
	UserName   string `json:"user_name"`
	Status     string `json:"status"`
	Definition string `json:"definition"`
	Result     string `json:"result"`
	LastError  string `json:"last_known_error"`
	Endpoint   string `json:"endpoint"`
	Vpc        string `json:"vpc_id"`
	Networks   []struct {
		Name             string `json:"name"`
		Subnet           string `json:"network_aws_id"`
		AvailabilityZone string `json:"availability_zone"`
	} `json:"networks"`
	Instances []struct {
		Name          string `json:"name"`
		InstanceAWSID string `json:"instance_aws_id"`
		PublicIP      string `json:"public_ip"`
		IP            string `json:"ip"`
	} `json:"instances"`
	Nats []struct {
		Name            string `json:"name"`
		NatGatewayAWSID string `json:"nat_gateway_aws_id"`
	} `json:"nats"`
	SecurityGroups []struct {
		Name               string `json:"name"`
		SecurityGroupAWSID string `json:"security_group_aws_id"`
	} `json:"security_groups"`
	ELBs []struct {
		Name    string `json:"name"`
		DNSName string `json:"dns_name"`
	} `json:"elbs"`
	RDSClusters []struct {
		Name     string `json:"name"`
		Endpoint string `json:"endpoint"`
	} `json:"rds_clusters"`
	RDSInstances []struct {
		Name     string `json:"name"`
		Endpoint string `json:"endpoint"`
	} `json:"rds_instances"`
}
