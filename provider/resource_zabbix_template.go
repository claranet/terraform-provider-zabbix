package provider

import (
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
				Type: schema.TypeList,
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

	if err := api.TemplatesCreate(templates); err != nil {
		return err
	}

	items := make(zabbix.Items, d.Get("item.#").(int))
	itemsTerraform := d.Get("item").([]interface{})
	for i, item := range itemsTerraform {
		value := item.(map[string]interface{})
		log.Printf("[DEBUG] Creating item: %#v", value["name"])
		items[i].Delay = value["delay"].(int)
		items[i].HostID = templates[0].TemplateID
		items[i].Key = value["key"].(string)
		items[i].Name = value["name"].(string)
	}

	if err := api.ItemsCreate(items); err != nil {
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

	groups, err := api.HostGroupsGet(zabbix.Params{
		"output":  "extend",
		"hostids": []string{d.Id()},
	})
	if err != nil {
		return err
	}

	groupNames := make([]string, len(groups))
	for i, g := range groups {
		groupNames[i] = g.Name
	}
	d.Set("groups", groupNames)

	items, err := api.ItemsGet(zabbix.Params{
		"output":      "extend",
		"templateids": []string{d.Id()},
	})
	if err != nil {
		return err
	}

	itemTerraList := make([]interface{}, len(items), len(items))
	for i, item := range items {
		itemTerra := make(map[string]interface{})

		log.Printf("[DEBUG] Item data %#v", item)
		itemTerra["delay"] = item.Delay
		itemTerra["name"] = item.Name
		itemTerra["key"] = item.Key
		itemTerra["item_id"] = item.ItemID
		itemTerraList[i] = itemTerra
	}
	if err := d.Set("item", itemTerraList); err != nil {
		return err
	}

	return nil
}

func resourceZabbixTemplateExist(d *schema.ResourceData, meta interface{}) (bool, error) {
	api := meta.(*zabbix.API)

	if _, err := api.TemplateGetByID(d.Id()); err != nil {
		log.Printf("[DEBUG] Template with id %s doesn't exist", d.Id())
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

	if err := api.TemplatesUpdate(templates); err != nil {
		return err
	}

	if d.HasChange("item") {
		log.Printf("[TRACE] template.item has changes")
		oldV, newV := d.GetChange("item")
		oldItemsTerraform := oldV.([]interface{})
		newItemsTerraform := newV.([]interface{})
		createdItems := zabbix.Items{}
		updatedItems := zabbix.Items{}
		deletedItems := zabbix.Items{}
		marked := make(map[string]bool)
		for _, newItemTerraform := range newItemsTerraform {
			log.Printf("[DEBUG] Checking item: %#v", newItemTerraform)

			item := zabbix.Item{
				HostID: templates[0].TemplateID,
				Delay:  newItemTerraform.(map[string]interface{})["delay"].(int),
				Key:    newItemTerraform.(map[string]interface{})["key"].(string),
				Name:   newItemTerraform.(map[string]interface{})["name"].(string),
			}

			// First pass on name
			for _, oldItemTerraform := range oldItemsTerraform {
				oldName := oldItemTerraform.(map[string]interface{})["name"].(string)
				if marked[oldName] {
					continue
				}
				if oldName == item.Name {
					log.Printf("[DEBUG] Marked item for update (matches name): %#v => %#v", oldItemTerraform, newItemTerraform)
					item.ItemID = oldItemTerraform.(map[string]interface{})["item_id"].(string)
					updatedItems = append(updatedItems, item)
					marked[item.Name] = true
					break
				}
			}

			// Second pass on key
			for _, oldItemTerraform := range oldItemsTerraform {
				oldName := oldItemTerraform.(map[string]interface{})["name"].(string)
				if marked[oldName] {
					continue
				}
				if oldItemTerraform.(map[string]interface{})["key"] == item.Key {
					log.Printf("[DEBUG] Marked item for update (matches key): %#v => %#v", oldItemTerraform, newItemTerraform)
					item.ItemID = oldItemTerraform.(map[string]interface{})["item_id"].(string)
					updatedItems = append(updatedItems, item)
					marked[item.Name] = true
					marked[oldName] = true
					break
				}
			}

			// New item not marked is considered created
			if !marked[item.Name] {
				log.Printf("[DEBUG] Marked item for creation: %#v", newItemTerraform)
				createdItems = append(createdItems, item)
			}
		}

		for _, oldItemTerraform := range oldItemsTerraform {
			// Old item not marked is considered removed
			if marked[oldItemTerraform.(map[string]interface{})["name"].(string)] {
				continue
			}
			log.Printf("[DEBUG] Marked item for deletion: %#v", oldItemTerraform)
			item := zabbix.Item{
				ItemID: oldItemTerraform.(map[string]interface{})["item_id"].(string),
				HostID: templates[0].TemplateID,
				Delay:  oldItemTerraform.(map[string]interface{})["delay"].(int),
				Key:    oldItemTerraform.(map[string]interface{})["key"].(string),
				Name:   oldItemTerraform.(map[string]interface{})["name"].(string),
			}
			deletedItems = append(deletedItems, item)
		}

		if len(createdItems) > 0 {
			log.Printf("[DEBUG] Item to create: %#v", createdItems)
			if err := api.ItemsCreate(createdItems); err != nil {
				return err
			}
		}
		if len(updatedItems) > 0 {
			log.Printf("[DEBUG] Item to update: %#v", updatedItems)
			if err := api.ItemsUpdate(updatedItems); err != nil {
				return err
			}
		}
		if len(deletedItems) > 0 {
			log.Printf("[DEBUG] Item to delete: %#v", deletedItems)
			if err := api.ItemsDelete(deletedItems); err != nil {
				return err
			}
		}
	}

	return resourceZabbixTemplateRead(d, meta)
}

func resourceZabbixTemplateDelete(d *schema.ResourceData, meta interface{}) error {
	api := meta.(*zabbix.API)

	return api.TemplatesDeleteByIds([]string{d.Id()})
}
