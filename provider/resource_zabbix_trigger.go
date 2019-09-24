package provider

import (
	"fmt"
	"strings"

	"github.com/claranet/go-zabbix-api"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceZabbixTrigger() *schema.Resource {
	return &schema.Resource{
		Create: resourceZabbixTriggerCreate,
		Read:   resourceZabbixTriggerRead,
		Exists: resourceZabbixTriggerExist,
		Update: resourceZabbixTriggerUpdate,
		Delete: resourceZabbixTriggerDelete,

		Schema: map[string]*schema.Schema{
			"trigger_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "(readonly) ID of the trigger",
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"expression": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"comment": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"priority": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},
			"status": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},
		},
	}
}

func resourceZabbixTriggerCreate(d *schema.ResourceData, meta interface{}) error {
	api := meta.(*zabbix.API)

	triggers := zabbix.Triggers{createTriggerObj(d)}
	err := api.TriggersCreate(triggers)
	if err != nil {
		return err
	}
	d.SetId(triggers[0].TriggerID)
	return resourceZabbixTriggerRead(d, meta)
}

func resourceZabbixTriggerRead(d *schema.ResourceData, meta interface{}) error {
	api := meta.(*zabbix.API)

	params := zabbix.Params{
		"output":           "extend",
		"expandExpression": true,
		"triggerids":       d.Id(),
	}
	res, err := api.TriggersGet(params)
	if err != nil {
		return err
	}
	if len(res) != 1 {
		return fmt.Errorf("Expected one result got : %d", len(res))
	}
	item := res[0]
	d.Set("trigger_id", item.TriggerID)
	d.Set("description", item.Description)
	d.Set("expression", item.Expression)
	d.Set("comment", item.Comments)
	d.Set("priority", item.Priority)
	d.Set("status", item.Status)
	return nil
}

func resourceZabbixTriggerExist(d *schema.ResourceData, meta interface{}) (bool, error) {
	api := meta.(*zabbix.API)

	_, err := api.TriggerGetByID(d.Id())
	if err != nil {
		if strings.Contains(err.Error(), "Expected exactly one result") {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func resourceZabbixTriggerUpdate(d *schema.ResourceData, meta interface{}) error {
	api := meta.(*zabbix.API)

	triggers := zabbix.Triggers{createTriggerObj(d)}
	err := api.TriggersUpdate(triggers)
	if err != nil {
		return err
	}
	return resourceZabbixTriggerRead(d, meta)
}

func resourceZabbixTriggerDelete(d *schema.ResourceData, meta interface{}) error {
	api := meta.(*zabbix.API)
	return api.TriggersDeleteByIds([]string{d.Id()})
}

func createTriggerObj(d *schema.ResourceData) zabbix.Trigger {
	return zabbix.Trigger{
		TriggerID:   d.Get("trigger_id").(string),
		Description: d.Get("description").(string),
		Expression:  d.Get("expression").(string),
		Comments:    d.Get("comment").(string),
		Priority:    zabbix.SeverityType(d.Get("priority").(int)),
		Status:      zabbix.StatusType(d.Get("status").(int)),
	}
}
