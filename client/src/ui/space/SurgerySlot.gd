extends PanelContainer

signal organ_dropped(inventory_index: int, slot_type: String, slot_index: int)

var slot_type: String = ""
var slot_index: int = 0

func setup(type: String, index: int, label_text: String):
	slot_type = type
	slot_index = index

	# Basic Visuals
	custom_minimum_size = Vector2(200, 40)

	var label = Label.new()
	label.text = label_text + " (Drop Here)"
	add_child(label)

func _can_drop_data(at_position, data):
	return data is Dictionary and data.has("item_id")

func _drop_data(at_position, data):
	if _can_drop_data(at_position, data):
		emit_signal("organ_dropped", data["item_index"], slot_type, slot_index)
