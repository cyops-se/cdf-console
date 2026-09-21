#!/bin/sh
# arg 1 ($1) - Peer identity (9133)
# arg 2 ($2) - Public WG key of peer
# arg 3 ($3) - Wan IP address of peer
# arg 4 ($4) - Allowed IP addresses (typically only the peer WG address)

peerid=$1
peerpubkey=$2
peerip=$3
allowedip=$4

ifdown wg1
uci -q delete network.ctrvpn
uci set network.$peerid=wireguard_wg1
uci set network.$peerid.allowed_ips=$allowedip
uci set network.$peerid.public_key=$peerpubkey
uci set network.$peerid.force_tunlink='0'
uci set network.$peerid.endpoint_host=$peerip
uci set network.$peerid.route_allowed_ips='0'
uci set network.$peerid.tunlink='any'
uci commit

ifup wg1
exit 0