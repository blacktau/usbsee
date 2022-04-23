package oldgui

import (
	"fyne.io/fyne/v2"
	"github.com/blacktau/usbsee/internal/localizations"
)

func MakeMainMenu(a *fyne.App, l *localizations.Localizer) *fyne.MainMenu {
	file := fyne.NewMenu(l.Get("mainmenu.file"))

	return fyne.NewMainMenu(
		file,
	)
}
