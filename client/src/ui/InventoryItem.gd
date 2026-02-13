extends PanelContainer

var item_index: int = -1
var item_id: String = ""

func setup(index: int, id: String, count: int):
	item_index = index
	item_id = id

	# Clear existing children
	for child in get_children():
		child.queue_free()

	# Add Label
	var label = Label.new()
	label.text = str(count) + "x\n" + id
	add_child(label)

	# Visual Style (Basic)
	custom_minimum_size = Vector2(50, 50)

func _get_drag_data(at_position):
	if item_index == -1:
		return null

	var data = {
		"item_index": item_index,
		"item_id": item_id
	}

	# Create Preview
	var preview = Label.new()
	preview.text = item_id
	preview.modulate = Color(1, 1, 1, 0.8)
	set_drag_preview(preview)

	return data
