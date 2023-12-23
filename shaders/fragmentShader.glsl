#version 410

#define  MAX_LIGHTS 15

//attributes from vertex shader
in vec2 texCoord;

out vec4 frag_colour;

uniform sampler2D u_texture;   //diffuse map
uniform sampler2D u_normals;   //normal map

uniform bool useNormals;

//values used for shading algorithm...
uniform vec4 ambientLight;    //ambient RGBA -- alpha is intensity 
uniform vec3 lightPos[MAX_LIGHTS];        //light position, normalized
uniform vec4 lightColour[MAX_LIGHTS];      //light RGBA -- alpha is intensity
uniform vec3 falloff[MAX_LIGHTS];         //attenuation coefficients

void main() {
	vec4 diffuseColour = texture(u_texture, vTexCoord);
	if (diffuseColour.a < 0.1)
		discard;

	vec3 normalMap;
	if (useNormals) {
		normalMap = texture(u_normals, vTexCoord).rgb;
	} else {
		normalMap = vec3(0.5, 0.5, 1.0);
	}

	vec3 N = normalize(normalMap * 2.0 - 1.0);

    vec3 diffuse = vec3(0.0);
	//pre-multiply ambient color with intensity

	for (int i = 0; i < MAX_LIGHTS; i++) {
		if (lightColour[i].a != 0) {
			//The delta position of light
			vec3 lightDir = vec3(lightPos[i].xy - gl_FragCoord.xy, lightPos[i].z);

			//Determine distance (used for attenuation) BEFORE we normalize our lightDir
			float D = length(lightDir);

			//normalize our vectors
			vec3 L = normalize(lightDir);
			
			//Pre-multiply light color with intensity
			//Then perform "N dot L" to determine our diffuse term

			vec3 reducedFalloff = falloff[i] / 1000;
			//calculate attenuation
			float attenuation = 1.0 / (reducedFalloff.x + (reducedFalloff.y * D) + (reducedFalloff.z * D * D));
			
			//the calculation which brings it all together
			diffuse += attenuation * (lightColour[i].rgb * lightColour[i].a) * max(dot(N, L), 0.0);
		}
	}

    //pre-multiply ambient color with intensity
	vec3 ambient = ambientLight.rgb * ambientLight.a;
	//the calculation which brings it all together
	vec3 intensity = min(vec3(1.0), ambient + diffuse); // don't remember if min is critical, but I think it might be to avoid shifting the hue when multiple lights add up to something very bright.
	vec3 finalColor = diffuseColour.rgb * intensity;
	frag_colour = vec4(finalColor, diffuseColour.a);
}
