package ic

import (
	"bytes"
	"encoding/base64"
	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"io/ioutil"
	"math/rand"
	"time"
)

var IC = identifyCode{}

type identifyCode struct {}

func (ic *identifyCode) GetBase64Encode(codeLength int) (code string, result string) {
	code = ic.getRandStr(codeLength)

	data := ic.imgText(200, 100, code)
	result = base64.StdEncoding.EncodeToString(data)
	result = "data:image/jpg;base64," + result
	return
}

func (ic *identifyCode) SaveFile(filepath string, codeLength int) {
	code := ic.getRandStr(codeLength)
	data := ic.imgText(200, 100, code)
	ioutil.WriteFile(filepath, data, 0644)
}

func (i *identifyCode) getRandStr(n int) (randStr string) {
	chars := "ABCDEFGHIJKMNPQRSTUVWXYZabcdefghijkmnpqrstuvwxyz23456789"
	charsLen := len(chars)
	if n > 10 {
		n = 10
	}

	rand.Seed(time.Now().UnixNano())
	for i := 0; i < n; i++ {
		randIndex := rand.Intn(charsLen)
		randStr += chars[randIndex : randIndex+1]
	}
	return randStr
}

func (ic *identifyCode) imgText(width, height int, text string) (b []byte) {
	textLen := len(text)
	dc := gg.NewContext(width, height)
	bgR, bgG, bgB, bgA := ic.getRandColorRange(240, 255)
	dc.SetRGBA255(bgR, bgG, bgB, bgA)
	dc.Clear()

	// 干扰线
	for i := 0; i < 30; i++ {
		x1, y1 := ic.getRandPos(width, height)
		x2, y2 := ic.getRandPos(width, height)
		r, g, b, a := ic.getRandColor(255)
		w := float64(rand.Intn(6) + 1)
		dc.SetRGBA255(r, g, b, a)
		dc.SetLineWidth(w)
		dc.DrawLine(x1, y1, x2, y2)
		dc.Stroke()
	}

	fontSize := float64(height/2) + 5
	face := ic.loadFontFace(fontSize)
	dc.SetFontFace(face)

	for i := 0; i < len(text); i++ {
		r, g, b, _ := ic.getRandColor(100)
		dc.SetRGBA255(r, g, b, 255)
		fontPosX := float64(width/textLen*i) + fontSize*0.6

		ic.writeText(dc, text[i:i+1], float64(fontPosX), float64(height/2))
	}

	buffer := bytes.NewBuffer(nil)
	dc.EncodePNG(buffer)
	b = buffer.Bytes()
	return
}

// 渲染文字
func (ic *identifyCode) writeText(dc *gg.Context, text string, x, y float64) {
	xfload := 5 - rand.Float64()*10 + x
	yfload := 5 - rand.Float64()*10 + y

	radians := 40 - rand.Float64()*80
	dc.RotateAbout(gg.Radians(radians), x, y)
	dc.DrawStringAnchored(text, xfload, yfload, 0.2, 0.5)
	dc.RotateAbout(-1*gg.Radians(radians), x, y)
	dc.Stroke()
}

// 随机坐标
func (ic *identifyCode) getRandPos(width, height int) (x float64, y float64) {
	x = rand.Float64() * float64(width)
	y = rand.Float64() * float64(height)
	return x, y
}

// 随机颜色
func (ic *identifyCode) getRandColor(maxColor int) (r, g, b, a int) {
	r = int(uint8(rand.Intn(maxColor)))
	g = int(uint8(rand.Intn(maxColor)))
	b = int(uint8(rand.Intn(maxColor)))
	a = int(uint8(rand.Intn(255)))
	return r, g, b, a
}

// 随机颜色范围
func (ic *identifyCode) getRandColorRange(miniColor, maxColor int) (r, g, b, a int) {
	if miniColor > maxColor {
		miniColor = 0
		maxColor = 255
	}
	r = int(uint8(rand.Intn(maxColor-miniColor) + miniColor))
	g = int(uint8(rand.Intn(maxColor-miniColor) + miniColor))
	b = int(uint8(rand.Intn(maxColor-miniColor) + miniColor))
	a = int(uint8(rand.Intn(maxColor-miniColor) + miniColor))
	return r, g, b, a
}

// 加载字体
func (ic *identifyCode) loadFontFace(points float64) font.Face {
	data, _ := ioutil.ReadFile("my.ttf")
	f, err := truetype.Parse(data)

	if err != nil {
		panic(err)
	}
	face := truetype.NewFace(f, &truetype.Options{
		Size: points,
	})
	return face
}
