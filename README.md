# Track system calls for a go exec program

This repository contains a simple Go program along with bash scripts. Using
these scripts, someone can track the system calls that the Go program uses to
exec an another application. With this information the user can also add a
seccomp filter to the Go program, allowing only the tracked system calls.

## Building

The Go program can be built by running the following commands:
```
go mod tidy
go build
```

## How to use

The Go program accepts two options:
- seccomp: Makes the Go program to use a secomp filter
- strict: Makes the seccomp filter strict, meaning that not allowed system
  calls will trap.

The `-strict` option does not have any effect, if the `-seccomp` option has
not been set. By default, the `-seccomp` option will make the Go program to use
a seccomp filter that logs all system calls. The default allowed set of system
calls consists of only the `exit` system call. However, this entry just exists
to showcase where to add new system calls.

### Running the Go program

After building the Go program, someone can execute it by running:
```
./goexec [-seccomp] [-strict] -- <app> [app arguments]
```

For instance to run the ls command:
```
./goexec -- ls / ## This should list all files and directories in /
```

To apply the seccomp filter:
```
./goexec -seccomp -- ls / ## This should list all files and directories in /
```

To apply the seccomp filter in strict mode:
```
./goexec -seccomp -strict -- ls / ## This will fail without modifications
```

> **_NOTE:_**  The above command will fail if no more system calls have been
> allowed in the seccomp filter.

### Tracking the system calls

To track the system calls of the Go program, someone can make use of the
`tracksyscalls.sh` bash script under the `scripts` directory. THe script simply
takes as an argument the application to track the system calls. Therefore for
the case of the Go program:

```
bash scripts/tracksyscalls.sh ./goexec -seccomp -- ls /
```

THe script will list all system calls used from the application using strace and
the output will be shown in stdout.

### Enabling the necessary system calls

TO automatically replace the list of allowed system calls in the Go program,
someone can make use of the `replace.sh` bash script under the `scripts/`
directory. The script takes two arguments:
1. A file containing all the system calls per line
2. The Go file that contains the system calls string array for the seccomp
   filter.

Therefore, in our case we can redirect the output of `tracksyscalls.sh` script
to a file and then use this file as the first argument for the `replace.sh`
script.

```
bash scripts/tracksyscalls.sh ./goexec -seccomp -- ls / > syscalls_list.txt
bash scripts/replace.sh syscalls_list.txt main.go
```

The above pair of commands wil ltrack the system calls of the Go program and
then it will replace the list of allowed system calls.

After rebuilding the Go program, the strict mode should run without issues:

```
./goexec -seccomp -strict -- ls / 
```

> **_NOTE:_**  The `replace.sh` script will try to find the definition of the
> allowed system calls, by searching for a pattern like the one in main.go file.
> The system calls are represented as a string array with the `syscalls` name.
> `var syscalls = []string {`
