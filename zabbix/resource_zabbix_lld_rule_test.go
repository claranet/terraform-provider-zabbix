package provider

import (
	"fmt"
	"testing"

	"github.com/claranet/go-zabbix-api"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccZabbixlldRule_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckZabbixlldRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccZabbixlldRuleConfig(),
				Check: resource.ComposeTestCheckFunc(
					// printastat(),
					resource.TestCheckResourceAttr("zabbix_lld_rule.lld_rule_test", "delay", "60"),
					resource.TestCheckResourceAttr("zabbix_lld_rule.lld_rule_test", "interface_id", "0"),
					resource.TestCheckResourceAttr("zabbix_lld_rule.lld_rule_test", "key", "key.lolo"),
					resource.TestCheckResourceAttr("zabbix_lld_rule.lld_rule_test", "name", "test_low_level_discovery_rule"),
					resource.TestCheckResourceAttr("zabbix_lld_rule.lld_rule_test", "type", "0"),
					resource.TestCheckResourceAttr("zabbix_lld_rule.lld_rule_test", "filter.#", "1"),
					resource.TestCheckResourceAttr("zabbix_lld_rule.lld_rule_test", "filter.3189296381.#", "0"),
					resource.TestCheckResourceAttr("zabbix_lld_rule.lld_rule_test", "filter.3189296381.condition.23998414.macro", "{#TESTMACRO}"),
					resource.TestCheckResourceAttr("zabbix_lld_rule.lld_rule_test", "filter.3189296381.condition.23998414.value", "^lo$"),
					resource.TestCheckResourceAttr("zabbix_lld_rule.lld_rule_test", "filter.3189296381.condition.23998414.operator", "8"),
				),
			},
		},
	})
}

func testAccCheckZabbixlldRuleDestroy(s *terraform.State) error {
	api := testAccProvider.Meta().(*zabbix.API)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "zabbix_item" {
			continue
		}

		_, err := api.ItemGetByID(rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("Item still exist %s", rs.Primary.ID)
		}

		expectedError := "Expected exactly one result, got 0."
		if err.Error() != expectedError {
			return fmt.Errorf("expected error : %s, got : %s", expectedError, err.Error())
		}
	}
	return nil
}

func testAccZabbixlldRuleConfig() string {
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
	`)
}

func printastat() resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources["zabbix_lld_rule.lld_rule_test"]

		if !ok {
			return fmt.Errorf("pas ok")
		}
		return fmt.Errorf("resourde %#v", rs.Primary)
	}
}
