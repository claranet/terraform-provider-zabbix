package provider

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider

func TestProvider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ terraform.ResourceProvider = Provider()
}

func init() {
	testAccProvider = Provider().(*schema.Provider)
	testAccProviders = map[string]terraform.ResourceProvider{
		"zabbix": testAccProvider,
	}
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("TEST_ZABBIX_URL"); v == "" {
		t.Fatal("TEST_ZABBIX_URL must be set for acceptance tests")
	}
	if v := os.Getenv("TEST_ZABBIX_USER"); v == "" {
		t.Fatal("TEST_ZABBIX_USER must be set for acceptance tests")
	}
	if v := os.Getenv("TEST_ZABBIX_PASSWORD"); v == "" {
		t.Fatal("TEST_ZABBIX_PASSWORD must be set for acceptance tests")
	}
	if v := os.Getenv("TEST_ZABBIX_VERBOSE"); v == "" {
		t.Fatal("TEST_ZABBIX_VERBOSE must be set for acceptance tests")
	}
	if v := os.Getenv("ZABBIX_SERVER_URL"); v == "" {
		t.Fatal("ZABBIX_SERVER_URL must be set for acceptance tests")
	}
	if v := os.Getenv("ZABBIX_USER"); v == "" {
		t.Fatal("ZABBIX_USER must be set for acceptance tests")
	}
	if v := os.Getenv("ZABBIX_PASSWORD"); v == "" {
		t.Fatal("ZABBIX_PASSWORD must be set for acceptance tests")
	}
}
