package provider

import (
	"fmt"
	"strings"

	"github.com/claranet/go-zabbix-api"
	"github.com/hashicorp/terraform/helper/schema"
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
				Elem:     &schema.Schema{Type: schema.TypeString},
				Required: true,
			},
			"trigger": &schema.Schema{
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Required: true,
			},
		},
	}
}

func resourceZabbixTemplateLinkCreate(d *schema.ResourceData, meta interface{}) error {
	d.SetId(randStringNumber(5))
	return nil
}

func resourceZabbixTemplateLinkRead(d *schema.ResourceData, meta interface{}) error {
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
	return resourceZabbixTemplateLinkRead(d, meta)
}

func resourceZabbixTemplateLinkDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func getTerraformTemplateItems(d *schema.ResourceData, api *zabbix.API) ([]string, error) {
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

	itemsTerraform := make([]string, len(items))
	for i, item := range items {
		itemsTerraform[i] = item.ItemID
	}
	return itemsTerraform, nil
}

func getTerraformTemplateTriggers(d *schema.ResourceData, api *zabbix.API) ([]string, error) {
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

	TriggersTerraform := make([]string, len(triggers))
	for i, trigger := range triggers {
		TriggersTerraform[i] = trigger.TriggerID
	}
	return TriggersTerraform, nil
}

func updateZabbixTemplateItem(d *schema.ResourceData, api *zabbix.API) error {
	localItems := d.Get("item").(*schema.Set)

	params := zabbix.Params{
		"output": "extend",
		"templateids": []string{
			d.Get("template_id").(string),
		},
		"inherited": false,
	}
	serverItems, err := api.ItemsGet(params)
	if err != nil {
		return err
	}

	for _, serverItem := range serverItems {
		exist := false

		for _, localItem := range localItems.List() {
			if localItem.(string) == serverItem.ItemID {
				exist = true
			}
		}

		if !exist {
			err = api.ItemsDelete(zabbix.Items{serverItem})
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func updateZabbixTemplateTrigger(d *schema.ResourceData, api *zabbix.API) error {
	localTriggers := d.Get("trigger").(*schema.Set)

	params := zabbix.Params{
		"output": "extend",
		"templateids": []string{
			d.Get("template_id").(string),
		},
		"inherited": false,
	}
	serverTriggers, err := api.TriggersGet(params)
	if err != nil {
		return err
	}

	for _, serverTrigger := range serverTriggers {
		exist := false

		for _, localItem := range localTriggers.List() {
			if localItem.(string) == serverTrigger.TriggerID {
				exist = true
			}
		}

		if !exist {
			err = api.TriggersDelete(zabbix.Triggers{serverTrigger})
			if err != nil {
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
