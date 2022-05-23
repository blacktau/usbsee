package gui

import (
	"github.com/gotk3/gotk3/glib"
)

type ActionActivate func()

type ActionEntry struct {
	Name          string
	Enabled       bool
	ParameterType *glib.VariantType
	OnActivate    ActionActivate
}

func mapActionEntries(actionMap glib.IActionMap, entries []ActionEntry) {
	for _, e := range entries {
		act := glib.SimpleActionNew(e.Name, e.ParameterType)
		act.Connect("activate", e.OnActivate)
		act.SetEnabled(e.Enabled)
		actionMap.AddAction(act)
	}
}
