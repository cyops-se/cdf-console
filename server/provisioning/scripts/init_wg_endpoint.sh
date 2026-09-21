#!/bin/sh
# arg 1 ($1) - IP address of WG server
# arg 2 ($2) - IP address of WG peer
# arg 3 ($3) - Public key WG peer

wgip=$1
peerip=$2
peerpuk=$3

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
uci set network.wg1.addresses=$wgip'/22'
uci set network.wghub1=wireguard_wg1
uci set network.wghub1.public_key=$peerpuk
uci set network.wghub1.force_tunlink='0'
uci set network.wghub1.route_allowed_ips='0'
uci set network.wghub1.tunlink='any'
uci set network.wghub1.allowed_ips='0.0.0.0/0'
uci set network.wghub1.endpoint_host=$peerip
uci commit
exit 0