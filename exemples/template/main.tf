provider "zabbix" {
  user = "${var.user}"
  password = "${var.password}"
  server_url = "http://localhost/api_jsonrpc.php" 
}

resource "zabbix_host_group" "TEMPLATE_LINUX" {
  name = "TEMPLATE_LINUX"
}


resource "zabbix_template" "Base_Linux_General" {
  host = "Base_Linux_General"
  groups = ["${zabbix_host_group.TEMPLATE_LINUX.name}"]
  description = "Complex template exemple"
  macro = {
    CPU_AVG = "85"
    CPU_DISASTER = "95"
    CPU_HIGH = "90"
    CPU_INTERVAL = "60m"
    CPU_LOAD_RATIO_AVG = "2"
    CPU_LOAD_RATIO_DISASTER = "3"
    CPU_LOAD_RATIO_HIGH = "2.5"
    CPU_LOAD_RATIO_INTERVAL = "30m"
    CPU_LOAD_RATIO_WARN = "1.5"
    CPU_WARN = "80"
    MEMORY_PERCENTAGE_AVG = "10"
    MEMORY_PERCENTAGE_DISABLE = "2"
    MEMORY_PERCENTAGE_HIGH = "5"
    MEMORY_PERCENTAGE_WARN = "15"
  }
}

resource "zabbix_item" "CPU_Load_avg_1min" {
  name = "CPU Load AVG 1min"
  key = "system.cpu.load[,avg1]"
  delay = 60
  history = 90
  trends = 90
  host_id = "${zabbix_template.Base_Linux_General.template_id}"
}

resource "zabbix_item" "CPU_PERCENTAGE_IDLE" {
  name = "CPU % Idle"
  key = "system.cpu.util[,idle,]"
  delay = 60
  history = 90
  trends = 365
  host_id = "${zabbix_template.Base_Linux_General.template_id}"
}

resource "zabbix_item" "CPU_number" {
  name = "CPU_number"
  key = "system.cpu.num[online]"
  delay = 300
  history = 1
  trends = 7
  host_id = "${zabbix_template.Base_Linux_General.template_id}"
}

resource "zabbix_item" "Memory_percent_available" {
  name = "Memory_percent_available"
  key = "vm.memory.size[pavailable]"
  delay = 60
  history = 7
  trends = 365
  host_id = "${zabbix_template.Base_Linux_General.template_id}"
}

resource "zabbix_item" "aaaa" {
  name = "cxccx"
  key = "aaa.aaa"
  delay = 60
  history = 7
  trends = 365
  host_id = "${zabbix_template.Base_Linux_General.template_id}"
}

resource "zabbix_trigger" "CPU_load_ratio_disaster" {
  description = "CPU: Load Ratio ({ITEM.LASTVALUE}) > {$CPU_LOAD_RATIO_DISASTER} during the last {$CPU_LOAD_RATIO_INTERVAL}"
  expression = "{${zabbix_template.Base_Linux_General.host}:${zabbix_item.CPU_Load_avg_1min.key}.min({$CPU_LOAD_RATIO_INTERVAL})} / {${zabbix_template.Base_Linux_General.host}:${zabbix_item.CPU_number.key}.min({$CPU_LOAD_RATIO_INTERVAL})} > {$CPU_LOAD_RATIO_DISASTER}"
  priority = 5
}

resource "zabbix_trigger" "CPU_load_ratio_high" {
  description = "CPU: Load Ratio ({ITEM.LASTVALUE}) > {$CPU_LOAD_RATIO_HIGH} during the last {$CPU_LOAD_RATIO_INTERVAL}"
  expression = "{${zabbix_template.Base_Linux_General.host}:${zabbix_item.CPU_Load_avg_1min.key}.min({$CPU_LOAD_RATIO_INTERVAL})} / {${zabbix_template.Base_Linux_General.host}:${zabbix_item.CPU_number.key}.min({$CPU_LOAD_RATIO_INTERVAL})} > {$CPU_LOAD_RATIO_HIGH}"
  priority = 4
  dependencies = [
    zabbix_trigger.CPU_load_ratio_disaster.trigger_id
  ]
}

resource "zabbix_trigger" "CPU_load_ratio_avg" {
  description = "CPU: Load Ratio ({ITEM.LASTVALUE}) > {$CPU_LOAD_RATIO_AVG} during the last {$CPU_LOAD_RATIO_INTERVAL}"
  expression = "{${zabbix_template.Base_Linux_General.host}:${zabbix_item.CPU_Load_avg_1min.key}.min({$CPU_LOAD_RATIO_INTERVAL})} / {${zabbix_template.Base_Linux_General.host}:${zabbix_item.CPU_number.key}.min({$CPU_LOAD_RATIO_INTERVAL})} > {$CPU_LOAD_RATIO_AVG}"
  priority = 3
  dependencies = [
    zabbix_trigger.CPU_load_ratio_high.trigger_id
  ]
}

resource "zabbix_trigger" "CPU_load_ratio_warn" {
  description = "CPU: Load Ratio ({ITEM.LASTVALUE}) > {$CPU_LOAD_RATIO_WAN} during the last {$CPU_LOAD_RATIO_INTERVAL}"
  expression = "{${zabbix_template.Base_Linux_General.host}:${zabbix_item.CPU_Load_avg_1min.key}.min({$CPU_LOAD_RATIO_INTERVAL})} / {${zabbix_template.Base_Linux_General.host}:${zabbix_item.CPU_number.key}.min({$CPU_LOAD_RATIO_INTERVAL})} > {$CPU_LOAD_RATIO_WARN}"
  priority = 2
  dependencies = [
    zabbix_trigger.CPU_load_ratio_avg.trigger_id
  ]
}

