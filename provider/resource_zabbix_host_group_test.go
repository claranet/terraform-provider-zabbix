package provider

import (
	"fmt"
	"testing"

	"github.com/claranet/go-zabbix-api"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccZabbixHostGroup_Basic(t *testing.T) {
	groupName := fmt.Sprintf("host_groud_%s", acctest.RandString(5))
	var hostGroup zabbix.HostGroup
	expectedHostGroup := zabbix.HostGroup{Name: groupName}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckZabbixHostGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccZabbixHostGroupConfig(groupName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckZabbixHostGroupExist("zabbix_host_group.zabbix", &hostGroup),
					testAccCheckZabbixHostGroupAttributes(&hostGroup, expectedHostGroup),
					resource.TestCheckResourceAttr("zabbix_host_group.zabbix", "name", groupName),
				),
			},
		},
	})
}

func testAccCheckZabbixHostGroupDestroy(s *terraform.State) error {
	api := testAccProvider.Meta().(*zabbix.API)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "zabbix_host_group" {
			continue
		}

		_, err := api.HostGroupGetByID(rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("Host group still exists")
		}
		expectedError := "Expected exactly one result, got 0."
		if err.Error() != expectedError {
			return fmt.Errorf("expected error : %s, got : %s", expectedError, err.Error())
		}
	}
	return nil
}

func testAccZabbixHostGroupConfig(groupName string) string {
	return fmt.Sprintf(`
		resource "zabbix_host_group" "zabbix" {
			name = "%s"
		}`, groupName,
	)
}

func testAccCheckZabbixHostGroupExist(resource string, hostGroup *zabbix.HostGroup) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("Not found; %s", resource)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No record ID set")
		}

		api := testAccProvider.Meta().(*zabbix.API)
		group, err := api.HostGroupGetByID(rs.Primary.ID)
		if err != nil {
			return err
		}
		*hostGroup = *group
		return nil
	}
}

func testAccCheckZabbixHostGroupAttributes(hostGroup *zabbix.HostGroup, want zabbix.HostGroup) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		if hostGroup.Name != want.Name {
			return fmt.Errorf("got host name : %q, expected : %q", hostGroup.Name, want.Name)
		}
		return nil
	}
}
