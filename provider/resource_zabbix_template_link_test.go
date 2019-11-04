package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
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
					resource.TestCheckResourceAttr("zabbix_template_link.template_link_test", "trigger.#", "1"),
				),
			},
			{
				Config: testAccZabbixTemplateLinkDeleteTrigger(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("zabbix_template_link.template_link_test", "item.#", "1"),
					resource.TestCheckResourceAttr("zabbix_template_link.template_link_test", "trigger.#", "0"),
				),
			},
			{
				Config: testAccZabbixTemplateLinkDeleteItem(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("zabbix_template_link.template_link_test", "item.#", "0"),
					resource.TestCheckResourceAttr("zabbix_template_link.template_link_test", "trigger.#", "0"),
				),
			},
			{
				Config: testAccZabbixTemplateLinkConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("zabbix_template_link.template_link_test", "item.#", "1"),
					resource.TestCheckResourceAttr("zabbix_template_link.template_link_test", "trigger.#", "1"),
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
			trends = "300"
			history = "25"
			host_id = "${zabbix_template.template_test.template_id}"
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
			trends = "300"
			history = "25"
			host_id = "${zabbix_template.template_test.template_id}"
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

func testAccCheckZabbixTemplateLinkDestroy(s *terraform.State) error {
	return nil
}
