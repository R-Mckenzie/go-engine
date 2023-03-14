package engine

import (
	"image"
	"os"

	_ "image/png"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type Image struct {
	id     uint32
	width  float32
	height float32
}

func LoadImage(filepath string) (Image, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return Image{}, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return Image{}, err
	}

	w := img.Bounds().Max.X
	h := img.Bounds().Max.Y

	pixels := make([]byte, w*h*4)
	bIndex := 0
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			pixels[bIndex] = byte(r / 256)
			bIndex++
			pixels[bIndex] = byte(g / 256)
			bIndex++
			pixels[bIndex] = byte(b / 256)
			bIndex++
			pixels[bIndex] = byte(a / 256)
			bIndex++
		}
	}

	var tex uint32
	gl.GenTextures(1, &tex)
	gl.BindTexture(gl.TEXTURE_2D, tex)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, int32(w), int32(h), 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(pixels))
	return Image{
		id:     tex,
		width:  float32(w),
		height: float32(h),
	}, nil
}

func (t Image) Use() {
	gl.BindTexture(gl.TEXTURE_2D, t.id)
}

type Texture struct {
	image     Image
	texCoords mgl32.Vec4 // {u min, u max, v min, v max}
}

func NewTexture(filepath string) Texture {
	img, err := LoadImage(filepath)
	if err != nil {
		panic(err)
	}

	return Texture{
		image:     img,
		texCoords: mgl32.Vec4{0, 1, 0, 1},
	}
}

func NewTextureFromAtlas(image Image, xOffset, yOffset, width, height float32) Texture {
	umin := xOffset / image.width
	umax := (xOffset + width) / image.width
	vmin := yOffset / image.height
	vmax := (yOffset + height) / image.height
	return Texture{
		image:     image,
		texCoords: mgl32.Vec4{umin, umax, vmin, vmax},
	}
}

func (t Texture) use() {
	t.image.Use()
}
