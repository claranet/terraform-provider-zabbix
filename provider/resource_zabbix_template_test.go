package provider

import (
	"fmt"
	"testing"

	"github.com/claranet/go-zabbix-api"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccZabbixTemplate_Basic(t *testing.T) {
	resourceName := "zabbix_template.template_test"
	strID := acctest.RandString(5)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckZabbixTemplateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccZabbixTemplateSimpleConfig(strID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "description", "test_template_description"),
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("template_%s", strID)),
					resource.TestCheckResourceAttr(resourceName, "host", "template_test"),
					resource.TestCheckResourceAttr(resourceName, "groups.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "item.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "item.1763311883.name", "name1"),
					resource.TestCheckResourceAttr(resourceName, "item.1763311883.key", "key1"),
					resource.TestCheckResourceAttr(resourceName, "item.1763311883.delay", "15"),
					resource.TestCheckResourceAttr(resourceName, "item.1507990063.name", "name2"),
					resource.TestCheckResourceAttr(resourceName, "item.1507990063.key", "key2"),
					resource.TestCheckResourceAttr(resourceName, "item.1507990063.delay", "30"),
					// resource.TestCheckResourceAttr(resourceName, "trigger.2837083952.description", "trigger_name_A"),
					// resource.TestCheckResourceAttr(resourceName, "trigger.2837083952.expression", "{template_test:key1.last()}=0"),
					// resource.TestCheckResourceAttr(resourceName, "trigger.2837083952.priority", "5"),
					// resource.TestCheckResourceAttr(resourceName, "trigger.2837083952.status", "0"),
					// resource.TestCheckResourceAttr(resourceName, "trigger.2837083952.description", "trigger_name_B"),
					// resource.TestCheckResourceAttr(resourceName, "trigger.2837083952.expression", "{template_test:key1.last()}=1"),
					// resource.TestCheckResourceAttr(resourceName, "trigger.2837083952.priority", "3"),
					// resource.TestCheckResourceAttr(resourceName, "trigger.2837083952.status", "1"),
				),
			},
			{
				Config: testAccZabbixTemplateSimpleUpdate(strID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "description", "update_test_template_description"),
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("update_template_%s", strID)),
					resource.TestCheckResourceAttr(resourceName, "host", "template_test"),
					resource.TestCheckResourceAttr(resourceName, "groups.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "item.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "item.2064979859.name", "nameA"),
					resource.TestCheckResourceAttr(resourceName, "item.2064979859.key", "key1"),
					resource.TestCheckResourceAttr(resourceName, "item.2064979859.delay", "60"),
					resource.TestCheckResourceAttr(resourceName, "item.576049446.name", "name2"),
					resource.TestCheckResourceAttr(resourceName, "item.576049446.key", "keyB"),
					resource.TestCheckResourceAttr(resourceName, "item.576049446.delay", "120"),
					resource.TestCheckResourceAttr(resourceName, "item.617140443.name", "nameC"),
					resource.TestCheckResourceAttr(resourceName, "item.617140443.key", "keyC"),
					resource.TestCheckResourceAttr(resourceName, "item.617140443.delay", "180"),
					// resource.TestCheckResourceAttr(resourceName, "trigger.1101767588.description", "trigger_name_A"),
					// resource.TestCheckResourceAttr(resourceName, "trigger.1101767588.expression", "{template_test:key1.last()}=3"),
					// resource.TestCheckResourceAttr(resourceName, "trigger.1101767588.priority", "4"),
					// resource.TestCheckResourceAttr(resourceName, "trigger.1101767588.status", "1"),
					// resource.TestCheckResourceAttr(resourceName, "trigger.1604588631.description", "trigger_name_2"),
					// resource.TestCheckResourceAttr(resourceName, "trigger.1604588631.expression", "{template_test:key1.last()}=1"),
					// resource.TestCheckResourceAttr(resourceName, "trigger.1604588631.priority", "2"),
					// resource.TestCheckResourceAttr(resourceName, "trigger.1604588631.status", "0"),
					// resource.TestCheckResourceAttr(resourceName, "trigger.2624594860.description", "trigger_name_C"),
					// resource.TestCheckResourceAttr(resourceName, "trigger.2624594860.expression", "{template_test:key1.last()}=6"),
					// resource.TestCheckResourceAttr(resourceName, "trigger.2624594860.priority", "3"),
					// resource.TestCheckResourceAttr(resourceName, "trigger.2624594860.status", "1"),
				),
			},
			{
				Config: testAccZabbixTemplateSimpleConfig(strID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "description", "test_template_description"),
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("template_%s", strID)),
					resource.TestCheckResourceAttr(resourceName, "host", "template_test"),
					resource.TestCheckResourceAttr(resourceName, "groups.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "item.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "item.1763311883.name", "name1"),
					resource.TestCheckResourceAttr(resourceName, "item.1763311883.key", "key1"),
					resource.TestCheckResourceAttr(resourceName, "item.1763311883.delay", "15"),
					resource.TestCheckResourceAttr(resourceName, "item.1507990063.name", "name2"),
					resource.TestCheckResourceAttr(resourceName, "item.1507990063.key", "key2"),
					resource.TestCheckResourceAttr(resourceName, "item.1507990063.delay", "30"),
					// resource.TestCheckResourceAttr(resourceName, "trigger.2837083952.description", "trigger_name_A"),
					// resource.TestCheckResourceAttr(resourceName, "trigger.2837083952.expression", "{template_test:key1.last()}=0"),
					// resource.TestCheckResourceAttr(resourceName, "trigger.2837083952.priority", "5"),
					// resource.TestCheckResourceAttr(resourceName, "trigger.2837083952.status", "0"),
					// resource.TestCheckResourceAttr(resourceName, "trigger.3192477590.description", "trigger_name_B"),
					// resource.TestCheckResourceAttr(resourceName, "trigger.3192477590.expression", "{template_test:key1.last()}=1"),
					// resource.TestCheckResourceAttr(resourceName, "trigger.3192477590.priority", "3"),
					// resource.TestCheckResourceAttr(resourceName, "trigger.3192477590.status", "1"),
				),
			},
		},
	})
}

