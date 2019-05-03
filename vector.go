package main

import "math"

type vec3 struct {
	x, y, z float64
}

func newVec3(x, y, z float64) (v vec3) {
	v.x = x
	v.y = y
	v.z = z
	return
}

func (v1 vec3) add(v2 vec3) (v vec3) {
	v.x = v1.x + v2.x
	v.y = v1.y + v2.y
	v.z = v1.z + v2.z
	return
}

func (v1 vec3) sub(v2 vec3) (v vec3) {
	v.x = v1.x - v2.x
	v.y = v1.y - v2.y
	v.z = v1.z - v2.z
	return
}

func (v1 vec3) inv() (v vec3) {
	v.x = -v1.x
	v.y = -v1.y
	v.z = -v1.z
	return
}

func (v1 vec3) mul(scalar float64) (v vec3) {
	v.x = v1.x * scalar
	v.y = v1.y * scalar
	v.z = v1.z * scalar
	return
}

func (v1 vec3) dot(v2 vec3) float64 {
	return v1.x*v2.x + v1.y*v2.y + v1.z*v2.z
}

func (v1 vec3) cross(v2 vec3) (v vec3) {
	v.x = v1.y*v2.z - v1.z*v2.y
	v.y = v1.z*v2.x - v1.x*v2.z
	v.z = v1.x*v2.y - v1.y*v2.x
	return
}

func (v1 vec3) norm() float64 {
	return math.Sqrt(v1.x*v1.x + v1.y*v1.y + v1.z*v1.z)
}

func (v1 vec3) normalize() (v vec3) {
	n := v1.norm()
	if n != 0 {
		v.x = v1.x / n
		v.y = v1.y / n
		v.z = v1.z / n
	} else {
		v.x, v.y, v.z = 0, 0, 0
	}
	return
}
