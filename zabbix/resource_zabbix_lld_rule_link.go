package provider

import (
	"fmt"
	"log"
	"strings"

	"github.com/claranet/go-zabbix-api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceZabbixlldRuleLink() *schema.Resource {
	return &schema.Resource{
		Create: resourceZabbixlldRuleLinkCreate,
		Read:   resourceZabbixlldRuleLinkRead,
		Exists: resourceZabbixlldRuleLinkExist,
		Update: resourceZabbixlldRuleLinkUpdate,
		Delete: resourceZabbixlldRuleLinkDelete,
		Schema: map[string]*schema.Schema{
			"lld_rule_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"item_prototype": &schema.Schema{
				Type:     schema.TypeSet,
				Elem:     schemaTemplateItemPrototype(),
				Optional: true,
			},
			"trigger_prototype": &schema.Schema{
				Type:     schema.TypeSet,
				Elem:     schemaTemplateTriggerPrototype(),
				Optional: true,
			},
		},
	}
}

func schemaTemplateItemPrototype() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"local": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
			"item_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func schemaTemplateTriggerPrototype() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"local": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
			"trigger_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceZabbixlldRuleLinkCreate(d *schema.ResourceData, meta interface{}) error {
	d.SetId(randStringNumber(5))
	return resourceZabbixlldRuleLinkReadTrusted(d, meta)
}

func resourceZabbixlldRuleLinkRead(d *schema.ResourceData, meta interface{}) error {
	api := meta.(*zabbix.API)

	itemsTerraform, err := getTerraformTemplateItemPrototypesForPlan(d, api)
	if err != nil {
		return err
	}
	d.Set("item_prototype", itemsTerraform)

	triggersTerraform, err := getTerraformTemplateTriggerPrototypesForPlan(d, api)
	if err != nil {
		return err
	}
	d.Set("trigger_prototype", triggersTerraform)
	return nil
}

func resourceZabbixlldRuleLinkReadTrusted(d *schema.ResourceData, meta interface{}) error {
	api := meta.(*zabbix.API)

	itemsTerraform, err := getTerraformTemplateItemPrototypes(d, api)
	if err != nil {
		return err
	}
	d.Set("item_prototype", itemsTerraform)

	triggersTerraform, err := getTerraformTemplateTriggerPrototypes(d, api)
	if err != nil {
		return err
	}
	d.Set("trigger_prototype", triggersTerraform)
	return nil
}

func resourceZabbixlldRuleLinkExist(d *schema.ResourceData, meta interface{}) (bool, error) {
	return true, nil
}

func resourceZabbixlldRuleLinkUpdate(d *schema.ResourceData, meta interface{}) error {
	api := meta.(*zabbix.API)

	err := updateZabbixTemplateItemPrototypes(d, api)
	if err != nil {
		return err
	}

	err = updateZabbixTemplateTriggerPrototypes(d, api)
	if err != nil {
		return err
	}
	return resourceZabbixlldRuleLinkReadTrusted(d, meta)
}

func resourceZabbixlldRuleLinkDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func getTerraformTemplateItemPrototypesForPlan(d *schema.ResourceData, api *zabbix.API) ([]interface{}, error) {
	params := zabbix.Params{
		"output": "extend",
		"discoveryids": []string{
			d.Get("lld_rule_id").(string),
		},
		"inherited": false,
	}
	items, err := api.ItemPrototypesGet(params)
	if err != nil {
		return nil, err
	}

	itemList := d.Get("item_prototype").(*schema.Set).List()
	itemLocal := make(map[string]bool)
	var itemsTerraform []interface{}

	for _, item := range itemList {
		var itemTerraform = make(map[string]interface{})
		value := item.(map[string]interface{})

		log.Printf("Found local item prototype with id : %s", value["item_id"].(string))
		itemLocal[value["item_id"].(string)] = true
		itemTerraform["local"] = true
		itemTerraform["item_id"] = value["item_id"].(string)
		itemsTerraform = append(itemsTerraform, itemTerraform)
	}
	for _, item := range items {
		var itemTerraform = make(map[string]interface{})

		if itemLocal[item.ItemID] {
			continue
		}
		log.Printf("Found server item prototype with id : %s", item.ItemID)
		itemTerraform["local"] = false
		itemTerraform["item_id"] = item.ItemID
		itemsTerraform = append(itemsTerraform, itemTerraform)
	}
	return itemsTerraform, nil
}

func getTerraformTemplateItemPrototypes(d *schema.ResourceData, api *zabbix.API) ([]interface{}, error) {
	params := zabbix.Params{
		"output": "extend",
		"discoveryids": []string{
			d.Get("lld_rule_id").(string),
		},
		"inherited": false,
	}
	items, err := api.ItemPrototypesGet(params)
	if err != nil {
		return nil, err
	}

	itemsTerraform := make([]interface{}, len(items))
	for i, item := range items {
		var itemTerraform = make(map[string]interface{})

		itemTerraform["local"] = true
		itemTerraform["item_id"] = item.ItemID
		itemsTerraform[i] = itemTerraform
	}
	return itemsTerraform, nil
}

