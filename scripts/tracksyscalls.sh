#! /bin/bash

strace -f -o syscall.trace "$@" &> app.out
cut -d' ' -f2- syscall.trace > no_pids
sed -n '/seccomp(/,$p' no_pids | tail -n +2 > after_seccomp
grep -o '^[a-zA-Z0-9_]\+' after_seccomp | sort | uniq
rm syscall.trace after_seccomp no_pids
