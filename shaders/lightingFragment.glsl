#version 410

uniform vec4 u_lightColour;
uniform float u_lightRadius;

out vec4 frag_colour;

void main() {
	frag_colour = vec4(0.0, 0.0, 0.0, 1.0);
}
