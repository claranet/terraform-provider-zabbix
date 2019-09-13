# Terraform zabbix provider

Allows to manage zabbix hosts

```
provider "zabbix" {
  user = "Admin"
  password = "zabbix"
  server_url = "http://localhost/api_jsonrpc.php"
}

resource "zabbix_host" "zabbix1" {
  host = "127.0.0.1"
  name = "the best name"
  interfaces {
    ip = "127.0.0.1"
    main = true
  }
  groups = ["Linux servers", "${zabbix_host_group.zabbix.name}"]
  templates = ["Template ICMP Ping"]
}

resource "zabbix_host_group" "zabbix" {
  name = "something"
}
```
