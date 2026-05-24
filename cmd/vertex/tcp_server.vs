package main
import "linux/libc"

class C : libc {
    func socket(domain: int, sockType: int, protocol: int) -> int
    func bind(sockfd: int, addr: any void, addrlen: int) -> int  // any void → ptr ✓
    func listen(sockfd: int, backlog: int) -> int
    func accept(sockfd: int, addr: mut any void, addrlen: mut any void) -> int
    func write(fd: int, buf: any void, count: int) -> int
    func close(fd: int) -> int
    func puts(str: any char) -> int
}

// ── TCP structs ───────────────────────────────────────────────────────────────
// Mirrors the C sockaddr_in layout (16 bytes total)
struct SockAddrIn {
    var sin_family: int16
    var sin_port:   uint16
    var sin_addr:   uint32
    var sin_zero1:  uint32
    var sin_zero2:  uint32
}

// ── entry point ───────────────────────────────────────────────────────────────
func main() -> int {
    var c = C()

    // 1. Create socket (AF_INET = 2, SOCK_STREAM = 1, IPPROTO_TCP = 0)
    var server_fd = c.socket(2, 1, 0)
    if server_fd < 0 {
        c.puts("Failed to create socket".any())
        return 1
    }

    // 2. Configure address
    // Port 8080 = 0x1F90. Byte-swapped to network order = 0x901F.
    var port: uint16 = 0x901F
    var addr = SockAddrIn{
        sin_family: int16(2),
        sin_port:   port,
        sin_addr:   uint32(0),
        sin_zero1:  uint32(0),
        sin_zero2:  uint32(0),
    }

    // 3. Bind — SockAddrIn.sizeof() resolves at compile time
    var bind_res = c.bind(server_fd, addr.any(), SockAddrIn.sizeof())
    if bind_res < 0 {
        c.puts("Bind failed. Port 8080 might be in use.".any())
        return 1
    }

    // 4. Listen
    var listen_res = c.listen(server_fd, 10)
    if listen_res < 0 {
        c.puts("Listen failed".any())
        return 1
    }

    c.puts("Server listening on port 8080...".any())

    // 5. Accept loop
    // client_addr and client_len are var bindings so .any() yields mutable pointers
    var client_addr = SockAddrIn{
        sin_family: int16(0),
        sin_port:   uint16(0),
        sin_addr:   uint32(0),
        sin_zero1:  uint32(0),
        sin_zero2:  uint32(0),
    }
    var client_len: int = SockAddrIn.sizeof()

    while true {
        var client_fd = c.accept(server_fd, client_addr.any(), client_len.any())
        if client_fd < 0 {
            c.puts("Accept failed".any())
            continue
        }

        c.puts("Client connected!".any())

        var response = "HTTP/1.1 200 OK\r\nContent-Length: 13\r\n\r\nHello Vertex!"
        c.write(client_fd, response.any(), 52)
        c.close(client_fd)
    }

    c.close(server_fd)
    return 0
}