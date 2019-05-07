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

type triangle struct {
	p1, p2, p3 vec3
	material   material
}

// Moller-Trumbore algorithm
// https://www.scratchapixel.com/lessons/3d-basic-rendering/ray-tracing-rendering-a-triangle/moller-trumbore-ray-triangle-intersection
func rayTriangleIntersect(orig, dir, p, v1, v2 vec3) (bool, float64) {
	pvec := dir.cross(v2)
	det := v1.dot(pvec)
	if det < 1e-5 {
		return false, 0
	}

	tvec := orig.sub(p)
	u := tvec.dot(pvec)
	if u < 0 || u > det {
		return false, 0
	}

	qvec := tvec.cross(v1)
	v := dir.dot(qvec)
	if v < 0 || u+v > det {
		return false, 0
	}

	tnear := v2.dot(qvec) * (1 / det)
	return tnear > 1e-5, tnear
}

func (t triangle) rayIntersect(orig, dir vec3) (bool, float64) {
	return rayTriangleIntersect(orig, dir, t.p1, t.p2.sub(t.p1), t.p3.sub(t.p1))
}

func (t triangle) normal(_ vec3) vec3 {
	return t.p1.sub(t.p2).cross(t.p1.sub(t.p3)).normalize()
}

func (t triangle) mat(_ vec3) material {
	return t.material
}

type plane struct {
	p, v1, v2 vec3
	material  material
}

/*
// https://stackoverflow.com/questions/21114796/3d-ray-quad-intersection-test-in-java
func (p plane) rayIntersect(orig, dir vec3) (bool, float64) {
	n := p.v1.cross(p.v2).normalize()
	nd := n.dot(dir)
	if nd > -1e-6 && nd < 1e-6 {
		return false, 0
	}

	t := -n.dot(orig.sub(p.p)) / nd
	m := orig.add(dir.mul(t))

	mp := m.sub(p.p)

	u := mp.dot(p.v1)
	if u < 0 || u > p.v1.dot(p.v1) {
		return false, 0
	}
	v := mp.dot(p.v2)
	if v < 0 || v > p.v2.dot(p.v2) {
		return false, 0
	}

	return true, t
}
*/

func (p plane) rayIntersect(orig, dir vec3) (bool, float64) {
	intersect, distance := rayTriangleIntersect(orig, dir, p.p, p.v1, p.v2)
	if intersect {
		return true, distance
	}
	p4 := p.p.add(p.v2).add(p.v1)
	return rayTriangleIntersect(orig, dir, p4, p.v1.inv(), p.v2.inv())
}

func (p plane) normal(hit vec3) vec3 {
	return p.v1.cross(p.v2).normalize()
}

func (p plane) mat(_ vec3) material {
	return p.material
}
