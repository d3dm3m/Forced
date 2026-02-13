extends Node

# Controls the Sanity Distortion Shader based on network events
@export var shader_rect_path: NodePath
var shader_mat: ShaderMaterial

func _ready():
	NetworkManager.connect("sanity_changed", _on_sanity_changed)

	if shader_rect_path:
		var rect = get_node(shader_rect_path)
		if rect and rect.material is ShaderMaterial:
			shader_mat = rect.material

func _on_sanity_changed(new_level: float):
	if !shader_mat:
		return

	# Sanity is 0-100 usually, shader might expect 0.0-1.0
	var normalized = clamp(new_level / 100.0, 0.0, 1.0)
	shader_mat.set_shader_parameter("sanity_level", normalized)
