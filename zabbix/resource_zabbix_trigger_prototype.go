package provider

import (
	"fmt"
	"log"
	"strings"

	"github.com/claranet/go-zabbix-api"
	"github.com/hashicorp/terraform/helper/schema"
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
	api := meta.(*zabbix.API)

	triggers := zabbix.TriggerPrototypes{createTriggerPrototypeObj(d)}
	err := api.TriggerPrototypesCreate(triggers)
	if err != nil {
		return err
	}
	d.SetId(triggers[0].TriggerID)
	return resourceZabbixTriggerPrototypeRead(d, meta)
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
	d.Set("trigger_id", trigger.TriggerID)
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
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func resourceZabbixTriggerPrototypeUpdate(d *schema.ResourceData, meta interface{}) error {
	api := meta.(*zabbix.API)

	triggers := zabbix.TriggerPrototypes{createTriggerPrototypeObj(d)}
	if !d.HasChange("dependencies") {
		triggers[0].Dependencies = nil
	}
	err := api.TriggerPrototypesUpdate(triggers)
	if err != nil {
		return err
	}
	return resourceZabbixTriggerPrototypeRead(d, meta)
}

func resourceZabbixTriggerPrototypeDelete(d *schema.ResourceData, meta interface{}) error {
	api := meta.(*zabbix.API)

	// triggers, err := api.TriggerPrototypesGet(zabbix.Params{
	// 	"ouput":       "extend",
	// 	"selectHosts": "extend",
	// 	"triggerids":  d.Id(),
	// })
	// if err != nil {
	// 	return fmt.Errorf("%s, with trigger %s", err.Error(), d.Id())
	// }
	// if len(triggers) != 1 {
	// 	return fmt.Errorf("Expected one item and got %d items", len(triggers))
	// }
	// trigger := triggers[0]

	// templates, err := api.TemplatesGet(zabbix.Params{
	// 	"output":            "extend",
	// 	"parentTemplateids": trigger.ParentHosts[0].HostID,
	// })

	// triggerids, err := api.TriggersDeleteIDs([]string{d.Id()})
	// if err != nil {
	// 	return fmt.Errorf("%s, with trigger %s", err.Error(), d.Id())
	// }
	// if len(triggerids) != len(templates)+1 {
	// 	return fmt.Errorf("Expected to delete %d trigger and %d were delete", len(templates)+1, len(triggerids))
	// }
	_, err := api.TriggerPrototypesDeleteIDs([]string{d.Id()})
	return err
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
		TriggerID:    d.Get("trigger_id").(string),
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