func testAccCheckZabbixTemplateDestroy(s *terraform.State) error {
	api := testAccProvider.Meta().(*zabbix.API)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "zabbix_template" {
			continue
		}

		_, err := api.TemplateGetByID(rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("Template still exist %s", rs.Primary.ID)
		}

		expectedError := "Expected exactly one result, got 0."
		if err.Error() != expectedError {
			return fmt.Errorf("expected error : %s, got : %s", expectedError, err.Error())
		}
	}
	return nil

}

func testAccZabbixTemplateSimpleConfig(strID string) string {
	return fmt.Sprintf(`
	resource "zabbix_host_group" "host_group_test" {
		name = "host_group_%s"
	}

	resource "zabbix_template" "template_test" {
		host = "template_test"
		groups = ["${zabbix_host_group.host_group_test.name}"]
		name = "template_%s"
		description = "test_template_description"

		item {
			name = "name1"
			key = "key1"
			delay = "15"
		}
		item {
			name = "name2"
			key = "key2"
			delay = "30"
		}
		trigger {
			description = "trigger_name_A"
			expression = "{template_test:key1.last()}=0"
			priority = 5
			status = 0
		}
		trigger {
			description = "trigger_name_B"
			expression = "{template_test:key1.last()}=1"
			priority = 3
			status = 1
		}
	}
	`, strID, strID)
}

func testAccZabbixTemplateSimpleUpdate(strID string) string {
	return fmt.Sprintf(`
	resource "zabbix_host_group" "host_group_test" {
		name = "host_group_%s"
	}

	resource "zabbix_template" "template_test" {
		host = "template_test"
		groups = ["${zabbix_host_group.host_group_test.name}"]
		name = "update_template_%s"
		description = "update_test_template_description"

		item {
			name = "nameA"
			key = "key1"
			delay = "60"
		}
		item {
			name = "name2"
			key = "keyB"
			delay = "120"
		}
		item {
			name = "nameC"
			key = "keyC"
			delay = "180"
		}
		trigger {
			description = "trigger_name_A"
			expression = "{template_test:key1.last()}=3"
			priority = 4
			status = 1
		}
		trigger {
			description = "trigger_name_2"
			expression = "{template_test:key1.last()}=1"
			priority = 2
			status = 0
		}
		trigger {
			description = "trigger_name_C"
			expression = "{template_test:key1.last()}=6"
			priority = 3
			status = 1
		}
	}
	`, strID, strID)
}
