#version 410

in vec2 texCoord;
out vec4 frag_colour;

uniform sampler2D u_texture;
uniform float exposure;

void main()
{
	vec4 frag = texture(u_texture, texCoord);
	frag_colour = vec4(vec3(1.0) - exp(-frag.rgb * exposure), frag.a);
}  
