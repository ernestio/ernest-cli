package main

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
}
