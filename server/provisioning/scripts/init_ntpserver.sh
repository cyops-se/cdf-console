#!/bin/sh

uci set ntpclient.1.hostname=''
uci set ntpclient.ntpclient.force='1'
uci set ntpclient.ntpclient.enabled='0'

uci set ntpserver.general=ntpserver
uci set ntpserver.general.enabled='1'

uci set firewall.26=rule
uci set firewall.26.src='wan'
uci set firewall.26.name='NTP'
uci set firewall.26.target='ACCEPT'
uci set firewall.26.priority='12'
uci set firewall.26.dest_port='123'
uci set firewall.26.proto='tcp' 'udp'
uci set firewall.26.enabled='1'

uci commit

/etc/init.d/ntpserver restart 2>/dev/null >/dev/null
/etc/init.d/firewall restart 2>/dev/null >/dev/null
exit 0