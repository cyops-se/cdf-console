#!/bin/sh

uci set firewall.20=rule
uci set firewall.20.dest_port='51820'
uci set firewall.20.src='wan'
uci set firewall.20.name='Allow-wireguard_wg1-traffic'
uci set firewall.20.target='ACCEPT'
uci set firewall.20.vpn_type='wireguard'
uci set firewall.20.proto='udp'
uci set firewall.20.family='ipv4'
uci set firewall.21=zone
uci set firewall.21.name='wireguard'
uci set firewall.21.forward='REJECT'
uci set firewall.21.network='wg1'
uci set firewall.21.output='ACCEPT'
uci set firewall.21.masq='0'
uci set firewall.21.input='ACCEPT'
uci set firewall.22=forwarding
uci set firewall.22.dest='lan'
uci set firewall.22.src='wireguard'
uci set firewall.23=forwarding
uci set firewall.23.dest='wireguard'
uci set firewall.23.src='lan'
uci set firewall.24=rule
uci set firewall.24.name='Allow VXLAN on port 4789'
uci set firewall.24._vxlan='1'
uci set firewall.24.target='ACCEPT'
uci set firewall.24.enabled='1'
uci set firewall.24.proto='udp'
uci set firewall.24.dest_port='4789'
uci set firewall.24.src='wireguard'
uci set firewall.24.family='ipv4'
uci commit

/etc/init.d/firewall restart 2>/dev/null >/dev/null
exit 0