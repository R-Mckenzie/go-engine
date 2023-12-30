#version 410

//attributes from vertex shader
in vec2 texCoord;

out vec4 frag_colour;

uniform sampler2D u_texture;   //diffuse map
uniform vec4 u_colour;

void main() {
	vec4 diffuseColour = texture(u_texture, texCoord) * u_colour;
	frag_colour = vec4(diffuseColour);
}
