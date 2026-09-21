#!/bin/sh


lanip=$1
wanip=$2
mask=$3

uci set network.lan.ipaddr=$lanip
uci set network.lan.netmask=$mask
uci set network.wan.ipaddr=$wanip
uci set network.wan.netmask=$mask
uci set network.wan.proto=static
uci commit

ifdown lan
ifup lan
ifdown wan
ifup wan
ifdown wg1
ifup wg1
exit 0