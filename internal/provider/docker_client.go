package provider

import (
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type ClientConfig struct {
	host string
	registryAddress string
	registryUsername string
	registryPassword string
	dockerClient *client.Client
}

func (c *ClientConfig) getAuthConfig() *types.AuthConfig {
	return &types.AuthConfig{
		ServerAddress: normalizeRegistryAddress(c.registryAddress),
		Username: c.registryUsername,
		Password: c.registryPassword,
	}
}

func getOrCreateDockerClient(d *schema.ResourceData, meta *map[string]*ClientConfig) (*ClientConfig, error) {
	configs := *meta
	c := d.Get("client").([]interface{})[0].(map[string]interface{})
	hostUri := c["host"].(string)

	if clientConfig, ok := configs[hostUri]; ok {
		return clientConfig, nil
	} else {
		config := Config{
			Host: hostUri,
		}

		dockerClient, err := config.NewClient()
		if err != nil {
			return nil, err
		}

		clientConfig := ClientConfig{
			host:             hostUri,
			registryAddress:  c["registry_address"].(string),
			registryUsername: c["registry_username"].(string),
			registryPassword: c["registry_password"].(string),
			dockerClient: dockerClient,
		}

		configs[hostUri] = &clientConfig
		return &clientConfig, err
	}
}

func generateClientSchema() *schema.Schema {
	return &schema.Schema {
		Type:     schema.TypeList,
		Description:   "Configuration for Docker API client",
		Required: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"host": {
					Type:        schema.TypeString,
					Required:    true,
					DefaultFunc: schema.EnvDefaultFunc("DOCKER_HOST", "unix:///var/run/docker.sock"),
					Description: "The Docker daemon address",
				},
				"registry_address": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "Address of the registry",
				},

				"registry_username": {
					Type:        schema.TypeString,
					Optional:    true,
					DefaultFunc: schema.EnvDefaultFunc("DOCKER_REGISTRY_USER", ""),
					Description: "Username for the registry",
				},

				"registry_password": {
					Type:        schema.TypeString,
					Optional:    true,
					Sensitive:   true,
					DefaultFunc: schema.EnvDefaultFunc("DOCKER_REGISTRY_PASS", ""),
					Description: "Password for the registry",
				},
			},
		},
	}
}


