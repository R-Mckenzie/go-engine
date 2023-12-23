#version 410

//attributes from vertex shader
in vec2 vTexCoord;
in mat4 vMV;

#define  MAX_LIGHTS 15

out vec4 frag_colour;

uniform sampler2D u_texture;   //diffuse map
uniform sampler2D u_normals;   //normal map

uniform bool UseNormals;

//values used for shading algorithm...
uniform vec4 AmbientColor;    //ambient RGBA -- alpha is intensity 

uniform vec3 LightPos[MAX_LIGHTS];        //light position, normalized
uniform vec4 LightColor[MAX_LIGHTS];      //light RGBA -- alpha is intensity
uniform vec3 Falloff[MAX_LIGHTS];         //attenuation coefficients

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

	vec3 N = normalize(NormalMap * 2.0 - 1.0);

    vec3 diffuse = vec3(0.0);
	//pre-multiply ambient color with intensity

	for (int i = 0; i < MAX_LIGHTS; i++) {
		if (LightColor[i].a != 0) {
			//The delta position of light
			vec3 LightDir = vec3(LightPos[i].xy - gl_FragCoord.xy, LightPos[i].z);

			//Determine distance (used for attenuation) BEFORE we normalize our LightDir
			float D = length(LightDir);

			//normalize our vectors
			vec3 L = normalize(LightDir);
			
			//Pre-multiply light color with intensity
			//Then perform "N dot L" to determine our diffuse term

			vec3 ReducedFalloff = Falloff[i] / 10000;
			//calculate attenuation
			float Attenuation = 1.0 / (ReducedFalloff.x + (ReducedFalloff.y * D) + (ReducedFalloff.z * D * D));
			
			//the calculation which brings it all together
			diffuse += Attenuation * (LightColor[i].rgb * LightColor[i].a) * max(dot(N, L), 0.0);
		}
	}

    //pre-multiply ambient color with intensity
	vec3 Ambient = AmbientColor.rgb * AmbientColor.a;

		//the calculation which brings it all together
		vec3 intensity = min(vec3(1.0), Ambient + diffuse); // don't remember if min is critical, but I think it might be to avoid shifting the hue when multiple lights add up to something very bright.
		vec3 finalColor = DiffuseColor.rgb * intensity;
		frag_colour = vec4(finalColor, DiffuseColor.a);
}
