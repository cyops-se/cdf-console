#!/bin/sh

uci set ntpclient.1.hostname='[hubip]'
uci set ntpclient.ntpclient.force='1'
uci commit

/etc/init.d/ntpclient restart 2>/dev/null >/dev/null
exit 0