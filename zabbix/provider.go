package provider

import (
	"fmt"
	"net/http"

	"github.com/claranet/go-zabbix-api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/logging"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

var zabbixAPIVersion = ""

// Provider define the provider and his resources
func Provider() terraform.ResourceProvider {
	p := &schema.Provider{
		Schema: map[string]*schema.Schema{
			"user": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("ZABBIX_USER", nil),
			},
			"password": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("ZABBIX_PASSWORD", nil),
			},
			"server_url": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("ZABBIX_SERVER_URL", nil),
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"zabbix_host":              resourceZabbixHost(),
			"zabbix_host_group":        resourceZabbixHostGroup(),
			"zabbix_item":              resourceZabbixItem(),
			"zabbix_trigger":           resourceZabbixTrigger(),
			"zabbix_template":          resourceZabbixTemplate(),
			"zabbix_template_link":     resourceZabbixTemplateLink(),
			"zabbix_lld_rule":          resourceZabbixlldRule(),
			"zabbix_item_prototype":    resourceZabbixItemPrototype(),
			"zabbix_trigger_prototype": resourceZabbixTriggerPrototype(),
		},
	}

	p.ConfigureFunc = func(d *schema.ResourceData) (interface{}, error) {
		terraformVersion := p.TerraformVersion
		if terraformVersion == "" {
			// Terraform 0.12 introduced this field to the protocol
			// We can therefore assume that if it's missing it's 0.10 or 0.11
			terraformVersion = "0.11+compatible"
		}
		return providerConfigure(d, terraformVersion)
	}

	return p
}

func providerConfigure(d *schema.ResourceData, terraformVersion string) (interface{}, error) {
	api := zabbix.NewAPI(d.Get("server_url").(string))

	api.UserAgent = fmt.Sprintf("HashiCorp/1.0 Terraform/%s", terraformVersion)

	if logging.IsDebugOrHigher() {
		httpClient := http.Client{}
		httpClient.Transport = logging.NewTransport("Zabbix", http.DefaultTransport)
		api.SetClient(&httpClient)
	}

	if _, err := api.Login(d.Get("user").(string), d.Get("password").(string)); err != nil {
		return nil, err
	}
	v, err := api.Version()
	if err != nil {
		return nil, err
	}
	zabbixAPIVersion = v

	return api, nil
}
