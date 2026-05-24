package main
build linux
import "lib/datachannel"

class Datachannel : datachannel {
    func rtcInitLogger(level: int, cb: func(int, any char))
    func rtcPreload()
    func rtcCleanup()

    func rtcCreatePeerConnection(conf: any opaque) -> int
    func rtcDeletePeerConnection(pc: int) -> int

    func rtcSetStateChangeCallback(
        pc: int, cb: func(int, int, mut any void), ptr: mut any void) -> int
    func rtcSetGatheringStateChangeCallback(
        pc: int, cb: func(int, int, mut any void), ptr: mut any void) -> int
    func rtcSetLocalDescriptionCallback(
        pc: int, cb: func(int, any char, any char, mut any void), ptr: mut any void) -> int
    func rtcSetLocalCandidateCallback(
        pc: int, cb: func(int, any char, any char, mut any void), ptr: mut any void) -> int
    func rtcSetDataChannelCallback(
        pc: int, cb: func(int, int, mut any void), ptr: mut any void) -> int

    func rtcCreateDataChannel(pc: int, label: any char) -> int
    func rtcDeleteDataChannel(dc: int) -> int

    func rtcSetOpenCallback(
        id: int, cb: func(int, mut any void), ptr: mut any void) -> int
    func rtcSetClosedCallback(
        id: int, cb: func(int, mut any void), ptr: mut any void) -> int
    func rtcSetErrorCallback(
        id: int, cb: func(int, any char, mut any void), ptr: mut any void) -> int
    func rtcSetMessageCallback(
        id: int, cb: func(int, any char, int, mut any void), ptr: mut any void) -> int

    func rtcSendMessage(id: int, data: any char, size: int) -> int
    func rtcSetRemoteDescription(pc: int, sdp: any char, sdpType: any char) -> int
    func rtcAddRemoteCandidate(pc: int, cand: any char, mid: any char) -> int

    func printf(fmt: any char, ...) -> int
    func puts(s: any char) -> int
    func sleep(seconds: uint) -> uint
}

// ─── Constants ────────────────────────────────────────────────────────────────

let RTC_CONNECTING   = 0
let RTC_CONNECTED    = 1
let RTC_DISCONNECTED = 2
let RTC_FAILED       = 3
let RTC_CLOSED       = 4

let RTC_GATHERING_NEW        = 0
let RTC_GATHERING_INPROGRESS = 1
let RTC_GATHERING_COMPLETE   = 2

let RTC_LOG_WARNING = 3

// ─── Wrapper classes ──────────────────────────────────────────────────────────

class PeerConnection {
    var handle:        int
    var onStateFn:     func(int)?
    var onGatherFn:    func(int)?
    var onSdpFn:       func(any char, any char)?
    var onCandidateFn: func(any char, any char)?
    var onChannelFn:   func(int)?
}

class DataChannel {
    var handle:    int
    var onOpenFn:  func()?
    var onCloseFn: func()?
    var onErrorFn: func(any char)?
    var onMsgFn:   func(any char, int)?
}

func init(pc: mut PeerConnection) {
    pc.handle        = -1
    pc.onStateFn     = nil
    pc.onGatherFn    = nil
    pc.onSdpFn       = nil
    pc.onCandidateFn = nil
    pc.onChannelFn   = nil
}

func init(dc: mut DataChannel) {
    dc.handle    = -1
    dc.onOpenFn  = nil
    dc.onCloseFn = nil
    dc.onErrorFn = nil
    dc.onMsgFn   = nil
}

// ─── Package-level instances ──────────────────────────────────────────────────

var g_pc = PeerConnection().new()
var g_dc = DataChannel().new()

// ─── Global logger ────────────────────────────────────────────────────────────

func onRtcLog(level: int, msg: any char) {
    var dc = Datachannel()
    dc.printf("[rtc] %s\n".any(), msg)
}

// ─── C → Vertex trampolines ───────────────────────────────────────────────────

