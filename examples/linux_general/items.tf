resource "zabbix_item" "cpu_load_avg1" {
  name    = "CPU Load AVG 1min"
  key     = "system.cpu.load[,avg1]"
  delay   = 60
  history = 90
  trends  = 90
  host_id = zabbix_template.base_linux_general.template_id
}

resource "zabbix_item" "cpu_util_idle" {
  name    = "CPU % Idle"
  key     = "system.cpu.util[,idle,]"
  delay   = 60
  history = 90
  trends  = 365
  host_id = zabbix_template.base_linux_general.template_id
}

resource "zabbix_item" "cpu_num_online" {
  name    = "CPU_number"
  key     = "system.cpu.num[online]"
  delay   = 300
  history = 1
  trends  = 7
  host_id = zabbix_template.base_linux_general.template_id
}

resource "zabbix_item" "memory_size_pavailable" {
  name    = "Memory_percent_available"
  key     = "vm.memory.size[pavailable]"
  delay   = 60
  history = 7
  trends  = 365
  host_id = zabbix_template.base_linux_general.template_id
}

