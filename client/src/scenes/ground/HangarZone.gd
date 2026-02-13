extends Area3D

func _ready():
	# Visuals: Transparent Blue Cylinder
	var mesh_inst = MeshInstance3D.new()
	var cylinder = CylinderMesh.new()
	cylinder.height = 4.0
	cylinder.radius = 2.0
	mesh_inst.mesh = cylinder

	var mat = StandardMaterial3D.new()
	mat.albedo_color = Color(0, 0.5, 1.0, 0.3) # Blue Transparent
	mat.transparency = BaseMaterial3D.TRANSPARENCY_ALPHA
	mat.cull_mode = BaseMaterial3D.CULL_DISABLED
	mesh_inst.material_override = mat

	add_child(mesh_inst)

	# Collision Shape
	var shape = CollisionShape3D.new()
	var cyl_shape = CylinderShape3D.new()
	cyl_shape.height = 4.0
	cyl_shape.radius = 2.0
	shape.shape = cyl_shape
	add_child(shape)

	# Signal
	connect("body_entered", _on_body_entered)

func _on_body_entered(body):
	# Check if body is the player (basic check for MVP)
	if body.has_method("_send_movement"): # Duck typing check for PlayerController
		print("Hangar Zone Entered! Requesting Launch...")
		NetworkManager.send_udp_packet("REQUEST_LAUNCH", {})
		# Note: The server will respond with PACKET_TYPE_LAUNCH_GRANTED.
		# NetworkManager should handle that to switch scenes.
