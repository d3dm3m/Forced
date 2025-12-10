extends Control

# Simple UI Grid for Inventory
# Expects a GridContainer child named "Grid"

func _ready():
	if not has_node("Grid"):
		var grid = GridContainer.new()
		grid.name = "Grid"
		grid.columns = 5
		add_child(grid)

	# Start listening for inventory updates (mock connection)
	# NetworkManager.connect("inventory_updated", _on_inventory_updated)

func _on_inventory_updated(items: Array):
	var grid = $Grid
	# Clear existing children
	for child in grid.get_children():
		child.queue_free()

	# Populate new items
	for item in items:
		var slot = PanelContainer.new()
		var label = Label.new()
		# item is expected to be { "item_id": "...", "count": ... }
		label.text = str(item.count) + "x\n" + item.item_id
		slot.add_child(label)
		grid.add_child(slot)
