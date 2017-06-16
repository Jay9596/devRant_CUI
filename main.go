package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/Jay9596/goRant"
	"github.com/jroimartin/gocui"
)

type setting struct {
	sort  string
	limit int
	skip  int
}

var (
	devRant     *goRant.Client
	cui         *gocui.Gui
	views       = []string{"input", "sort", "limit"}
	active      = 0
	printLimit  = 5
	lastLim     = 0
	rantSetting setting
	rants       []goRant.Rant
	rantOpen    = false
	current     string
	listEnd     = false
)

func layout(g *gocui.Gui) error {
	mX, mY := g.Size()
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
	m, mErr := g.SetView("main", 0, 0, mX/2-1, mY-1)
	if mErr != nil {
		if mErr != gocui.ErrUnknownView {
			return mErr
		}
		m.Autoscroll = false
		m.Wrap = true
		m.Editable = false
	}
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
	ds, dErr := cui.SetView("output", mX/2, mY/2-7, mX-1, mY/2-2)
	if dErr != nil {
		if dErr != gocui.ErrUnknownView {
			return dErr
		}
		ds.Editable = false
		ds.Wrap = true
		ds.Autoscroll = true
	}

	logo, artErr := cui.SetView("banner", mX/2, 0, mX-2, mY/4-1)
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

func getMain() *gocui.View {
	v, err := cui.View("main")
	if err != nil {
		log.Panic(err.Error())
	}
	return v
}

func main() {
	devRant = goRant.New()
	rantSetting.sort = "algo"
	rantSetting.limit = 20
	rantSetting.skip = 0
	var err error
	cui, err = gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer cui.Close()

	cui.SelFgColor = gocui.ColorYellow
	cui.Highlight = true

	cui.SetManagerFunc(layout)

	if err := initKeyBinding(cui); err != nil {
		log.Panic("Key binding err", err)
	}

	if err := cui.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err.Error())
	}

	if err := cui.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln("Error in Main loop ", err)
	}
}

//Initialise Key bindings
func initKeyBinding(g *gocui.Gui) error {
	if err := cui.SetKeybinding("input", gocui.KeyEnter, gocui.ModNone, enterCom); err != nil {
		return err
	}
	if err := cui.SetKeybinding("input", gocui.KeyArrowUp, gocui.ModNone, upKey); err != nil {
		return err
	}
	if err := cui.SetKeybinding("input", gocui.KeyBackspace, gocui.ModNone, backSp); err != nil {
		return err
	}
	if err := cui.SetKeybinding("input", gocui.KeyArrowLeft, gocui.ModNone, leftKey); err != nil {
		return err
	}
	if err := cui.SetKeybinding("input", gocui.KeyArrowRight, gocui.ModNone, rightKey); err != nil {
		return err
	}
	if err := cui.SetKeybinding("input", gocui.KeyArrowDown, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error { return nil }); err != nil {
		return err
	}

	if err := cui.SetKeybinding("sort", gocui.KeyEnter, gocui.ModNone, setSort); err != nil {
		return err
	}
	if err := cui.SetKeybinding("limit", gocui.KeyEnter, gocui.ModNone, setLimit); err != nil {
		return err
	}
	if err := cui.SetKeybinding("", gocui.KeyTab, gocui.ModNone, nextTab); err != nil {
		return err
	}
	return nil
}

//Set current View on top
func setCurrentViewOnTop(g *gocui.Gui, name string) (*gocui.View, error) {
	if _, err := g.SetCurrentView(name); err != nil {
		return nil, err
	}
	return g.SetViewOnTop(name)
}

//Func to switch between frames
func nextTab(g *gocui.Gui, v *gocui.View) error {
	nextIndex := (active + 1) % len(views)
	name := views[nextIndex]

	if _, err := setCurrentViewOnTop(g, name); err != nil {
		return err
	}

	if nextIndex == 0 || nextIndex == 3 {
		g.Cursor = true
	} else {
		g.Cursor = false
	}

	active = nextIndex
	return nil
}

//func to handle Inputs
func enterCom(g *gocui.Gui, v *gocui.View) error {
	x, _ := v.Size()
	_, lineY := v.Cursor()
	v.SetCursor(x, lineY)
	command, err := v.Line(lineY)
	if err != nil {
		return err
	}
	v.EditNewLine()
	fmt.Fprintf(v, ":/ >")
	v.SetCursor(5, lineY+1)
	checkCommand(command)
	return nil
}

