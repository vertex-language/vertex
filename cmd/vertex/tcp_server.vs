package main

import "linux/libc"

// Zero-cost native bindings to POSIX socket API
class Net : libc {
    func socket(domain: int32, stype: int32, proto: int32) -> int32
    func setsockopt(fd: int32, level: int32, opt: int32, val: string, len: int32) -> int32
    func bind(fd: int32, addr: string, len: int32) -> int32
    func listen(fd: int32, backlog: int32) -> int32
    func accept(fd: int32, addr: string, addrlen: string) -> int32
    func recv(fd: int32, buf: string, n: int32, flags: int32) -> int32
    func send(fd: int32, buf: string, n: int32, flags: int32) -> int32
    func close(fd: int32) -> int32
    func htons(port: uint16) -> uint16
}

class C : libc {
    func printf(fmt: string, ...) -> int32
    func exit(code: int32)
}

let AF_INET:      int32  = 2
let SOCK_STREAM:  int32  = 1
let SOL_SOCKET:   int32  = 1
let SO_REUSEADDR: int32  = 2
let PORT:         uint16 = 8080
let BACKLOG:      int32  = 128

// each accepted client spins up here in its own thread — echo server
func serve(fd: int32) thread {
    var net = Net()
    var c   = C()
    defer net.close(fd)

    // 4096-byte stack buffer — no heap, no gc
    var buf = [uint8](4096)

    while true {
        let n = net.recv(fd, buf, 4096, 0)
        if n <= 0 { break }
        net.send(fd, buf, n, 0)
        c.printf("fd=%-4d  echoed %d bytes\n", fd, n)
    }

    c.printf("fd=%-4d  disconnected\n", fd)
}

func main() -> int {
    var net = Net()
    var c   = C()

    let serverFd = net.socket(AF_INET, SOCK_STREAM, 0)
    if serverFd < 0 {
        c.printf("error: socket() failed\n")
        c.exit(1)
    }
    defer net.close(serverFd)

    // allow port reuse across quick restarts
    var opt: int32 = 1
    net.setsockopt(serverFd, SOL_SOCKET, SO_REUSEADDR, opt, 4)

    // build sockaddr_in inline: sin_family(2) | sin_port(2) | sin_addr(4) | zero(8)
    var addr = [uint8](16)
    addr[0] = uint8(AF_INET)
    addr[2] = uint8(net.htons(PORT) >> 8)
    addr[3] = uint8(net.htons(PORT) & 0xFF)
    // sin_addr stays 0 — INADDR_ANY, sin_zero stays 0 — already zeroed

    if net.bind(serverFd, addr, 16) < 0 {
        c.printf("error: bind() failed\n")
        c.exit(1)
    }

    if net.listen(serverFd, BACKLOG) < 0 {
        c.printf("error: listen() failed\n")
        c.exit(1)
    }

    c.printf("tcp echo  0.0.0.0:%d\n", PORT)

    // accept loop — non-blocking hand-off, each client owns its thread
    while true {
        let clientFd = net.accept(serverFd, nil, nil)
        if clientFd < 0 { break }
        c.printf("accepted  fd=%d\n", clientFd)
        serve(fd: clientFd).spawn()
    }

    return 0
}