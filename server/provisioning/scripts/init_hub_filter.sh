#!/bin/sh

file=/etc/firewall.user

cat > "$file" << EOF
# Flush/Reset
etables -F

# Accept ARP
ebtables -A FORWARD -p 0x0806 -j ACCEPT
EOF

# interface;dstip;target;log
for arg in "$@"; do
    iface=(echo "$arg" | cut -d ";" -f 1)
    ip=(echo "$arg" | cut -d ";" -f 2)
    target=(echo "$arg" | cut -d ";" -f 3)
    log=(echo "$arg" | cut -d ";" -f 4)

    if [[ log == 1 ]]; then
      echo "iptables -A FORWARD -p IPv4 -i $iface --ip-dst $ip -j LOG" >> "$file"
    fi
    echo "ebtales -A FORWARD -p IPv4 -i $iface --ip-dst $ip -j $target" >> "$file"
done

# cat >> /etc/firewall.user << EOF
# EOF

# Accessible from endpoints
# ebtables -A FORWARD -p IPv4 -i eth0 --ip-dst 172.16.90.25 -j ACCEPT
# ebtables -A FORWARD -p IPv4 -i eth0 --ip-dst 172.16.90.26 -j ACCEPT

# Only send to needed vxlan
# ebtables -A FORWARD -p IPv4 -i vxlan9130 --ip-dst 172.16.91.30 -j ACCEPT
# ebtables -A FORWARD -p IPv4 -i vxlan9131 --ip-dst 172.16.91.31 -j ACCEPT
# ebtables -A FORWARD -p IPv4 -i vxlan9132 --ip-dst 172.16.91.32 -j ACCEPT
# ebtables -A FORWARD -p IPv4 -i vxlan9133 --ip-dst 172.16.91.33 -j ACCEPT

# Get the value of N from the user (or set it directly if known)
# N=2  # Example: Skip the first 2 arguments

# Iterate over arguments starting from the (N+1)th one
# for arg in "${@:$((N + 1))}"; do
#     echo "$arg"
# done

# Drop all others
echo "ebtables -A FORWARD -i vxlan+ -j DROP" >> "$file"

/etc/init.d/firewall restart 2>/dev/null >/dev/null
