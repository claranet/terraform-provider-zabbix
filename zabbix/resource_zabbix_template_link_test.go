package provider

import (
	"fmt"
	"log"
	"testing"

	"github.com/claranet/go-zabbix-api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccZabbixTemplateLink_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckZabbixTemplateLinkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccZabbixTemplateLinkConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("zabbix_template_link.template_link_test", "item.#", "1"),
					resource.TestCheckResourceAttr("zabbix_template_link.template_link_test", "server_item.#", "0"),
					resource.TestCheckResourceAttr("zabbix_template_link.template_link_test", "trigger.#", "1"),
					resource.TestCheckResourceAttr("zabbix_template_link.template_link_test", "server_trigger.#", "0"),
				),
			},
			{
				Config: testAccZabbixTemplateLinkDeleteTrigger(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("zabbix_template_link.template_link_test", "item.#", "1"),
					resource.TestCheckResourceAttr("zabbix_template_link.template_link_test", "server_item.#", "0"),
					resource.TestCheckResourceAttr("zabbix_template_link.template_link_test", "trigger.#", "0"),
					resource.TestCheckResourceAttr("zabbix_template_link.template_link_test", "server_trigger.#", "0"),
				),
			},
			{
				Config: testAccZabbixTemplateLinkDeleteItem(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("zabbix_template_link.template_link_test", "item.#", "0"),
					resource.TestCheckResourceAttr("zabbix_template_link.template_link_test", "server_item.#", "0"),
					resource.TestCheckResourceAttr("zabbix_template_link.template_link_test", "trigger.#", "0"),
					resource.TestCheckResourceAttr("zabbix_template_link.template_link_test", "server_trigger.#", "0"),
				),
			},
			{
				Config: testAccZabbixTemplateLinkConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("zabbix_template_link.template_link_test", "item.#", "1"),
					resource.TestCheckResourceAttr("zabbix_template_link.template_link_test", "server_item.#", "0"),
					resource.TestCheckResourceAttr("zabbix_template_link.template_link_test", "trigger.#", "1"),
					resource.TestCheckResourceAttr("zabbix_template_link.template_link_test", "server_trigger.#", "0"),
				),
			},
		},
	})
}

func TestAccZabbixTemplateLink_DeleteServerItem(t *testing.T) {
	var template zabbix.Template
	item := zabbix.Item{
		Name:  "server_item",
		Key:   "server.key",
		Type:  zabbix.ZabbixAgent,
		Delay: "30",
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckZabbixTemplateLinkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccZabbixTemplateLinkConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTemplateExist("zabbix_template.template_test", &template),
					resource.TestCheckResourceAttr("zabbix_template_link.template_link_test", "item.#", "1"),
					resource.TestCheckResourceAttr("zabbix_template_link.template_link_test", "server_item.#", "0"),
					resource.TestCheckResourceAttr("zabbix_template_link.template_link_test", "trigger.#", "1"),
					resource.TestCheckResourceAttr("zabbix_template_link.template_link_test", "server_trigger.#", "0"),
				),
			},
			{
				PreConfig: testAccZabbixTemplateLinkCreateServerItem(template, &item),
				Config:    testAccZabbixTemplateLinkConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTemplateServerItemDelete(&item),
					resource.TestCheckResourceAttr("zabbix_template_link.template_link_test", "item.#", "1"),
					resource.TestCheckResourceAttr("zabbix_template_link.template_link_test", "server_item.#", "0"),
					resource.TestCheckResourceAttr("zabbix_template_link.template_link_test", "trigger.#", "1"),
					resource.TestCheckResourceAttr("zabbix_template_link.template_link_test", "server_trigger.#", "0"),
				),
			},
		},
	})
}

func TestAccZabbixTemplateLink_DeleteServerTrigger(t *testing.T) {
	var template zabbix.Template
	item := zabbix.Item{
		Name:  "server_item",
		Key:   "server.key",
		Type:  zabbix.ZabbixAgent,
		Delay: "30",
	}
	trigger := zabbix.Trigger{
		Description: "server_trigger",
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckZabbixTemplateLinkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccZabbixTemplateLinkConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTemplateExist("zabbix_template.template_test", &template),
					resource.TestCheckResourceAttr("zabbix_template_link.template_link_test", "item.#", "1"),
					resource.TestCheckResourceAttr("zabbix_template_link.template_link_test", "server_item.#", "0"),
					resource.TestCheckResourceAttr("zabbix_template_link.template_link_test", "trigger.#", "1"),
					resource.TestCheckResourceAttr("zabbix_template_link.template_link_test", "server_trigger.#", "0"),
				),
			},
			{
				PreConfig: testAccZabbixTemplateLinkCreateServerTrigger(template, item, &trigger),
				Config:    testAccZabbixTemplateLinkConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTemplateServerTriggerDelete(&trigger),
					resource.TestCheckResourceAttr("zabbix_template_link.template_link_test", "item.#", "1"),
					resource.TestCheckResourceAttr("zabbix_template_link.template_link_test", "server_item.#", "0"),
					resource.TestCheckResourceAttr("zabbix_template_link.template_link_test", "trigger.#", "1"),
					resource.TestCheckResourceAttr("zabbix_template_link.template_link_test", "server_trigger.#", "0"),
				),
			},
		},
	})
}

