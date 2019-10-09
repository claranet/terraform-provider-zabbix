package provider

import (
	"fmt"
	"testing"

	"github.com/claranet/go-zabbix-api"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccZabbixItem_Basic(t *testing.T) {
	strID := acctest.RandString(5)
	groupName := fmt.Sprintf("host_group_%s", strID)
	templateName := fmt.Sprintf("template_%s", strID)
	itemName := fmt.Sprintf("item_%s", strID)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckZabbixItemDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccZabbixItemConfig(groupName, templateName, itemName),
				Check: resource.ComposeTestCheckFunc(
					testAccZabbixItemExist("zabbix_item.my_item1"),
					resource.TestCheckResourceAttr("zabbix_item.my_item1", "name", itemName),
					resource.TestCheckResourceAttr("zabbix_item.my_item1", "key", "bilou.bilou"),
					resource.TestCheckResourceAttr("zabbix_item.my_item1", "delay", "34"),
					resource.TestCheckResourceAttr("zabbix_item.my_item1", "description", fmt.Sprintf("description for item : %s", itemName)),
					resource.TestCheckResourceAttr("zabbix_item.my_item1", "trends", "300"),
					resource.TestCheckResourceAttr("zabbix_item.my_item1", "history", "25"),
				),
			},
			{
				Config: testAccZabbixItemUpdate(groupName, templateName, itemName),
				Check: resource.ComposeTestCheckFunc(
					testAccZabbixItemExist("zabbix_item.my_item1"),
					resource.TestCheckResourceAttr("zabbix_item.my_item1", "name", fmt.Sprintf("update_%s", itemName)),
					resource.TestCheckResourceAttr("zabbix_item.my_item1", "key", "update.bilou.bilou"),
					resource.TestCheckResourceAttr("zabbix_item.my_item1", "delay", "23"),
					resource.TestCheckResourceAttr("zabbix_item.my_item1", "description", fmt.Sprintf("update description for item : %s", itemName)),
					resource.TestCheckResourceAttr("zabbix_item.my_item1", "trends", "3"),
					resource.TestCheckResourceAttr("zabbix_item.my_item1", "history", "2"),
				),
			},
		},
	})
}

func testAccZabbixItemConfig(groupName string, templateName string, itemName string) string {
	return fmt.Sprintf(`
		resource "zabbix_host_group" "zabbix" {
			name = "%s"
		}

		resource "zabbix_template" "my_zbx_template" {
			host = "%s"
			groups = ["${zabbix_host_group.zabbix.name}"]
			name = "display name %s"
			description = "description for template %s"
	  	}
	  
		resource "zabbix_item" "my_item1" {
			name = "%s"
			key = "bilou.bilou"
			delay = "34"
			description = "description for item : %s"
			trends = "300"
			history = "25"
			host_id = "${zabbix_template.my_zbx_template.template_id}"
	  	}
	`, groupName, templateName, templateName, templateName, itemName, itemName)
}

func testAccZabbixItemUpdate(groupName string, templateName string, itemName string) string {
	return fmt.Sprintf(`
		resource "zabbix_host_group" "zabbix" {
			name = "%s"
		}

		resource "zabbix_template" "my_zbx_template" {
			host = "%s"
			groups = ["${zabbix_host_group.zabbix.name}"]
			name = "display name %s"
			description = "description for template %s"
	  	}
	  
		resource "zabbix_item" "my_item1" {
			name = "update_%s"
			key = "update.bilou.bilou"
			delay = "23"
			description = "update description for item : %s"
			trends = "3"
			history = "2"
			host_id = "${zabbix_template.my_zbx_template.template_id}"
	  	}
	`, groupName, templateName, templateName, templateName, itemName, itemName)
}

func testAccZabbixItemExist(resource string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("Not found : %s", resource)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No record ID set")
		}
		return nil
	}
}

func testAccCheckZabbixItemDestroy(s *terraform.State) error {
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
