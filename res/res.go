package res

import (
	"bytes"
	"image"
	"io"
	"io/fs"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"

	"github.com/hajimehoshi/ebiten/v2/audio/vorbis"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
	"golang.org/x/image/font/sfnt"

	"embed"
)

//go:embed font/* img/*.png audio/* shader/*
var assets embed.FS

var fonts map[string]*sfnt.Font = make(map[string]*sfnt.Font)

var shaderDict = make(map[string]*ebiten.Shader)

func init() {
	LoadFonts()
	LoadShaders()
}

func LoadShaders() {
	files := make([]string, 0, 10)
	if err := fs.WalkDir(&assets, "shader", func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}
		files = append(files, path)
		return nil
	}); err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		bytes, err := assets.ReadFile(f)
		if err != nil {
			log.Fatal(err)
		}
		shader, err := ebiten.NewShader(bytes)
		if err != nil {
			log.Fatal(err)
		}
		base := strings.TrimSuffix(filepath.Base(f), filepath.Ext(f))
		shaderDict[base] = shader
	}
}

func LoadFonts() {
	bytes, err := assets.ReadFile("font/Roboto-Medium.ttf")
	if err != nil {
		log.Fatal(err)
	}
	f, err := sfnt.Parse(bytes)

	if err != nil {
		log.Fatal(err)
	}

	fonts["Roboto-Medium"] = f
}

func GetFont(n string) *sfnt.Font {
	return fonts[n]
}

func ReadImage(p string) (image.Image, error) {
	data, err := assets.ReadFile(path.Join("img", p))
	if err != nil {
		return nil, err
	}
	img, _, err := image.Decode(bytes.NewReader(data))
	return img, err
}

// Images is the map of all loaded images.
var Images map[string]*ebiten.Image = make(map[string]*ebiten.Image)

// GetImage returns the image matching the given file name. IT ALSO LOADS IT.
func GetImage(p string) *ebiten.Image {
	if v, ok := Images[p]; ok {
		return v
	}
	img, err := ReadImage(p + ".png")
	if err != nil {
		log.Println("error reading image " + p)
		return nil
	}
	eimg := ebiten.NewImageFromImage(img)
	Images[p] = eimg
	return eimg
}

func DecodeWavToBytes(audioContext *audio.Context, fileName string) []byte {
	data, err := assets.ReadFile(path.Join("audio", fileName))
	if err != nil {
		log.Fatal(err)
		return nil
	}
	s, err := wav.DecodeWithSampleRate(48000, bytes.NewReader(data))
	if err != nil {
		log.Fatal(err)
		return nil
	}
	b, err := io.ReadAll(s)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	return b
}

type AudioStream interface {
	io.ReadSeeker
	Length() int64
}

func OggToStream(audioContext *audio.Context, fileName string) AudioStream {
	data, err := assets.ReadFile(path.Join("audio", fileName))
	if err != nil {
		log.Fatal(err)
		return nil
	}
	s, err := vorbis.DecodeWithoutResampling(bytes.NewReader(data))
	if err != nil {
		log.Fatal(err)
		return nil
	}
	return s
}

func GetShader(name string) *ebiten.Shader {
	return shaderDict[name]
}

func ReloadShader(name string) {
	file, err := os.ReadFile(name)
	if err != nil {
		log.Fatal(err)
		return
	}
	shader, err := ebiten.NewShader(file)
	if err != nil {
		log.Println(err)
		return
	}
	base := strings.TrimSuffix(filepath.Base(name), filepath.Ext(name))
	shaderDict[base] = shader
}
