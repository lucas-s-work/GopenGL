package opengl

import (
	"Gopengl/util"
	"fmt"
	"image"
	"image/draw"
	_ "image/png" //needed to load png file
	"os"

	"github.com/go-gl/gl/v4.1-core/gl"
)

type Texture struct {
	id     uint32
	width  int
	height int
}

func LoadTexture(file string) *Texture {
	imgFile, err := os.Open(util.RelativePath(file))
	if err != nil {
		panic(fmt.Errorf("texture %q not found on disk: %v", file, err))
	}

	// Get imagine data
	img, _, err := image.Decode(imgFile)
	if err != nil {
		panic(fmt.Errorf("Image load error, error: %v", err))
	}

	bounds := img.Bounds()
	rgba := image.NewRGBA(bounds)
	if rgba.Stride != rgba.Rect.Size().X*4 {
		panic(fmt.Errorf("unsupported stride"))
	}
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)

	var texture uint32
	gl.GenTextures(1, &texture)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, texture)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGBA,
		int32(rgba.Rect.Size().X),
		int32(rgba.Rect.Size().Y),
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		gl.Ptr(rgba.Pix))

	return &Texture{
		texture,
		bounds.Max.X,
		bounds.Max.Y,
	}
}

/*
Texture usage methods
*/

// TODO add support for multiple texture units if available
// this works good enough for simple 2d games with only a few sprite sheets
func (t *Texture) Use() {
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, t.id)
}

// NormCoords ... normalize pixture texture coordinates
func (t *Texture) PixToTex(texs []float32) []float32 {
	normedTexs := make([]float32, len(texs))
	even := false
	width := float32(t.width)
	height := float32(t.height)

	for i, coord := range texs {
		even = !even

		if even == true {
			normedTexs[i] = coord / width
			continue
		}

		normedTexs[i] = coord / height
	}

	return normedTexs
}
