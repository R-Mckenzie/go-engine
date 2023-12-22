#version 410

in vec2 texCoord;
out vec4 frag_colour;

uniform sampler2D u_texture;

void main()
{
	frag_colour = texture(u_texture, texCoord);
}  
