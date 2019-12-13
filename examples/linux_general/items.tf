resource "zabbix_item" "cpu_load_avg1" {
  name    = "CPU Load AVG 1min"
  key     = "system.cpu.load[,avg1]"
  delay   = 60
  history = "360d"
  trends  = "34d"
  host_id = zabbix_template.base_linux_general.id
}

resource "zabbix_item" "cpu_util_idle" {
  name    = "CPU % Idle"
  key     = "system.cpu.util[,idle,]"
  delay   = 60
  history = "90d"
  trends  = "365d"
  host_id = zabbix_template.base_linux_general.id
}

resource "zabbix_item" "cpu_num_online" {
  name    = "CPU_number"
  key     = "system.cpu.num[online]"
  delay   = 300
  history = "1d"
  trends  = "7d"
  host_id = zabbix_template.base_linux_general.id
}

resource "zabbix_item" "memory_size_pavailable" {
  name    = "Memory_percent_available"
  key     = "vm.memory.size[pavailable]"
  delay   = 60
  history = "7d"
  trends  = "365d"
  host_id = zabbix_template.base_linux_general.id
}

resource "zabbix_item" "ssh_server_is_running" {
  name = "SSH server is running"
  key = "net.tcp.service[tcp]"
  delay = 30
  history = "7d"
  trends = "365d"
  host_id = zabbix_template.base_linux_network.id
}
