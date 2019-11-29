package provider

import (
	"fmt"
	"log"
	"strings"

	"github.com/claranet/go-zabbix-api"
	"github.com/hashicorp/terraform/helper/schema"
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
			"item_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "(readonly) ID of the item prototype.",
			},
			"delay": &schema.Schema{
				Type:     schema.TypeInt,
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
			},
			"delta": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Value that will be stored. ",
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
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Number of days to keep item's history data. Default: 90. ",
				Default:     "90",
			},
			"trends": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Number of days to keep item's trends data. Default: 365. ",
				Default:     "365",
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
		ItemID:       d.Get("item_id").(string),
		Delay:        d.Get("delay").(int),
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
		History:      d.Get("history").(int),
		Trends:       d.Get("trends").(int),
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
	items := zabbix.ItemPrototypes{*item}

	err = api.ItemPrototypesCreate(items)
	if err != nil {
		return err
	}

	d.Set("item_id", items[0].ItemID)
	d.SetId(items[0].ItemID)
	return resourceZabbixItemPrototypeRead(d, meta)
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

	d.Set("item_id", item.ItemID)
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
		log.Printf("Item exist error : %s", err.Error())
		if strings.Contains(err.Error(), "Expected exactly one result") {
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
	items := zabbix.ItemPrototypes{*item}

	log.Printf("[DEBUG] Update item prototype %#v", items)
	err = api.ItemPrototypesUpdate(items)
	if err != nil {
		return err
	}

	return resourceZabbixItemPrototypeRead(d, meta)
}

func resourceZabbixItemPrototypeDelete(d *schema.ResourceData, meta interface{}) error {
	api := meta.(*zabbix.API)

	// items, err := api.ItemPrototypesGet(zabbix.Params{
	// 	"output":      "extend",
	// 	"selectHosts": "extend",
	// 	"itemids":     d.Id(),
	// })
	// if err != nil {
	// 	return fmt.Errorf("%s, with item %s", err.Error(), d.Id())
	// }
	// if len(items) != 1 {
	// 	return fmt.Errorf("Expected one item and got %d items", len(items))
	// }
	// item := items[0]

	// templates, err := api.TemplatesGet(zabbix.Params{
	// 	"ouput":             "extend",
	// 	"parentTemplateids": item.ItemParent[0].HostID,
	// })

	// itemids, err := api.ItemPrototypesDeleteIDs([]string{d.Id()})
	// if err != nil {
	// 	return err
	// }
	// if len(itemids) != len(templates)+1 {
	// 	return fmt.Errorf("Expected to delete %d item and %d were delete", len(templates)+1, len(itemids))
	// }
	_, err := api.ItemPrototypesDeleteIDs([]string{d.Id()})
	return err
}
