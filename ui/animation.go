package ui

import (
	"math"
)

type Location struct {
	X, Y float64
}

type PositionSettable interface {
	SetPos(x, y float64)
}

type LinearPathAnimator struct {
	Name         string
	Target       PositionSettable
	Finished     bool
	Path         []Location
	CurrentFrame int
	Complete     func()
	Step         func()
}

func LinearEasing(x float64) float64 {
	return x
}

func EaseOutQuadratic(x float64) float64 {
	return 1 - math.Pow(1-x, 2)
}

func EaseOutCubic(x float64) float64 {
	return 1 - math.Pow(1-x, 3)
}

func EaseOutQuintic(x float64) float64 {
	return 1 - math.Pow(1-x, 5)
}

// TODO: think about what to do if framecount is 1 or 0
func NewLinearPathAnimator(t PositionSettable, frameCount int, start Location, end Location, easing func(float64) float64, complete func()) *LinearPathAnimator {
	dx := end.X - start.X
	dy := end.Y - start.Y
	path := make([]Location, frameCount-1)
	for i := 1; i < frameCount; i++ {
		t := float64(i) / float64(frameCount-1)
		et := easing(t)
		path[i-1] = Location{X: start.X + dx*et, Y: start.Y + dy*et}
	}

	return &LinearPathAnimator{
		Target:       t,
		Finished:     false,
		Path:         path,
		CurrentFrame: 0,
		Complete:     complete,
	}
}

func (a *LinearPathAnimator) Update() {
	l := a.Path[a.CurrentFrame]
	a.Target.SetPos(l.X, l.Y)
	a.CurrentFrame += 1
	//fmt.Printf("linear path frame %d of %d\n", a.CurrentFrame, len(a.Path))

	if a.CurrentFrame >= len(a.Path) {
		a.Finished = true
		if a.Complete != nil {
			a.Complete()
		}
		return
	}
}

func (a *LinearPathAnimator) IsFinished() bool {
	return a.Finished
}

func (a *LinearPathAnimator) IsBlocking() bool {
	// TODO: this may want to be configurable
	return !a.Finished
}

type Anim interface {
	Update()
	IsBlocking() bool
	IsFinished() bool
}

type BlockingAnimation struct {
	TicksLeft int
	Finished  bool
	Name      string
}

func (b *BlockingAnimation) Update() {
	b.TicksLeft -= 1
	if b.TicksLeft == 0 {
		b.Finished = true
	}
}

func (b *BlockingAnimation) IsFinished() bool {
	return b.Finished
}

func (b *BlockingAnimation) IsBlocking() bool {
	// finishes at the same time as blocking
	return !b.IsFinished()
}

func NewBlockingAnim(ticks int) Anim {
	return &BlockingAnimation{
		TicksLeft: ticks,
		Finished:  false,
	}
}