//func to QUIT
func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

//func to handle upKey on Input frame
func upKey(g *gocui.Gui, v *gocui.View) error {
	_, lineY := v.Cursor()
	if lineY == 0 {
		return nil
	}
	preCom, err := v.Line(lineY - 1)
	if err != nil {
		return err
	}
	com := cleanComm(preCom)
	fmt.Fprintf(v, "%s", com)
	v.SetCursor(5+len(com), lineY)
	return nil
}

func leftKey(g *gocui.Gui, v *gocui.View) error {
	x, y := v.Cursor()
	if x > 5 {
		v.SetCursor(x-1, y)
	}
	return nil
}

func rightKey(g *gocui.Gui, v *gocui.View) error {
	x, y := v.Cursor()
	str, _ := v.Line(y)
	if x > len(cleanComm(str)) {
		v.SetCursor(x+1, y)
	}
	return nil
}

func backSp(g *gocui.Gui, v *gocui.View) error {
	x, _ := v.Cursor()
	if x > 5 {
		v.EditDelete(true)
	}
	return nil
}

//func to print to output frame
func output(clr bool, str string) {
	v, _ := cui.View("output")
	if clr {
		v.Clear()
	}
	fmt.Fprintf(v, "%s\n", str)
}

/*//redundant func, same as duput()
func printThis(clr bool, str string) {
	v, _ := cui.View("output")
	if clr {
		v.Clear()
	}
	fmt.Fprintf(v, "%s\n", str)
}*/

//func to get Rants
func getRants(algo string, limit int, skip int) []goRant.Rant {
	rs, err := devRant.Rants(rantSetting.sort, rantSetting.limit, rantSetting.skip)
	if err != nil {
		output(false, err.Error())
	}
	return rs
}

//func to get single Rant
func fetchRant(num int) {
	output(true, "Fetching Rant")
	r, comms, err := devRant.GetRant(rants[num].ID)
	if err != nil {
		v, _ := cui.View("output")
		v.Clear()
		fmt.Fprint(v, err.Error())
	}
	printRant(r, comms)
}

//func to print single Rant
func printRant(r goRant.Rant, coms []goRant.Comment) {
	var rant, tags string
	view := getMain()
	if r.AttachedImage.URL == "" {
		rant = fmt.Sprintf("%s     %d+\n\n%s", r.Username, r.UserScore, r.Text)
	} else {
		rant = fmt.Sprintf("%s     %d+\n\n%s\nImage: \"%s\"", r.Username, r.UserScore, r.Text, r.AttachedImage.URL)
	}
	for _, t := range r.Tags {
		tags += t + ", "
	}
	view.Clear()
	fmt.Fprintf(view, "%s\nScore: %dn\nTags: %s\n\nComments:%d\n", rant, r.Score, tags, r.NumComments)
	for _, cm := range coms {
		body := fmt.Sprintf("%s     %d+\n%s\nScore: %d", cm.Username, cm.UserScore, cm.Body, cm.Score)
		fmt.Fprintf(view, "%s\n\n", body)
	}
}

//func rants command calls
func fetchRants() {
	output(true, "Fetching rants....")
	lastLim = 0
	view := getMain()
	view.Clear()
	ch := make(chan bool)
	rCh := make(chan []goRant.Rant)
	go func(ch chan bool, r chan []goRant.Rant) {
		rs := getRants(rantSetting.sort, rantSetting.limit, rantSetting.skip)
		r <- rs
		ch <- true
	}(ch, rCh)
	rs := <-rCh
	if got := <-ch; got {
		rants = rs
		current = "rants"
		printRants(rs, 0)
		output(false, "Done!!")
	}
}

//func to print []Rant
func printRants(rs []goRant.Rant, start int) {
	view := getMain()
	view.Clear()
	l := len(rs)
	lim := printLimit
	if start >= l {
		lim = l - (start - printLimit)
	}
	for i := start; i < start+lim; i++ {
		if i >= l {
			output(false, "End of List")
			listEnd = true
			break
		}
		var rant string
		r := rs[i]
		if r.AttachedImage.URL != "" {
			rant = fmt.Sprintf("%s\nImage : \"%s\"", r.Text, r.AttachedImage.URL)
		} else {
			rant = fmt.Sprintf("%s", r.Text)
		}
		fmt.Fprintf(view, ">>%d\n%s\n ++ %d -- \n\n", i, rant, r.Score)
	}
}

