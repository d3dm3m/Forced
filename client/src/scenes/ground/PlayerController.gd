extends CharacterBody3D

# Settings
@export var speed: float = 10.0
@export var acceleration: float = 50.0
@export var friction: float = 60.0

# Network
var update_timer: float = 0.0
const UPDATE_RATE: float = 0.05 # 20Hz

# Combat
var current_target_id: String = ""
var camera: Camera3D

# Visuals
var torso: Node3D
var slew_rate: float = 2.0 # Radians/sec, dynamic based on class

func _ready():
	# Temporary: Spawn Hangar Zone for testing
	var hangar = preload("res://src/scenes/ground/HangarZone.gd").new()
	hangar.position = Vector3(10, 0, 10)
	get_parent().call_deferred("add_child", hangar)

	# Find Camera (Assuming CameraRig is sibling or child, for now grab viewport camera)
	camera = get_viewport().get_camera_3d()

	# Create Torso (Visual Representation for Turret Slew)
	torso = MeshInstance3D.new()
	var box = BoxMesh.new()
	box.size = Vector3(0.5, 1.5, 0.5)
	torso.mesh = box
	torso.position.y = 0.75
	add_child(torso)

	# Listen for Gear
	NetworkManager.connect("packet_received", _on_packet_received)

func _on_packet_received(type: String, payload: Dictionary):
	if type == "LOGIN_SUCCESS":
		var gear = payload.get("ground_gear", {})
		if gear.has("primary_weapon"):
			var weapon = gear["primary_weapon"]
			if weapon:
				print("PlayerController: Equipped Weapon: ", weapon.get("item_id", "Unknown"))

func _physics_process(delta):
	# Movement (Inertia)
	var input_dir = Input.get_vector("ui_left", "ui_right", "ui_up", "ui_down")
	var direction = Vector3(input_dir.x, 0, input_dir.y).normalized()

	# Note: move_toward handles linear acceleration/friction automatically
	if direction:
		velocity.x = move_toward(velocity.x, direction.x * speed, acceleration * delta)
		velocity.z = move_toward(velocity.z, direction.z * speed, acceleration * delta)
	else:
		velocity.x = move_toward(velocity.x, 0, friction * delta)
		velocity.z = move_toward(velocity.z, 0, friction * delta)

	move_and_slide()

	# Turret Slew (Torso Tracking)
	_handle_turret_slew(delta)

	# Network Update
	update_timer += delta
	if update_timer >= UPDATE_RATE:
		update_timer = 0
		_send_movement()

func _handle_turret_slew(delta):
	if !camera: return

	var mouse_pos = get_viewport().get_mouse_position()
	var from = camera.project_ray_origin(mouse_pos)
	var to = from + camera.project_ray_normal(mouse_pos) * 1000.0

	# Raycast to ground plane (Y=0)
	# Simple plane intersection math since physics raycast needs collision
	var t = -from.y / (to.y - from.y)
	if t >= 0:
		var target_point = from + (to - from) * t
		target_point.y = torso.global_position.y # Look level

		# Calculate angle difference
		var current_quat = torso.global_transform.basis.get_rotation_quaternion()
		var target_transform = torso.global_transform.looking_at(target_point, Vector3.UP)
		var target_quat = target_transform.basis.get_rotation_quaternion()

		# Rotate towards
		var new_quat = current_quat.slerp(target_quat, slew_rate * delta)
		torso.global_transform.basis = Basis(new_quat)

		# Firing Solution Check
		var forward = -torso.global_transform.basis.z
		var to_target = (target_point - torso.global_position).normalized()
		var dot = forward.dot(to_target)

		if dot > 0.99:
			# Aligned
			pass # Green Reticle logic would go here
		else:
			pass # Red Reticle logic

func _input(event):
	# Tab Targeting
	if event.is_action_pressed("ui_focus_next"): # Tab
		_cycle_target()

	# Mouse Click Targeting
	if event is InputEventMouseButton and event.pressed and event.button_index == MOUSE_BUTTON_LEFT:
		pass # Raycast handled in _handle_turret_slew implicitly for looking, selection separate

	# F1 Self Target (Placeholder for Party)
	if event is InputEventKey and event.pressed and event.keycode == KEY_F1:
		current_target_id = "" # Clear or Set to Self ID if known
		print("Target Self/None")

	# Fire
	if event.is_action_pressed("ui_accept"): # Space or Enter, map "Fire" later
		_fire_weapon()

func _cycle_target():
	# Access GroundEntityManager. Assuming it's a sibling or singleton.
	# For MVP, let's assume we can find it in the Scene Tree.
	var manager = get_node_or_null("../GroundEntityManager")
	if manager:
		current_target_id = manager.get_next_target(current_target_id)
		print("Target Locked: ", current_target_id)
		_update_reticle()

func _raycast_target(mouse_pos):
	if !camera:
		camera = get_viewport().get_camera_3d()
		if !camera: return

	var from = camera.project_ray_origin(mouse_pos)
	var to = from + camera.project_ray_normal(mouse_pos) * 1000.0

	var space_state = get_world_3d().direct_space_state
	var query = PhysicsRayQueryParameters3D.create(from, to)
	query.collision_mask = 2 # Layer 2 (Entities)

	var result = space_state.intersect_ray(query)
	if result:
		var collider = result.collider
		if collider.has_meta("entity_id"):
			current_target_id = collider.get_meta("entity_id")
			print("Click Target: ", current_target_id)
			_update_reticle()

func _fire_weapon():
	if current_target_id == "":
		print("No Target!")
		return

	print("Firing at ", current_target_id)

	# Visuals: Muzzle Flash (Placeholder)

	# Network
	var payload = { "target_id": current_target_id }
	NetworkManager.send_udp_packet("PACKET_TYPE_GROUND_ATTACK", payload)

func _update_reticle():
	# Visual feedback for target (Simple print for now or highlight shader)
	# In real imp, we would move a Sprite3D to the target's position.
	pass

func _send_movement():
	var payload = {
		"x": global_position.x,
		"y": global_position.z, # Mapping 3D Z to 2D Y
		"vx": velocity.x,
		"vy": velocity.z,
		"ts": Time.get_unix_time_from_system() * 1000
	}
	NetworkManager.send_udp_packet("PACKET_TYPE_GROUND_MOVEMENT", payload)