func _dcOpenTrampoline(id: int, ptr: mut any void) {
    if let f = g_dc.onOpenFn { f() }
}

func _dcCloseTrampoline(id: int, ptr: mut any void) {
    if let f = g_dc.onCloseFn { f() }
}

func _dcErrorTrampoline(id: int, err: any char, ptr: mut any void) {
    if let f = g_dc.onErrorFn { f(err) }
}

func _dcMsgTrampoline(id: int, msg: any char, size: int, ptr: mut any void) {
    if let f = g_dc.onMsgFn { f(msg, size) }
}

func _pcStateTrampoline(pc: int, state: int, ptr: mut any void) {
    if let f = g_pc.onStateFn { f(state) }
}

func _pcGatherTrampoline(pc: int, state: int, ptr: mut any void) {
    if let f = g_pc.onGatherFn { f(state) }
}

func _pcSdpTrampoline(pc: int, sdp: any char, sdpType: any char, ptr: mut any void) {
    if let f = g_pc.onSdpFn { f(sdp, sdpType) }
}

func _pcCandidateTrampoline(pc: int, cand: any char, mid: any char, ptr: mut any void) {
    if let f = g_pc.onCandidateFn { f(cand, mid) }
}

func _pcChannelTrampoline(pc: int, dc: int, ptr: mut any void) {
    if let f = g_pc.onChannelFn { f(dc) }
}

// ─── PeerConnection — .on*() event API ───────────────────────────────────────

func onState(pc: mut PeerConnection, handler: func(int)) {
    var dc = Datachannel()
    pc.onStateFn = handler
    dc.rtcSetStateChangeCallback(pc.handle, _pcStateTrampoline, nil)
}

func onGathering(pc: mut PeerConnection, handler: func(int)) {
    var dc = Datachannel()
    pc.onGatherFn = handler
    dc.rtcSetGatheringStateChangeCallback(pc.handle, _pcGatherTrampoline, nil)
}

func onLocalDescription(pc: mut PeerConnection, handler: func(any char, any char)) {
    var dc = Datachannel()
    pc.onSdpFn = handler
    dc.rtcSetLocalDescriptionCallback(pc.handle, _pcSdpTrampoline, nil)
}

func onLocalCandidate(pc: mut PeerConnection, handler: func(any char, any char)) {
    var dc = Datachannel()
    pc.onCandidateFn = handler
    dc.rtcSetLocalCandidateCallback(pc.handle, _pcCandidateTrampoline, nil)
}

func onDataChannel(pc: mut PeerConnection, handler: func(int)) {
    var dc = Datachannel()
    pc.onChannelFn = handler
    dc.rtcSetDataChannelCallback(pc.handle, _pcChannelTrampoline, nil)
}

// ─── DataChannel — .on*() event API + send ───────────────────────────────────

func onOpen(dc: mut DataChannel, handler: func()) {
    var rtc = Datachannel()
    dc.onOpenFn = handler
    rtc.rtcSetOpenCallback(dc.handle, _dcOpenTrampoline, nil)
}

func onClose(dc: mut DataChannel, handler: func()) {
    var rtc = Datachannel()
    dc.onCloseFn = handler
    rtc.rtcSetClosedCallback(dc.handle, _dcCloseTrampoline, nil)
}

func onError(dc: mut DataChannel, handler: func(any char)) {
    var rtc = Datachannel()
    dc.onErrorFn = handler
    rtc.rtcSetErrorCallback(dc.handle, _dcErrorTrampoline, nil)
}

func onMessage(dc: mut DataChannel, handler: func(any char, int)) {
    var rtc = Datachannel()
    dc.onMsgFn = handler
    rtc.rtcSetMessageCallback(dc.handle, _dcMsgTrampoline, nil)
}

func send(dc: DataChannel, data: string, size: int) -> int {
    var rtc = Datachannel()
    return rtc.rtcSendMessage(dc.handle, data.any(), size)
}

// ─── Entry point ─────────────────────────────────────────────────────────────

