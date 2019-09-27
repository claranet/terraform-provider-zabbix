package provider

import (
	"encoding/json"
	"log"
	"strconv"
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
			"item": &schema.Schema{
				Type: schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"item_id": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"delay": &schema.Schema{
							Type:     schema.TypeInt,
							Required: true,
						},
						"key": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "Item key.",
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "Name of the item.",
						},
					},
				},
				Optional: true,
				Set:      getUniqueItemID,
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

	items := make(zabbix.Items, d.Get("item.#").(int))
	itemsTerraform := d.Get("item")
	for i, item := range itemsTerraform.(*schema.Set).List() {
		value := item.(map[string]interface{})
		log.Printf("Creating item with name : %s", value["name"])
		items[i].Delay = value["delay"].(int)
		items[i].HostID = templates[0].TemplateID
		items[i].Key = value["key"].(string)
		items[i].Name = value["name"].(string)
	}
	err = api.ItemsCreate(items)
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

	params2 := zabbix.Params{
		"output":      "extend",
		"templateids": d.Id(),
	}
	items, err := api.ItemsGet(params2)
	if err != nil {
		return err
	}

	var itemTerraList []interface{}
	for _, item := range items {
		itemTerra := map[string]interface{}{}

		str, _ := json.Marshal(item)
		log.Printf("my item data %s", string(str))
		itemTerra["delay"] = item.Delay
		itemTerra["name"] = item.Name
		itemTerra["key"] = item.Key
		itemTerra["item_id"] = item.ItemID
		itemTerraList = append(itemTerraList, itemTerra)
	}
	d.Set("item", itemTerraList)
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
func getUniqueItemID(a interface{}) int {
	item := a.(map[string]interface{})

	i, err := strconv.Atoi(item["item_id"].(string))
	if err != nil {
		return schema.HashString(item["name"].(string))
	}
	return i
}
