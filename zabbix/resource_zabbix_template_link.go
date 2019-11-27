package provider

import (
	"fmt"
	"log"
	"strings"

	"github.com/claranet/go-zabbix-api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceZabbixTemplateLink() *schema.Resource {
	return &schema.Resource{
		Create: resourceZabbixTemplateLinkCreate,
		Read:   resourceZabbixTemplateLinkRead,
		Exists: resourceZabbixTemplateLinkExist,
		Update: resourceZabbixTemplateLinkUpdate,
		Delete: resourceZabbixTemplateLinkDelete,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				templateID, itemsID, triggersID, err := resourceZabbixTemplateLinkParseID(d.Id())
				if err != nil {
					return nil, err
				}
				d.Set("item", itemsID)
				d.Set("trigger", triggersID)
				d.Set("template_id", templateID)
				d.SetId(randStringNumber(5))
				return []*schema.ResourceData{d}, nil
			}},
		Schema: map[string]*schema.Schema{
			"template_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"item": &schema.Schema{
				Type:     schema.TypeSet,
				Elem:     schemaTemplateItem(),
				Optional: true,
			},
			"trigger": &schema.Schema{
				Type:     schema.TypeSet,
				Elem:     schemaTemplateTrigger(),
				Optional: true,
			},
			"lld_rule": &schema.Schema{
				Type:     schema.TypeSet,
				Elem:     schemaTemplatelldRule(),
				Optional: true,
			},
		},
	}
}

