extends CharacterBody3D

# Newtonian Physics Parameters
@export var acceleration_force: float = 50.0
@export var rotation_speed: float = 2.0
@export var max_speed: float = 100.0
@export var braking_factor: float = 2.0

var network_tick_timer: float = 0.0
const NETWORK_TICK_RATE: float = 0.1 # 10Hz

func _physics_process(delta):
	handle_input(delta)
	move_and_slide()
	handle_networking(delta)

func handle_input(delta):
	var rot_dir = 0
	if Input.is_action_pressed("ui_left") or Input.is_key_pressed(KEY_A):
		rot_dir += 1
	if Input.is_action_pressed("ui_right") or Input.is_key_pressed(KEY_D):
		rot_dir -= 1

	rotate_y(rot_dir * rotation_speed * delta)

	var forward_dir = -global_transform.basis.z
	var input_accel = Vector3.ZERO

	if Input.is_action_pressed("ui_up") or Input.is_key_pressed(KEY_W):
		input_accel += forward_dir * acceleration_force * delta
	if Input.is_action_pressed("ui_down") or Input.is_key_pressed(KEY_S):
		input_accel -= forward_dir * acceleration_force * delta

	if Input.is_action_pressed("ui_select") or Input.is_key_pressed(KEY_SPACE):
		var current_speed = velocity.length()
		if current_speed > 0:
			var brake_force = velocity.normalized() * -1 * acceleration_force * braking_factor * delta
			if brake_force.length() > current_speed:
				velocity = Vector3.ZERO
			else:
				velocity += brake_force

	velocity += input_accel

	if velocity.length() > max_speed:
		velocity = velocity.normalized() * max_speed

func handle_networking(delta):
	# Use NetworkManager instead of direct WS peer
	network_tick_timer += delta
	if network_tick_timer >= NETWORK_TICK_RATE:
		network_tick_timer = 0
		send_state()

func send_state():
	var payload = {
		"pos_x": global_position.x,
		"pos_y": global_position.y,
		"pos_z": global_position.z,
		"vel_x": velocity.x,
		"vel_y": velocity.y,
		"vel_z": velocity.z,
		"rot_y": rotation.y
	}
	NetworkManager.send_packet("PACKET_TYPE_SPACE_STATE", payload)
