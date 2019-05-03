package main

type albedoCoeff struct {
	diffusive   float64
	specular    float64
	refelective float64
	refractive  float64
}

type material interface {
	diffuseColor() vec3
	albedo() albedoCoeff
	specularExponent() float64
	refractiveIndex() float64
}

type simpleMaterial struct {
	color vec3
}

func (m simpleMaterial) diffuseColor() vec3 {
	return m.color
}

func (m simpleMaterial) albedo() albedoCoeff {
	return albedoCoeff{1, 0, 0, 0}
}

func (m simpleMaterial) specularExponent() float64 {
	return 0
}

func (m simpleMaterial) refractiveIndex() float64 {
	return 0
}

type phongMaterial struct {
	color        vec3
	albedoCf     albedoCoeff
	specular     float64
	refractIndex float64
}

func (m phongMaterial) albedo() albedoCoeff {
	return m.albedoCf
}

func (m phongMaterial) diffuseColor() vec3 {
	return m.color
}

func (m phongMaterial) specularExponent() float64 {
	return m.specular
}

func (m phongMaterial) refractiveIndex() float64 {
	return m.refractIndex
}
