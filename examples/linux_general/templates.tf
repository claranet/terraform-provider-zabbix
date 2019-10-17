resource "zabbix_template" "base_linux_general" {
  host        = "Base_Linux_General"
  groups      = [zabbix_host_group.template_linux.name]
  description = "Linux general template without network and disk support"
  macro = {
    CPU_AVG                   = "85"
    CPU_DISASTER              = "95"
    CPU_HIGH                  = "90"
    CPU_INTERVAL              = "60m"
    CPU_LOAD_RATIO_AVG        = "2"
    CPU_LOAD_RATIO_DISASTER   = "3"
    CPU_LOAD_RATIO_HIGH       = "2.5"
    CPU_LOAD_RATIO_INTERVAL   = "30m"
    CPU_LOAD_RATIO_WARN       = "1.5"
    CPU_WARN                  = "80"
    MEMORY_PERCENTAGE_AVG     = "10"
    MEMORY_PERCENTAGE_DISABLE = "2"
    MEMORY_PERCENTAGE_HIGH    = "5"
    MEMORY_PERCENTAGE_WARN    = "15"
  }
}

# This virtual resource is responsible of ensuring no other items are associated to the template
resource "zabbix_template_link" "base_linux_general_link" {
  template_id = zabbix_template.base_linux_general.id
  item = [
    zabbix_item.cpu_load_avg1.id,
    zabbix_item.cpu_util_idle.id,
    zabbix_item.cpu_num_online.id,
    zabbix_item.memory_size_pavailable.id,
  ]
  trigger = [
    zabbix_trigger.cpu_load_disaster.id,
    zabbix_trigger.cpu_load_high.id,
    zabbix_trigger.cpu_load_avg.id,
    zabbix_trigger.cpu_load_warn.id,
    zabbix_trigger.cpu_utilization_disaster.id,
    zabbix_trigger.cpu_utilization_high.id,
    zabbix_trigger.cpu_utilization_avg.id,
    zabbix_trigger.cpu_utilization_warn.id,
    zabbix_trigger.memory_space_disaster.id,
    zabbix_trigger.memory_space_high.id,
    zabbix_trigger.memory_space_avg.id,
    zabbix_trigger.memory_space_warn.id,
  ]
}


