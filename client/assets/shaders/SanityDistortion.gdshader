shader_type canvas_item;

uniform float sanity_level : hint_range(0.0, 1.0) = 1.0;
uniform sampler2D screen_texture : hint_screen_texture, repeat_disable, filter_nearest;

void fragment() {
    // 0.0 Sanity = Max Distortion, 1.0 = Clean
    float distortion_strength = 1.0 - sanity_level;

    // Chromatic Aberration
    // Shift channels based on strength
    float r_offset = distortion_strength * 0.01;
    float b_offset = distortion_strength * -0.01;

    vec4 red_channel = texture(screen_texture, SCREEN_UV + vec2(r_offset, 0.0));
    vec4 green_channel = texture(screen_texture, SCREEN_UV);
    vec4 blue_channel = texture(screen_texture, SCREEN_UV + vec2(b_offset, 0.0));

    COLOR = vec4(red_channel.r, green_channel.g, blue_channel.b, 1.0);

    // Breathing / Pulsing Effect
    // Only at low sanity (< 0.5)
    if (sanity_level < 0.5) {
        float time = TIME * 2.0;
        float pulse = sin(time) * 0.02 * (0.5 - sanity_level); // stronger as sanity drops

        // Simple UV warp for breathing
        vec2 centered_uv = SCREEN_UV - 0.5;
        float dist = length(centered_uv);
        vec2 warped_uv = SCREEN_UV + (centered_uv / dist) * pulse * dist;

        // Re-sample with warp
        // (For simplicity in this snippet, effectively combining them involves complex sampling,
        // so we'll just mix the color tint here or assume the UV warp applies to the texture lookup above)
        // Let's just apply a slight dark vignette pulse for "Breathing"
        float vignette = smoothstep(0.8, 0.2, length(centered_uv) + pulse);
        COLOR.rgb *= vignette;
    }
}
