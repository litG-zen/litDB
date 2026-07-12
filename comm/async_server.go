package comm

import (
	"fmt"
	"github/litG-zen/litDB/conf"
	"github/litG-zen/litDB/core"
	"net"
	"os"
	"syscall"
	"time"
)

var cronFrequency time.Duration = 1 * time.Second
var lastCronRun time.Time = time.Now()

func RunAsyncServer() error {
	connected_clients := make(map[int]*FDComm)

	var max_client int = 20000

	// Epoll event array. This buffer will be holdind FD that are ready for i/o operations by system call.
	var events []syscall.EpollEvent = make([]syscall.EpollEvent, max_client)

	// Create a TCP socket.

	// AF_INET: IPv4 protocol
	// SOCK_STREAM: This tells the system that we want to keep the connection open and send data in a stream format.
	// O_NONBLOCK: This flag makes the socket non-blocking,
	// which means that system calls on this socket will return immediately if they cannot be completed,
	// rather than blocking the execution of the program until they can be completed.

	serverFD, err := syscall.Socket(syscall.AF_INET, syscall.O_NONBLOCK|syscall.SOCK_STREAM, 0)
	if err != nil {
		fmt.Printf("socket error: %s\n", err)
		os.Exit(1)
	}
	defer syscall.Close(serverFD) // close the socket when the connection is closed.

	// Set the socket to non-blocking mode.
	// This allows the server to handle multiple clients concurrently without blocking on any single client.
	if err = syscall.SetNonblock(serverFD, true); err != nil {
		return err
	}

	// Bind the IP and port to the socket.
	ip := net.ParseIP(conf.HOST)
	fmt.Printf("Server is listening on %s:%d\n", conf.HOST, conf.PORT)
	if err = syscall.Bind(serverFD, &syscall.SockaddrInet4{
		Port: conf.PORT,
		Addr: [4]byte{ip[0], ip[1], ip[2], ip[3]},
	}); err != nil {
		return err
	}

	// Listen for incoming connections.
	if err = syscall.Listen(serverFD, max_client); err != nil {
		return err
	}

	// Asyncio server loop begins here.
	epollFD, err := syscall.EpollCreate1(0)
	if err != nil {
		fmt.Printf("epoll create error: %s\n", err)
	}
	defer syscall.Close(epollFD)

	// We now want the server(serverFD) to be monitored by epoll for incoming connections.
	var socketServerEvent syscall.EpollEvent = syscall.EpollEvent{
		Events: syscall.EPOLLIN, // EPOLLIN: This event is triggered when there is data to read on the socket.
		Fd:     int32(serverFD), // The file descriptor of the server socket that we want to monitor for incoming connections.
	}

	// We add the server socket to the epoll instance using the EpollCtl function.
	if err = syscall.EpollCtl(epollFD, syscall.EPOLL_CTL_ADD, serverFD, &socketServerEvent); err != nil {
		return err
	}

	for {
		if time.Now().After(lastCronRun.Add(cronFrequency)) {
			core.DeleteExpiredKeys()
			lastCronRun = time.Now()
		}
		// Wait for events on the monitored file descriptors. This call will block until at least one event occurs.
		eventCount, err := syscall.EpollWait(epollFD, events, -1)
		if err != nil {
			return err
		} else {
			for i := 0; i < eventCount; i++ {
				event := events[i]
				var clientAddr string

				if event.Fd == int32(serverFD) {
					// This means that there is a new incoming connection on the server socket.
					connFD, client_addr, err := syscall.Accept(serverFD)
					if err != nil {
						fmt.Printf("accept error: %s\n", err)
						continue
					}

					switch addr := client_addr.(type) {
					case *syscall.SockaddrInet4:
						clientAddr = fmt.Sprintf("%s:%d", net.IP(addr.Addr[:]), addr.Port)
						fmt.Printf("New connection from %s\n", clientAddr)
					case *syscall.SockaddrInet6:
						clientAddr = fmt.Sprintf("%s:%d", net.IP(addr.Addr[:]), addr.Port)
						fmt.Printf("New connection from %s\n", clientAddr)
					default:
						fmt.Printf("Unknown client address type\n")
					}
					syscall.SetNonblock(serverFD, true)

					connected_clients[connFD] = &FDComm{Fd: connFD, ClientAddr: clientAddr}
					fmt.Printf("New client connected. Total clients: %d\n", len(connected_clients))

					var clientEvent syscall.EpollEvent = syscall.EpollEvent{
						Events: syscall.EPOLLIN, // EPOLLIN: This event is triggered when there is data to read on the socket.
						Fd:     int32(connFD),   // The file descriptor of the client socket that we want to monitor for incoming data.
					}
					// register the client socket with epoll to monitor for incoming data.
					if err = syscall.EpollCtl(epollFD, syscall.EPOLL_CTL_ADD, connFD, &clientEvent); err != nil {
						fmt.Printf("epoll ctl error: %s\n", err)
						continue
					}
				} else {
					comm := *connected_clients[int(event.Fd)]
					cmds, err := ReadCommands(&comm, comm.ClientAddr)
					if err != nil {
						fmt.Printf("%s client disconnected. Total clients left: %d\n", comm.ClientAddr, len(connected_clients)-1)
						syscall.Close(int(event.Fd))
						delete(connected_clients, int(event.Fd))
						continue
					}
					Reply(&comm, cmds)
				}

			}
		}
	}
}
