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
	var getHost zabbix.Host
	randName := acctest.RandString(5)
	host := fmt.Sprintf("host_%s", randName)
	name := fmt.Sprintf("name_%s", randName)
	hostGroup := fmt.Sprintf("host_group_%s", randName)
	expectedHost := zabbix.Host{
		Host:       host,
		Name:       name,
		Interfaces: zabbix.HostInterfaces{zabbix.HostInterface{IP: "127.0.0.1", Main: 1}},
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckZabbixHostDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccZabbixHostConfig(host, name, hostGroup),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckZabbixHostExist("zabbix_host.zabbix1", &getHost),
					testAccCheckZabbixHostAttributes(&getHost, expectedHost, []string{"Linux servers", hostGroup}, []string{"Template ICMP Ping"}),
					resource.TestCheckResourceAttr("zabbix_host.zabbix1", "host", host),
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

func testAccZabbixHostConfig(host string, name string, hostGroup string) string {
	return fmt.Sprintf(`
	  	resource "zabbix_host" "zabbix1" {
			host = "%s"
			name = "%s"
			interfaces {
		  		ip = "127.0.0.1"
				main = true
			}
			groups = ["Linux servers", "${zabbix_host_group.zabbix.name}"]
			templates = ["Template ICMP Ping"]
	  	}
	  
	  	resource "zabbix_host_group" "zabbix" {
			name = "%s"
	  	}`, host, name, hostGroup,
	)
}

func testAccCheckZabbixHostExist(resource string, host *zabbix.Host) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("Not found: %s", resource)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No record ID id set")
		}

		api := testAccProvider.Meta().(*zabbix.API)
		getHost, err := api.HostGetByID(rs.Primary.ID)
		if err != nil {
			return err
		}
		*host = *getHost
		return nil
	}
}

func testAccCheckZabbixHostAttributes(host *zabbix.Host, want zabbix.Host, groupNames []string, templateNames []string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		api := testAccProvider.Meta().(*zabbix.API)

		if host.Host != want.Host {
			return fmt.Errorf("Got host name: %q, expected: %q", host.Host, want.Host)
		}
		if host.Name != want.Name {
			return fmt.Errorf("Got name: %q, expected: %q", host.Name, want.Name)
		}

		param := zabbix.Params{
			"output": "extend",
			"hostids": []string{
				host.HostID,
			},
		}
		groups, err := api.HostGroupsGet(param)
		if err != nil {
			return err
		}
		if len(groups) != len(groupNames) {
			return fmt.Errorf("Got %d groups, but expected %d groups", len(groups), len(groupNames))
		}
		for _, groupName := range groupNames {
			if !containGroup(groups, groupName) {
				return fmt.Errorf("Group not found: %s", groupName)
			}
		}

		templates, err := api.TemplatesGet(param)
		if err != nil {
			return err
		}
		for _, templateName := range templateNames {
			if !containTemplate(templates, templateName) {
				return fmt.Errorf("Template not found : %s", templateName)
			}
		}
		return nil
	}
}

func containGroup(groupNames zabbix.HostGroups, name string) bool {
	for _, group := range groupNames {
		if name == group.Name {
			return true
		}
	}
	return false
}

func containTemplate(templateNames zabbix.Templates, name string) bool {
	for _, template := range templateNames {
		if name == template.Name {
			return true
		}
	}
	return false
}
