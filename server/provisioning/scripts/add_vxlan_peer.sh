#!/bin/sh
# arg 1 ($1) - WG address of peer, e.g. 10.48...
# arg 2 ($2) - id, eg 9133

wgip=$1
vxlanid=$2

ifdown vxlan$vxlanid

uci -q delete network.vxlan$vxlanid
uci -q del_list network.br_lan.ports=vxlan$vxlanid
uci set network.vxlan$vxlanid=device
uci set network.vxlan$vxlanid.type='vxlan'
uci set network.vxlan$vxlanid.name=vxlan$vxlanid
uci set network.vxlan$vxlanid.port='4789'
uci set network.vxlan$vxlanid.proxy='0'
uci set network.vxlan$vxlanid.id=$vxlanid
uci set network.vxlan$vxlanid.learning='0'
uci set network.vxlan$vxlanid.udp6zerocsumtx='0'
uci set network.vxlan$vxlanid.rsc='0'
uci set network.vxlan$vxlanid.udp6zerocsumrx='0'
uci set network.vxlan$vxlanid.l3miss='0'
uci set network.vxlan$vxlanid.l2miss='0'
uci set network.vxlan$vxlanid.udpcsum='0'
uci set network.vxlan$vxlanid.remote=$wgip
uci add_list network.br_lan.ports=vxlan$vxlanid
uci commit

ifup vxlan$vxlanid
exit 0