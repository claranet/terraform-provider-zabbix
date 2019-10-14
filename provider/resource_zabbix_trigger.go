package provider

import (
	"fmt"
	"log"
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
	d.Set("trigger_id", trigger.TriggerID)
	log.Printf("trigger expressiob %s", trigger.Expression)
	d.Set("description", trigger.Description)
	d.Set("expression", trigger.Expression)
	d.Set("comment", trigger.Comments)
	d.Set("priority", trigger.Priority)
	d.Set("status", trigger.Status)

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

	triggers, err := api.TriggersGet(zabbix.Params{
		"ouput":       "extend",
		"selectHosts": "extend",
		"triggerids":  d.Id(),
	})
	if err != nil {
		return fmt.Errorf("%s, with trigger %s", err.Error(), d.Id())
	}
	if len(triggers) != 1 {
		return fmt.Errorf("Expected one item and got %d items", len(triggers))
	}
	trigger := triggers[0]

	templates, err := api.TemplatesGet(zabbix.Params{
		"output":            "extend",
		"parentTemplateids": trigger.TriggerParent[0].HostID,
	})

	triggerids, err := api.TriggersDeleteNoError([]string{d.Id()})
	if err != nil {
		return fmt.Errorf("%s, with trigger %s", err.Error(), d.Id())
	}
	if len(triggerids) != len(templates)+1 {
		return fmt.Errorf("Expected to delete %d trigger and %d were delete", len(templates)+1, len(triggerids))
	}
	return nil
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
		TriggerID:    d.Get("trigger_id").(string),
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
