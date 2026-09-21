ebtables -F

ebtables -A OUTPUT -p 0x0806 -j ACCEPT
# Loop for each, add logging
ebtables -A OUTPUT -p IPv4 -o eth0 --ip-dst 172.16.90.25 -j ACCEPT
ebtables -A OUTPUT -p IPv4 -o eth0 --ip-dst 172.16.90.26 -j ACCEPT

# Loop for each vxlan, add logging
ebtables -A OUTPUT -p IPv4 -o vxlan9130 --ip-dst 172.16.91.30 -j ACCEPT
ebtables -A OUTPUT -p IPv4 -o vxlan9131 --ip-dst 172.16.91.31 -j ACCEPT
ebtables -A OUTPUT -p IPv4 -o vxlan9132 --ip-dst 172.16.91.32 -j ACCEPT
ebtables -A OUTPUT -p IPv4 -o vxlan9133 --ip-dst 172.16.91.33 -j ACCEPT

# iptables -A OUTPUT -4 -o vxlan+ -j LOG
ebtables -A OUTPUT -o vxlan+ -j DROP
