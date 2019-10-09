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
					resource.TestCheckResourceAttr(resourceName, "macro.MACRO1", "value1"),
					resource.TestCheckResourceAttr(resourceName, "macro.MACRO2", "value2"),
				),
			},
			{
				Config: testAccZabbixTemplateSimpleUpdate(strID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "description", "update_test_template_description"),
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("update_template_%s", strID)),
					resource.TestCheckResourceAttr(resourceName, "host", fmt.Sprintf("update_template_%s", strID)),
					resource.TestCheckResourceAttr(resourceName, "groups.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "macro.MACRO1", "update_value1"),
					resource.TestCheckResourceAttr(resourceName, "macro.UPDATE_MACRO2", "value2"),
				),
			},
		},
	})
}

func testAccCheckZabbixTemplateDestroy(s *terraform.State) error {
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
		macro = {
		  MACRO1 = "value1"
		  MACRO2 = "value2"
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
		macro = {
		  MACRO1 = "update_value1"
		  UPDATE_MACRO2 = "value2"
		}
	}
	`, strID, strID, strID)
}
