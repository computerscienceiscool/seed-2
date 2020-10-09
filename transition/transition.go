package transition

import (
	"fmt"
	"time"

	"qlova.org/seed"
	"qlova.org/seed/client"
	"qlova.org/seed/css"
	"qlova.org/seed/js"
	"qlova.org/seed/page"
	"qlova.org/seed/popup"
	"qlova.org/seed/vfx/animation"
	"qlova.org/seed/view"
)

var fadeIn = animation.New(
	animation.Frames{
		0:   css.SetOpacity(css.Zero),
		100: css.SetOpacity(css.Number(1)),
	},
	animation.Duration(400*time.Millisecond),
)

var fadeOut = animation.New(
	animation.Frames{
		0:   css.SetOpacity(css.Number(1)),
		100: css.SetOpacity(css.Zero),
	},
	animation.Duration(400*time.Millisecond),
)

type Transition struct {
	seed.Option

	In, Out animation.Animation
}

type Option func(*Transition)

func New(options ...Option) Transition {
	var t Transition
	for _, o := range options {
		o(&t)
	}

	t.Option = seed.NewOption(func(c seed.Seed) {

		enter := js.Script(func(q js.Ctx) {
			t.In.AddTo(client.Seed{c, q})
			fmt.Fprintf(q, `seed.in(%v, 0.4);`, client.Seed{c, q}.Element())
		})

		exit := js.Script(func(q js.Ctx) {
			t.Out.AddTo(client.Seed{c, q})
			fmt.Fprintf(q, `seed.out(%v, 0.4);`, client.Seed{c, q}.Element())
		})

		switch c.(type) {
		case page.Seed:
			c.With(
				page.OnEnter(enter),
				page.OnExit(exit),
			)
		case popup.Seed:
			c.With(
				popup.OnShow(enter),
				popup.OnHide(exit),
			)
		case view.Seed:
			c.With(
				view.OnEnter(enter),
				view.OnExit(exit),
			)
		default:
			c.With(
				client.On("visible", js.Script(func(q js.Ctx) {
					t.In.AddTo(client.Seed{c, q})

				})),
				client.On("hidden", js.Script(func(q js.Ctx) {
					t.Out.AddTo(client.Seed{c, q})
				})),
			)
		}

	})

	return t
}

func In(in animation.Animation) Option {
	return func(t *Transition) {
		t.In = in
	}
}

func Out(out animation.Animation) Option {
	return func(t *Transition) {
		t.Out = out
	}
}
