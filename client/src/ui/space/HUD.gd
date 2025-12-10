extends Control

# Preload the Draggable Window scene or script
var window_script = preload("res://src/ui/space/DraggableWindow.gd")

# HUD Elements
var health_bar: ProgressBar
var bio_load_bar: ProgressBar
var warning_label: Label
var purge_button: Button

func _ready():
	# Setup Network Connection
	NetworkManager.connect("ship_stats_updated", _on_ship_stats_updated)

	# Create HUD Layout
	setup_status_bars()

	# Example: Create a Radar Window (Existing logic)
	create_radar_window()

func setup_status_bars():
	var container = VBoxContainer.new()
	container.position = Vector2(20, 20)
	container.size = Vector2(200, 100)
	add_child(container)

	# Health Bar
	var hp_label = Label.new()
	hp_label.text = "HULL INTEGRITY"
	container.add_child(hp_label)

	health_bar = ProgressBar.new()
	health_bar.max_value = 1000
	health_bar.value = 1000
	# Set color to Green roughly via modulation or theme (omitted for brevity)
	container.add_child(health_bar)

	# Spacer
	container.add_child(HSeparator.new())

	# Bio-Load Bar
	var bio_label = Label.new()
	bio_label.text = "BIO-LOAD"
	container.add_child(bio_label)

	bio_load_bar = ProgressBar.new()
	bio_load_bar.max_value = 100
	bio_load_bar.value = 0
	container.add_child(bio_load_bar)

	# Warning Label
	warning_label = Label.new()
	warning_label.text = "REJECTION IMMINENT"
	warning_label.modulate = Color(1, 0, 0)
	warning_label.visible = false
	container.add_child(warning_label)

	# Purge Button
	purge_button = Button.new()
	purge_button.text = "EMERGENCY PURGE"
	purge_button.modulate = Color(1, 0, 0) # Red button
	purge_button.connect("pressed", _on_purge_pressed)
	container.add_child(purge_button)

	# Spacer
	container.add_child(HSeparator.new())

	# Toggle Buttons
	var market_btn = Button.new()
	market_btn.text = "OPEN MARKET"
	market_btn.connect("pressed", _on_toggle_market)
	container.add_child(market_btn)

	var surgery_btn = Button.new()
	surgery_btn.text = "OPEN SURGERY"
	surgery_btn.connect("pressed", _on_toggle_surgery)
	container.add_child(surgery_btn)

	# Load Windows
	var market = preload("res://src/ui/space/MarketWindow.gd").new()
	market.visible = false
	add_child(market)

	var surgery = preload("res://src/ui/space/SurgeryWindow.gd").new()
	surgery.visible = false
	add_child(surgery)

func _on_toggle_market():
	var win = get_node("MarketWindow")
	if win: win.visible = !win.visible

func _on_toggle_surgery():
	var win = get_node("SurgeryWindow")
	if win: win.visible = !win.visible

func _on_ship_stats_updated(stats: Dictionary):
	# Update Health
	var cur_hp = stats.get("current_health", 1000.0)
	var max_hp = stats.get("max_health", 1000.0)

	health_bar.max_value = max_hp
	health_bar.value = cur_hp

	# Visual Feedback for Health
	if cur_hp < (max_hp * 0.2):
		health_bar.modulate = Color(1, 0, 0) # Red Critical
		# Panic Signal to SanityController could be emitted here or centralized
	else:
		health_bar.modulate = Color(1, 1, 1) # Normal

	# Update Bio-Load
	var cur_load = stats.get("bio_load", 0)
	var max_load = stats.get("bio_capacity", 50)

	bio_load_bar.max_value = max_load
	bio_load_bar.value = cur_load

	# Color Coding Bio-Load
	var load_percent = float(cur_load) / float(max_load) if max_load > 0 else 0.0

	if load_percent > 1.0:
		bio_load_bar.modulate = Color(1, 0, 0) # Red (Overload)
		warning_label.visible = true
		# Blink effect could be done via tween or process
	elif load_percent > 0.8:
		bio_load_bar.modulate = Color(1, 0.6, 0) # Orange (Warning)
		warning_label.visible = false
	else:
		bio_load_bar.modulate = Color(0, 1, 0) # Green (Safe)
		warning_label.visible = false

func _on_purge_pressed():
	print("Purge Initiated!")
	# Logic to send PURGE packet would go here.
	# NetworkManager.send_packet("PACKET_TYPE_PURGE_REQUEST", {})

# ... Existing Radar Logic ...
func create_radar_window():
	var radar_win = Panel.new()
	radar_win.set_script(window_script)
	radar_win.name = "RadarWindow"
	radar_win.size = Vector2(200, 200)
	radar_win.position = Vector2(50, 150) # Moved down to avoid bars

	# Add some content to look like a radar
	var label = Label.new()
	label.text = "D-SCAN RADAR"
	label.position = Vector2(10, 5)
	radar_win.add_child(label)

	var radar_display = ColorRect.new()
	radar_display.color = Color(0, 0.2, 0, 0.8)
	radar_display.position = Vector2(10, 30)
	radar_display.size = Vector2(180, 160)
	radar_win.add_child(radar_display)

	# Add a center blip (Self)
	var self_blip = ColorRect.new()
	self_blip.color = Color.WHITE
	self_blip.size = Vector2(4, 4)
	self_blip.position = Vector2(90 - 2, 80 - 2)
	radar_display.add_child(self_blip)

	add_child(radar_win)
