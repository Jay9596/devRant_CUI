package UI

import (
	"fmt"

	"github.com/jroimartin/gocui"
)

// Layout func has defination for all views
func Layout(g *gocui.Gui) error {
	mX, mY := g.Size()
	// Input area, on lower right side
	inp, err := g.SetView("input", mX/2, mY/2-1, mX-1, mY-1)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		inp.Editable = true
		inp.Overwrite = false
		inp.Wrap = true
		inp.Title = "Input"
		if _, err := g.SetCurrentView("input"); err != nil {
			return err
		}
		fmt.Fprintf(inp, " :/ >")
		inp.SetCursor(5, 0)
	}
	// Main view, the left portion
	m, mErr := g.SetView("main", 0, 0, mX/2-1, mY-1)
	if mErr != nil {
		if mErr != gocui.ErrUnknownView {
			return mErr
		}
		m.Autoscroll = false
		m.Wrap = true
		m.Editable = false
	}

	// The sort view
	alg, aErr := g.SetView("sort", mX/2, mY/4-1, mX-1, mY/4+1)
	if aErr != nil {
		if aErr != gocui.ErrUnknownView {
			return aErr
		}
		alg.Title = "Sort"
		alg.Editable = true
		alg.Frame = true
		fmt.Fprintf(alg, "algo")
		alg.SetCursor(0, 0)
	}

	// The limit view
	lim, lErr := g.SetView("limit", mX/2, mY/4+2, mX-1, mY/4+4)
	if lErr != nil {
		if lErr != gocui.ErrUnknownView {
			return lErr
		}
		lim.Title = "Limit"
		lim.Editable = true
		lim.Frame = true
		fmt.Fprintf(lim, "5")
		lim.SetCursor(0, 0)
	}

	// The output view
	ds, dErr := g.SetView("output", mX/2, mY/2-7, mX-1, mY/2-2)
	if dErr != nil {
		if dErr != gocui.ErrUnknownView {
			return dErr
		}
		ds.Editable = false
		ds.Wrap = true
		ds.Autoscroll = true
	}

	//Bannder with devRant written
	logo, artErr := g.SetView("banner", mX/2, 0, mX-2, mY/4-1)
	if artErr != nil {
		if artErr != gocui.ErrUnknownView {
			return err
		}
		logo.Wrap = false
		logo.Editable = false
		logo.Frame = false

		art := ` ########  ######## ##     ## ########     ###    ##    ## ######## 
 ##     ## ##       ##     ## ##     ##   ## ##   ###   ##    ##    
 ##     ## ##       ##     ## ##     ##  ##   ##  ####  ##    ##    
 ##     ## ######   ##     ## ########  ##     ## ## ## ##    ##    
 ##     ## ##        ##   ##  ##   ##   ######### ##  ####    ##    
 ##     ## ##         ## ##   ##    ##  ##     ## ##   ###    ##    
 ########  ########    ###    ##     ## ##     ## ##    ##    ##   `
		fmt.Fprintf(logo, "%s", art)
	}
	return nil
}

//CommentView creates a new View for comments
func CommentView(g *gocui.Gui) error {
	mX, mY := g.Size()
	v, err := g.SetView("comment", 0, mY/2-1, mX/2-1, mY-1)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Wrap = true
		v.Title = "Comment"
		v.Frame = true
		v.Autoscroll = false
		v.Editable = false
	}
	if _, err := g.SetCurrentView("comment"); err != nil {
		return err
	}
	return nil
}
