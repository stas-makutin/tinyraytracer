package main

import "math"

type rayCatcher interface {
	rayIntersect(orig, dir vec3) (bool, float64)
	normal(hit vec3) vec3
	mat(hit vec3) material
}

type sphere struct {
	center   vec3
	radius   float64
	material material
}

func (s sphere) rayIntersect(orig, dir vec3) (bool, float64) {
	l := s.center.sub(orig)
	tca := l.dot(dir)
	d2 := l.dot(l) - tca*tca
	r2 := s.radius * s.radius
	if d2 > r2 {
		return false, 0
	}
	thc := math.Sqrt(r2 - d2)
	t0 := tca - thc
	t1 := tca + thc
	if t0 < 0 {
		if t1 < 0 {
			return false, 0
		}
		return true, t1
	}
	return true, t0
}

func (s sphere) normal(hit vec3) vec3 {
	return hit.sub(s.center).normalize()
}

func (s sphere) mat(_ vec3) material {
	return s.material
}

type plane struct {
	p1, p2, p3 vec3
	material   material
}

func (p plane) rayIntersect(orig, dir vec3) (bool, float64) {
	return false, 0
}

func (p plane) normal(_ vec3) vec3 {
	return p.p1.sub(p.p2).cross(p.p1.sub(p.p3)).normalize()
}

func (p plane) mat(_ vec3) material {
	return p.material
}
