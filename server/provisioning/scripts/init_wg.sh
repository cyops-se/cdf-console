#!/bin/sh
# arg 1 ($1) - IP address of WG server

wgip=$1

ifdown wg1

uci -q delete network.wg1
uci -q delete network.wghub1
uci -q delete network.ctrvpn
uci commit

prk=`wg genkey`
puk=`echo $prk | wg pubkey`

uci set network.wg1=interface
uci set network.wg1.private_key=$prk
uci set network.wg1.public_key=$puk
uci set network.wg1.listen_port='51820'
uci set network.wg1.proto='wireguard'
uci set network.wg1.disabled='0'
uci set network.wg1.addresses=$wgip
uci commit

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
uci commit

ifup wg1

/etc/init.d/firewall restart 2>/dev/null >/dev/null
exit 0