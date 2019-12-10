package provider

import (
	"fmt"
	"log"
	"strings"

	"github.com/claranet/go-zabbix-api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceZabbixTriggerPrototype() *schema.Resource {
	return &schema.Resource{
		Create: resourceZabbixTriggerPrototypeCreate,
		Read:   resourceZabbixTriggerPrototypeRead,
		Exists: resourceZabbixTriggerPrototypeExist,
		Update: resourceZabbixTriggerPrototypeUpdate,
		Delete: resourceZabbixTriggerPrototypeDelete,
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

func resourceZabbixTriggerPrototypeCreate(d *schema.ResourceData, meta interface{}) error {
	trigger := createTriggerPrototypeObj(d)

	return createRetry(d, meta, createTriggerPrototype, trigger, resourceZabbixTriggerPrototypeRead)
}

func resourceZabbixTriggerPrototypeRead(d *schema.ResourceData, meta interface{}) error {
	api := meta.(*zabbix.API)

	params := zabbix.Params{
		"output":             "extend",
		"selectDependencies": "extend",
		"selectFunctions":    "extend",
		"selectItems":        "extend",
		"triggerids":         d.Id(),
	}
	res, err := api.TriggerPrototypesGet(params)
	if err != nil {
		return err
	}
	if len(res) != 1 {
		return fmt.Errorf("Expected one result got : %d", len(res))
	}
	trigger := res[0]
	err = getTriggerPrototypeExpression(&trigger, api)
	log.Printf("[DEBUG] trigger expression: %s", trigger.Expression)
	d.Set("description", trigger.Description)
	d.Set("expression", trigger.Expression)
	d.Set("priority", trigger.Priority)
	d.Set("status", trigger.Status)

	var dependencies []string
	for _, dependencie := range trigger.Dependencies {
		dependencies = append(dependencies, dependencie.TriggerID)
	}
	d.Set("dependencies", dependencies)
	return nil
}

func resourceZabbixTriggerPrototypeExist(d *schema.ResourceData, meta interface{}) (bool, error) {
	api := meta.(*zabbix.API)

	_, err := api.TriggerPrototypeGetByID(d.Id())
	if err != nil {
		if strings.Contains(err.Error(), "Expected exactly one result") {
			log.Printf("Trigger prototype with id %s doesn't exist", d.Id())
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func resourceZabbixTriggerPrototypeUpdate(d *schema.ResourceData, meta interface{}) error {
	trigger := createTriggerPrototypeObj(d)
	trigger.TriggerID = d.Id()
	if !d.HasChange("dependencies") {
		trigger.Dependencies = nil
	}
	return createRetry(d, meta, updateTriggerPrototype, trigger, resourceZabbixTriggerPrototypeRead)
}

func resourceZabbixTriggerPrototypeDelete(d *schema.ResourceData, meta interface{}) error {
	api := meta.(*zabbix.API)

	return deleteRetry(d.Id(), getTriggerPrototypeParentID, api.TriggerPrototypesDeleteIDs, api)
}

func createTriggerPrototypeDependencies(d *schema.ResourceData) zabbix.TriggerPrototypes {
	size := d.Get("dependencies.#").(int)
	dependencies := make(zabbix.TriggerPrototypes, size)

	terraformDependencies := d.Get("dependencies").(*schema.Set)
	for i, terraformDependencie := range terraformDependencies.List() {
		dependencies[i].TriggerID = terraformDependencie.(string)
	}
	return dependencies
}

func createTriggerPrototypeObj(d *schema.ResourceData) zabbix.TriggerPrototype {
	return zabbix.TriggerPrototype{
		Description:  d.Get("description").(string),
		Expression:   d.Get("expression").(string),
		Priority:     zabbix.SeverityType(d.Get("priority").(int)),
		Status:       zabbix.StatusType(d.Get("status").(int)),
		Dependencies: createTriggerPrototypeDependencies(d),
	}
}

func getTriggerPrototypeExpression(trigger *zabbix.TriggerPrototype, api *zabbix.API) error {
	for _, function := range trigger.Functions {
		var item zabbix.ItemPrototype

		items, err := api.ItemPrototypesGet(zabbix.Params{
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
		if len(item.Hosts) != 1 {
			return fmt.Errorf("Expected one parent host for item with id %s, and got : %d", function.ItemID, len(item.Hosts))
		}
		idstr := fmt.Sprintf("{%s}", function.FunctionID)
		expendValue := fmt.Sprintf("{%s:%s.%s(%s)}", item.Hosts[0].Host, item.Key, function.Function, function.Parameter)
		trigger.Expression = strings.Replace(trigger.Expression, idstr, expendValue, 1)
	}
	return nil
}

func getTriggerPrototypeParentID(api *zabbix.API, id string) (string, error) {
	triggers, err := api.TriggerPrototypesGet(zabbix.Params{
		"ouput":       "extend",
		"selectHosts": "extend",
		"triggerids":  id,
	})
	if err != nil {
		return "", err
	}
	if len(triggers) != 1 {
		return "", fmt.Errorf("Expected one trigger prototype and got %d trigger prototype", len(triggers))
	}
	if len(triggers[0].ParentHosts) != 1 {
		return "", fmt.Errorf("Expected one parent for trigger prototype %s and got %d", id, len(triggers[0].ParentHosts))
	}
	return triggers[0].ParentHosts[0].HostID, nil
}

func createTriggerPrototype(trigger interface{}, api *zabbix.API) (id string, err error) {
	triggers := zabbix.TriggerPrototypes{trigger.(zabbix.TriggerPrototype)}

	err = api.TriggerPrototypesCreate(triggers)
	if err != nil {
		return
	}
	id = triggers[0].TriggerID
	return
}

func updateTriggerPrototype(trigger interface{}, api *zabbix.API) (id string, err error) {
	triggers := zabbix.TriggerPrototypes{trigger.(zabbix.TriggerPrototype)}

	err = api.TriggerPrototypesUpdate(triggers)
	if err != nil {
		return
	}
	id = triggers[0].TriggerID
	return
}
