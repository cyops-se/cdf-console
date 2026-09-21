#!/bin/sh

newpwd=$1

passwd << EOF
$newpwd
$newpwd
EOF

passwd admin << EOF
$newpwd
$newpwd
EOF
exit 0