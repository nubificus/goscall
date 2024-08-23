package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"syscall"
	seccomp "github.com/elastic/go-seccomp-bpf"
)

func usage() {
	fmt.Println("Usage of goexec: goexec [flags] -- <app> [app_args]")
	fmt.Println("Flags:")
	flag.PrintDefaults()
}

func setSeccompFilter(strictMode bool) error {
	//Create a filter.
	var syscalls = []string {
		"exit",
	}
	var secPolicy seccomp.Policy
	if strictMode {
		secPolicy.DefaultAction = seccomp.ActionTrap
	} else {
		secPolicy.DefaultAction = seccomp.ActionLog
	}
	secPolicy.Syscalls = []seccomp.SyscallGroup {
		{
			Action: seccomp.ActionAllow,
			Names: syscalls,
		},
	}
	filter := seccomp.Filter{
	        NoNewPrivs: true,
	        Flag:       seccomp.FilterFlagTSync,
	        Policy:     secPolicy,
	}

	// Load it. This will set no_new_privs before loading.
	return seccomp.LoadFilter(filter)
}

func main() {
	var seccomp bool
	var secMode bool

	flag.BoolVar(&seccomp, "seccomp", false, "Use seccomp filter")
	flag.BoolVar(&secMode, "strict", false, "Strict policy for seccomp. Not allowed syscalls will trap")
	flag.Usage = usage
	flag.Parse()

	if flag.NArg() < 1 {
		flag.Usage()
		return
	}

	appArgs := flag.Args()
	binary, err := exec.LookPath(appArgs[0])
	if err != nil {
		fmt.Println(appArgs[0], ": was not found in PATH")
		return
	}
	if seccomp {
		err := setSeccompFilter(secMode)
		if err != nil {
	        	fmt.Println("failed to load filter: ", err)
	        	return
		}
	}

	err = syscall.Exec(binary, appArgs, os.Environ()) //nolint: gosec
	if err != nil {
		fmt.Println("failed to exec ", appArgs[0], ":",  err)
	        return
	}

	return
}
