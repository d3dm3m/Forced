extends Node3D

# Visualizes other players on ground
var entities: Dictionary = {} # ID -> Node3D

func _ready():
	NetworkManager.connect("packet_received", _on_packet_received)

func _on_packet_received(type: String, payload: Dictionary):
	if type == "PACKET_TYPE_GROUND_STATE":
		_update_entities(payload.get("entities", []))

func _update_entities(list: Array):
	for data in list:
		var id = data.get("id")

		if entities.has(id):
			var node = entities[id]
			var x = data.get("x", 0.0)
			var y = data.get("y", 0.0)
			node.global_position = Vector3(x, 0, y)
		else:
			_spawn_entity(id, data)

func _spawn_entity(id: String, data: Dictionary):
	var mesh_inst = MeshInstance3D.new()
	var capsule = CapsuleMesh.new()
	capsule.height = 2.0
	mesh_inst.mesh = capsule
	mesh_inst.position.y = 1.0 # Offset up
	add_child(mesh_inst)

	# Add Collider for Raycast Targeting
	var static_body = StaticBody3D.new()
	var shape = CollisionShape3D.new()
	var cap_shape = CapsuleShape3D.new()
	cap_shape.height = 2.0
	shape.shape = cap_shape
	static_body.add_child(shape)

	# Metadata for identification
	static_body.set_meta("entity_id", id)

	# Offset collider to match mesh
	static_body.position.y = 1.0
	add_child(static_body)

	var x = data.get("x", 0.0)
	var y = data.get("y", 0.0)
	# mesh_inst is a child, but we need to move the parent node if we wrapped it?
	# Ah, I added mesh_inst as child of self (GroundEntityManager).
	# I should probably wrap them in a Node3D "Entity" so I move one thing.

	# Refactor for cleaner movement
	mesh_inst.queue_free()
	static_body.queue_free()

	var wrapper = Node3D.new()
	add_child(wrapper)
	wrapper.global_position = Vector3(x, 0, y)

	# Re-add visuals to wrapper
	var vis = MeshInstance3D.new()
	vis.mesh = capsule
	vis.position.y = 1.0
	wrapper.add_child(vis)

	# Re-add collider to wrapper
	var body = StaticBody3D.new()
	body.add_child(shape.duplicate()) # Shape resource is shared
	body.set_meta("entity_id", id)
	body.position.y = 1.0
	body.collision_layer = 2 # Enemy/Entity Layer
	wrapper.add_child(body)

	entities[id] = wrapper

func get_next_target(current_id: String) -> String:
	var keys = entities.keys()
	if keys.is_empty():
		return ""

	var index = keys.find(current_id)
	if index == -1 or index == keys.size() - 1:
		return keys[0]

	return keys[index + 1]
