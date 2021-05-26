package provider

import (
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"strings"
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

func getOrCreateDockerClient(d *schema.ResourceData, meta map[string]*ClientConfig) (*ClientConfig, error) {
	hostUri := d.Get("client.host").(string)
	if clientConfig, ok := meta[hostUri]; ok {
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
			host:             d.Get("client.host").(string),
			registryAddress:  d.Get("client.registry_address").(string),
			registryUsername: d.Get("client.registry_username").(string),
			registryPassword: d.Get("client.registry_password").(string),
			dockerClient: dockerClient,
		}

		meta[hostUri] = &clientConfig
		return &clientConfig, err
	}
}

func generateClientSchema() *schema.Schema {
	return &schema.Schema {
		Type:          schema.TypeSet,
		Description:   "Configuration for Docker API client",
		Optional:      false,
		MaxItems:      1,
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

// ConvertToHostname converts a registry url which has http|https prepended
// to just an hostname.
// Copied from github.com/docker/docker/registry.ConvertToHostname to reduce dependencies.
func convertToHostname(url string) string {
	stripped := url
	// DevSkim: ignore DS137138
	if strings.HasPrefix(url, "http://") {
		// DevSkim: ignore DS137138
		stripped = strings.TrimPrefix(url, "http://")
	} else if strings.HasPrefix(url, "https://") {
		stripped = strings.TrimPrefix(url, "https://")
	}

	nameParts := strings.SplitN(stripped, "/", 2)

	return nameParts[0]
}


