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
					resource.TestCheckResourceAttr(resourceName, "host", fmt.Sprintf("template_%s", strID)),
					resource.TestCheckResourceAttr(resourceName, "groups.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "item.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "item.0.name", "name1"),
					resource.TestCheckResourceAttr(resourceName, "item.0.key", "key1"),
					resource.TestCheckResourceAttr(resourceName, "item.0.delay", "15"),
					resource.TestCheckResourceAttr(resourceName, "item.1.name", "name2"),
					resource.TestCheckResourceAttr(resourceName, "item.1.key", "key2"),
					resource.TestCheckResourceAttr(resourceName, "item.1.delay", "30"),
				),
			},
			{
				Config: testAccZabbixTemplateSimpleUpdate(strID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "description", "update_test_template_description"),
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("update_template_%s", strID)),
					resource.TestCheckResourceAttr(resourceName, "host", fmt.Sprintf("update_template_%s", strID)),
					resource.TestCheckResourceAttr(resourceName, "groups.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "item.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "item.0.name", "nameA"),
					resource.TestCheckResourceAttr(resourceName, "item.0.key", "key1"),
					resource.TestCheckResourceAttr(resourceName, "item.0.delay", "60"),
					resource.TestCheckResourceAttr(resourceName, "item.1.name", "name2"),
					resource.TestCheckResourceAttr(resourceName, "item.1.key", "keyB"),
					resource.TestCheckResourceAttr(resourceName, "item.1.delay", "120"),
					resource.TestCheckResourceAttr(resourceName, "item.2.name", "nameC"),
					resource.TestCheckResourceAttr(resourceName, "item.2.key", "keyC"),
					resource.TestCheckResourceAttr(resourceName, "item.2.delay", "180"),
				),
			},
			{
				Config: testAccZabbixTemplateSimpleConfig(strID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "description", "test_template_description"),
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("template_%s", strID)),
					resource.TestCheckResourceAttr(resourceName, "host", fmt.Sprintf("template_%s", strID)),
					resource.TestCheckResourceAttr(resourceName, "groups.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "item.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "item.0.name", "name1"),
					resource.TestCheckResourceAttr(resourceName, "item.0.key", "key1"),
					resource.TestCheckResourceAttr(resourceName, "item.0.delay", "15"),
					resource.TestCheckResourceAttr(resourceName, "item.1.name", "name2"),
					resource.TestCheckResourceAttr(resourceName, "item.1.key", "key2"),
					resource.TestCheckResourceAttr(resourceName, "item.1.delay", "30"),
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
		host = "template_%s"
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
	}
	`, strID, strID, strID)
}

func testAccZabbixTemplateSimpleUpdate(strID string) string {
	return fmt.Sprintf(`
	resource "zabbix_host_group" "host_group_test" {
		name = "host_group_%s"
	}

	resource "zabbix_template" "template_test" {
		host = "update_template_%s"
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
	}
	`, strID, strID, strID)
}
