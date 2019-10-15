provider "zabbix" {
  user = "${var.user}"
  password = "${var.password}"
  server_url = "http://localhost/api_jsonrpc.php" 
}

resource "zabbix_host_group" "demo_group" {
  name = "Template demo group"
}

resource "zabbix_template" "demo_template" {
  host = "template"
  name = "template demo"
  description = "An exemple of template with item and trigger"
  groups = ["${zabbix_host_group.demo_group.name}"]
  macro = {
      MACRO_TEMPLATE = "12"
  }
}

resource "zabbix_item" "demo_item" {
  name = "demo item"
  key = "demo.key"
  delay = "34"
  description = "Item for the demo template"
  trends = "300"
  history = "25"
  host_id = "${zabbix_template.demo_template.template_id}"
}

resource "zabbix_trigger" "demo_trigger" {
  description = "demo trigger"
  expression = "{${zabbix_template.demo_template.host}:${zabbix_item.demo_item.key}.last()}={$MACRO_TEMPLATE}"
  priority = 5
  status = 0
}

# This virtual resource is responsible of ensuring no other items are associated to the template
resource "zabbix_template_link" "demo_template_link" {
  template_id = zabbix_template.demo_template.id
  item = [
    zabbix_item.demo_item.id,
  ]
  trigger = [
    zabbix_trigger.demo_trigger.id,
  ]
}