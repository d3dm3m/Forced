extends Node3D

# Visualizes Projectiles based on Server Events
# Projectiles are just visual meshes interpolated between Start and End.

func _ready():
	NetworkManager.connect("packet_received", _on_packet_received)

func _on_packet_received(type: String, payload: Dictionary):
	if type == "PACKET_TYPE_PROJECTILE_SPAWN":
		_spawn_projectile(payload)

func _spawn_projectile(data: Dictionary):
	var start = Vector3(data.get("start_x"), 1.0, data.get("start_y"))
	var end = Vector3(data.get("end_x"), 1.0, data.get("end_y"))
	var speed = data.get("speed", 20.0)
	var outcome = data.get("outcome", "Hit")

	var mesh = MeshInstance3D.new()
	var cylinder = CylinderMesh.new()
	cylinder.height = 0.5
	cylinder.radius = 0.1
	mesh.mesh = cylinder

	# Rotate to face target
	mesh.rotation.x = PI / 2 # Lay flat

	add_child(mesh)
	mesh.global_position = start
	mesh.look_at(end)

	# Tween movement
	var dist = start.distance_to(end)
	var duration = dist / speed

	var tween = create_tween()
	tween.tween_property(mesh, "global_position", end, duration)

	if outcome == "Miss":
		# overshoot
		var miss_pos = end + (end - start).normalized() * 5.0
		tween.stop() # Reset
		dist = start.distance_to(miss_pos)
		duration = dist / speed
		tween = create_tween()
		tween.tween_property(mesh, "global_position", miss_pos, duration)

	tween.tween_callback(mesh.queue_free)
