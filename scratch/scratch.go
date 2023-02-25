package scratch

import (
	"fmt"

	"github.com/visualfc/atk/tk"
)

type Window struct {
	*tk.Window
}

func NewWindow() *Window {
	mw := &Window{tk.RootWindow()}
	lbl := tk.NewLabel(mw, "Hello ATK")
	frm := tk.NewFrame(mw)
	btn := tk.NewButton(frm, "Quit")

	btn.OnCommand(func() {
		tk.Quit()
	})
	frm.Attach(btn.Id())
	frm.SetBorderWidth(5)
	frm.SetWidth(300)
	frm.SetHeight(300)
	fmt.Printf("Form Widget Width %d", frm.Width())
	tk.NewVPackLayout(mw).AddWidgets(lbl, tk.NewLayoutSpacer(mw, 0, true), frm)

	mw.ResizeN(600, 500)
	return mw
}

func Run() {
	tk.MainLoop(func() {
		mw := NewWindow()
		mw.SetTitle("ATK Sample")
		mw.Center(nil)
		mw.ShowNormal()
	})
}
