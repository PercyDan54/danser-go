package play

import (
	"github.com/wieku/danser-go/app/settings"
	"github.com/wieku/danser-go/framework/graphics/batch"
	"github.com/wieku/danser-go/framework/graphics/font"
	color2 "github.com/wieku/danser-go/framework/math/color"
	"github.com/wieku/danser-go/framework/math/vector"
	"strconv"
)

type KpsCounter struct {
	font  *font.Font
	text  string
	value int
	hits  []float64
}

func NewKpsCounter() *KpsCounter {
	return &KpsCounter{
		font:  font.GetFont("HUDFont"),
		text:  "0 kps",
		value: 0,
	}
}

func (kpsCounter *KpsCounter) Add(time float64) {
	kpsCounter.hits = append(kpsCounter.hits, time)
}

func (kpsCounter *KpsCounter) Update(time float64) {
	var newHits []float64
	for _, hit := range kpsCounter.hits {
		window := 1000.0
		relativeTime := time - hit
		if relativeTime <= window {
			newHits = append(newHits, hit)
		}
	}
	kpsCounter.hits = newHits
	kpsCounter.value = len(newHits)
	kpsCounter.text = strconv.Itoa(kpsCounter.value) + " kps"
}

func (kpsCounter *KpsCounter) Draw(batch *batch.QuadBatch, alpha float64) {
	batch.ResetTransform()

	kpsAlpha := settings.Gameplay.KpsCounter.Opacity * alpha

	if kpsAlpha < 0.001 || !settings.Gameplay.KpsCounter.Show {
		return
	}

	scale := settings.Gameplay.KpsCounter.Scale
	position := vector.NewVec2d(settings.Gameplay.KpsCounter.XPosition, settings.Gameplay.KpsCounter.YPosition)
	origin := vector.ParseOrigin(settings.Gameplay.KpsCounter.Align)
	cS := settings.Gameplay.KpsCounter.Color

	color := color2.NewHSVA(float32(cS.Hue), float32(cS.Saturation), float32(cS.Value), float32(kpsAlpha))
	kpsCounter.draw(batch, kpsCounter.text, position, 0, scale, color, origin)

	batch.ResetTransform()
}

func (kpsCounter *KpsCounter) draw(batch *batch.QuadBatch, text string, position vector.Vector2d, length float64, scale float64, color color2.Color, origin vector.Vector2d) {
	batch.SetColor(0, 0, 0, float64(color.A)*0.8)
	kpsCounter.font.DrawOriginV(batch, position.AddS(scale+length, scale), origin, 40*scale, true, text)

	batch.SetColorM(color)
	kpsCounter.font.DrawOriginV(batch, position.AddS(length, 0), origin, 40*scale, true, text)
}
