package provider

import (
	"fmt"
	"testing"

	"github.com/claranet/go-zabbix-api"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccZabbixHost_Basic(t *testing.T) {
	hostName := fmt.Sprintf("host_name_%s", acctest.RandString(5))
	hostGroup := fmt.Sprintf("host_group_%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckZabbixHostDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccZabbixHostConfig(hostName, hostGroup),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckZabbixHostExit("zabbix_host.zabbix1"),
					resource.TestCheckResourceAttr("zabbix_host.zabbix1", "host", hostName),
				),
			},
		},
	})
}

func testAccCheckZabbixHostDestroy(s *terraform.State) error {
	api := testAccProvider.Meta().(*zabbix.API)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "zabbix_host" {
			continue
		}

		_, err := api.HostGroupGetByID(rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("Host still exists")
		}
		expectedError := "Expected exactly one result, got 0."
		if err.Error() != expectedError {
			return fmt.Errorf("expected error : %s, got : %s", expectedError, err.Error())
		}
	}
	return nil
}

func testAccZabbixHostConfig(hostName string, hostGroup string) string {
	return fmt.Sprintf(`
	  	resource "zabbix_host" "zabbix1" {
			host = "%s"
			interfaces {
		  		ip = "127.0.0.1"
				main = true
			}
			groups = ["Linux servers", "${zabbix_host_group.zabbix.name}"]
			templates = ["Template ICMP Ping"]
	  	}
	  
	  	resource "zabbix_host_group" "zabbix" {
			name = "%s"
	  	}`, hostName, hostGroup,
	)
}

func testAccCheckZabbixHostExit(resource string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("Not found: %s", resource)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No record ID id set")
		}

		api := testAccProvider.Meta().(*zabbix.API)
		_, err := api.HostGetByID(rs.Primary.ID)
		if err != nil {
			return err
		}
		return nil
	}
}
