resource "zabbix_lld_rule" "test_lld_rule" {
    delay = 60
    host_id = zabbix_template.base_linux_general.id
    interface_id = "0"
    key = "ll.lolo"
    name = "test_low_level_discovery_rulemm"
    type = 0
    filter {
        condition {
            macro = "{#SLT}"
            value = "^lo$"
        }
        eval_type = 0
    }
}

resource "zabbix_item_prototype" "test_item_prototype" {
  delay = 60
  host_id  = zabbix_template.base_linux_general.id
  rule_id = zabbix_lld_rule.test_lld_rule.id
  interface_id = "0"
  key = "sort"
  name = "hmmmmmmmmmmmmm"
  type = 0
  status = 0
}

resource "zabbix_trigger_prototype" "test_trigger_prototype" {
  description = "basic name"
  expression = "{${zabbix_template.base_linux_general.host}:${zabbix_item_prototype.test_item_prototype.key}.last()}=0"
  priority = 5
}
