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
func (t triangle) rayIntersect(orig, dir vec3) (bool, float64) {
	edge1 := t.p2.sub(t.p1)
	edge2 := t.p3.sub(t.p1)
	pvec := dir.cross(edge2)
	det := edge1.dot(pvec)
	if det < 1e-5 {
		return false, 0
	}

	tvec := orig.sub(t.p1)
	u := tvec.dot(pvec)
	if u < 0 || u > det {
		return false, 0
	}

	qvec := tvec.cross(edge1)
	v := dir.dot(qvec)
	if v < 0 || u+v > det {
		return false, 0
	}

	tnear := edge2.dot(qvec) * (1 / det)
	return tnear > 1e-5, tnear
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

func (p plane) rayIntersect(orig, dir vec3) (bool, float64) {
	p2 := p.p.add(p.v1)
	p3 := p.p.add(p.v2)
	intersect, distance := triangle{p.p, p2, p3, nil}.rayIntersect(orig, dir)
	if intersect {
		return true, distance
	}
	p4 := p3.add(p.v1)
	return triangle{p2, p4, p3, nil}.rayIntersect(orig, dir)
}

func (p plane) normal(hit vec3) vec3 {
	return p.v1.cross(p.v2).normalize()
}

func (p plane) mat(_ vec3) material {
	return p.material
}
