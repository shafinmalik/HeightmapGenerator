package main

import (
	"image"
	"image/color"
	"image/png"
	"math"
	"math/rand"
	"os"
	"time"
)

var SET_SIZE int = 256
var OCTAVE float32 = 0.45

type heightmap struct {
	filename string
	width    int
	height   int
	cMap     *image.RGBA
}

func heightmap_init(name string, width int, height int) *heightmap {
	frame := heightmap{filename: name}
	frame.width = width
	frame.height = height
	topleft := image.Point{0, 0}
	lowRight := image.Point{width, height}

	frame.cMap = image.NewRGBA(image.Rectangle{topleft, lowRight})
	return &frame
}

func populate_map(h *heightmap, p []int) {

	for i := 0; i < h.width; i++ {
		for j := 0; j < h.height; j++ {
			horizontal := float32(j) / float32(h.width)
			vertical := float32(i) / float32(h.height)
			pixel_value := noise(horizontal, vertical, OCTAVE, p)

			adj_pixel := uint8(math.Floor(float64(255 * pixel_value)))
			// fmt.Println(adj_pixel)
			value := color.RGBA{adj_pixel, adj_pixel, adj_pixel, 0xff}
			// fmt.Println(value)

			h.cMap.Set(i, j, value)
		}
	}
}

func randomize_permutations(size int) []int {
	set := make([]int, size)
	max := 255
	min := 0
	for i := 0; i < 256; i++ {
		set[i] = rand.Intn(max-min) + min
	}

	return set
}

func noise(x_unit float32, y_unit float32, z_unit float32, p []int) float32 {
	X_corner := int(math.Floor(float64(x_unit))) & 255
	Y_corner := int(math.Floor(float64(y_unit))) & 255
	Z_corner := int(math.Floor(float64(z_unit))) & 255

	x_pos := x_unit - float32(math.Floor(float64(x_unit)))
	y_pos := x_unit - float32(math.Floor(float64(y_unit)))
	z_pos := x_unit - float32(math.Floor(float64(z_unit)))

	// Fade
	u := fade(x_pos)
	v := fade(y_pos)
	w := fade(z_pos)

	// Coordinates
	A := (p[X_corner] + Y_corner)
	AA := (p[A] + Z_corner)
	AB := (p[A+1] + Z_corner)
	B := (p[X_corner+1] + Y_corner)
	BA := (p[B] + Z_corner)
	BB := (p[B+1] + Z_corner)

	grad_A0 := gradient(p[AB+1], x_pos, y_pos-1, z_pos-1)
	grad_A1 := gradient(p[BB+1], x_pos-1, y_pos-1, z_pos-1)
	lerp_A := LERP(u, grad_A0, grad_A1)

	grad_B0 := gradient(p[BA+1], x_pos-1, y_pos, z_pos-1)
	grad_B1 := gradient(p[AA+1], x_pos, y_pos, z_pos-1)
	lerp_B := LERP(u, grad_B1, grad_B0)

	lerp_C := LERP(v, lerp_B, lerp_A) // Goes into final result

	grad_D0 := gradient(p[BB], x_pos-1, y_pos-1, z_pos)
	grad_D1 := gradient(p[AB], x_pos, y_pos-1, z_pos)
	lerp_D := LERP(u, grad_D1, grad_D0)

	grad_E0 := gradient(p[BA], x_pos-1, y_pos, z_pos)
	grad_E1 := gradient(p[AA], x_pos, y_pos, z_pos)
	lerp_E := LERP(u, grad_E1, grad_E0)

	lerp_F := LERP(v, lerp_E, lerp_D)

	result := LERP(w, lerp_F, lerp_C)
	noise := (result + 1.0) / 2.0

	return noise
}

// Based on: Perlin, K. (2002, July). Improving noise.
func fade(input float32) float32 {
	term := float32(math.Pow(float64(input), float64(3)))
	effector := float32(input*(input*6-15) + 10)

	return term * effector
}

func LERP(t float32, a float32, b float32) float32 {
	return a + t*(b-a)
}

// Based on: Perlin, K. (2002, July). Improving noise.
func gradient(hash int, x float32, y float32, z float32) float32 {
	h := hash & 15
	u := x
	if h < 8 {
		u = y
	}

	v := u
	if h < 4 {
		v = y
	} else if h == 12 || h == 14 {
		v = x
	} else {
		v = z
	}

	lower := -u
	if (h & 1) == 0 {
		lower = u
	}
	upper := -v
	if (h & 2) == 0 {
		upper = v
	}

	return (lower + upper)
}

func main() {
	rand.Seed(int64(time.Now().Nanosecond()))
	p := randomize_permutations(SET_SIZE)

	map_image := heightmap_init("img.png", 300, 300)
	populate_map(map_image, p)
	f, _ := os.Create(map_image.filename)
	png.Encode(f, map_image.cMap)
}
