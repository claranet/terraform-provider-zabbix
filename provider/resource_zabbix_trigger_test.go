package provider

import (
	"fmt"
	"testing"

	"github.com/claranet/go-zabbix-api"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccZabbixTrigger_Basic(t *testing.T) {
	resourceName := "zabbix_trigger.trigger_test"
	strID := acctest.RandString(5)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckZabbixTriggerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccZabbixTriggerSimpleConfig(strID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "description", fmt.Sprintf("trigger_%s", strID)),
					resource.TestCheckResourceAttr(resourceName, "expression", fmt.Sprintf("{template_%s:lili.lala.last()}=0", strID)),
					resource.TestCheckResourceAttr(resourceName, "comment", "trigger_comment"),
					resource.TestCheckResourceAttr(resourceName, "priority", "5"),
					resource.TestCheckResourceAttr(resourceName, "status", "1"),
				),
			},
			{
				Config: testAccZabbixTriggerSimpleConfigUpdate(strID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "description", fmt.Sprintf("update_trigger_%s", strID)),
					resource.TestCheckResourceAttr(resourceName, "expression", fmt.Sprintf("{template_%s:lili.lala.min(1)}=0", strID)),
					resource.TestCheckResourceAttr(resourceName, "comment", "update_trigger_comment"),
					resource.TestCheckResourceAttr(resourceName, "priority", "0"),
					resource.TestCheckResourceAttr(resourceName, "status", "0"),
				),
			},
			{
				Config: testAccZabbixTriggerOmitEmpty(strID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "description", fmt.Sprintf("update_trigger_%s", strID)),
					resource.TestCheckResourceAttr(resourceName, "expression", fmt.Sprintf("{template_%s:lili.lala.min(1)}=0", strID)),
					resource.TestCheckResourceAttr(resourceName, "comment", ""),
					resource.TestCheckResourceAttr(resourceName, "priority", "0"),
					resource.TestCheckResourceAttr(resourceName, "status", "0"),
					resource.TestCheckResourceAttr(resourceName, "dependencies.#", "0"),
				),
			},
		},
	})
}

func TestAccZabbixTrigger_BasicMacro(t *testing.T) {
	resourceName := "zabbix_trigger.trigger_test"
	strID := acctest.RandString(5)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckZabbixTriggerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccZabbixTriggerMacroConfig(strID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "description", fmt.Sprintf("trigger_%s", strID)),
					resource.TestCheckResourceAttr(resourceName, "expression", fmt.Sprintf("{template_%s:lili.lala.min({$MACRO_TRIGGER})}=0", strID)),
					resource.TestCheckResourceAttr(resourceName, "comment", "trigger_comment"),
					resource.TestCheckResourceAttr(resourceName, "priority", "5"),
					resource.TestCheckResourceAttr(resourceName, "status", "1"),
				),
			},
			{
				Config: testAccZabbixTriggerMacroConfigUpdate(strID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "description", fmt.Sprintf("update_trigger_%s", strID)),
					resource.TestCheckResourceAttr(resourceName, "expression", fmt.Sprintf("{template_%s:lili.lala.min({$MACRO_UPDATE})}=0", strID)),
					resource.TestCheckResourceAttr(resourceName, "comment", "update_trigger_comment"),
					resource.TestCheckResourceAttr(resourceName, "priority", "3"),
					resource.TestCheckResourceAttr(resourceName, "status", "0"),
				),
			},
		},
	})
}

func TestAccZabbixTrigger_BasicDependencies(t *testing.T) {
	resourceName := "zabbix_trigger.trigger_test_3"
	strID := acctest.RandString(5)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckZabbixTriggerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccZabbixTriggerDependencies(strID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "description", fmt.Sprintf("trigger_3_%s", strID)),
					resource.TestCheckResourceAttr(resourceName, "expression", fmt.Sprintf("{template_%s:lili.lala.last()}=0", strID)),
					resource.TestCheckResourceAttr(resourceName, "dependencies.#", "2"),
				),
			},
		},
	})
}

