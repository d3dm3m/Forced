extends Node3D

var light: SpotLight3D

func _ready():
	light = SpotLight3D.new()
	add_child(light)

	# Visual Configuration
	light.shadow_enabled = true
	light.spot_angle = 60.0
	light.light_energy = 5.0 # High contrast against dark world
	light.light_color = Color(0.9, 0.95, 1.0) # Cold industrial white

	# Optimization
	light.shadow_bias = 0.05

	# Atmosphere
	# We assume WorldEnvironment is handled by the scene, but we can enforce local darkness if needed.
	# For MVP, the light itself provides the "see" part.

func setup(range_val: float):
	if light:
		light.spot_range = range_val
