#version 410

//attributes from vertex shader
in vec2 vTexCoord;
in mat4 vMV;

out vec4 frag_colour;

uniform sampler2D u_texture;   //diffuse map
uniform sampler2D u_normals;   //normal map

uniform bool UseNormals;

//values used for shading algorithm...
uniform vec2 Resolution;      //resolution of screen
uniform vec3 LightPos;        //light position, normalized
uniform vec4 LightColor;      //light RGBA -- alpha is intensity
uniform vec4 AmbientColor;    //ambient RGBA -- alpha is intensity 
uniform vec3 Falloff;         //attenuation coefficients

void main() {

	vec4 DiffuseColor = texture(u_texture, vTexCoord);
	if (DiffuseColor.a < 0.1)
		discard;

	vec3 NormalMap;
	if (UseNormals) {
		NormalMap = texture(u_normals, vTexCoord).rgb;
	} else {
		NormalMap = vec3(0.5, 0.5, 1.0);
	}

	//The delta position of light
	vec3 LightDir = vec3(LightPos.xy - gl_FragCoord.xy, LightPos.z);

	//Determine distance (used for attenuation) BEFORE we normalize our LightDir
	float D = length(LightDir);

	//normalize our vectors
	vec3 N = normalize(NormalMap * 2.0 - 1.0);
	vec3 L = normalize(LightDir);
	
	//Pre-multiply light color with intensity
	//Then perform "N dot L" to determine our diffuse term
	vec3 Diffuse = (LightColor.rgb * LightColor.a) * max(dot(N, L), 0.0);

	//pre-multiply ambient color with intensity
	vec3 Ambient = AmbientColor.rgb * AmbientColor.a;
	
	vec3 ReducedFalloff = Falloff / vec3(Resolution, 1000);
	//calculate attenuation
	float Attenuation = 1.0 / ( ReducedFalloff.x + (ReducedFalloff.y*D) + (ReducedFalloff.z*D*D) );
	
	//the calculation which brings it all together
	vec3 Intensity = Ambient + Diffuse * Attenuation;
	vec3 FinalColor = DiffuseColor.rgb * Intensity;

	frag_colour = vec4(FinalColor, DiffuseColor.a);
}
