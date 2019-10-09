package provider

import (
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

	itemsTerraform, err := getTerraformTemplateLinks(d, api)
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

	err := updateZabbixTemplateLink(d, api)
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

func getTerraformTemplateLinks(d *schema.ResourceData, api *zabbix.API) ([]string, error) {
	params := zabbix.Params{
		"output": "extend",
		"templateids": []string{
			d.Get("template_id").(string),
		},
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

func updateZabbixTemplateLink(d *schema.ResourceData, api *zabbix.API) error {
	localItems := d.Get("item").(*schema.Set)

	params := zabbix.Params{
		"output": "extend",
		"templateids": []string{
			d.Get("template_id").(string),
		},
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
