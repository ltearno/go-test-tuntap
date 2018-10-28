# Tests in GO lang

This one is a little program creating a TUN/TAP device (so requires a linux kernel) and reads incoming bytes through a go routine.

The bytes are then sent through a Go channel for the application to use.

The bytes are then displayed on the console.

Uses packages "syscall" and "x/sys/unix".