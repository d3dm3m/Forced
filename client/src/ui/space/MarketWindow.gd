extends "res://src/ui/space/DraggableWindow.gd"

# Hardcoded Market for Sprint 15
var items = [
	{ "id": "synthetic_heart", "name": "Synthetic Heart", "price": 500 },
	{ "id": "ocular_implant", "name": "Ocular Implant", "price": 250 },
	{ "id": "void_heart", "name": "Void Heart (Illegal)", "price": 1000 },
	{ "id": "mining_laser", "name": "Mining Laser", "price": 100 },
]

func _ready():
	name = "MarketWindow"
	size = Vector2(250, 300)
	position = Vector2(300, 50)

	# Title
	var title = Label.new()
	title.text = "THE MEAT MARKET"
	title.position = Vector2(10, 5)
	add_child(title)

	var container = VBoxContainer.new()
	container.position = Vector2(10, 30)
	container.size = Vector2(230, 260)
	add_child(container)

	for item in items:
		var row = HBoxContainer.new()

		var label = Label.new()
		label.text = item.name + "\n" + str(item.price) + " Solium"
		label.size_flags_horizontal = Control.SIZE_EXPAND_FILL
		row.add_child(label)

		var buy_btn = Button.new()
		buy_btn.text = "BUY"
		buy_btn.connect("pressed", func(): _on_buy_pressed(item.id))
		row.add_child(buy_btn)

		container.add_child(row)
		container.add_child(HSeparator.new())

func _on_buy_pressed(item_id: String):
	print("Buying: ", item_id)
	var payload = { "item_id": item_id }
	NetworkManager.send_packet("PACKET_TYPE_BUY_ITEM", payload)
