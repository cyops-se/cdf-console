#!/bin/sh

for i in $(seq 1 25); do uci -q delete firewall.$i; done

uci set firewall.1=defaults
uci set firewall.1.flow_offloading='1'
uci set firewall.1.flow_offloading_hw='1'
uci set firewall.1.syn_flood='1'
uci set firewall.1.output='ACCEPT'
uci set firewall.1.forward='REJECT'
uci set firewall.1.input='ACCEPT'
uci set firewall.2=zone
uci set firewall.2.name='lan'
uci set firewall.2.network='lan'
uci set firewall.2.output='ACCEPT'
uci set firewall.2.forward='ACCEPT'
uci set firewall.2.input='ACCEPT'
uci set firewall.3=zone
uci set firewall.3.name='wan'
uci set firewall.3.output='ACCEPT'
uci set firewall.3.forward='REJECT'
uci set firewall.3.mtu_fix='1'
uci set firewall.3.port_scan='0'
uci set firewall.3.network='wan'
uci set firewall.3.masq='0'
uci set firewall.3.input='REJECT'
uci set firewall.4=forwarding
uci set firewall.4.src='lan'
uci set firewall.4.dest='wan'
uci set firewall.6=rule
uci set firewall.6.name='Allow-Ping'
uci set firewall.6.proto='icmp'
uci set firewall.6.family='ipv4'
uci set firewall.6.target='ACCEPT'
uci set firewall.6.priority='2'
uci set firewall.6.src='*'
uci set firewall.14=include
uci set firewall.14.path='/etc/firewall.user'
uci set firewall.15=rule
uci set firewall.15.proto='tcp'
uci set firewall.15.name='Enable_SSH_WAN'
uci set firewall.15.target='ACCEPT'
uci set firewall.15.src='wan'
uci set firewall.15.priority='10'
uci set firewall.15.enabled='1'
uci set firewall.15.dest_port='22'
uci set firewall.15.family='ipv4'
uci set firewall.16=rule
uci set firewall.16.proto='tcp'
uci set firewall.16.name='Enable_HTTP_WAN'
uci set firewall.16.target='ACCEPT'
uci set firewall.16.src='wan'
uci set firewall.16.priority='11'
uci set firewall.16.enabled='1'
uci set firewall.16.dest_port='80 443'
uci set firewall.16.family='ipv4'
uci commit

/etc/init.d/firewall restart 2>/dev/null >/dev/null
exit 0