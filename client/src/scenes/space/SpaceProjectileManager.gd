extends Node3D

# SpaceProjectileManager
# Visualizes 3D projectiles with trails and behaviors (Missile, Linear, Instant)

func _ready():
	NetworkManager.connect("packet_received", _on_packet_received)

func _on_packet_received(type: String, payload: Dictionary):
	if type == "PACKET_TYPE_PROJECTILE_SPAWN":
		_spawn_projectile(payload)

func _spawn_projectile(data: Dictionary):
	var behavior = data.get("behavior", "linear")

	if behavior == "instant":
		_draw_beam(data)
	else:
		_spawn_moving_projectile(data)

func _spawn_moving_projectile(data: Dictionary):
	var start = Vector3(data.get("start_x"), data.get("start_z"), data.get("start_y")) # Z-Y swap check?
	# Godot Y is up. Protocol Z is likely depth/up depending on Space convention.
	# Assuming Space Core uses X/Z plane and Y is up? Or full 3D?
	# Let's assume direct mapping: X->X, Y->Y, Z->Z.
	start = Vector3(data.get("start_x"), data.get("start_y"), data.get("start_z"))
	var end = Vector3(data.get("end_x"), data.get("end_y"), data.get("end_z"))

	var speed = data.get("speed", 50.0)
	var behavior = data.get("behavior", "linear")

	var mesh_inst = MeshInstance3D.new()

	if behavior == "missile":
		var sphere = SphereMesh.new()
		sphere.radius = 0.5
		mesh_inst.mesh = sphere
		# Add Trail (Placeholder)
		# var trail = GPUParticles3D.new() ...
	else:
		# Linear (Railgun)
		var cylinder = CylinderMesh.new()
		cylinder.height = 2.0
		cylinder.radius = 0.1
		mesh_inst.mesh = cylinder
		mesh_inst.rotation.x = PI/2 # Align forward

	add_child(mesh_inst)
	mesh_inst.global_position = start
	mesh_inst.look_at(end)

	var dist = start.distance_to(end)
	var duration = dist / speed

	if duration <= 0: duration = 0.1

	var tween = create_tween()

	if behavior == "missile":
		# Quadratic Curve interpolation could go here for "launch" feel
		# For MVP, linear move to target
		tween.tween_property(mesh_inst, "global_position", end, duration)
	else:
		tween.tween_property(mesh_inst, "global_position", end, duration)

	tween.tween_callback(mesh_inst.queue_free)

func _draw_beam(data: Dictionary):
	# Draw Line3D logic
	pass
