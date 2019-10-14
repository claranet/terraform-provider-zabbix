package provider

import (
	"fmt"
	"log"
	"strings"

	"github.com/claranet/go-zabbix-api"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceZabbixTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceZabbixTemplateCreate,
		Read:   resourceZabbixTemplateRead,
		Exists: resourceZabbixTemplateExist,
		Update: resourceZabbixTemplateUpdate,
		Delete: resourceZabbixTemplateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
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
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Required:    true,
				Description: "ID of the Host Group",
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
			"macro": &schema.Schema{
				Type:        schema.TypeMap,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "User macros for the template",
			},
			"linked_template": &schema.Schema{
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
			"linked_host": &schema.Schema{
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
		},
	}
}

func createZabbixMacro(d *schema.ResourceData) zabbix.Macros {
	var macros zabbix.Macros

	terraformMacros := d.Get("macro").(map[string]interface{})
	for i, terraformMacro := range terraformMacros {
		macro := zabbix.Macro{
			MacroName: fmt.Sprintf("{$%s}", i),
			Value:     terraformMacro.(string),
		}
		macros = append(macros, macro)
	}
	return macros
}

func createLinkedTemplate(d *schema.ResourceData) zabbix.Templates {
	var templates zabbix.Templates

	terraformTemplates := d.Get("linked_template").(*schema.Set)
	for _, terraformTemplate := range terraformTemplates.List() {
		zabbixTemplate := zabbix.Template{
			TemplateID: terraformTemplate.(string),
		}
		templates = append(templates, zabbixTemplate)
	}
	return templates
}

func createLinkedHost(d *schema.ResourceData) []string {
	var hosts []string

	terraformHosts := d.Get("linked_host").(*schema.Set)
	for _, terraformHost := range terraformHosts.List() {
		hosts = append(hosts, terraformHost.(string))
	}
	return hosts
}

func createTemplateObj(d *schema.ResourceData, api *zabbix.API) (*zabbix.Template, error) {
	template := zabbix.Template{
		Host:            d.Get("host").(string),
		Name:            d.Get("name").(string),
		Description:     d.Get("description").(string),
		UserMacros:      createZabbixMacro(d),
		LinkedTemplates: createLinkedTemplate(d),
		LinkedHosts:     createLinkedHost(d),
	}
	hostGroupIDs, err := getHostGroups(d, api)
	if err != nil {
		return nil, err
	}
	template.Groups = make([]zabbix.HostGroup, len(hostGroupIDs))
	for i, ID := range hostGroupIDs {
		template.Groups[i].GroupID = ID.GroupID
	}
	if template.UserMacros == nil {
		template.UserMacros = zabbix.Macros{}
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

	params := zabbix.Params{
		"templateids":     d.Id(),
		"selectHosts":     "extend",
		"selectTemplates": "extend",
		"output":          "extend",
		"selectMacros":    "extend",
	}
	templates, err := api.TemplatesGet(params)
	if err != nil {
		return err
	}

	template := templates[0]
	d.Set("host", template.Host)
	if template.Host != template.Name && d.Get("name").(string) == "" {
		d.Set("name", template.Name)
	}
	d.Set("description", template.Description)

	terraformMacros, err := createTerraformMacro(template)
	if err != nil {
		return err
	}
	d.Set("macro", terraformMacros)

	terraformGroups, err := createTerraformTemplateGroup(d, api)
	if err != nil {
		return err
	}
	d.Set("groups", terraformGroups)
	return nil
}

func resourceZabbixTemplateExist(d *schema.ResourceData, meta interface{}) (bool, error) {
	api := meta.(*zabbix.API)

	_, err := api.TemplateGetByID(d.Id())
	if err != nil {
		log.Printf("Template with id %s doesn t exist", d.Id())
		if strings.Contains(err.Error(), "Expected exactly one result") {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func resourceZabbixTemplateUpdate(d *schema.ResourceData, meta interface{}) error {
	api := meta.(*zabbix.API)

	template, err := createTemplateObj(d, api)
	if err != nil {
		return err
	}
	template.TemplatesClear = getUnlinkedTemplate(d)
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

func createTerraformMacro(template zabbix.Template) (map[string]interface{}, error) {
	terraformMacros := make(map[string]interface{}, len(template.UserMacros))

	for _, macro := range template.UserMacros {
		var name string
		if noPrefix := strings.Split(macro.MacroName, "{$"); len(noPrefix) == 2 {
			name = noPrefix[1]
		} else {
			return nil, fmt.Errorf("Invalid macro name \"%s\"", macro.MacroName)
		}
		if noSuffix := strings.Split(name, "}"); len(noSuffix) == 2 {
			name = noSuffix[0]
		} else {
			return nil, fmt.Errorf("Invalid macro name \"%s\"", macro.MacroName)
		}
		terraformMacros[name] = macro.Value
	}
	return terraformMacros, nil
}

func createTerraformTemplateGroup(d *schema.ResourceData, api *zabbix.API) ([]string, error) {
	params := zabbix.Params{
		"output": "extend",
		"hostids": []string{
			d.Id(),
		},
	}

	groups, err := api.HostGroupsGet(params)
	if err != nil {
		return nil, err
	}

	groupNames := make([]string, len(groups))
	for i, g := range groups {
		groupNames[i] = g.Name
	}
	return groupNames, nil
}

func createTerraformLinkedTemplate(template zabbix.Template) []string {
	var terraformTemplates []string

	for _, linkedTemplate := range template.LinkedTemplates {
		terraformTemplates = append(terraformTemplates, linkedTemplate.TemplateID)
	}
	return terraformTemplates
}

func createTerraformLinkedHost(template zabbix.Template) []string {
	var terraformHosts []string

	for _, linkedHost := range template.LinkedHosts {
		terraformHosts = append(terraformHosts, linkedHost)
	}
	return terraformHosts
}

func getUnlinkedTemplate(d *schema.ResourceData) zabbix.Templates {
	before, after := d.GetChange("linked_template")
	beforeID := before.(*schema.Set).List()
	afterID := after.(*schema.Set).List()
	var unlinkID zabbix.Templates

	for _, l := range beforeID {
		present := false
		for _, k := range afterID {
			if l == k {
				present = true
			}
		}
		if !present {
			unlinkID = append(unlinkID, zabbix.Template{TemplateID: l.(string)})
		}
	}
	return unlinkID
}
