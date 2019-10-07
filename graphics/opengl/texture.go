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

var (
	textureUnits = []uint32{
		gl.TEXTURE0,
		gl.TEXTURE1,
		gl.TEXTURE2,
		gl.TEXTURE3,
		gl.TEXTURE4,
		gl.TEXTURE5,
		gl.TEXTURE6,
		gl.TEXTURE7,
		gl.TEXTURE8,
		gl.TEXTURE9,
		gl.TEXTURE10,
		gl.TEXTURE11,
		gl.TEXTURE12,
		gl.TEXTURE13,
		gl.TEXTURE14,
	}
	storedTextures       []*Texture
	currentTextureUnitId uint32 = 0
	textureUnitUsed      uint32 = 0
)

const textureIdsBeforeChange = 32

type Texture struct {
	id          uint32
	width       int
	height      int
	file        string
	textureUnit uint32
}

/**
Loads a texture, does not reload it if already created
*/
func LoadTexture(file string) *Texture {
	// Load existing texture
	existingTex := FindTex(file)

	if existingTex != nil {
		return existingTex
	}

	// Create new texture if it doesn't exist
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
	gl.ActiveTexture(currentTextureUnit())
	gl.GenTextures(1, &texture)
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

	textureObj := &Texture{
		texture,
		bounds.Max.X,
		bounds.Max.Y,
		file,
		currentTextureUnitId,
	}

	gl.BindTexture(gl.TEXTURE_2D, 0)

	//Add texture to texture store
	storedTextures = append(storedTextures, textureObj)

	return textureObj
}

func FindTex(file string) *Texture {
	for _, tex := range storedTextures {
		if tex.file == file {
			return tex
		}
	}

	return nil
}

/*
Texture usage methods
*/

func (t *Texture) Use() {
	gl.ActiveTexture(t.textureUnit)
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

// Util

func currentTextureUnit() uint32 {
	if textureUnitUsed > textureIdsBeforeChange {
		currentTextureUnitId++
		textureUnitUsed = 0

		if currentTextureUnitId >= uint32(len(textureUnits)) {
			panic("No free texture units")
		}
	}

	textureUnitUsed++

	return textureUnits[currentTextureUnitId]
}
