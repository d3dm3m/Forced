extends Node

# Signals
signal connected_to_server()
signal disconnected_from_server()
signal packet_received(type: String, payload: Dictionary)
signal sanity_changed(new_level: float)
signal ship_stats_updated(stats: Dictionary)

# Connection State
var ws_peer: WebSocketPeer = WebSocketPeer.new()
var udp_peer: PacketPeerUDP = PacketPeerUDP.new()
var udp_host: String = "127.0.0.1"
var udp_port: int = 5000

var is_connected_to_space: bool = false
var space_host: String = "ws://localhost:8080/ws"

func _ready():
	# UDP Setup
	udp_peer.connect_to_host(udp_host, udp_port)
	# For prototype, we might start disconnected or auto-connect.
	pass

func connect_to_space(token: String):
	print("NetworkManager: Connecting to Space Core at ", space_host)
	var err = ws_peer.connect_to_url(space_host)
	if err != OK:
		print("NetworkManager: Failed to connect to Space Core")
		return

	is_connected_to_space = true
	# We need to wait for connection to be open before sending token.
	# This is handled in _process.
	# For simplicity in this script, we store the token to send later.
	_pending_token = token

var _pending_token: String = ""
var _handshake_sent: bool = false

func _send_login_packet():
	var pkt = {
		"type": "LOGIN_WITH_TOKEN",
		"payload": {
			"token": _pending_token
		}
	}
	send_packet(pkt.type, pkt.payload)

func send_packet(type: String, payload: Dictionary):
	if ws_peer.get_ready_state() == WebSocketPeer.STATE_OPEN:
		var pkt = {
			"type": type,
			"payload": payload
		}
		ws_peer.send_text(JSON.stringify(pkt))

func _handle_packet(data: Dictionary):
	var type = data.get("type", "")
	var payload = data.get("payload", {})

	emit_signal("packet_received", type, payload)

	match type:
		"LOGIN_SUCCESS":
			print("NetworkManager: Login Success")
			emit_signal("connected_to_server")
			emit_signal("packet_received", type, payload) # Ensure logic sees this too
		"PACKET_TYPE_SANITY_UPDATE":
			var val = payload.get("sanity", 100.0)
			emit_signal("sanity_changed", val)
		"PACKET_TYPE_SHIP_STATS":
			emit_signal("ship_stats_updated", payload)
		"PACKET_TYPE_GROUND_STATE":
			emit_signal("packet_received", type, payload)

func send_udp_packet(type: String, payload: Dictionary):
	var pkt = {
		"type": type,
		"payload": payload
	}
	var packet_bytes = JSON.stringify(pkt).to_utf8_buffer()
	udp_peer.put_packet(packet_bytes)

func _process(delta):
	# WebSocket Polling
	ws_peer.poll()
	var state = ws_peer.get_ready_state()

	if state == WebSocketPeer.STATE_OPEN:
		if !_handshake_sent and _pending_token != "":
			_send_login_packet()
			_handshake_sent = true

		while ws_peer.get_available_packet_count() > 0:
			var pkt = ws_peer.get_packet()
			var txt = pkt.get_string_from_utf8()
			var json = JSON.parse_string(txt)
			if json:
				_handle_packet(json)
	elif state == WebSocketPeer.STATE_CLOSED:
		if is_connected_to_space:
			print("NetworkManager: Disconnected from Space Core")
			is_connected_to_space = false
			emit_signal("disconnected_from_server")

	# UDP Polling
	while udp_peer.get_available_packet_count() > 0:
		var pkt = udp_peer.get_packet()
		var txt = pkt.get_string_from_utf8()
		var json = JSON.parse_string(txt)
		if json:
			_handle_packet(json)