func fetchProfile(name string) {
	p, err := devRant.Profile(name)
	if err != nil {
		output(false, err.Error())
		return
	}
	printProfile(p)
}

//func to print profile
func printProfile(user goRant.User) {
	v := getMain()
	v.Clear()

	info := fmt.Sprintf("Username: %s     +%d\nAbout: %s\nSkills: %s\nLocation: %s\nGithub: %s\n", user.Username, user.Score, user.About, user.Skills, user.Location, user.Github)
	counts := user.Content.Counts
	fmt.Fprintf(v, "%s\n\nRants      : %d\n++'s       : %d\nComments   : %d\nFavourites : %d", info, counts.Rants, counts.Upvotes, counts.Comments, counts.Favourites)
}

//func to print search result
func printRes(term string) {
	lastLim = 0
	rs, err := devRant.Search(term)
	if err != nil {
		output(false, err.Error())
	}
	rants = rs
	current = "search"
	printRants(rs, 0)
}

func fetchStories() {
	lastLim = 0
	listEnd = false
	ss, err := devRant.Stories()
	if err != nil {
		output(false, err.Error())
		return
	}
	rants = ss
	current = "stories"
	printRants(ss, 0)
}

func fetchWeeklyRants() {
	lastLim = 0
	listEnd = false
	wrs, err := devRant.WeeklyRant()
	if err != nil {
		output(false, err.Error())
		return
	}
	rants = wrs
	current = "weekly"
	printRants(wrs, 0)
}

func fetchSurprise() {
	r, coms, err := devRant.Surprise()
	if err != nil {
		output(false, err.Error())
		return
	}
	printRant(r, coms)
}

func fetchCollabs() {
	lastLim = 0
	listEnd = false
	rs, err := devRant.Collabs()
	if err != nil {
		output(false, err.Error())
		return
	}
	rants = rs
	current = "collabs"
	printRants(rs, 0)
}

//duh! func to print help
func printHelp() {
	v := getMain()
	v.Clear()
	output(true, "")
	help := `Command: Output
help: help info
rants: Gets all rants
next or :n : print next rants (at once 20 rants are fetched) 
rant <number>: View the single rant
cd.. or back : to go back to rants view
profile <username>: Prints user info of given username
search <term>: Prints result of search of given term
weekly : Prints weekly rants
surprise or :!! : Prints a random rant
collabs : Prints Collabs
stories : Prints stories
exit or ":q": Exit
Ctrl + C: Quit`
	fmt.Fprintf(v, help)

}

func clearConsole() {
	v, err := cui.View("input")
	if err != nil {
		log.Panic(err.Error())
	}
	v.Clear()
	o, _ := cui.View("output")
	o.Clear()
	fmt.Fprint(v, " :/ >")
	v.SetCursor(5, 0)
}

func exit() {
	getMain().Clear()
	output(true, "Ctrl + C to Quit")
}

//Clean input by removing the leading ":/ >" and trailing extra chars
func cleanComm(com string) string {
	ind := strings.IndexRune(com, '>')
	cleanInp := strings.ToLower(com[ind+1 : len(com)])
	return stripCtlAndExtFromBytes(cleanInp)
}

