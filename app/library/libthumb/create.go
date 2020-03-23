package libthumb

import (
	"io/ioutil"
	"math/rand"
	"strings"

	libimage "hanyny/app/library/libimage"

	"github.com/fogleman/gg"
	"github.com/gosimple/slug"
	"github.com/nfnt/resize"
)

// Thumbnail ...
type Thumbnail struct {
	Name       string
	Input      string
	Output     string
	Font1      string
	Font2      string
	Color1     string
	Color2     string
	Title1     string
	Title2     string
	Background string
}

// New  ...
func New(title string, width int, height int) (string, error) {
	thumbnail := Thumbnail{}
	err := thumbnail.getInfo(title)
	if err != nil {
		return "", err
	}
	img, err := gg.LoadImage(thumbnail.Background)
	if err != nil {
		return "", err
	}

	im := resize.Resize(uint(width), uint(height), img, resize.Lanczos3)

	dc := gg.NewContext(width, height)
	dc.SetHexColor("#1e272e")
	dc.Clear()
	dc.DrawRoundedRectangle(0, 0, 512, 512, 0)
	dc.DrawImage(im, 0, 0)

	// Draw font 3
	dc.DrawStringAnchored("https://hanyny.com", float64(width), float64(height), 1.1, -1)

	// Draw font 1
	dc.SetHexColor(thumbnail.Color1)
	dc.LoadFontFace(thumbnail.Font1, 30)
	dc.DrawStringAnchored(thumbnail.Title1, float64(width)/2, float64(height)/2, 0.5, -0.8)

	// Draw font 2
	dc.SetHexColor(thumbnail.Color2)
	dc.LoadFontFace(thumbnail.Font2, 40)
	dc.DrawStringAnchored(thumbnail.Title2, float64(width)/2, float64(height)/2, 0.5, 0.8)

	// OUT
	dc.Clip()
	dc.SavePNG(thumbnail.Output)

	go libimage.OptimizePath(thumbnail.Output, thumbnail.Output)

	return thumbnail.Output, nil
}

func (t *Thumbnail) getInfo(title string) error {
	dirBg := "public/thumbnail/backgrounds/"
	dirFont := "public/thumbnail/fonts/"
	dirOut := "public/thumbnail/store/"

	// Lấy tên background
	filesBg, err := ioutil.ReadDir(dirBg)
	if err != nil {
		return err
	}
	randomBg := rand.Intn((len(filesBg)-1)-0) + 0
	fileBgName := filesBg[randomBg].Name()

	// Lấy tên fonts
	filesFont, err := ioutil.ReadDir(dirFont)
	if err != nil {
		return err
	}
	randomFont := rand.Intn((len(filesFont)-1)-0) + 0
	fileFontName1 := filesFont[randomFont].Name()
	randomFont2 := rand.Intn((len(filesFont)-1)-0) + 0
	fileFontName2 := filesFont[randomFont2].Name()

	fileNameOutput := slug.Make(title)

	t.Background = dirBg + fileBgName
	// t.Background = dirBg + "3923277-bright-wallpaper.jpg"
	t.Output = dirOut + fileNameOutput + ".png"
	// t.Output = dirOut + "aaaaaaaaaaaaaaaaaaa" + ".png"
	t.Font1 = dirFont + fileFontName1
	t.Font2 = dirFont + fileFontName2

	// Random color
	colors := []string{"#1e272e", "#ff3f34", "#ffa801", "#f53b57", "#3c40c6", "#05c46b", "#16a085", "#27ae60", "#2980b9", "#8e44ad", "#2c3e50", "#f39c12", "#f1c40f", "#e67e22", "#d35400", "#c0392b", "#34495e", "#0a3d62", "#0c2461", "#079992", "#b71540", "#1e3799", "#1B1464", "#006266", "#EA2027", "#5758BB", "#6F1E51", "#009432", "#0652DD"}
	randomColor1 := rand.Intn((len(colors)-1)-0) + 0
	t.Color1 = colors[randomColor1]
	randomColor2 := rand.Intn((len(colors)-1)-0) + 0
	t.Color2 = colors[randomColor2]

	// Custom title
	if strings.Index(title, "non") > 0 {
		arrayChar := strings.Split(title, "non")
		if len(arrayChar) > 2 {
			t.Title1 = arrayChar[0] + "non"
			for i := 1; i < len(arrayChar); i++ {
				if i < len(arrayChar) {
					t.Title2 = arrayChar[i] + "non"
				} else {
					t.Title2 = arrayChar[i]
				}
			}
		}
	} else {
		arrayChar := strings.Split(title, " ")
		lengthTitle := len(arrayChar)

		length1 := (lengthTitle / 2) - 1
		length2 := (lengthTitle / 2) + 1
		for index := 0; index < length1; index++ {
			t.Title1 = t.Title1 + " " + arrayChar[index]
		}
		for index := 0; index < length2; index++ {
			t.Title2 = t.Title2 + " " + arrayChar[(index+length1)]
		}
	}
	return nil
}
