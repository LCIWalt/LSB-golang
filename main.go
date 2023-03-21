package main

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"io"
	"log"
	"os"

	"golang.org/x/image/bmp"
)

type Pixel struct {
	R int
	G int
	B int
	A int
}

func main() {
	in, err := os.Open("./t.bmp")
	if err != nil {
		panic(err)
	}
	img, err := bmp.Decode(in)
	if err != nil {
		log.Fatal(err)
	}

	defer in.Close()
	var encrypt string
	fmt.Scanf("%s", &encrypt)
	var length int
	length = len(encrypt)
	//fmt.Println(encrypt)
	//encrypt = "NUIST"
	var bounds = img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	var tar = []byte(encrypt)
	//fmt.Println(tar)
	var s = BytesToBinaryString(tar)

	var mark int
	mark = 0
	var pixels [][]Pixel
	for y := 0; y < height; y++ {
		var row []Pixel
		for x := 0; x < width; x++ {
			var ps [4]uint32
			ps[0], ps[1], ps[2], ps[3] = img.At(x, y).RGBA()
			if mark < len(s) {
				ps[2] = codingTheLastRGBA(ps[2]/257, s[mark]) //！！！/257！
				ps[2] *= 257
				mark++
			}
			row = append(row, rgbaToPixel(ps[0], ps[1], ps[2], ps[3]))
		}
		pixels = append(pixels, row)
	}
	r := image.NewNRGBA(image.Rect(0, 0, width, height))

	for j := 0; j < height; j++ {
		for k := 0; k < width; k++ {
			r.Set(k, j, color.RGBA{uint8(pixels[j][k].R), uint8(pixels[j][k].G), uint8(pixels[j][k].B), uint8(pixels[j][k].A)})
			//有坑 uint32和uint8的转换会导致你改动的最后一位丢失
		}
	}

	w, _ := os.Create("output.bmp")

	bmp.Encode(w, r)
	if err != nil {
		log.Fatal(err)
	}

	de, _ := os.Open("./output.bmp")
	if err != nil {
		log.Fatal(err)
	}

	decodeLSB(de, length)

}

func rgbaToPixel(r uint32, g uint32, b uint32, a uint32) Pixel {
	return Pixel{int(r / 257), int(g / 257), int(b / 257), int(a / 257)}
}

func BytesToBinaryString(bs []byte) string {
	buf := bytes.NewBuffer([]byte{})
	for _, v := range bs {
		buf.WriteString(fmt.Sprintf("%08b", v))
	}
	return buf.String()
}

func codingTheLastRGBA(x uint32, s byte) uint32 {
	if s == '1' && x%2 == 0 {
		x++
	}
	if s == '0' && x%2 == 1 {
		x--
	}
	return x
}
func decodeLSB(w io.Reader, l int) { //l:密码长度
	img, _ := bmp.Decode(w)
	var bounds = img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y
	var ans []byte
	var mark int
	mark = 0
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			var ps [4]uint32
			ps[0], ps[1], ps[2], ps[3] = img.At(x, y).RGBA()
			if mark < (l * 8) {
				ans = append(ans, cho(ps[2]/257))
				mark++
			}
		}

	}
	CoutBytes(ans)
}
func cho(x uint32) byte {
	//fmt.Print(x % 2)
	if x%2 == 0 {
		return 0
	} else {
		return 1
	}
}
func CoutBytes(b []byte) {
	var bb = make([]byte, len(b)/8)
	mark := 0
	var t byte //byte:1 --true
	t = '0' - 47

	for i := 0; i < len(b); i += 8 {
		for j := i; j < i+8; j++ {
			if b[j] == t {
				//fmt.Println("O")
				//fmt.Print(8+i-j)
				//fmt.Println(1<<(7+i-j))
				bb[mark] += 1 << (7 + i - j)
			}
		}
		mark++
	}
	fmt.Println(string(bb))
}
