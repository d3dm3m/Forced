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

# Action State Machine
enum State { IDLE, WINDUP, BACKSWING }
var current_state: int = State.IDLE
var state_timer: float = 0.0
# Stats (Should be loaded from class)
var cast_point: float = 0.3
var backswing: float = 0.5

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

	# Fog of War (Vision Sensor)
	var fog_script = preload("res://src/vfx/FogOfWar.gd")
	var fog = fog_script.new()
	torso.add_child(fog)
	# Default range, ideally updated from Class Stats in Login
	fog.setup(20.0)
	fog.position.y = 0.5 # High on torso
	fog.rotation.x = -0.1 # Slight tilt down

	# Atmosphere (WorldEnvironment)
	# Force pitch black ambient to make Fog of War work
	var world_env = WorldEnvironment.new()
	var env = Environment.new()
	env.background_mode = Environment.BG_COLOR
	env.background_color = Color.BLACK
	env.ambient_light_source = Environment.AMBIENT_SOURCE_COLOR
	env.ambient_light_color = Color.BLACK
	env.ambient_light_energy = 0.0
	world_env.environment = env
	add_child(world_env)

	# Listen for Gear
	NetworkManager.connect("packet_received", _on_packet_received)

func _on_packet_received(type: String, payload: Dictionary):
	if type == "LOGIN_SUCCESS":
		var gear = payload.get("ground_gear", {})
		if gear.has("primary_weapon"):
			var weapon = gear["primary_weapon"]
			if weapon:
				print("PlayerController: Equipped Weapon: ", weapon.get("item_id", "Unknown"))
	elif type == "PACKET_TYPE_COMBAT_HIT":
		var dmg = payload.get("damage", 0.0)
		print("TOOK DAMAGE: ", dmg)
	elif type == "PACKET_TYPE_DEATH":
		print("YOU DIED. RESPAWNING...")

func _physics_process(delta):
	# State Machine Logic
	if current_state != State.IDLE:
		state_timer -= delta
		if state_timer <= 0:
			if current_state == State.WINDUP:
				_perform_attack()
			elif current_state == State.BACKSWING:
				current_state = State.IDLE
				print("Ready.")

	# Tank / RTS Control
	var input_dir = Input.get_vector("ui_left", "ui_right", "ui_up", "ui_down")
	var target_dir = Vector3(input_dir.x, 0, input_dir.y).normalized()

	# Action Canceling / Orb Walking
	if target_dir:
		if current_state == State.WINDUP:
			current_state = State.IDLE
			print("Attack Cancelled!")
		elif current_state == State.BACKSWING:
			current_state = State.IDLE
			print("Backswing Cancelled (Orb Walk)!")

	if target_dir:
		# 1. Rotate Body towards Target
		# Using a fixed turn rate (could be class based)
		var current_transform = global_transform
		var target_pos = global_position + target_dir
		var new_transform = current_transform.looking_at(target_pos, Vector3.UP)

		# Rotate towards target quaternion
		var current_quat = current_transform.basis.get_rotation_quaternion()
		var target_quat = new_transform.basis.get_rotation_quaternion()
		var next_quat = current_quat.slerp(target_quat, 5.0 * delta) # 5.0 = Turn Rate

		global_transform.basis = Basis(next_quat)

		# 2. Check Alignment
		# Only move if facing roughly the right way (~15 deg)
		var forward = -global_transform.basis.z
		var dot = forward.dot(target_dir)

		if dot > 0.9:
			# Aligned enough to move
			velocity.x = move_toward(velocity.x, target_dir.x * speed, acceleration * delta)
			velocity.z = move_toward(velocity.z, target_dir.z * speed, acceleration * delta)
		else:
			# Not aligned, just braking
			velocity.x = move_toward(velocity.x, 0, friction * delta)
			velocity.z = move_toward(velocity.z, 0, friction * delta)
	else:
		# No Input, Stop
		velocity.x = move_toward(velocity.x, 0, friction * delta)
		velocity.z = move_toward(velocity.z, 0, friction * delta)

	move_and_slide()

	# Turret Slew (Torso Tracking - Independent of Body)
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
	if current_state != State.IDLE:
		return

	if current_target_id == "":
		print("No Target!")
		return

	# Start Windup
	current_state = State.WINDUP
	state_timer = cast_point
	print("Winding up...")

func _perform_attack():
	print("Fired at ", current_target_id)

	# Visuals: Muzzle Flash (Placeholder)

	# Network
	var payload = { "target_id": current_target_id }
	NetworkManager.send_udp_packet("PACKET_TYPE_GROUND_ATTACK", payload)

	# Enter Backswing
	current_state = State.BACKSWING
	state_timer = backswing

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
