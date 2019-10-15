provider "zabbix" {
  user = "${var.user}"
  password = "${var.password}"
  server_url = "http://localhost/api_jsonrpc.php" 
}

resource "zabbix_host_group" "demo_group" {
  name = "Group demo"
}

resource "zabbix_template" "demo_template" {
  host = "simple template"
  name = "simple template demo"
  description = "A simple template exemple"
  groups = ["${zabbix_host_group.demo_group.name}"]
  macro = {
      MACRO_TEMPLATE = "12"
  }
}
