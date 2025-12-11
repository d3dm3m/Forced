extends Node3D

@export var target_path: NodePath
@export var smooth_speed: float = 5.0
@export var offset: Vector3 = Vector3(20, 30, 20) # High and steep (RTS)

var target: Node3D
var camera: Camera3D

func _ready():
	if target_path:
		target = get_node(target_path)

	# Create Camera if not present as child (for script attachment usage)
	# But typically this script goes on a Node3D holding the camera
	camera = $Camera3D
	if not camera:
		camera = Camera3D.new()
		add_child(camera)
		camera.name = "Camera3D"

	# Set Orthogonal Projection for 2.5D look
	camera.projection = Camera3D.PROJECTION_ORTHOGONAL
	camera.size = 20.0 # Adjust size for zoom level

	# Set isometric rotation (Look at center from offset)
	# Typically isometric is 45 degrees Y, and roughly 30-35 degrees X down
	# But looking from (20, 20, 20) to (0,0,0) does this automatically if we use look_at

	# Initial position
	if target:
		global_position = target.global_position + offset
		look_at(target.global_position)

func _physics_process(delta):
	if !target:
		return

	var desired_position = target.global_position + offset
	var smoothed_position = global_position.lerp(desired_position, smooth_speed * delta)
	global_position = smoothed_position

	# Ensure we always look at the target, or a fixed point relative to the camera
	# For strict isometric, rotation usually shouldn't change, just position.
	# So we just move the rig.