func testAccCheckZabbixTriggerDestroy(s *terraform.State) error {
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

func testAccZabbixTriggerSimpleConfig(strID string) string {
	return fmt.Sprintf(`		
	resource "zabbix_host_group" "host_group_test" {
		name = "host_group_%s"
	}

	resource "zabbix_template" "template_test" {
		host = "template_%s"
		groups = ["${zabbix_host_group.host_group_test.name}"]
		description = "description for template"
	  }
  
	resource "zabbix_item" "item_test" {
		name = "name_%s"
		key = "lili.lala"
		delay = "34"
		description = "description for item"
		trends = "300"
		history = "25"
		delta = 1
		type = 2
		host_id = "${zabbix_template.template_test.template_id}"
	}
	
	resource "zabbix_trigger" "trigger_test" {
		description = "trigger_%s"
		expression = "{${zabbix_template.template_test.host}:${zabbix_item.item_test.key}.last()}=0"
		comment = "trigger_comment"
		priority = 5
		status = 1
	}`, strID, strID, strID, strID)
}

func testAccZabbixTriggerSimpleConfigUpdate(strID string) string {
	return fmt.Sprintf(`		
	resource "zabbix_host_group" "host_group_test" {
		name = "host_group_%s"
	}

	resource "zabbix_template" "template_test" {
		host = "template_%s"
		groups = ["${zabbix_host_group.host_group_test.name}"]
		description = "description for template"
	  }
  
	resource "zabbix_item" "item_test" {
		name = "name_%s"
		key = "lili.lala"
		delay = "34"
		description = "description for item"
		trends = "300"
		history = "25"
		delta = 1
		type = 2
		host_id = "${zabbix_template.template_test.template_id}"
	}
	
	resource "zabbix_trigger" "trigger_test" {
		description = "update_trigger_%s"
		expression = "{${zabbix_template.template_test.host}:${zabbix_item.item_test.key}.min(1)}=0"
		comment = "update_trigger_comment"
		priority = 0
		status = 0
	}`, strID, strID, strID, strID)
}

func testAccZabbixTriggerOmitEmpty(strID string) string {
	return fmt.Sprintf(`		
	resource "zabbix_host_group" "host_group_test" {
		name = "host_group_%s"
	}

	resource "zabbix_template" "template_test" {
		host = "template_%s"
		groups = ["${zabbix_host_group.host_group_test.name}"]
		description = "description for template"
	  }
  
	resource "zabbix_item" "item_test" {
		name = "name_%s"
		key = "lili.lala"
		delay = "34"
		description = "description for item"
		trends = "300"
		history = "25"
		delta = 1
		type = 2
		host_id = "${zabbix_template.template_test.template_id}"
	}
	
	resource "zabbix_trigger" "trigger_test" {
		description = "update_trigger_%s"
		expression = "{${zabbix_template.template_test.host}:${zabbix_item.item_test.key}.min(1)}=0"
	}`, strID, strID, strID, strID)
}

func testAccZabbixTriggerMacroConfig(strID string) string {
	return fmt.Sprintf(`		
	resource "zabbix_host_group" "host_group_test" {
		name = "host_group_%s"
	}

	resource "zabbix_template" "template_test" {
		host = "template_%s"
		groups = ["${zabbix_host_group.host_group_test.name}"]
		description = "description for template"
		macro = {
			MACRO_TRIGGER = "12m",
			MACRO_UPDATE = "21m",
		}
	  }
  
	resource "zabbix_item" "item_test" {
		name = "name_%s"
		key = "lili.lala"
		delay = "34"
		description = "description for item"
		trends = "300"
		history = "25"
		delta = 1
		type = 2
		host_id = "${zabbix_template.template_test.template_id}"
	}
	
	resource "zabbix_trigger" "trigger_test" {
		description = "trigger_%s"
		expression = "{${zabbix_template.template_test.host}:${zabbix_item.item_test.key}.min({$MACRO_TRIGGER})}=0"
		comment = "trigger_comment"
		priority = 5
		status = 1
	}`, strID, strID, strID, strID)
}

func testAccZabbixTriggerMacroConfigUpdate(strID string) string {
	return fmt.Sprintf(`		
	resource "zabbix_host_group" "host_group_test" {
		name = "host_group_%s"
	}

	resource "zabbix_template" "template_test" {
		host = "template_%s"
		groups = ["${zabbix_host_group.host_group_test.name}"]
		description = "description for template"
		macro = {
			MACRO_TRIGGER = "12m",
			MACRO_UPDATE = "21m",
		}
	  }
  
	resource "zabbix_item" "item_test" {
		name = "name_%s"
		key = "lili.lala"
		delay = "34"
		description = "description for item"
		trends = "300"
		history = "25"
		delta = 1
		type = 2
		host_id = "${zabbix_template.template_test.template_id}"
	}
	
	resource "zabbix_trigger" "trigger_test" {
		description = "update_trigger_%s"
		expression = "{${zabbix_template.template_test.host}:${zabbix_item.item_test.key}.min({$MACRO_UPDATE})}=0"
		comment = "update_trigger_comment"
		priority = 3
		status = 0
	}`, strID, strID, strID, strID)
}

func testAccZabbixTriggerDependencies(strID string) string {
	return fmt.Sprintf(`		
	resource "zabbix_host_group" "host_group_test" {
		name = "host_group_%s"
	}

	resource "zabbix_template" "template_test" {
		host = "template_%s"
		groups = ["${zabbix_host_group.host_group_test.name}"]
		description = "description for template"
	  }
  
	resource "zabbix_item" "item_test" {
		name = "name_%s"
		key = "lili.lala"
		delay = "34"
		description = "description for item"
		trends = "300"
		history = "25"
		delta = 1
		type = 2
		host_id = "${zabbix_template.template_test.template_id}"
	}
	
	resource "zabbix_trigger" "trigger_test" {
		description = "trigger_%s"
		expression = "{${zabbix_template.template_test.host}:${zabbix_item.item_test.key}.last()}=0"
	}

	resource "zabbix_trigger" "trigger_test_2" {
		description = "trigger_2_%s"
		expression = "{${zabbix_template.template_test.host}:${zabbix_item.item_test.key}.last()}=0"
	}
	
	resource "zabbix_trigger" "trigger_test_3" {
		description = "trigger_3_%s"
		expression = "{${zabbix_template.template_test.host}:${zabbix_item.item_test.key}.last()}=0"
		dependencies = [
			zabbix_trigger.trigger_test.id,
			zabbix_trigger.trigger_test_2.id,
		]
	}`, strID, strID, strID, strID, strID, strID)
}
