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
	create_slot(container, "High Slot 1", "high_slots", 0)
	# Slot: Mid 1
	create_slot(container, "Mid Slot 1", "mid_slots", 0)
	# Slot: Low 1
	create_slot(container, "Low Slot 1", "low_slots", 0)

	var note = Label.new()
	note.text = "Drag Items from Inventory"
	note.modulate = Color(0.7, 0.7, 0.7)
	container.add_child(note)

func create_slot(container, label_text, slot_type, slot_index):
	var slot_script = preload("res://src/ui/space/SurgerySlot.gd")
	var slot = slot_script.new()
	slot.setup(slot_type, slot_index, label_text)
	slot.connect("organ_dropped", _on_organ_dropped)
	container.add_child(slot)

func _on_organ_dropped(inventory_index, slot_type, slot_index):
	print("Grafting Inventory[", inventory_index, "] into ", slot_type, "[", slot_index, "]")
	var payload = {
		"inventory_index": inventory_index,
		"slot_type": slot_type,
		"slot_index": slot_index
	}
	NetworkManager.send_packet("PACKET_TYPE_GRAFT_ORGAN", payload)
