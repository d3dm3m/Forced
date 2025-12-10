extends "res://src/ui/space/DraggableWindow.gd"

# Surgery Window: Shows Slots and allows installation
# For Sprint 15 MVP, we just show buttons to install into specific slots from "Inventory Index 0"
# to demonstrate the loop. A full Inventory Picker is complex.

func _ready():
	name = "SurgeryWindow"
	size = Vector2(250, 300)
	position = Vector2(600, 50)

	var title = Label.new()
	title.text = "SURGERY BAY"
	title.position = Vector2(10, 5)
	add_child(title)

	var container = VBoxContainer.new()
	container.position = Vector2(10, 30)
	container.size = Vector2(230, 260)
	add_child(container)

	# Slot: High 1
	add_slot_row(container, "High Slot 1", "high_slots", 0)
	# Slot: Mid 1
	add_slot_row(container, "Mid Slot 1", "mid_slots", 0)
	# Slot: Low 1
	add_slot_row(container, "Low Slot 1", "low_slots", 0)

	# Note: This UI assumes you want to install the FIRST item in your inventory.
	# Real UI needs drag-and-drop from Inventory Window.
	var note = Label.new()
	note.text = "NOTE: Installs Item #0 from Inventory"
	note.modulate = Color(0.7, 0.7, 0.7)
	container.add_child(note)

func add_slot_row(container, label_text, slot_type, slot_index):
	var row = HBoxContainer.new()

	var label = Label.new()
	label.text = label_text
	label.size_flags_horizontal = Control.SIZE_EXPAND_FILL
	row.add_child(label)

	var btn = Button.new()
	btn.text = "GRAFT (Idx 0)"
	btn.connect("pressed", func(): _on_graft_pressed(slot_type, slot_index))
	row.add_child(btn)

	container.add_child(row)

func _on_graft_pressed(slot_type, slot_index):
	print("Grafting Inventory[0] into ", slot_type, "[", slot_index, "]")
	# Hardcoded to index 0 for MVP testing
	var payload = {
		"inventory_index": 0,
		"slot_type": slot_type,
		"slot_index": slot_index
	}
	NetworkManager.send_packet("PACKET_TYPE_GRAFT_ORGAN", payload)
