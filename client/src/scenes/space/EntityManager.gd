extends Node3D

# Simple Entity Manager to visualize other ships/objects
# Mapped by Entity ID -> Node3D
var entities: Dictionary = {}

func _ready():
	NetworkManager.connect("packet_received", _on_packet_received)

func _on_packet_received(type: String, payload: Dictionary):
	match type:
		"PACKET_TYPE_SPACE_STATE":
			_handle_space_state(payload)
		"PACKET_TYPE_ENTITY_DESTROYED":
			_handle_entity_destroyed(payload)

func _handle_space_state(payload: Dictionary):
	# Expected Payload: { "entities": [ { "id": "...", "pos_x": ... }, ... ] }
	# Note: The current server implementation might just send individual updates or a list.
	# The mechanics.go/main.go implementation of Space Core mainly echoes or validates,
	# but `ENTITY_TRACKING` was mentioned in context.
	# For Sprint 14, we assume we might receive a list of nearby entities.

	var entity_list = payload.get("entities", [])
	var seen_ids = []

	for ent_data in entity_list:
		var id = ent_data.get("id")
		seen_ids.append(id)

		if entities.has(id):
			_update_entity(id, ent_data)
		else:
			_spawn_entity(id, ent_data)

func _spawn_entity(id: String, data: Dictionary):
	var mesh_inst = MeshInstance3D.new()
	var sphere = SphereMesh.new()
	sphere.radius = 1.0
	mesh_inst.mesh = sphere
	add_child(mesh_inst)

	entities[id] = mesh_inst
	_update_entity(id, data)

func _update_entity(id: String, data: Dictionary):
	var node = entities[id]
	var x = data.get("pos_x", 0.0)
	var y = data.get("pos_y", 0.0)
	var z = data.get("pos_z", 0.0)
	node.global_position = Vector3(x, y, z)

func _handle_entity_destroyed(payload: Dictionary):
	var id = payload.get("id")
	if entities.has(id):
		entities[id].queue_free()
		entities.erase(id)
