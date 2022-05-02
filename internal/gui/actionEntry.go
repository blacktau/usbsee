package gui

import (
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

type ActionActivate func()

type ActionEntry struct {
	Name          string
	Enabled       bool
	ParameterType *glib.VariantType
	OnActivate    ActionActivate
}

func MapActionEntries(app *gtk.Application, entries []ActionEntry) {
	for _, e := range entries {
		act := glib.SimpleActionNew(e.Name, e.ParameterType)
		act.Connect("activate", e.OnActivate)
		act.SetEnabled(e.Enabled)
		app.AddAction(act)
	}
}