func main() -> int {
    var rtc = Datachannel()

    rtc.puts("=== Vertex WebRTC DataChannel demo ===".any())

    rtc.rtcInitLogger(RTC_LOG_WARNING, onRtcLog)
    rtc.rtcPreload()
    defer rtc.rtcCleanup()

    // ── Peer connection ───────────────────────────────────────────────────────

    g_pc.handle = rtc.rtcCreatePeerConnection(nil)
    if g_pc.handle < 0 {
        rtc.puts("fatal: rtcCreatePeerConnection failed".any())
        return 1
    }
    defer rtc.rtcDeletePeerConnection(g_pc.handle)

    g_pc.onState(func(state: int) {
        var r = Datachannel()
        switch state {
        case RTC_CONNECTING:
            r.puts("[pc] connecting...".any())
        case RTC_CONNECTED:
            r.puts("[pc] connected!".any())
        case RTC_DISCONNECTED:
            r.puts("[pc] disconnected".any())
        case RTC_FAILED:
            r.puts("[pc] connection failed".any())
        case RTC_CLOSED:
            r.puts("[pc] closed".any())
        default:
        }
    })

    g_pc.onGathering(func(state: int) {
        var r = Datachannel()
        switch state {
        case RTC_GATHERING_INPROGRESS:
            r.puts("[ice] gathering candidates...".any())
        case RTC_GATHERING_COMPLETE:
            r.puts("[ice] gathering complete — SDP ready to signal".any())
        default:
        }
    })

    g_pc.onLocalDescription(func(sdp: any char, sdpType: any char) {
        var r = Datachannel()
        r.printf("[sdp] local %s ready → forward to remote peer\n".any(), sdpType)
    })

    g_pc.onLocalCandidate(func(cand: any char, mid: any char) {
        var r = Datachannel()
        r.printf("[ice] candidate on mid=%s → forward to remote peer\n".any(), mid)
    })

    g_pc.onDataChannel(func(dcHandle: int) {
        var r = Datachannel()
        r.puts("[dc] incoming data channel from remote peer".any())
        g_dc.handle = dcHandle
        r.rtcSetOpenCallback(g_dc.handle,    _dcOpenTrampoline,  nil)
        r.rtcSetClosedCallback(g_dc.handle,  _dcCloseTrampoline, nil)
        r.rtcSetErrorCallback(g_dc.handle,   _dcErrorTrampoline, nil)
        r.rtcSetMessageCallback(g_dc.handle, _dcMsgTrampoline,   nil)
    })

    // ── Data channel ──────────────────────────────────────────────────────────

    g_dc.handle = rtc.rtcCreateDataChannel(g_pc.handle, "vertex-chat".any())
    if g_dc.handle < 0 {
        rtc.puts("fatal: rtcCreateDataChannel failed".any())
        return 1
    }
    defer rtc.rtcDeleteDataChannel(g_dc.handle)

    g_dc.onOpen(func() {
        var r = Datachannel()
        r.puts("[dc] channel open — sending greeting".any())
        g_dc.send(data: "Hello from Vertex over WebRTC!", size: 30)
    })

    g_dc.onMessage(func(msg: any char, size: int) {
        var r = Datachannel()
        r.printf("[dc] recv (%d bytes): %s\n".any(), size, msg)
        g_dc.send(data: "echo: got it", size: 12)
    })

    g_dc.onClose(func() {
        var r = Datachannel()
        r.puts("[dc] channel closed".any())
    })

    g_dc.onError(func(err: any char) {
        var r = Datachannel()
        r.printf("[dc] error: %s\n".any(), err)
    })

    // ── Run ───────────────────────────────────────────────────────────────────

    rtc.puts("[main] gathering ICE candidates (5 s)...".any())
    rtc.sleep(uint(5))

    rtc.puts("[main] (signaling exchange goes here)".any())

    rtc.sleep(uint(10))
    rtc.puts("[main] shutting down".any())
    return 0
}