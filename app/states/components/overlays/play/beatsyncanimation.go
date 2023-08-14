package play

import (
	"github.com/wieku/danser-go/app/beatmap"
	"github.com/wieku/danser-go/app/settings"
	"github.com/wieku/danser-go/app/skin"
	"github.com/wieku/danser-go/framework/graphics/batch"
	"github.com/wieku/danser-go/framework/graphics/sprite"
	"github.com/wieku/danser-go/framework/math/vector"
)

type BeatSyncAnimation struct {
	*sprite.Animation
	beatmap         *beatmap.BeatMap
	animationLength int
	lastTiming      float64
	beatLength      float64
	time            float64
}

func NewBeatSyncAnimation(beatmap *beatmap.BeatMap) *BeatSyncAnimation {
	animation := &BeatSyncAnimation{
		beatmap: beatmap,
	}
	animation.Animation = sprite.NewAnimation(skin.GetFrames("animation", true), 50, true, 0, vector.NewVec2d(0, 0), vector.Centre)
	animation.animationLength = animation.GetTextureCount()
	animation.SetStartTime(beatmap.Timings.GetOriginalPointAt(0).Time)

	return animation
}

func (animation *BeatSyncAnimation) Update(time float64) {
	animation.time = time
	animation.Animation.Update(time)

	timing := animation.beatmap.Timings.Current
	if &timing != nil && timing.GetBaseBeatLength() != animation.beatLength {
		animation.beatLength = timing.GetBaseBeatLength()
		animation.lastTiming = time
		animation.SetFrameDelay(animation.beatLength / float64(animation.animationLength) * float64(timing.Signature/2) / settings.Gameplay.BeatSyncAnimation.Speed)
	}
}

func (animation *BeatSyncAnimation) Draw(batch *batch.QuadBatch, alpha float64) {
	opacity := settings.Gameplay.BeatSyncAnimation.Opacity * alpha

	if opacity < 0.001 || !settings.Gameplay.BeatSyncAnimation.Show {
		return
	}

	batch.ResetTransform()
	batch.SetColor(1, 1, 1, opacity)

	xPos := settings.Gameplay.BeatSyncAnimation.XPosition
	yPos := settings.Gameplay.BeatSyncAnimation.YPosition

	batch.SetTranslation(vector.NewVec2d(xPos, yPos))

	scl := settings.Gameplay.BeatSyncAnimation.Scale

	batch.SetScale(scl, scl)
	animation.Animation.Draw(animation.time, batch)

	batch.SetColor(1, 1, 1, 1)
	batch.ResetTransform()
}