resource "zabbix_trigger" "CPU_utilization_disaster" {
  description = "CPU: Utilization ({ITEM.LASTVALUE}) > {$CPU_DISASTER}% during the last {$CPU_INTERVAL}"
  expression = "100 - {${zabbix_template.Base_Linux_General.host}:${zabbix_item.CPU_PERCENTAGE_IDLE.key}.max({$CPU_INTERVAL})} > {$CPU_DISASTER}"
  priority = 5
}

resource "zabbix_trigger" "CPU_utilization_high" {
  description = "CPU: Utilization ({ITEM.LASTVALUE}) > {$CPU_AVG}% during the last {$CPU_INTERVAL}"
  expression = "100 - {${zabbix_template.Base_Linux_General.host}:${zabbix_item.CPU_PERCENTAGE_IDLE.key}.max({$CPU_INTERVAL})} > {$CPU_HIGH}"
  priority = 4
  dependencies = [
    zabbix_trigger.CPU_utilization_disaster.trigger_id
  ]
}

resource "zabbix_trigger" "CPU_utilization_avg" {
  description = "	CPU: Utilization ({ITEM.LASTVALUE}) > {$CPU_HIGH}% during the last {$CPU_INTERVAL}"
  expression = "100 - {${zabbix_template.Base_Linux_General.host}:${zabbix_item.CPU_PERCENTAGE_IDLE.key}.max({$CPU_INTERVAL})} > {$CPU_AVG}"
  priority = 3
  dependencies = [
    zabbix_trigger.CPU_utilization_high.trigger_id
  ]
}

resource "zabbix_trigger" "CPU_utilization_warn" {
  description = "CPU: Utilization ({ITEM.LASTVALUE}) > {$CPU_WARN}% during the last {$CPU_INTERVAL}"
  expression = "100 - {${zabbix_template.Base_Linux_General.host}:${zabbix_item.CPU_PERCENTAGE_IDLE.key}.max({$CPU_INTERVAL})} > {$CPU_WARN}"
  priority = 2
  dependencies = [
    zabbix_trigger.CPU_utilization_avg.trigger_id
  ]
}

resource "zabbix_trigger" "Memory_free_space_disaster" {
  description = "Memory: Free space ({ITEM.LASTVALUE}) < {$MEMORY_PERCENTAGE_DISASTER}%"
  expression = "{${zabbix_template.Base_Linux_General.host}:${zabbix_item.Memory_percent_available.key}.last()} < {$MEMORY_PERCENTAGE_DISASTER}"
  priority = 5
}

resource "zabbix_trigger" "Memory_free_space_high" {
  description = "Memory: Free space ({ITEM.LASTVALUE}) < {$MEMORY_PERCENTAGE_HIGH}%"
  expression = "{${zabbix_template.Base_Linux_General.host}:${zabbix_item.Memory_percent_available.key}.last()} < {$MEMORY_PERCENTAGE_HIGH}"
  priority = 4
  dependencies = [
    zabbix_trigger.Memory_free_space_disaster.trigger_id
  ]
}

resource "zabbix_trigger" "Memory_free_space_avg" {
  description = "Memory: Free space ({ITEM.LASTVALUE}) < {$MEMORY_PERCENTAGE_AVG}%"
  expression = "{${zabbix_template.Base_Linux_General.host}:${zabbix_item.Memory_percent_available.key}.last()} < {$MEMORY_PERCENTAGE_AVG}"
  priority = 3
  dependencies = [
    zabbix_trigger.Memory_free_space_high.trigger_id
  ]
}

resource "zabbix_trigger" "Memory_free_space_warn" {
  description = "Memory: Free space ({ITEM.LASTVALUE}) < {$MEMORY_PERCENTAGE_WARN}%"
  expression = "{${zabbix_template.Base_Linux_General.host}:${zabbix_item.Memory_percent_available.key}.last()} < {$MEMORY_PERCENTAGE_WARN}"
  priority = 2
  dependencies = [
    zabbix_trigger.Memory_free_space_avg.trigger_id
  ]
}

# This virtual resource is responsible of ensuring no other items are associated to the template
resource "zabbix_template_link" "my_zbx_template_items" {
  template_id = zabbix_template.Base_Linux_General.id
  item = [
    zabbix_item.CPU_Load_avg_1min.id,
    zabbix_item.CPU_PERCENTAGE_IDLE.id,
    zabbix_item.CPU_number.id,
    zabbix_item.Memory_percent_available.id,
    zabbix_item.aaaa.id
  ]
  trigger = [
    zabbix_trigger.CPU_load_ratio_disaster.id,
    zabbix_trigger.CPU_load_ratio_high.id,
    zabbix_trigger.CPU_load_ratio_avg.id,
    zabbix_trigger.CPU_load_ratio_warn.id,
    zabbix_trigger.CPU_utilization_disaster.id,
    zabbix_trigger.CPU_utilization_high.id,
    zabbix_trigger.CPU_utilization_avg.id,
    zabbix_trigger.CPU_utilization_warn.id,
    zabbix_trigger.Memory_free_space_disaster.id,
    zabbix_trigger.Memory_free_space_high.id,
    zabbix_trigger.Memory_free_space_avg.id,
    zabbix_trigger.Memory_free_space_warn.id,
  ]
}