func testAccZabbixTemplateLinkConfig() string {
	return fmt.Sprintf(`
		resource "zabbix_host_group" "zabbix" {
			name = "host group test"
		}

		resource "zabbix_template" "template_test" {
			host = "template_test"
			groups = ["${zabbix_host_group.zabbix.name}"]
			name = "display name for template test"
	  	}
	  
		resource "zabbix_item" "item_test_0" {
			name = "item_test_0"
			key = "bilou.bilou"
			delay = "34"
			host_id = "${zabbix_template.template_test.id}"
		}
		
		resource "zabbix_trigger" "trigger_test_0" {
			description = "trigger_test_0"
			expression  = "{${zabbix_template.template_test.host}:${zabbix_item.item_test_0.key}.last()} = 0"
			priority    = 5
		}

		resource "zabbix_template_link" "template_link_test" {
			template_id = zabbix_template.template_test.id
			item {
				item_id = zabbix_item.item_test_0.id
			}
			trigger {
				trigger_id = zabbix_trigger.trigger_test_0.id
			}
		}
	`)
}

func testAccZabbixTemplateLinkDeleteTrigger() string {
	return fmt.Sprintf(`
		resource "zabbix_host_group" "zabbix" {
			name = "host group test"
		}

		resource "zabbix_template" "template_test" {
			host = "template_test"
			groups = ["${zabbix_host_group.zabbix.name}"]
			name = "display name for template test"
	  	}
	  
		resource "zabbix_item" "item_test_0" {
			name = "item_test_0"
			key = "bilou.bilou"
			delay = "34"
			host_id = "${zabbix_template.template_test.id}"
		}

		resource "zabbix_template_link" "template_link_test" {
			template_id = zabbix_template.template_test.id
			item {
				item_id = zabbix_item.item_test_0.id
			}
		}
	`)
}

func testAccZabbixTemplateLinkDeleteItem() string {
	return fmt.Sprintf(`
		resource "zabbix_host_group" "zabbix" {
			name = "host group test"
		}

		resource "zabbix_template" "template_test" {
			host = "template_test"
			groups = ["${zabbix_host_group.zabbix.name}"]
			name = "display name for template test"
		  }
		  
		resource "zabbix_template_link" "template_link_test" {
			template_id = zabbix_template.template_test.id
		}
	`)
}

func testAccZabbixTemplateLinkCreateServerItem(template zabbix.Template, item *zabbix.Item) func() {
	return func() {
		api := testAccProvider.Meta().(*zabbix.API)

		item.HostID = template.TemplateID
		items := zabbix.Items{*item}
		err := api.ItemsCreate(items)
		if err != nil {
			return
		}
		item.ItemID = items[0].ItemID
	}
}

func testAccZabbixTemplateLinkCreateServerTrigger(template zabbix.Template, item zabbix.Item, trigger *zabbix.Trigger) func() {
	return func() {
		api := testAccProvider.Meta().(*zabbix.API)

		trigger.Expression = fmt.Sprintf("{%s:%s.last()} = 0", template.Host, item.Key)
		triggers := zabbix.Triggers{*trigger}
		err := api.TriggersCreate(triggers)
		if err != nil {
			log.Print(err)
			return
		}
		trigger.TriggerID = triggers[0].TriggerID
	}
}

func testAccCheckZabbixTemplateLinkDestroy(s *terraform.State) error {
	return nil
}

func testAccCheckTemplateExist(n string, template *zabbix.Template) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		api := testAccProvider.Meta().(*zabbix.API)
		templates, err := api.TemplateGetByID(rs.Primary.ID)
		if err != nil {
			return err
		}
		template = templates
		return nil
	}
}

func testAccCheckTemplateServerItemDelete(item *zabbix.Item) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		api := testAccProvider.Meta().(*zabbix.API)

		_, err := api.ItemGetByID(item.ItemID)
		if err == nil {
			return fmt.Errorf("Expected an error")
		}

		expectedErr := "Expected exactly one result, got 0."
		if err.Error() != expectedErr {
			return fmt.Errorf("expected error : %s, got : %s", expectedErr, err.Error())
		}
		return nil
	}
}

func testAccCheckTemplateServerTriggerDelete(trigger *zabbix.Trigger) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		api := testAccProvider.Meta().(*zabbix.API)

		_, err := api.TriggerGetByID(trigger.TriggerID)
		if err == nil {
			return fmt.Errorf("Expected an error")
		}

		expectedErr := "Expected exactly one result, got 0."
		if err.Error() != expectedErr {
			return fmt.Errorf("expected error : %s, got : %s", expectedErr, err.Error())
		}
		return nil
	}
}
