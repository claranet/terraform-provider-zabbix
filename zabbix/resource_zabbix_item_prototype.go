package provider

import (
	"fmt"
	"log"
	"strings"

	"github.com/claranet/go-zabbix-api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/mcuadros/go-version"
)

func resourceZabbixItemPrototype() *schema.Resource {
	return &schema.Resource{
		Create: resourceZabbixItemPrototypeCreate,
		Read:   resourceZabbixItemPrototypeRead,
		Exists: resourceZabbixItemPrototypeExist,
		Update: resourceZabbixItemPrototypeUpdate,
		Delete: resourceZabbixItemPrototypeDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"delay": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"host_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the low level discovery that the item prototype belongs to.",
			},
			"interface_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "0",
			},
			"key": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Item prototype key.",
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the item prototype.",
			},
			"type": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(int)
					if v < 0 || v > 16 {
						errs = append(errs, fmt.Errorf("%q, must be between 0 and 16 inclusive, got %d", key, v))
					}
					return
				},
			},
			"value_type": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(int)
					if v < 0 || v > 4 {
						errs = append(errs, fmt.Errorf("%q, must be between 0 and 4 inclusive, got %d", key, v))
					}
					return
				},
			},
			"rule_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"data_type": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(int)
					if v < 0 || v > 3 {
						errs = append(errs, fmt.Errorf("%q, must be between 0 and 3 inclusive, got %d", key, v))
					}
					return
				},
				Description: "Remove in zabbix 3.4. Data type of the item prototype.",
			},
			"delta": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Remove in zabbix 3.4. Value that will be stored.",
				Default:     0,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(int)
					if v < 0 || v > 2 {
						errs = append(errs, fmt.Errorf("%q, must be between 0 and 2 inclusive, got %d", key, v))
					}
					return
				},
			},
			"description": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description of the item prototype.",
				Default:     "",
			},
			"history": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Number of days to keep item's history data. Default: 90. ",
				DefaultFunc: func() (interface{}, error) {
					if version.Compare(zabbixAPIVersion, "3.4.0", ">=") {
						return "90d", nil
					}
					return "90", nil
				},
			},
			"trends": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Number of days to keep item's trends data. Default: 365. ",
				DefaultFunc: func() (interface{}, error) {
					if version.Compare(zabbixAPIVersion, "3.4.0", ">=") {
						return "365d", nil
					}
					return "365", nil
				},
			},
			"trapper_host": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Allowed hosts. Used only by trapper items. ",
			},
			"status": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     "0",
				Description: "Allowed hosts. Used only by trapper items. ",
			},
		},
	}
}

func createItemPrototypeObject(d *schema.ResourceData, api *zabbix.API) (*zabbix.ItemPrototype, error) {

	item := zabbix.ItemPrototype{
		Delay:        d.Get("delay").(string),
		HostID:       d.Get("host_id").(string),
		InterfaceID:  d.Get("interface_id").(string),
		Key:          d.Get("key").(string),
		Name:         d.Get("name").(string),
		Type:         zabbix.ItemType(d.Get("type").(int)),
		ValueType:    zabbix.ValueType(d.Get("value_type").(int)),
		RuleID:       d.Get("rule_id").(string),
		DataType:     zabbix.DataType(d.Get("data_type").(int)),
		Delta:        zabbix.DeltaType(d.Get("delta").(int)),
		Description:  d.Get("description").(string),
		History:      d.Get("history").(string),
		Trends:       d.Get("trends").(string),
		TrapperHosts: d.Get("trapper_host").(string),
		Status:       d.Get("status").(int),
	}
	return &item, nil
}

func resourceZabbixItemPrototypeCreate(d *schema.ResourceData, meta interface{}) error {
	api := meta.(*zabbix.API)

	item, err := createItemPrototypeObject(d, api)
	if err != nil {
		return err
	}

	return createRetry(d, meta, createItemPrototype, *item, resourceZabbixItemPrototypeRead)
}

func resourceZabbixItemPrototypeRead(d *schema.ResourceData, meta interface{}) error {
	api := meta.(*zabbix.API)

	items, err := api.ItemPrototypesGet(zabbix.Params{
		"itemids":             d.Id(),
		"output":              "extend",
		"selectDiscoveryRule": "extend",
	})
	if err != nil {
		return err
	}
	if len(items) != 1 {
		return fmt.Errorf("Expected one item prototype and got : %d ", len(items))
	}
	item := items[0]

	d.Set("delay", item.Delay)
	d.Set("host_id", item.HostID)
	d.Set("interface_id", item.InterfaceID)
	d.Set("key", item.Key)
	d.Set("name", item.Name)
	d.Set("type", item.Type)
	d.Set("value_type", item.ValueType)
	d.Set("rule_id", item.DiscoveryRule.ItemID)
	d.Set("data_type", item.DataType)
	d.Set("delta", item.Delta)
	d.Set("description", item.Description)
	d.Set("history", item.History)
	d.Set("trends", item.Trends)
	d.Set("trapper_host", item.TrapperHosts)
	d.Set("status", item.Status)

	log.Printf("Item prototype name is %s\n", item.Name)
	return nil
}

func resourceZabbixItemPrototypeExist(d *schema.ResourceData, meta interface{}) (bool, error) {
	api := meta.(*zabbix.API)

	_, err := api.ItemPrototypeGetByID(d.Id())
	if err != nil {
		if strings.Contains(err.Error(), "Expected exactly one result") {
			log.Printf("Item prototype with id %s doesn't exist", d.Id())
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func resourceZabbixItemPrototypeUpdate(d *schema.ResourceData, meta interface{}) error {
	api := meta.(*zabbix.API)

	item, err := createItemPrototypeObject(d, api)
	if err != nil {
		return err
	}

	item.ItemID = d.Id()
	log.Printf("[DEBUG] Update item prototype %#v", item)
	return createRetry(d, meta, updateItemPrototype, *item, resourceZabbixItemPrototypeRead)
}

func resourceZabbixItemPrototypeDelete(d *schema.ResourceData, meta interface{}) error {
	api := meta.(*zabbix.API)

	return deleteRetry(d.Id(), getItemPrototypeParentID, api.ItemPrototypesDeleteIDs, api)
}

func getItemPrototypeParentID(api *zabbix.API, id string) (string, error) {
	items, err := api.ItemPrototypesGet(zabbix.Params{
		"output":      "extend",
		"selectHosts": "extend",
		"itemids":     id,
	})
	if err != nil {
		return "", fmt.Errorf("%s, with item %s", err.Error(), id)
	}
	if len(items) != 1 {
		return "", fmt.Errorf("Expected one item and got %d items", len(items))
	}
	if len(items[0].Hosts) != 1 {
		return "", fmt.Errorf("Expected one parent for item %s and got %d", id, len(items[0].Hosts))
	}
	return items[0].Hosts[0].HostID, nil
}

func createItemPrototype(item interface{}, api *zabbix.API) (id string, err error) {
	items := zabbix.ItemPrototypes{item.(zabbix.ItemPrototype)}

	err = api.ItemPrototypesCreate(items)
	if err != nil {
		return
	}
	id = items[0].ItemID
	return
}

func updateItemPrototype(item interface{}, api *zabbix.API) (id string, err error) {
	items := zabbix.ItemPrototypes{item.(zabbix.ItemPrototype)}

	err = api.ItemPrototypesUpdate(items)
	if err != nil {
		return
	}
	id = items[0].ItemID
	return
}