func getTerraformTemplateTriggerPrototypesForPlan(d *schema.ResourceData, api *zabbix.API) ([]interface{}, error) {
	params := zabbix.Params{
		"output": "extend",
		"discoveryids": []string{
			d.Get("lld_rule_id").(string),
		},
		"inherited": false,
	}
	triggers, err := api.TriggerPrototypesGet(params)
	if err != nil {
		return nil, err
	}

	triggerList := d.Get("trigger_prototype").(*schema.Set).List()
	triggerLocal := make(map[string]bool)
	var triggersTerraform []interface{}
	for _, trigger := range triggerList {
		triggerTerraform := make(map[string]interface{})
		value := trigger.(map[string]interface{})

		log.Printf("Found local trigger prototype with id : %s", value["trigger_id"].(string))
		triggerLocal[value["trigger_id"].(string)] = true
		triggerTerraform["trigger_id"] = value["trigger_id"].(string)
		triggerTerraform["local"] = true
		triggersTerraform = append(triggersTerraform, triggerTerraform)
	}
	for _, trigger := range triggers {
		var triggerTerraform = make(map[string]interface{})

		if triggerLocal[trigger.TriggerID] {
			continue
		}
		log.Printf("Found server trigger prototype with id : %s", trigger.TriggerID)
		triggerTerraform["local"] = false
		triggerTerraform["trigger_id"] = trigger.TriggerID
		triggersTerraform = append(triggersTerraform, triggerTerraform)
	}
	return triggersTerraform, nil
}

func getTerraformTemplateTriggerPrototypes(d *schema.ResourceData, api *zabbix.API) ([]interface{}, error) {
	params := zabbix.Params{
		"output": "extend",
		"discoveryids": []string{
			d.Get("lld_rule_id").(string),
		},
		"inherited": false,
	}
	triggers, err := api.TriggerPrototypesGet(params)
	if err != nil {
		return nil, err
	}

	triggersTerraform := make([]interface{}, len(triggers))
	for i, trigger := range triggers {
		var triggerTerraform = make(map[string]interface{})

		triggerTerraform["local"] = true
		triggerTerraform["trigger_id"] = trigger.TriggerID
		triggersTerraform[i] = triggerTerraform
	}
	return triggersTerraform, nil
}

func updateZabbixTemplateItemPrototypes(d *schema.ResourceData, api *zabbix.API) error {
	if d.HasChange("item_prototype") {
		oldV, newV := d.GetChange("item_prototype")
		oldItems := oldV.(*schema.Set).List()
		newItems := newV.(*schema.Set).List()
		var deletedItems []string
		templatedItems, err := api.ItemPrototypesGet(zabbix.Params{
			"discoveryids": []string{
				d.Get("lld_rule_id").(string),
			},
			"inherited": true,
		})

		if err != nil {
			return err
		}
		log.Printf("[DEBUG] Found templated item prototype %#v", templatedItems)
		for _, oldItem := range oldItems {
			oldItemValue := oldItem.(map[string]interface{})
			exist := false

			if oldItemValue["local"] == true {
				continue
			}

			for _, newItem := range newItems {
				newItemValue := newItem.(map[string]interface{})
				if newItemValue["item_id"].(string) == oldItemValue["item_id"].(string) {
					exist = true
				}
			}

			if !exist {
				templated := false

				for _, templatedItem := range templatedItems {
					if templatedItem.ItemID == oldItemValue["item_id"].(string) {
						templated = true
						break
					}
				}
				if !templated {
					deletedItems = append(deletedItems, oldItemValue["item_id"].(string))
				}
			}
		}
		if len(deletedItems) > 0 {
			log.Printf("[DEBUG] template link will delete item prototype with ids : %#v", deletedItems)
			_, err := api.ItemPrototypesDeleteIDs(deletedItems)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func updateZabbixTemplateTriggerPrototypes(d *schema.ResourceData, api *zabbix.API) error {
	if d.HasChange("trigger_prototype") {
		oldV, newV := d.GetChange("trigger_prototype")
		oldTriggers := oldV.(*schema.Set).List()
		newTriggers := newV.(*schema.Set).List()
		var deletedTriggers []string
		templatedTriggers, err := api.TriggerPrototypesGet(zabbix.Params{
			"output": "extend",
			"discoveryids": []string{
				d.Get("lld_rule_id").(string),
			},
			"inherited": true,
		})

		if err != nil {
			return err
		}
		log.Printf("[DEBUG] found templated trigger prototype %#v", templatedTriggers)
		for _, oldTrigger := range oldTriggers {
			oldTriggerValue := oldTrigger.(map[string]interface{})
			exist := false

			if oldTriggerValue["local"] == true {
				continue
			}

			for _, newTrigger := range newTriggers {
				newTriggerValue := newTrigger.(map[string]interface{})
				if oldTriggerValue["trigger_id"].(string) == newTriggerValue["trigger_id"].(string) {
					exist = true
				}
			}

			if !exist {
				templated := false

				for _, templatedTrigger := range templatedTriggers {
					if templatedTrigger.TriggerID == oldTriggerValue["trigger_id"].(string) {
						templated = true
						break
					}
				}
				if !templated {
					deletedTriggers = append(deletedTriggers, oldTriggerValue["trigger_id"].(string))
				}
			}
		}
		if len(deletedTriggers) > 0 {
			log.Printf("[DEBUG] template link will delete trigger prototype with ids : %#v", deletedTriggers)
			_, err := api.TriggerPrototypesDeleteIDs(deletedTriggers)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func resourceZabbixlldRuleLinkParseID(ID string) (templateID string, itemID []string, triggerID []string, err error) {
	parseID := strings.Split(ID, "_")
	if len(parseID) != 3 {
		err = fmt.Errorf(`Expected id format TEMPLATEID_ITEMID_TRIGGERID,
		if you have multiple ITEMID and TRIGGERID use "." to separate the id`)
		return
	}
	templateID = parseID[0]
	itemID = strings.Split(parseID[1], ".")
	triggerID = strings.Split(parseID[2], ".")
	return
}
