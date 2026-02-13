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
	var item_script = preload("res://src/ui/InventoryItem.gd")
	var idx = 0
	for item in items:
		var slot = item_script.new()
		slot.setup(idx, item.item_id, int(item.count))
		grid.add_child(slot)
		idx += 1
