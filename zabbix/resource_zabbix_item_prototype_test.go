package provider

import (
	"fmt"
	"testing"

	"github.com/claranet/go-zabbix-api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccZabbixItemPrototype_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckZabbixItemPrototypeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccZabbixItemPrototypeConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("zabbix_item_prototype.item_prototype_test", "delay", "60"),
					resource.TestCheckResourceAttr("zabbix_item_prototype.item_prototype_test", "interface_id", "0"),
					resource.TestCheckResourceAttr("zabbix_item_prototype.item_prototype_test", "key", "test.key"),
					resource.TestCheckResourceAttr("zabbix_item_prototype.item_prototype_test", "name", "item_prototype_test"),
					resource.TestCheckResourceAttr("zabbix_item_prototype.item_prototype_test", "type", "0"),
					resource.TestCheckResourceAttr("zabbix_item_prototype.item_prototype_test", "status", "0"),
				),
			},
			{
				Config: testAccZabbixItemPrototypeUpdateConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("zabbix_item_prototype.item_prototype_test", "delay", "90"),
					resource.TestCheckResourceAttr("zabbix_item_prototype.item_prototype_test", "interface_id", "0"),
					resource.TestCheckResourceAttr("zabbix_item_prototype.item_prototype_test", "key", "test.key.update"),
					resource.TestCheckResourceAttr("zabbix_item_prototype.item_prototype_test", "name", "item_prototype_test_update"),
					resource.TestCheckResourceAttr("zabbix_item_prototype.item_prototype_test", "type", "2"),
					resource.TestCheckResourceAttr("zabbix_item_prototype.item_prototype_test", "status", "1"),
				),
			},
		},
	})
}

func testAccCheckZabbixItemPrototypeDestroy(s *terraform.State) error {
	api := testAccProvider.Meta().(*zabbix.API)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "zabbix_item_prototype" {
			continue
		}

		_, err := api.ItemPrototypeGetByID(rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("Item prototype still exist %s", rs.Primary.ID)
		}

		expectedError := "Expected exactly one result, got 0."
		if err.Error() != expectedError {
			return fmt.Errorf("expected error : %s, got : %s", expectedError, err.Error())
		}
	}
	return nil
}

func testAccZabbixItemPrototypeConfig() string {
	return fmt.Sprintf(`
		resource "zabbix_host_group" "zabbix" {
			name = "host group test"
		}

		resource "zabbix_template" "template_test" {
			host = "template_test"
			groups = ["${zabbix_host_group.zabbix.name}"]
			name = "display name for template test"
	  	}

		resource "zabbix_lld_rule" "lld_rule_test" {
			delay = 60
			host_id = zabbix_template.template_test.id
			interface_id = "0"
			key = "key.lolo"
			name = "test_low_level_discovery_rule"
			type = 0
			filter {
				condition {
					macro = "{#TESTMACRO}"
					value = "^lo$"
				}
				eval_type = 0
			}
		}

		resource "zabbix_item_prototype" "item_prototype_test" {
			delay = 60
			host_id  = zabbix_template.template_test.id
			rule_id = zabbix_lld_rule.lld_rule_test.id
			interface_id = "0"
			key = "test.key"
			name = "item_prototype_test"
			type = 0
			status = 0
		}
	`)
}

func testAccZabbixItemPrototypeUpdateConfig() string {
	return fmt.Sprintf(`
		resource "zabbix_host_group" "zabbix" {
			name = "host group test"
		}

		resource "zabbix_template" "template_test" {
			host = "template_test"
			groups = ["${zabbix_host_group.zabbix.name}"]
			name = "display name for template test"
	  	}

		resource "zabbix_lld_rule" "lld_rule_test" {
			delay = 60
			host_id = zabbix_template.template_test.id
			interface_id = "0"
			key = "key.lolo"
			name = "test_low_level_discovery_rule"
			type = 0
			filter {
				condition {
					macro = "{#TESTMACRO}"
					value = "^lo$"
				}
				eval_type = 0
			}
		}

		resource "zabbix_item_prototype" "item_prototype_test" {
			delay = 90
			host_id  = zabbix_template.template_test.id
			rule_id = zabbix_lld_rule.lld_rule_test.id
			interface_id = "0"
			key = "test.key.update"
			name = "item_prototype_test_update"
			type = 2
			status = 1
		}
	`)
}
