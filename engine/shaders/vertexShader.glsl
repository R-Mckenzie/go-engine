#version 410

layout (location = 0) in vec3 pos;
layout (location = 1) in vec2 i_texCoord;

out vec2 texCoord;

uniform mat4 u_model;
uniform mat4 u_view;
uniform mat4 u_projection;

void main() {
	texCoord = i_texCoord;
	gl_Position =  u_projection * u_view * u_model * vec4(pos, 1.0);
}
