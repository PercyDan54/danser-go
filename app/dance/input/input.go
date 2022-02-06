package input

import (
	"github.com/wieku/danser-go/app/beatmap/objects"
	"github.com/wieku/danser-go/app/graphics"
	"github.com/wieku/danser-go/framework/math/mutils"
	"math"
)

const singleTapThreshold = 140

type NaturalInputProcessor struct {
	queue  []objects.IHitObject
	cursor *graphics.Cursor

	lastTime float64

	wasLeftBefore  bool
	index int
	previousEnd    float64
	releaseAt []float64
	releaseLeftAt  float64
	releaseRightAt float64
}

func NewNaturalInputProcessor(objs []objects.IHitObject, cursor *graphics.Cursor) *NaturalInputProcessor {
	processor := new(NaturalInputProcessor)
	processor.cursor = cursor
	processor.queue = make([]objects.IHitObject, len(objs))
	processor.releaseAt = make([]float64, 4)
	processor.releaseAt[0] = -10000000
	processor.releaseAt[1] = -10000000
	processor.releaseAt[2] = -10000000
	processor.releaseAt[3] = -10000000

	copy(processor.queue, objs)

	return processor
}

func (processor *NaturalInputProcessor) Update(time float64) {
	if len(processor.queue) > 0 {
		for i := 0; i < len(processor.queue); i++ {
			g := processor.queue[i]
			if g.GetStartTime() > time {
				break
			}

			if processor.lastTime < g.GetStartTime() && time >= g.GetStartTime() {
				//startTime := g.GetStartTime()
				endTime := g.GetEndTime()

				releaseAt := endTime + 50.0

				if i+1 < len(processor.queue) {
					nTime := processor.queue[mutils.MinI(i+2, len(processor.queue)-1)].GetStartTime()

					releaseAt = mutils.ClampF64(nTime-2, endTime+1, releaseAt)
				}

				processor.releaseAt[int(math.Abs(float64(processor.index))) % 4] = releaseAt

				processor.previousEnd = endTime
				if(processor.index / 4 >= 1){
				   processor.index = -2
				}

				processor.queue = append(processor.queue[:i], processor.queue[i+1:]...)
				processor.index++

				processor.lastTime = time
				i--
			}
		}
	}

	switch (int(math.Abs(float64(processor.index - 1))) % 4){
		case 0:
			processor.cursor.LeftKey = true
			processor.cursor.RightKey = false
			processor.cursor.LeftMouse = time < processor.releaseAt[2]
			processor.cursor.RightMouse = time < processor.releaseAt[3]
		case 1:
			processor.cursor.LeftKey = time < processor.releaseAt[0]
			processor.cursor.RightKey = true
			processor.cursor.LeftMouse = false
			processor.cursor.RightMouse = false
		case 2:
			processor.cursor.LeftKey = false
			processor.cursor.RightKey = time < processor.releaseAt[1]
			processor.cursor.LeftMouse = true
			processor.cursor.RightMouse = false
		case 3:
			processor.cursor.LeftKey = false
			processor.cursor.RightKey = time < processor.releaseAt[1]
			processor.cursor.LeftMouse = time < processor.releaseAt[2]
			processor.cursor.RightMouse = true
	}

	if(time - processor.previousEnd > 100){
		processor.cursor.LeftKey = false
		processor.cursor.RightKey = false
		processor.cursor.LeftMouse = false
		processor.cursor.RightMouse = false
	}
}
