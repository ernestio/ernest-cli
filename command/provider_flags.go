package command

import (
	"github.com/urfave/cli"
)

// ProviderVCloudFlags : ..
var ProviderVCloudFlags = []cli.Flag{
	stringFlag("user", "", "Your VCloud valid user name"),
	stringFlag("password", "", "Your VCloud valid password"),
	stringFlag("org", "", "Your vCloud Organization"),
	stringFlag("vse-url", "", "VSE URL"),
	stringFlag("vcloud-url", "", "VCloud URL"),
	stringFlag("public-network", "", "Public Network"),
	stringFlag("vcloud-region, reg", "", "Project region"),
}

// ProviderAWSFlags : ...
var ProviderAWSFlags = []cli.Flag{
	stringFlag("access_key_id, k", "", "AWS access key id"),
	stringFlag("secret_access_key, sak", "", "AWS Secret access key"),
	stringFlag("region, r", "", "Project region"),
}

// ProviderAzureFlags : ...
var ProviderAzureFlags = []cli.Flag{
	stringFlag("subscription_id, s", "", "Azure subscription id"),
	stringFlag("client_id, c", "", "Azure client id"),
	stringFlag("client_secret, p", "", "Azure client secret"),
	stringFlag("tenant_id, t", "", "Azure tenant_id"),
	stringFlag("environment, e", "", "Azure environment. Supported values are public(default), usgovernment, german and chine"),
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
