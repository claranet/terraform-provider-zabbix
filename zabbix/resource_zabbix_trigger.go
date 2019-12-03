package provider

import (
	"fmt"
	"log"
	"strings"

	"github.com/claranet/go-zabbix-api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceZabbixTrigger() *schema.Resource {
	return &schema.Resource{
		Create: resourceZabbixTriggerCreate,
		Read:   resourceZabbixTriggerRead,
		Exists: resourceZabbixTriggerExist,
		Update: resourceZabbixTriggerUpdate,
		Delete: resourceZabbixTriggerDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
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
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(int)
					if v < 0 || v > 5 {
						errs = append(errs, fmt.Errorf("%q, must be between 0 and 5 inclusive, got %d", key, v))
					}
					return
				},
			},
			"status": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(int)
					if v < 0 || v > 1 {
						errs = append(errs, fmt.Errorf("%q, must be between 0 and 1 inclusive, got %d", key, v))
					}
					return
				},
			},
			"dependencies": &schema.Schema{
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "ID of the trigger it depands",
			},
		},
	}
}

func resourceZabbixTriggerCreate(d *schema.ResourceData, meta interface{}) error {
	trigger := createTriggerObj(d)

	return createRetry(d, meta, createTrigger, trigger, resourceZabbixTriggerRead)
}

func resourceZabbixTriggerRead(d *schema.ResourceData, meta interface{}) error {
	api := meta.(*zabbix.API)

	params := zabbix.Params{
		"output":             "extend",
		"selectDependencies": "extend",
		"selectFunctions":    "extend",
		"selectItems":        "extend",
		"triggerids":         d.Id(),
	}
	res, err := api.TriggersGet(params)
	if err != nil {
		return err
	}
	if len(res) != 1 {
		return fmt.Errorf("Expected one result got : %d", len(res))
	}
	trigger := res[0]
	err = getTriggerExpression(&trigger, api)
	log.Printf("[DEBUG] trigger expression: %s", trigger.Expression)
	d.Set("description", trigger.Description)
	d.Set("expression", trigger.Expression)
	if trigger.Comments != "" {
		d.Set("comment", trigger.Comments)
	}
	d.Set("priority", trigger.Priority)
	if trigger.Status != 0 {
		d.Set("status", trigger.Status)
	} else {
		d.Set("value", 0)
	}

	var dependencies []string
	for _, dependencie := range trigger.Dependencies {
		dependencies = append(dependencies, dependencie.TriggerID)
	}
	d.Set("dependencies", dependencies)
	return nil
}

func resourceZabbixTriggerExist(d *schema.ResourceData, meta interface{}) (bool, error) {
	api := meta.(*zabbix.API)

	_, err := api.TriggerGetByID(d.Id())
	if err != nil {
		if strings.Contains(err.Error(), "Expected exactly one result") {
			log.Printf("Trigger with id %s doesn't exist", d.Id())
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func resourceZabbixTriggerUpdate(d *schema.ResourceData, meta interface{}) error {
	trigger := createTriggerObj(d)

	trigger.TriggerID = d.Id()
	if !d.HasChange("dependencies") {
		trigger.Dependencies = nil
	}
	return createRetry(d, meta, updateTrigger, trigger, resourceZabbixTriggerRead)
}

func resourceZabbixTriggerDelete(d *schema.ResourceData, meta interface{}) error {
	api := meta.(*zabbix.API)

	return deleteRetry(d.Id(), getTriggerParentID, api.TriggersDeleteIDs, api)
}

func createTriggerDependencies(d *schema.ResourceData) zabbix.Triggers {
	size := d.Get("dependencies.#").(int)
	dependencies := make(zabbix.Triggers, size)

	terraformDependencies := d.Get("dependencies").(*schema.Set)
	for i, terraformDependencie := range terraformDependencies.List() {
		dependencies[i].TriggerID = terraformDependencie.(string)
	}
	return dependencies
}

func createTriggerObj(d *schema.ResourceData) zabbix.Trigger {
	return zabbix.Trigger{
		Description:  d.Get("description").(string),
		Expression:   d.Get("expression").(string),
		Comments:     d.Get("comment").(string),
		Priority:     zabbix.SeverityType(d.Get("priority").(int)),
		Status:       zabbix.StatusType(d.Get("status").(int)),
		Dependencies: createTriggerDependencies(d),
	}
}

func getTriggerExpression(trigger *zabbix.Trigger, api *zabbix.API) error {
	for _, function := range trigger.Functions {
		var item zabbix.Item

		items, err := api.ItemsGet(zabbix.Params{
			"output":      "extend",
			"selectHosts": "extend",
			"itemids":     function.ItemID,
		})
		if err != nil {
			return err
		}
		if len(items) != 1 {
			return fmt.Errorf("Expected one item with id : %s and got : %d", function.ItemID, len(items))
		}
		item = items[0]
		if len(item.ItemParent) != 1 {
			return fmt.Errorf("Expected one parent host for item with id %s, and got : %d", function.ItemID, len(item.ItemParent))
		}
		idstr := fmt.Sprintf("{%s}", function.FunctionID)
		expendValue := fmt.Sprintf("{%s:%s.%s(%s)}", item.ItemParent[0].Host, item.Key, function.Function, function.Parameter)
		trigger.Expression = strings.Replace(trigger.Expression, idstr, expendValue, 1)
	}
	return nil
}

func getTriggerParentID(api *zabbix.API, id string) (string, error) {
	triggers, err := api.TriggersGet(zabbix.Params{
		"ouput":       "extend",
		"selectHosts": "extend",
		"triggerids":  id,
	})
	if err != nil {
		return "", err
	}
	if len(triggers) != 1 {
		return "", fmt.Errorf("Expected one item and got %d items", len(triggers))
	}
	if len(triggers[0].ParentHosts) != 1 {
		return "", fmt.Errorf("Expected one parent for item %s and got %d", id, len(triggers[0].ParentHosts))
	}
	return triggers[0].ParentHosts[0].HostID, nil
}

func createTrigger(trigger interface{}, api *zabbix.API) (id string, err error) {
	triggers := zabbix.Triggers{trigger.(zabbix.Trigger)}

	err = api.TriggersCreate(triggers)
	if err != nil {
		return
	}
	id = triggers[0].TriggerID
	return
}

func updateTrigger(trigger interface{}, api *zabbix.API) (id string, err error) {
	triggers := zabbix.Triggers{trigger.(zabbix.Trigger)}

	err = api.TriggersUpdate(triggers)
	if err != nil {
		return
	}
	id = triggers[0].TriggerID
	return
}
