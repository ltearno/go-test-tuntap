package main

// #include <stdio.h>
// #include <linux/if.h>
// #include <linux/if_tun.h>
// #include <sys/ioctl.h>
// #include <sys/types.h>
// #include <sys/stat.h>
// #include <fcntl.h>
// #include <string.h>
// int test() { printf("hello from C!\n"); }
import "C"
import "errors"
import "fmt"
import unix "golang.org/x/sys/unix"
import "syscall"
import "unsafe"
import "encoding/hex"

const SIZEOFIREQ = 40

type ifReq struct {
	Name  [unix.IFNAMSIZ]byte
	Flags uint16
	pad   [SIZEOFIREQ - unix.IFNAMSIZ - 2]byte
}

func connectTunTap(name string) (int, error) {
	fd, err := unix.Open("/dev/net/tun", unix.O_RDWR, 0)
	if err != nil {
		return 0, errors.New("canoot open tun device")
	}

	ifRequest := ifReq{}
	ifRequest.Flags = unix.IFF_TAP
	if name != "" {
		copy(ifRequest.Name[:], []byte(name))
	}

	_, _, errno := syscall.Syscall(
		syscall.SYS_IOCTL,
		uintptr(fd),
		uintptr(unix.TUNSETIFF),
		uintptr(unsafe.Pointer(&ifRequest)),
	)

	if errno != 0 {
		return 0, errors.New("failed ioctl call")
	}

	fmt.Printf("opened tuntap device %s\n", ifRequest.Name)

	return fd, nil
}

func tunIoLoop(fd int, c chan []byte) {
	defer fmt.Printf("finished io read loop on %d", fd)

	buffer := make([]byte, 2048)

	for {
		nbRead, err := unix.Read(fd, buffer)
		if err != nil {
			return
		}

		readden := buffer[0:nbRead]
		c <- readden
	}
}

func main() {
	fmt.Println("Hello")
	C.test()
	unix.Environ()

	tunFd, err := connectTunTap("")
	if err != nil {
		fmt.Println("cannot open tuntap device !")
		panic(err)
	}

	c := make(chan []byte)

	go tunIoLoop(tunFd, c)

	for {
		select {
		case received := <-c:
			fmt.Printf("%s\n", hex.EncodeToString(received))
		}
	}
}
