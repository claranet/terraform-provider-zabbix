package provider

import (
	"fmt"
	"log"
	"strings"

	"github.com/claranet/go-zabbix-api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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

func createTemplateObj(d *schema.ResourceData, api *zabbix.API) (*zabbix.Template, error) {
	template := zabbix.Template{
		Host:            d.Get("host").(string),
		Name:            d.Get("name").(string),
		Description:     d.Get("description").(string),
		UserMacros:      createZabbixMacro(d),
		LinkedTemplates: createLinkedTemplate(d),
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

	return createRetry(d, meta, createTemplate, *template, resourceZabbixTemplateRead)
}

func resourceZabbixTemplateRead(d *schema.ResourceData, meta interface{}) error {
	api := meta.(*zabbix.API)

	params := zabbix.Params{
		"templateids":  d.Id(),
		"output":       "extend",
		"selectMacros": "extend",
	}
	templates, err := api.TemplatesGet(params)
	if err != nil {
		return err
	}
	if len(templates) != 1 {
		log.Printf("[DEBUG] Expected one template with id %s and got %#v", d.Id(), templates)
		return fmt.Errorf("Expected one template with id %s and got %d templates", d.Id(), len(templates))
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
		if strings.Contains(err.Error(), "Expected exactly one result") {
			log.Printf("Template with id %s doesn't exist", d.Id())
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

	return createRetry(d, meta, updateTemplate, *template, resourceZabbixTemplateRead)
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

func createTemplate(template interface{}, api *zabbix.API) (id string, err error) {
	templates := zabbix.Templates{template.(zabbix.Template)}

	err = api.TemplatesCreate(templates)
	if err != nil {
		return
	}
	id = templates[0].TemplateID
	return
}

func updateTemplate(template interface{}, api *zabbix.API) (id string, err error) {
	templates := zabbix.Templates{template.(zabbix.Template)}

	err = api.TemplatesUpdate(templates)
	if err != nil {
		return
	}
	id = templates[0].TemplateID
	return
}
