package provider

import (
	"github.com/claranet/go-zabbix-api"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceZabbixTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceZabbixTemplateCreate,
		Read:   resourceZabbixTemplateRead,
		Update: resourceZabbixTemplateUpdate,
		Delete: resourceZabbixTemplateDelete,
		Schema: map[string]*schema.Schema{
			"template_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "(readonly) ID of the template. ",
			},
			"host": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Technical name of the template",
			},
			"groups": &schema.Schema{
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Required: true,
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Visible name of the template",
			},
			"description": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description of the template",
			},
		},
	}
}

func createTemplateObj(d *schema.ResourceData, api *zabbix.API) (*zabbix.Template, error) {
	template := zabbix.Template{
		Host:        d.Get("host").(string),
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
	}
	hostGroupIDs, err := getHostGroups(d, api)
	if err != nil {
		return nil, err
	}
	template.Groups = make([]zabbix.HostGroup, len(hostGroupIDs))
	for i, ID := range hostGroupIDs {
		template.Groups[i].GroupID = ID.GroupID
	}
	return &template, nil
}

func resourceZabbixTemplateCreate(d *schema.ResourceData, meta interface{}) error {
	api := meta.(*zabbix.API)

	template, err := createTemplateObj(d, api)
	if err != nil {
		return err
	}
	templates := zabbix.Templates{*template}

	err = api.TemplatesCreate(templates)
	if err != nil {
		return err
	}
	d.Set("template_id", templates[0].TemplateID)
	d.SetId(templates[0].TemplateID)
	return resourceZabbixTemplateRead(d, meta)
}

func resourceZabbixTemplateRead(d *schema.ResourceData, meta interface{}) error {
	api := meta.(*zabbix.API)

	template, err := api.TemplateGetByID(d.Id())
	if err != nil {
		return err
	}

	d.Set("host", template.Host)
	d.Set("name", template.Name)
	d.Set("description", template.Description)

	params := zabbix.Params{
		"output": "extend",
		"hostids": []string{
			d.Id(),
		},
	}

	groups, err := api.HostGroupsGet(params)
	if err != nil {
		return err
	}

	groupNames := make([]string, len(groups))
	for i, g := range groups {
		groupNames[i] = g.Name
	}
	d.Set("groups", groupNames)
	return nil
}

func resourceZabbixTemplateUpdate(d *schema.ResourceData, meta interface{}) error {
	api := meta.(*zabbix.API)

	template, err := createTemplateObj(d, api)
	if err != nil {
		return err
	}
	template.TemplateID = d.Id()
	templates := zabbix.Templates{*template}

	err = api.TemplatesUpdate(templates)
	if err != nil {
		return err
	}

	return resourceZabbixTemplateRead(d, meta)
}

func resourceZabbixTemplateDelete(d *schema.ResourceData, meta interface{}) error {
	api := meta.(*zabbix.API)

	return api.TemplatesDeleteByIds([]string{d.Id()})
}