func checkCommand(com string) {
	cleanInp := cleanComm(com)

	//rants Command
	if strings.Compare(cleanInp, "rants") == 0 {
		fetchRants()
		return
	}

	//rant Command
	if strings.Contains(cleanInp, "rant") {
		parts := strings.Fields(cleanInp)
		if len(parts) != 2 {
			output(false, "Invalid usage of rant command\nRefer help")
		} else {
			if strings.Compare(parts[0], "rant") == 0 {
				num, err := strconv.Atoi(parts[1])
				if err != nil {
					output(false, "Invalid rant number")
					return
				}
				rantOpen = true
				fetchRant(num)
			}
		}
		return
	}

	//Profile Command
	if strings.Contains(cleanInp, "profile") {
		parts := strings.Fields(cleanInp)
		if len(parts) != 2 {
			output(false, "Invalid usage of profile command\nRefer help")
		} else {
			if strings.Compare(parts[0], "profile") == 0 {
				name := parts[1]
				fetchProfile(name)
			}
		}
		return
	}

	//Search Command
	if strings.Contains(cleanInp, "search") {
		parts := strings.Fields(cleanInp)
		if len(parts) != 2 {
			output(false, "Invalid usage of search command\nRefer help")
		} else {
			if strings.Compare(parts[0], "search") == 0 {
				term := parts[1]
				output(false, fmt.Sprintf("Searched: %s", term))
				printRes(term)
			}
		}
		return
	}

	//stories Command
	if strings.Compare(cleanInp, "stories") == 0 {
		fetchStories()
		return
	}

	//weekly Command
	if strings.Compare(cleanInp, "weekly") == 0 || strings.Compare(cleanInp, "weekly rants") == 0 {
		fetchWeeklyRants()
		return
	}

	//surprise Command OR :!!
	if strings.Compare(cleanInp, "surprise") == 0 || strings.Compare(cleanInp, ":!!") == 0 {
		fetchSurprise()
		return
	}

	if strings.Compare(cleanInp, "collab") == 0 || strings.Compare(cleanInp, "collabs") == 0 {
		fetchCollabs()
		return
	}

	//Next Command
	if strings.Compare(cleanInp, ":n") == 0 || strings.Compare(cleanInp, "next") == 0 {
		lastLim += printLimit
		switch current {
		case "rants":
			{
				if lastLim < rantSetting.limit {
					printRants(rants, lastLim)
				} else {

					lastLim = 0
					rantSetting.skip += rantSetting.limit
					rants = getRants(rantSetting.sort, rantSetting.limit+20, rantSetting.skip)
					printRants(rants, lastLim)
				}
				break
			}
		case "weekly":
		case "search":
		case "stories":
		case "collabs":
			{
				if !listEnd {
					printRants(rants, lastLim)
				} else {
					output(false, "End of list\nFetch something else")
				}
			}
		}

		return
	}
	//Help Command
	if strings.Compare(cleanInp, "help") == 0 {
		printHelp()
		return
	}

	//back Command or cd ..
	if strings.Compare(cleanInp, "cd ..") == 0 || strings.Compare(cleanInp, "cd..") == 0 || strings.Compare(cleanInp, "back") == 0 {
		if rantOpen {
			rantOpen = false
			printRants(rants, lastLim)
		} else {
			output(false, "Cannot go back")
		}
		return
	}
	//Clear Command
	if strings.Compare(cleanInp, "clear") == 0 {
		clearConsole()
		return
	}

	//Clean Command
	if strings.Compare(cleanInp, "clean") == 0 {
		getMain().Clear()
		return
	}

	//Commands for exiting
	if strings.Compare(cleanInp, "exit") == 0 || strings.Compare(cleanInp, ":q") == 0 {
		exit()
		return
	}

	//Incase of invalid command
	output(true, "\nInvalid Command\nRefer \"help\"")

}

//Handler for Sort frame
func setSort(g *gocui.Gui, v *gocui.View) error {
	var sort string
	val, err := v.Line(0)
	if err != nil {
		return err
	}
	val = stripCtlAndExtFromBytes(val)
	if strings.Compare(val, "algo") == 0 {
		sort = "algo"
	} else if strings.Compare(val, "recent") == 0 {
		sort = "recent"
	} else if strings.Compare(val, "top") == 0 {
		sort = "top"
	} else {
		sort = "algo"
	}
	rantSetting.sort = sort
	output(false, fmt.Sprintf("Sort set to %s", rantSetting.sort))
	go func() {
		v, _ := cui.View("sort")
		v.Clear()
		fmt.Fprintf(v, "%s", rantSetting.sort)
		v.SetCursor(0, 0)
	}()
	return nil
}

//Handler for Limit frame
func setLimit(g *gocui.Gui, v *gocui.View) error {
	var limit int
	n, err := v.Line(0)
	if err != nil {
		return err
	}
	n = stripCtlAndExtFromBytes(n)
	num, e := strconv.Atoi(strings.TrimSpace(n))
	if e != nil {
		output(false, "Limit shout be a number >= 0")
		printLimit = 5
	}
	limit = num
	if limit > 0 {
		printLimit = limit
	} else {
		printLimit = 5
	}
	output(false, fmt.Sprintf("Limit set to %d", printLimit))
	go func() {
		v, _ := cui.View("limit")
		v.Clear()
		fmt.Fprintf(v, "%d", printLimit)
		v.SetCursor(0, 0)
	}()
	return nil
}

//Needed to clean the console Input
func stripCtlAndExtFromBytes(str string) string {
	b := make([]byte, len(str))
	var bl int
	for i := 0; i < len(str); i++ {
		c := str[i]
		if c >= 32 && c < 127 {
			b[bl] = c
			bl++
		}
	}
	return string(b[:bl])
}