func schemaTemplateItem() *schema.Resource {
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

func schemaTemplateTrigger() *schema.Resource {
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

func schemaTemplatelldRule() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"local": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
			"lld_rule_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceZabbixTemplateLinkCreate(d *schema.ResourceData, meta interface{}) error {
	d.SetId(randStringNumber(5))
	return resourceZabbixTemplateLinkReadTrusted(d, meta)
}

func resourceZabbixTemplateLinkRead(d *schema.ResourceData, meta interface{}) error {
	api := meta.(*zabbix.API)

	itemsTerraform, err := getTerraformTemplateItemsForPlan(d, api)
	if err != nil {
		return err
	}
	log.Print("Resource zabbix template link read :", itemsTerraform)
	d.Set("item", itemsTerraform)

	triggersTerraform, err := getTerraformTemplateTriggersForPlan(d, api)
	if err != nil {
		return err
	}
	d.Set("trigger", triggersTerraform)
	return nil
}

func resourceZabbixTemplateLinkReadTrusted(d *schema.ResourceData, meta interface{}) error {
	api := meta.(*zabbix.API)

	itemsTerraform, err := getTerraformTemplateItems(d, api)
	if err != nil {
		return err
	}
	d.Set("item", itemsTerraform)

	triggersTerraform, err := getTerraformTemplateTriggers(d, api)
	if err != nil {
		return err
	}
	d.Set("trigger", triggersTerraform)
	return nil
}

func resourceZabbixTemplateLinkExist(d *schema.ResourceData, meta interface{}) (bool, error) {
	return true, nil
}

func resourceZabbixTemplateLinkUpdate(d *schema.ResourceData, meta interface{}) error {
	api := meta.(*zabbix.API)

	err := updateZabbixTemplateItem(d, api)
	if err != nil {
		return err
	}
	err = updateZabbixTemplateTrigger(d, api)
	if err != nil {
		return err
	}
	return resourceZabbixTemplateLinkReadTrusted(d, meta)
}

func resourceZabbixTemplateLinkDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func getTerraformTemplateItemsForPlan(d *schema.ResourceData, api *zabbix.API) ([]interface{}, error) {
	params := zabbix.Params{
		"output": "extend",
		"templateids": []string{
			d.Get("template_id").(string),
		},
		"inherited": false,
	}
	items, err := api.ItemsGet(params)
	if err != nil {
		return nil, err
	}

	itemList := d.Get("item").(*schema.Set).List()
	itemLocal := make(map[string]bool)
	var itemsTerraform []interface{}

	for _, item := range itemList {
		var itemTerraform = make(map[string]interface{})
		value := item.(map[string]interface{})

		log.Printf("Found local item with id : %s", value["item_id"].(string))
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
		log.Printf("Found server item with id : %s", item.ItemID)
		itemTerraform["local"] = false
		itemTerraform["item_id"] = item.ItemID
		itemsTerraform = append(itemsTerraform, itemTerraform)
	}
	return itemsTerraform, nil
}

func getTerraformTemplateItems(d *schema.ResourceData, api *zabbix.API) ([]interface{}, error) {
	params := zabbix.Params{
		"output": "extend",
		"templateids": []string{
			d.Get("template_id").(string),
		},
		"inherited": false,
	}
	items, err := api.ItemsGet(params)
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

func getTerraformTemplateTriggersForPlan(d *schema.ResourceData, api *zabbix.API) ([]interface{}, error) {
	params := zabbix.Params{
		"output": "extend",
		"templateids": []string{
			d.Get("template_id").(string),
		},
		"inherited": false,
	}
	triggers, err := api.TriggersGet(params)
	if err != nil {
		return nil, err
	}

	triggerList := d.Get("trigger").(*schema.Set).List()
	triggerLocal := make(map[string]bool)
	var triggersTerraform []interface{}
	for _, trigger := range triggerList {
		triggerTerraform := make(map[string]interface{})
		value := trigger.(map[string]interface{})

		log.Printf("Found local trigger with id : %s", value["trigger_id"].(string))
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
		log.Printf("Found server trigger with id : %s", trigger.TriggerID)
		triggerTerraform["local"] = false
		triggerTerraform["trigger_id"] = trigger.TriggerID
		triggersTerraform = append(triggersTerraform, triggerTerraform)
	}
	return triggersTerraform, nil
}

func getTerraformTemplateTriggers(d *schema.ResourceData, api *zabbix.API) ([]interface{}, error) {
	params := zabbix.Params{
		"output": "extend",
		"templateids": []string{
			d.Get("template_id").(string),
		},
		"inherited": false,
	}
	triggers, err := api.TriggersGet(params)
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

func getTerraformTemplatelldRulesForPlan(d *schema.ResourceData, api *zabbix.API) ([]interface{}, error) {
	params := zabbix.Params{
		"output": "extend",
		"templateids": []string{
			d.Get("template_id").(string),
		},
		"inherited": false,
	}
	lldRules, err := api.DiscoveryRulesGet(params)
	if err != nil {
		return nil, err
	}

	lldRuleList := d.Get("lld_rule").(*schema.Set).List()
	lldRuleLocal := make(map[string]bool)
	var lldRulesTerraform []interface{}
	for _, lldRule := range lldRuleList {
		lldRuleTerraform := make(map[string]interface{})
		value := lldRule.(map[string]interface{})

		log.Printf("Found local lldRule with id : %s", value["lld_rule_id"].(string))
		lldRuleLocal[value["lldRule_id"].(string)] = true
		lldRuleTerraform["lld_rule_id"] = value["lld_rule_id"].(string)
		lldRuleTerraform["local"] = true
		lldRulesTerraform = append(lldRulesTerraform, lldRuleTerraform)
	}
	for _, lldRule := range lldRules {
		var lldRuleTerraform = make(map[string]interface{})

		if lldRuleLocal[lldRule.ItemID] {
			continue
		}
		log.Printf("Found server lldRule with id : %s", lldRule.ItemID)
		lldRuleTerraform["local"] = false
		lldRuleTerraform["lld_rule_id"] = lldRule.ItemID
		lldRulesTerraform = append(lldRulesTerraform, lldRuleTerraform)
	}
	return lldRulesTerraform, nil
}

func getTerraformTemplatelldRules(d *schema.ResourceData, api *zabbix.API) ([]interface{}, error) {
	params := zabbix.Params{
		"output": "extend",
		"templateids": []string{
			d.Get("template_id").(string),
		},
		"inherited": false,
	}
	lldRules, err := api.DiscoveryRulesGet(params)
	if err != nil {
		return nil, err
	}

	lldRulesTerraform := make([]interface{}, len(lldRules))
	for i, lldRule := range lldRules {
		var lldRuleTerraform = make(map[string]interface{})

		lldRuleTerraform["local"] = true
		lldRuleTerraform["lld_rule_id"] = lldRule.ItemID
		lldRulesTerraform[i] = lldRuleTerraform
	}
	return lldRulesTerraform, nil
}

func updateZabbixTemplateItem(d *schema.ResourceData, api *zabbix.API) error {
	if d.HasChange("item") {
		oldV, newV := d.GetChange("item")
		oldItems := oldV.(*schema.Set).List()
		newItems := newV.(*schema.Set).List()
		var deletedItems []string
		templatedItems, err := api.ItemsGet(zabbix.Params{
			"templateids": []string{
				d.Get("template_id").(string),
			},
			"inherited": true,
		})

		if err != nil {
			return err
		}
		log.Printf("[DEBUG] Found templated item %#v", templatedItems)
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
			log.Printf("[DEBUG] template link will delete item with ids : %#v", deletedItems)
			if err := api.ItemsDeleteByIds(deletedItems); err != nil {
				return err
			}
		}
	}
	return nil
}

func updateZabbixTemplateTrigger(d *schema.ResourceData, api *zabbix.API) error {
	if d.HasChange("trigger") {
		oldV, newV := d.GetChange("trigger")
		oldTriggers := oldV.(*schema.Set).List()
		newTriggers := newV.(*schema.Set).List()
		var deletedTriggers []string
		templatedTriggers, err := api.TriggersGet(zabbix.Params{
			"output": "extend",
			"templateids": []string{
				d.Get("template_id").(string),
			},
			"inherited": true,
		})

		if err != nil {
			return err
		}
		log.Printf("[DEBUG] found templated trigger %#v", templatedTriggers)
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
			log.Printf("[DEBUG] template link will delete trigger with ids : %#v", deletedTriggers)
			if err := api.ItemsDeleteByIds(deletedTriggers); err != nil {
				return err
			}
		}
	}
	return nil
}

func updateZabbixTemplateDiscoveryRule(d *schema.ResourceData, api *zabbix.API) error {
	if d.HasChange("lld_dule") {
		oldV, newV := d.GetChange("lld_rule")
		oldlldRules := oldV.(*schema.Set).List()
		newlldRules := newV.(*schema.Set).List()
		var deletedlldRules []string
		templatedlldRules, err := api.DiscoveryRulesGet(zabbix.Params{
			"output": "extend",
			"templateids": []string{
				d.Get("template_id").(string),
			},
			"inherited": true,
		})

		if err != nil {
			return err
		}
		log.Printf("[DEBUG] found templated lldRule %#v", templatedlldRules)
		for _, oldlldRule := range oldlldRules {
			oldlldRuleValue := oldlldRule.(map[string]interface{})
			exist := false

			if oldlldRuleValue["local"] == true {
				continue
			}

			for _, newlldRule := range newlldRules {
				newlldRuleValue := newlldRule.(map[string]interface{})
				if oldlldRuleValue["lld_rule_id"].(string) == newlldRuleValue["lld_rule_id"].(string) {
					exist = true
				}
			}

			if !exist {
				templated := false

				for _, templatedlldRule := range templatedlldRules {
					if templatedlldRule.ItemID == oldlldRuleValue["lld_rule_id"].(string) {
						templated = true
						break
					}
				}
				if !templated {
					deletedlldRules = append(deletedlldRules, oldlldRuleValue["lld_rule_id"].(string))
				}
			}
		}
		if len(deletedlldRules) > 0 {
			log.Printf("[DEBUG] template link will delete lldRule with ids : %#v", deletedlldRules)
			if err := api.DiscoveryRulesDeletesByIDs(deletedlldRules); err != nil {
				return err
			}
		}
	}
	return nil
}

func resourceZabbixTemplateLinkParseID(ID string) (templateID string, itemID []string, triggerID []string, err error) {
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
