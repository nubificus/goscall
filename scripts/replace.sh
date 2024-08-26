#! /bin/bash

new_array=$(sed 's/^/\t\t"/;s/$/",/' $1 | paste -sd '\n')
sed -i.bak '/var syscalls = \[\]string {/,/}/{
r /dev/stdin
d
}' main.go <<EOF
	var syscalls = []string {
$new_array
	}
EOF
