package command

import (
	"github.com/urfave/cli"
)

// ProviderVCloudFlags : ..
var ProviderVCloudFlags = []cli.Flag{
	tStringFlag("vcloud.create.flags.user"),
	tStringFlag("vcloud.create.flags.password"),
	tStringFlag("vcloud.create.flags.org"),
	tStringFlag("vcloud.create.flags.vse-url"),
	tStringFlag("vcloud.create.flags.vcloud-url"),
	tStringFlag("vcloud.create.flags.public-network"),
	tStringFlag("vcloud.create.flags.vcloud-region"),
}

// ProviderAWSFlags : ...
var ProviderAWSFlags = []cli.Flag{
	tStringFlag("aws.create.flags.region"),
	tStringFlag("aws.create.flags.access_key_id"),
	tStringFlag("aws.create.flags.secret_access_key"),
}

// ProviderAzureFlags : ...
var ProviderAzureFlags = []cli.Flag{
	tStringFlag("azure.update.flags.subscription_id"),
	tStringFlag("azure.update.flags.client_id"),
	tStringFlag("azure.update.flags.client_secret"),
	tStringFlag("azure.update.flags.tenant_id"),
	tStringFlag("azure.update.flags.environment"),
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
