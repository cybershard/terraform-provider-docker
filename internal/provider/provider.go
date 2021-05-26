package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func init() {
	// Set descriptions to support markdown syntax, this will be used in document generation
	// and the language server.
	schema.DescriptionKind = schema.StringMarkdown

	// Customize the content of descriptions when output. For example you can add defaults on
	// to the exported descriptions if present.
	// schema.SchemaDescriptionBuilder = func(s *schema.Schema) string {
	// 	desc := s.Description
	// 	if s.Default != nil {
	// 		desc += fmt.Sprintf(" Defaults to `%v`.", s.Default)
	// 	}
	// 	return strings.TrimSpace(desc)
	// }
}

// New creates the Docker provider
func New(version string) func() *schema.Provider {
	return func() *schema.Provider {
		p := &schema.Provider{
			ResourcesMap: map[string]*schema.Resource{
				//"docker_container":      resourceDockerContainer(),
				"docker_image":          resourceDockerImage(),
				//"docker_registry_image": resourceDockerRegistryImage(),
				//"docker_network":        resourceDockerNetwork(),
				//"docker_volume":         resourceDockerVolume(),
				//"docker_config":         resourceDockerConfig(),
				//"docker_secret":         resourceDockerSecret(),
				//"docker_service":        resourceDockerService(),
				//"docker_plugin":         resourceDockerPlugin(),
			},

			//DataSourcesMap: map[string]*schema.Resource{
				//"docker_registry_image": dataSourceDockerRegistryImage(),
				//"docker_network":        dataSourceDockerNetwork(),
				//"docker_plugin":         dataSourceDockerPlugin(),
			//},
		}

		p.ConfigureContextFunc = func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
			hostMap := make(map[string]*ClientConfig)
			return &hostMap, nil
		}

		return p
	}
}
