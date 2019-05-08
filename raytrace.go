package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
)

func clip(x float64) float64 {
	if x > 1.0 {
		return 1.0
	}
	if x < 0.0 {
		return 0.0
	}
	return x
}

func render() {
	var width, height int = 1024, 768
	var framebuffer = make([]vec3, width*height)

	ivory := phongMaterial{vec3{0.4, 0.4, 0.3}, albedoCoeff{0.6, 0.3, 0, 0}, 50, 1.0}
	redRubber := phongMaterial{vec3{0.3, 0.1, 0.1}, albedoCoeff{0.9, 0.1, 0, 0}, 10, 1.0}
	mirror := phongMaterial{vec3{1.0, 1.0, 1.0}, albedoCoeff{0.0, 10.0, 0.8, 0}, 1425., 1.0}
	glass := phongMaterial{vec3{0.6, 0.7, 0.8}, albedoCoeff{0.0, 0.5, 0.1, 0.8}, 125., 1.5}
	white := simpleMaterial{vec3{0.3, 0.3, 0.3}}
	yellow := simpleMaterial{vec3{0.3, 0.7 * 0.3, 0.3 * 0.3}}

	board := plane{vec3{-10, -4, -10}, vec3{20, 0, 0}, vec3{0, 0, -20}}

	scene := []rayCatcher{
		rayObject{sphere{vec3{-3.0, 0.0, -16}, 2}, simpleMaterialSelector{ivory}},
		rayObject{sphere{vec3{-1.0, -1.5, -12}, 2}, simpleMaterialSelector{glass}},
		rayObject{sphere{vec3{1.5, -0.5, -18}, 3}, simpleMaterialSelector{redRubber}},
		rayObject{sphere{vec3{7.0, 5.0, -18}, 4}, simpleMaterialSelector{mirror}},
		rayObject{board, checkerboardMaterialSelector{white, yellow, board}},
	}
	lights := []light{
		light{vec3{-20, 20, 20}, 1.5},
		light{vec3{30, 50, -25}, 1.8},
		light{vec3{30, 20, 30}, 1.7},
	}

	fov := 60 * math.Pi / 180
	orig := vec3{0, 0, 0}

	for j := 0; j < height; j++ {
		for i := 0; i < width; i++ {
			w := float64(width)
			h := float64(height)
			dirX := (float64(i) + 0.5) - w/2.
			dirY := -(float64(j) + 0.5) + h/2. // this flips the image at the same time
			dirZ := -h / (2. * math.Tan(fov/2.))

			dir := vec3{dirX, dirY, dirZ}.normalize()

			sceneIntersect := func(orig, dir vec3) (dist float64, hitPoint, normal vec3, mat material) {
				dist = math.MaxFloat64
				for _, obj := range scene {
					intersects, objDist := obj.rayIntersect(orig, dir)
					if intersects && objDist < dist {
						hitPoint = orig.add(dir.mul(objDist))
						mat = obj.mat(hitPoint)
						if mat != nil {
							dist = objDist
							normal = obj.normal(hitPoint)
						}
					}
				}
				return
			}

			var castRay func(orig, dir vec3, depth int) (color vec3)
			castRay = func(orig, dir vec3, depth int) (color vec3) {
				dist, hitPoint, normal, mat := sceneIntersect(orig, dir)

				if depth > 4 || dist >= 1000 {
					color = vec3{0.2, 0.7, 0.8}
					return
				}

				reflect := func(i, n vec3) vec3 {
					return i.sub(n.mul(2.0 * i.dot(n)))
				}

				// Snell's law
				refract := func(i, n vec3, refractiveIndex float64) vec3 {
					cosi := -math.Max(-1, math.Min(1, i.dot(n)))
					etai, etat := 1.0, refractiveIndex
					normal := n
					if cosi < 0 { // if the ray comes from the inside the object, swap the air and the media
						cosi = -cosi
						etai, etat = etat, etai
						normal = normal.inv()
					}
					eta := etai / etat
					k := 1 - eta*eta*(1-cosi*cosi)
					if k < 0 {
						// k<0 = total reflection, no ray to refract. I refract it anyways, this has no physical meaning
						return vec3{1, 0, 0}
					}
					return i.mul(eta).add(normal.mul(eta*cosi - math.Sqrt(k)))
				}

				origin := func(dir, normal, point vec3) vec3 {
					// offset the original point to avoid occlusion by the object itself
					if dir.dot(normal) < 0 {
						return point.sub(normal.mul(1e-3))
					}
					return point.add(normal.mul(1e-3))
				}

				reflectDir := reflect(dir, normal) //.normalize()
				refractDir := refract(dir, normal, mat.refractiveIndex()).normalize()
				reflectOrig := origin(reflectDir, normal, hitPoint)
				refractOrig := origin(refractDir, normal, hitPoint)
				reflectColor := castRay(reflectOrig, reflectDir, depth+1)
				refractColor := castRay(refractOrig, refractDir, depth+1)

				diffuseLightIntensity, specularLightIntensity := 0.0, 0.0
				for _, lt := range lights {
					lightDir := lt.position.sub(hitPoint).normalize()
					lightDistance := lt.position.sub(hitPoint).norm()

					shadowOrig := origin(lightDir, normal, hitPoint)

					shadowDist, shadowPoint, _, _ := sceneIntersect(shadowOrig, lightDir)
					if shadowDist < 1000 && shadowPoint.sub(shadowOrig).norm() < lightDistance {
						continue
					}

					diffuseLightIntensity += lt.intensivity * math.Max(0, lightDir.dot(normal))

					specularLightIntensity += math.Pow(math.Max(0, reflect(lightDir.inv(), normal).inv().dot(dir)), mat.specularExponent()) * lt.intensivity
				}

				color = mat.diffuseColor().mul(diffuseLightIntensity * mat.albedo().diffusive).
					add(vec3{1, 1, 1}.mul(specularLightIntensity * mat.albedo().specular)).
					add(reflectColor.mul(mat.albedo().refelective)).
					add(refractColor.mul(mat.albedo().refractive))

				return
			}

			framebuffer[i+j*width] = castRay(orig, dir, 0)
		}
	}

	imgRect := image.Rect(0, 0, width, height)
	img := image.NewRGBA(imgRect)

	for j := 0; j < height; j++ {
		for i := 0; i < width; i++ {
			v := framebuffer[i+j*width]
			max := math.Max(v.x, math.Max(v.y, v.z))
			if max > 1 {
				v = v.mul(1 / max)
			}
			img.Set(i, j, color.RGBA{uint8(clip(v.x) * 255), uint8(clip(v.y) * 255), uint8(clip(v.z) * 255), 255})
		}
	}

	out, err := os.Create("./output.png")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	png.Encode(out, img)
	fmt.Println("Image generated successfully.")
}

func main() {
	render()
}
