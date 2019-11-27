package provider

import (
	"fmt"
	"log"
	"strings"

	"github.com/claranet/go-zabbix-api"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceZabbixlldRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceZabbixlldRuleCreate,
		Read:   resourceZabbixlldRuleRead,
		Exists: resourceZabbixlldRuleExist,
		Update: resourceZabbixlldRuleUpdate,
		Delete: resourceZabbixlldRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"delay": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"host_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"interface_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"key": &schema.Schema{Type: schema.TypeString,
				Required: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"type": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"filter": &schema.Schema{
				Type:     schema.TypeSet,
				MaxItems: 1,
				Elem:     resourcelldRuleFilter(),
				Required: true,
			},
		},
	}
}

func resourcelldRuleFilter() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"condition": &schema.Schema{
				Type:     schema.TypeSet,
				Elem:     resourcelldRuleFilterCondition(),
				Required: true,
			},
			"eval_type": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"formula": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
		},
	}
}

func resourcelldRuleFilterCondition() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"macro": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"value": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"operator": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  "8",
			},
		},
	}
}

func resourceZabbixlldRuleCreate(d *schema.ResourceData, meta interface{}) error {
	api := meta.(*zabbix.API)
	rule := createlldRuleObject(d)
	rules := zabbix.LLDRules{rule}

	err := api.DiscoveryRulesCreate(rules)
	if err != nil {
		return err
	}
	d.SetId(rules[0].ItemID)
	return resourceZabbixlldRuleRead(d, meta)
}

func resourceZabbixlldRuleRead(d *schema.ResourceData, meta interface{}) error {
	api := meta.(*zabbix.API)
	params := zabbix.Params{
		"itemids":      d.Id(),
		"output":       "extend",
		"selectFilter": "extend",
		"inherited":    false,
	}

	lldRules, err := api.DiscoveryRulesGet(params)
	if err != nil {
		return err
	}
	if len(lldRules) != 1 {
		return fmt.Errorf("Expected one low level discovery rule with id %s and got %d rules", d.Id(), len(lldRules))
	}
	lldRule := lldRules[0]

	d.Set("delay", lldRule.Delay)
	d.Set("host_id", lldRule.HostID)
	d.Set("interface_id", lldRule.InterfaceID)
	d.Set("key", lldRule.Key)
	d.Set("name", lldRule.Name)
	d.Set("type", lldRule.Type)

	var terraformConditions []interface{}
	for _, condition := range lldRule.Filter.Conditions {
		terraformCondition := map[string]interface{}{}

		terraformCondition["macro"] = condition.LLDMacro
		terraformCondition["value"] = condition.Value
		terraformCondition["operator"] = condition.Operator
		terraformConditions = append(terraformConditions, terraformCondition)
	}

	filter := map[string]interface{}{}
	filter["condition"] = terraformConditions
	filter["eval_type"] = lldRule.Filter.EvalType
	filter["formula"] = lldRule.Filter.Formula

	d.Set("filter", []interface{}{filter})
	return nil
}

func resourceZabbixlldRuleExist(d *schema.ResourceData, meta interface{}) (bool, error) {
	api := meta.(*zabbix.API)

	_, err := api.DiscoveryRulesGetByID(d.Id())
	if err != nil {
		log.Printf("lld rule exist error : %s", err.Error())
		if strings.Contains(err.Error(), "Expected exactly one result") {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func resourceZabbixlldRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	api := meta.(*zabbix.API)
	rule := createlldRuleObject(d)
	rule.ItemID = d.Id()
	rules := zabbix.LLDRules{rule}

	err := api.DiscoveryRulesUpdate(rules)
	if err != nil {
		return err
	}
	return resourceZabbixlldRuleRead(d, meta)
}

func resourceZabbixlldRuleDelete(d *schema.ResourceData, meta interface{}) error {
	api := meta.(*zabbix.API)

	err := api.DiscoveryRulesDeletesByIDs([]string{d.Id()})
	return err
}

func createlldRuleObject(d *schema.ResourceData) zabbix.LLDRule {
	return zabbix.LLDRule{
		Delay:       d.Get("delay").(int),
		HostID:      d.Get("host_id").(string),
		InterfaceID: d.Get("interface_id").(string),
		Key:         d.Get("key").(string),
		Name:        d.Get("name").(string),
		Type:        zabbix.ItemType(d.Get("type").(int)),
		Filter:      createlldRuleConditionObject(d),
	}
}

func createlldRuleConditionObject(d *schema.ResourceData) zabbix.LLDRuleFilter {
	filters := d.Get("filter").(*schema.Set)
	filter := filters.List()[0].(map[string]interface{})
	conditions := filter["condition"].(*schema.Set)
	var filterObject zabbix.LLDRuleFilter

	filterObject.EvalType = filter["eval_type"].(int)
	filterObject.Formula = filter["formula"].(string)
	for _, condition := range conditions.List() {
		value := condition.(map[string]interface{})
		cond := zabbix.LLDRulesFilterCondition{
			LLDMacro: value["macro"].(string),
			Value:    value["value"].(string),
			Operator: value["operator"].(int),
		}
		filterObject.Conditions = append(filterObject.Conditions, cond)
	}
	return filterObject
}
