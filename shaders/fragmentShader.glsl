#version 410

in vec2 texCoord;
out vec4 frag_colour;
uniform sampler2D ourTexture;

void main() {
	vec4 texColour = texture(ourTexture, texCoord);
	if (texColour.a < 0.1)
		discard;
	frag_colour = texColour;
}