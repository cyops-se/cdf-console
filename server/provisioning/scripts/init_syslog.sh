#!/bin/sh

# arg 1 ($1) - Address of syslog server
# arg 2 ($2) - Port of syslog server
# arg 3 ($3) - Connection type

ip=$1
port=$2
proto=$3

uci set system.default=remote_logger
# uci set system.default.log_proto=$proto
uci set system.default.log_proto='udp'
# uci set system.default.log_port=$port
uci set system.default.log_port='8514'
uci set system.default.log_hostname='0'
# uci set system.default.log_ip=$ip
uci set system.default.log_ip='10.49.88.254'
uci commit

exit 0
