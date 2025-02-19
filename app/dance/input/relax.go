package input

import (
	"github.com/wieku/danser-go/app/beatmap/difficulty"
	"github.com/wieku/danser-go/app/graphics"
	"github.com/wieku/danser-go/app/rulesets/osu"
)

const (
	stableLeniency = 12
	lazerLeniency  = 3
)

type RelaxInputProcessor struct {
	cursor  *graphics.Cursor
	ruleset *osu.OsuRuleSet
	wasLeft bool
}

func NewRelaxInputProcessor(ruleset *osu.OsuRuleSet, cursor *graphics.Cursor) *RelaxInputProcessor {
	processor := new(RelaxInputProcessor)
	processor.cursor = cursor
	processor.ruleset = ruleset

	return processor
}

func (processor *RelaxInputProcessor) Update(time float64) {
	processed := processor.ruleset.GetProcessed()
	player := processor.ruleset.GetPlayer(processor.cursor)

	click := false

	isLazer := processor.ruleset.GetPlayerDifficulty(processor.cursor).CheckModActive(difficulty.Lazer)

	for _, o := range processed {
		circle, c1 := o.(*osu.Circle)
		slider, c2 := o.(*osu.Slider)

		objectStartTime := processor.ruleset.GetBeatMap().HitObjects[o.GetNumber()].GetStartTime()

		if ((c1 && !circle.IsHit(player)) || (c2 && !slider.IsStartHit(player))) &&
			((!isLazer && time > objectStartTime-stableLeniency) || (isLazer && time >= objectStartTime-lazerLeniency)) {
			click = true
		}
	}

	processor.cursor.LeftButton = click && !processor.wasLeft
	processor.cursor.RightButton = click && processor.wasLeft

	if click {
		processor.wasLeft = !processor.wasLeft
	}
}
