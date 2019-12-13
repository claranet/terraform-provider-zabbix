resource "zabbix_lld_rule" "filesystem_discovery_rule" {
    delay = 300
    host_id = zabbix_template.base_linux_filesystem.id
    interface_id = "0"
    key = "vfs.fs.discovery"
    name = "FileSystem discovery rule"
    type = 0
    filter {
        condition {
            macro = "{#FSTYPE}"
            value = "@fs"
        }
        eval_type = 0
    }
}

resource "zabbix_item_prototype" "free_disk_inodes" {
  delay = 60
  host_id  = zabbix_template.base_linux_filesystem.id
  rule_id = zabbix_lld_rule.filesystem_discovery_rule.id
  interface_id = "0"
  key = "vfs.fs.inode[{#FSNAME},pfree]"
  name = "Free disk inodes on {#FSNAME}"
  type = 0
  status = 0
}

resource "zabbix_item_prototype" "free_disk_space" {
  delay = 60
  host_id  = zabbix_template.base_linux_filesystem.id
  rule_id = zabbix_lld_rule.filesystem_discovery_rule.id
  interface_id = "0"
  key = "vfs.fs.size[{#FSNAME},pfree]"
  name = "Free disk space on {#FSNAME}"
  type = 0
  status = 0
}

resource "zabbix_trigger_prototype" "free_disk_inodes_disaster" {
  description = "${format("Disk [{#FSNAME}]: Free inodes ({ITEM.LASTVALUE}) < {$INODE_DISASTER:\"{#FSNAME}\"}%%")}"
  expression = "${format("{%s:%s.last()}<{$INODE_DISASTER:\"{#FSNAME}}\"}", zabbix_template.base_linux_filesystem.host, zabbix_item_prototype.free_disk_inodes.key)}"
  priority = 5
}

resource "zabbix_trigger_prototype" "free_disk_inodes_high" {
  description = "${format("Disk [{#FSNAME}]: Free inodes ({ITEM.LASTVALUE}) < {$INODE_HIGH:\"{#FSNAME}\"}%%")}"
  expression = "${format("{%s:%s.last()}<{$INODE_HIGH:\"{#FSNAME}}\"}", zabbix_template.base_linux_filesystem.host, zabbix_item_prototype.free_disk_inodes.key)}"
  priority = 4
  dependencies = [
    zabbix_trigger_prototype.free_disk_inodes_disaster.id
  ]
}

resource "zabbix_trigger_prototype" "free_disk_inodes_avg" {
  description = "${format("Disk [{#FSNAME}]: Free inodes ({ITEM.LASTVALUE}) < {$INODE_AVG:\"{#FSNAME}\"}%%")}"
  expression = "${format("{%s:%s.last()}<{$INODE_AVG:\"{#FSNAME}}\"}", zabbix_template.base_linux_filesystem.host, zabbix_item_prototype.free_disk_inodes.key)}"
  priority = 3
  dependencies = [
    zabbix_trigger_prototype.free_disk_inodes_high.id
  ]
}

resource "zabbix_trigger_prototype" "free_disk_inodes_warn" {
  description = "${format("Disk [{#FSNAME}]: Free inodes ({ITEM.LASTVALUE}) < {$INODE_WARN:\"{#FSNAME}\"}%%")}"
  expression = "${format("{%s:%s.last()}<{$INODE_WARN:\"{#FSNAME}}\"}", zabbix_template.base_linux_filesystem.host, zabbix_item_prototype.free_disk_inodes.key)}"
  priority = 2
  dependencies = [
    zabbix_trigger_prototype.free_disk_inodes_avg.id
  ]
}

resource "zabbix_trigger_prototype" "free_disk_space_disaster" {
  description = "${format("Disk [{#FSNAME}]: Free space ({ITEM.LASTVALUE}) < {$DISK_SPACE_DISASTER:\"{#FSNAME}\"}%%")}"
  expression = "${format("{%s:%s.last()}<{$DISK_SPACE_DISASTER:\"{#FSNAME}\"}", zabbix_template.base_linux_filesystem.host, zabbix_item_prototype.free_disk_space.key)}"
  priority = 5
}

resource "zabbix_trigger_prototype" "free_disk_space_high" {
  description = "${format("Disk [{#FSNAME}]: Free space ({ITEM.LASTVALUE}) < {$DISK_SPACE_HIGH:\"{#FSNAME}\"}%%")}"
  expression = "${format("{%s:%s.last()}<{$DISK_SPACE_HIGH:\"{#FSNAME}\"}", zabbix_template.base_linux_filesystem.host, zabbix_item_prototype.free_disk_space.key)}"
  priority = 4
  dependencies = [
    zabbix_trigger_prototype.free_disk_space_disaster.id
  ]
}

resource "zabbix_trigger_prototype" "free_disk_space_avg" {
  description = "${format("Disk [{#FSNAME}]: Free space ({ITEM.LASTVALUE}) < {$DISK_SPACE_AVG:\"{#FSNAME}\"}%%")}"
  expression = "${format("{%s:%s.last()}<{$DISK_SPACE_AVG:\"{#FSNAME}\"}", zabbix_template.base_linux_filesystem.host, zabbix_item_prototype.free_disk_space.key)}"
  priority = 3
  dependencies = [
    zabbix_trigger_prototype.free_disk_space_high.id
  ]
}

resource "zabbix_trigger_prototype" "free_disk_space_warn" {
  description = "${format("Disk [{#FSNAME}]: Free space ({ITEM.LASTVALUE}) < {$DISK_SPACE_WARN:\"{#FSNAME}\"}%%")}"
  expression = "${format("{%s:%s.last()}<{$DISK_SPACE_WARN:\"{#FSNAME}\"}", zabbix_template.base_linux_filesystem.host, zabbix_item_prototype.free_disk_space.key)}"
  priority = 2
  dependencies = [
    zabbix_trigger_prototype.free_disk_space_avg.id
  ]
}

resource "zabbix_lld_rule_link" "test_link" {
  lld_rule_id = zabbix_lld_rule.filesystem_discovery_rule.id

  dynamic "item_prototype" {
    for_each = [
      zabbix_item_prototype.free_disk_inodes.id,
      zabbix_item_prototype.free_disk_space.id
    ]

    content {
      item_id = item_prototype.value
    }
  }

  dynamic "trigger_prototype" {
    for_each = [
      zabbix_trigger_prototype.free_disk_inodes_disaster.id,
      zabbix_trigger_prototype.free_disk_inodes_high.id,
      zabbix_trigger_prototype.free_disk_inodes_avg.id,
      zabbix_trigger_prototype.free_disk_inodes_warn.id,
      zabbix_trigger_prototype.free_disk_space_disaster.id,
      zabbix_trigger_prototype.free_disk_space_high.id,
      zabbix_trigger_prototype.free_disk_space_avg.id,
      zabbix_trigger_prototype.free_disk_space_warn.id
    ]

    content {
      trigger_id = trigger_prototype.value
    }
  }
}