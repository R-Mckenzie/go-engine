#version 410

in vec2 texCoord;
out vec4 frag_colour;

uniform sampler2D screenTexture;

void main()
{
	frag_colour = texture(screenTexture, texCoord);
}  
