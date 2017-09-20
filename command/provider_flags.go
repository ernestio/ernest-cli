package command

import (
	"github.com/urfave/cli"
)

// ProviderVCloudFlags : ..
var ProviderVCloudFlags = []cli.Flag{
	cli.StringFlag{
		Name:  "user",
		Value: "",
		Usage: "Your VCloud valid user name",
	},
	cli.StringFlag{
		Name:  "password",
		Value: "",
		Usage: "Your VCloud valid password",
	},
	cli.StringFlag{
		Name:  "org",
		Value: "",
		Usage: "Your vCloud Organization",
	},
	cli.StringFlag{
		Name:  "vse-url",
		Value: "",
		Usage: "VSE URL",
	},
	cli.StringFlag{
		Name:  "vcloud-url",
		Value: "",
		Usage: "VCloud URL",
	},
	cli.StringFlag{
		Name:  "public-network",
		Value: "",
		Usage: "Public Network",
	},
	cli.StringFlag{
		Name:  "vcloud-region, reg",
		Value: "",
		Usage: "Project region",
	},
}

// ProviderAWSFlags : ...
var ProviderAWSFlags = []cli.Flag{
	cli.StringFlag{
		Name:  "access_key_id, k",
		Value: "",
		Usage: "AWS access key id",
	},
	cli.StringFlag{
		Name:  "secret_access_key, sak",
		Value: "",
		Usage: "AWS Secret access key",
	},
	cli.StringFlag{
		Name:  "region, r",
		Value: "",
		Usage: "Project region",
	},
}

// ProviderAzureFlags : ...
var ProviderAzureFlags = []cli.Flag{
	cli.StringFlag{
		Name:  "subscription_id, s",
		Value: "",
		Usage: "Azure subscription id",
	},
	cli.StringFlag{
		Name:  "client_id, c",
		Value: "",
		Usage: "Azure client id",
	},
	cli.StringFlag{
		Name:  "client_secret, p",
		Value: "",
		Usage: "Azure client secret",
	},
	cli.StringFlag{
		Name:  "tenant_id, t",
		Value: "",
		Usage: "Azure tenant_id",
	},
	cli.StringFlag{
		Name:  "environment, e",
		Value: "",
		Usage: "Azure environment. Supported values are public(default), usgovernment, german and chine",
	},
}

// AWSVCloudFlags : All aws vcloud provider flags
var AWSVCloudFlags = append(ProviderVCloudFlags, ProviderAWSFlags...)

// AllProviderFlags : All provider flags
var AllProviderFlags = append(AWSVCloudFlags, ProviderAzureFlags...)

// ProviderFlagsToSlice :
func ProviderFlagsToSlice(c *cli.Context) map[string]interface{} {
	keys := make(map[string]interface{}, 0)
	keys["subscription_id"] = "azure_subscription_id"
	keys["client_id"] = "azure_client_id"
	keys["client_secret"] = "azure_client_secret"
	keys["tenant_id"] = "azure_tenant_id"
	keys["access_key_id"] = "aws_access_key_id"
	keys["secret_access_key"] = "aws_secret_access_key"
	keys["region"] = "region"
	keys["vcloud-region"] = "region"
	keys["public-network"] = "external_network"
	keys["vcloud-url"] = "vcloud_url"
	keys["vse-url"] = "vse_url"
	keys["password"] = "password"
	keys["user"] = "username"

	flags := make(map[string]interface{}, 0)

	for key, val := range keys {
		v := c.String(key)
		if v != "" {
			flags[val.(string)] = v
		}
	}

	return flags
}
