resource "zabbix_trigger" "cpu_load_disaster" {
  description = "CPU: Load Ratio ({ITEM.LASTVALUE}) > {$CPU_LOAD_RATIO_DISASTER} during the last {$CPU_LOAD_RATIO_INTERVAL}"
  expression  = "{${zabbix_template.base_linux_general.host}:${zabbix_item.cpu_load_avg1.key}.min({$CPU_LOAD_RATIO_INTERVAL})} / {${zabbix_template.base_linux_general.host}:${zabbix_item.cpu_num_online.key}.min({$CPU_LOAD_RATIO_INTERVAL})} > {$CPU_LOAD_RATIO_DISASTER}"
  priority    = 5
}

resource "zabbix_trigger" "cpu_load_high" {
  description = "CPU: Load Ratio ({ITEM.LASTVALUE}) > {$CPU_LOAD_RATIO_HIGH} during the last {$CPU_LOAD_RATIO_INTERVAL}"
  expression  = "{${zabbix_template.base_linux_general.host}:${zabbix_item.cpu_load_avg1.key}.min({$CPU_LOAD_RATIO_INTERVAL})} / {${zabbix_template.base_linux_general.host}:${zabbix_item.cpu_num_online.key}.min({$CPU_LOAD_RATIO_INTERVAL})} > {$CPU_LOAD_RATIO_HIGH}"
  priority    = 4
  dependencies = [
    zabbix_trigger.cpu_load_disaster.trigger_id,
  ]
}

resource "zabbix_trigger" "cpu_load_avg" {
  description = "CPU: Load Ratio ({ITEM.LASTVALUE}) > {$CPU_LOAD_RATIO_AVG} during the last {$CPU_LOAD_RATIO_INTERVAL}"
  expression  = "{${zabbix_template.base_linux_general.host}:${zabbix_item.cpu_load_avg1.key}.min({$CPU_LOAD_RATIO_INTERVAL})} / {${zabbix_template.base_linux_general.host}:${zabbix_item.cpu_num_online.key}.min({$CPU_LOAD_RATIO_INTERVAL})} > {$CPU_LOAD_RATIO_AVG}"
  priority    = 3
  dependencies = [
    zabbix_trigger.cpu_load_high.trigger_id,
  ]
}

resource "zabbix_trigger" "cpu_load_warn" {
  description = "CPU: Load Ratio ({ITEM.LASTVALUE}) > {$CPU_LOAD_RATIO_WAN} during the last {$CPU_LOAD_RATIO_INTERVAL}"
  expression  = "{${zabbix_template.base_linux_general.host}:${zabbix_item.cpu_load_avg1.key}.min({$CPU_LOAD_RATIO_INTERVAL})} / {${zabbix_template.base_linux_general.host}:${zabbix_item.cpu_num_online.key}.min({$CPU_LOAD_RATIO_INTERVAL})} > {$CPU_LOAD_RATIO_WARN}"
  priority    = 2
  dependencies = [
    zabbix_trigger.cpu_load_avg.trigger_id,
  ]
}

resource "zabbix_trigger" "cpu_utilization_disaster" {
  description = "CPU: Utilization ({ITEM.LASTVALUE}) > {$CPU_DISASTER}% during the last {$CPU_INTERVAL}"
  expression  = "100 - {${zabbix_template.base_linux_general.host}:${zabbix_item.cpu_util_idle.key}.max({$CPU_INTERVAL})} > {$CPU_DISASTER}"
  priority    = 5
}

resource "zabbix_trigger" "cpu_utilization_high" {
  description = "CPU: Utilization ({ITEM.LASTVALUE}) > {$CPU_AVG}% during the last {$CPU_INTERVAL}"
  expression  = "100 - {${zabbix_template.base_linux_general.host}:${zabbix_item.cpu_util_idle.key}.max({$CPU_INTERVAL})} > {$CPU_HIGH}"
  priority    = 4
  dependencies = [
    zabbix_trigger.cpu_utilization_disaster.trigger_id,
  ]
}

resource "zabbix_trigger" "cpu_utilization_avg" {
  description = "	CPU: Utilization ({ITEM.LASTVALUE}) > {$CPU_HIGH}% during the last {$CPU_INTERVAL}"
  expression  = "100 - {${zabbix_template.base_linux_general.host}:${zabbix_item.cpu_util_idle.key}.max({$CPU_INTERVAL})} > {$CPU_AVG}"
  priority    = 3
  dependencies = [
    zabbix_trigger.cpu_utilization_high.trigger_id,
  ]
}

resource "zabbix_trigger" "cpu_utilization_warn" {
  description = "CPU: Utilization ({ITEM.LASTVALUE}) > {$CPU_WARN}% during the last {$CPU_INTERVAL}"
  expression  = "100 - {${zabbix_template.base_linux_general.host}:${zabbix_item.cpu_util_idle.key}.max({$CPU_INTERVAL})} > {$CPU_WARN}"
  priority    = 2
  dependencies = [
    zabbix_trigger.cpu_utilization_avg.trigger_id,
  ]
}

resource "zabbix_trigger" "memory_space_disaster" {
  description = "Memory: Free space ({ITEM.LASTVALUE}) < {$MEMORY_PERCENTAGE_DISASTER}%"
  expression  = "{${zabbix_template.base_linux_general.host}:${zabbix_item.memory_size_pavailable.key}.last()} < {$MEMORY_PERCENTAGE_DISASTER}"
  priority    = 5
}

resource "zabbix_trigger" "memory_space_high" {
  description = "Memory: Free space ({ITEM.LASTVALUE}) < {$MEMORY_PERCENTAGE_HIGH}%"
  expression  = "{${zabbix_template.base_linux_general.host}:${zabbix_item.memory_size_pavailable.key}.last()} < {$MEMORY_PERCENTAGE_HIGH}"
  priority    = 4
  dependencies = [
    zabbix_trigger.memory_space_disaster.trigger_id,
  ]
}

resource "zabbix_trigger" "memory_space_avg" {
  description = "Memory: Free space ({ITEM.LASTVALUE}) < {$MEMORY_PERCENTAGE_AVG}%"
  expression  = "{${zabbix_template.base_linux_general.host}:${zabbix_item.memory_size_pavailable.key}.last()} < {$MEMORY_PERCENTAGE_AVG}"
  priority    = 3
  dependencies = [
    zabbix_trigger.memory_space_high.trigger_id,
  ]
}

resource "zabbix_trigger" "memory_space_warn" {
  description = "Memory: Free space ({ITEM.LASTVALUE}) < {$MEMORY_PERCENTAGE_WARN}%"
  expression  = "{${zabbix_template.base_linux_general.host}:${zabbix_item.memory_size_pavailable.key}.last()} < {$MEMORY_PERCENTAGE_WARN}"
  priority    = 2
  dependencies = [
    zabbix_trigger.memory_space_avg.trigger_id,
  ]
}
