#!/bin/sh
# arg 1 ($1) - WG address of hub, e.g. 10.48...
# arg 2 ($2) - id, eg 9133

wgip=$1
vxlanid=$2

ifdown vxlan1

uci -q delete network.vxlan1
uci -q del_list network.br_lan.ports='vxlan1'
uci set network.vxlan1=device
uci set network.vxlan1.type='vxlan'
uci set network.vxlan1.name='vxlan1'
uci set network.vxlan1.port='4789'
uci set network.vxlan1.proxy='0'
uci set network.vxlan1.id=$vxlanid
uci set network.vxlan1.learning='0'
uci set network.vxlan1.udp6zerocsumtx='0'
uci set network.vxlan1.rsc='0'
uci set network.vxlan1.udp6zerocsumrx='0'
uci set network.vxlan1.l3miss='0'
uci set network.vxlan1.l2miss='0'
uci set network.vxlan1.udpcsum='0'
uci set network.vxlan1.remote=$wgip
uci add_list network.br_lan.ports='vxlan1'
uci commit

ifup vxlan1

/etc/init.d/firewall restart 2>/dev/null >/dev/null
exit 0