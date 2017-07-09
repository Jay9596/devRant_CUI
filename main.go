package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"time"

	"github.com/Jay9596/devRant_cui/UI"
	"github.com/Jay9596/goRant"
	"github.com/jroimartin/gocui"
)

// struct to store rant settings
type setting struct {
	sort  string
	limit int
	skip  int
}

// Global var declaration
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
	comments    []goRant.Comment
	currCom     = 0
	current     string
	listEnd     = false
	up          = 0
	opened      int
)

// Returs the main View
// Separate func, because this view is used many times
// To reduce writing boilerplate code
func getMain() *gocui.View {
	v, err := cui.View("main")
	if err != nil {
		log.Panic(err.Error())
	}
	return v
}

// MAIN func ()
func main() {
	//Instace of goRant API Wrapper
	devRant = goRant.New()
	rantSetting.sort = "algo"
	rantSetting.limit = 20
	rantSetting.skip = 0
	var err error

	//Instace of gocui
	cui, err = gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer cui.Close()

	cui.SelFgColor = gocui.ColorYellow
	cui.Highlight = true

	// Called layout func to draw UI
	cui.SetManagerFunc(UI.Layout)
	cui.Cursor = true
	// Initialise all keybindings, except quit
	if err := initKeyBinding(cui); err != nil {
		log.Panic("Key binding err", err)
	}

	// Initialise QUIT keybinding
	if err := cui.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err.Error())
	}

	// Main Loop for the CUI
	if err := cui.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln("Error in Main loop ", err)
	}
}

//Initialise Key bindings
func initKeyBinding(g *gocui.Gui) error {
	// Bindings for Input View
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
	// Bindings for Sort View
	if err := cui.SetKeybinding("sort", gocui.KeyEnter, gocui.ModNone, setSort); err != nil {
		return err
	}
	// Bindings for Limit View
	if err := cui.SetKeybinding("limit", gocui.KeyEnter, gocui.ModNone, setLimit); err != nil {
		return err
	}
	// Bindings for All View
	// Tab binding to cycle between views
	if err := cui.SetKeybinding("", gocui.KeyTab, gocui.ModNone, nextTab); err != nil {
		return err
	}

	//Comment view binding
	if err := cui.SetKeybinding("comment", gocui.KeyArrowRight, gocui.ModNone, comDone); err != nil {
		return err
	}

	if err := cui.SetKeybinding("comment", gocui.KeyArrowUp, gocui.ModNone, comUp); err != nil {
		return err
	}

	if err := cui.SetKeybinding("comment", gocui.KeyArrowDown, gocui.ModNone, comDown); err != nil {
		return err
	}

	//Mainview keybinding
	if err := cui.SetKeybinding("main", gocui.KeyArrowDown, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			scrollView(v, 1)
			return nil
		}); err != nil {
		return err
	}

	if err := cui.SetKeybinding("main", gocui.KeyArrowUp, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			scrollView(v, -1)
			return nil
		}); err != nil {
		return err
	}

	if err := cui.SetKeybinding("main", gocui.KeyArrowRight, gocui.ModNone, mainDone); err != nil {
		return err
	}

	return nil
}

// Set current View on top
// Used to switch active view
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
	up = 0
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
	if up == 0 {
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
		up++
	}
	return nil
}

//func to handle leftKey on Input frame
func leftKey(g *gocui.Gui, v *gocui.View) error {
	x, y := v.Cursor()
	if x > 5 {
		v.SetCursor(x-1, y)
	}
	return nil
}

//func to handle rightKey on Input frame
func rightKey(g *gocui.Gui, v *gocui.View) error {
	x, y := v.Cursor()
	str, _ := v.Line(y)
	if x > len(cleanComm(str)) {
		v.SetCursor(x+1, y)
	}
	return nil
}

//func to handle backspace on Input frame
func backSp(g *gocui.Gui, v *gocui.View) error {
	x, _ := v.Cursor()
	if x > 5 {
		v.EditDelete(true)
	}
	return nil
}

//handles Right key on comments View
func comDone(g *gocui.Gui, v *gocui.View) error {
	if err := cui.DeleteView("comment"); err != nil {
		return err
	}
	if _, err := setCurrentViewOnTop(cui, "input"); err != nil {
		return err
	}
	return nil
}

//handles Up key on comments View
func comUp(g *gocui.Gui, v *gocui.View) error {
	currCom--
	if currCom >= 0 {
		printComment()
	} else {
		currCom++
		output(false, "Reached start of comments")
	}
	return nil
}

//handles down key on comments View
func comDown(g *gocui.Gui, v *gocui.View) error {
	currCom++
	if currCom < len(comments) {
		printComment()
	} else {
		currCom--
		output(false, "Reached end of comments")
	}
	return nil
}

//handles exit from main view
func mainDone(g *gocui.Gui, v *gocui.View) error {
	v.SetOrigin(0, 0)
	if _, err := setCurrentViewOnTop(cui, "input"); err != nil {
		return err
	}
	return nil
}

//handle cursor on main view
func scrollView(v *gocui.View, dy int) error {
	if v != nil {
		v.Autoscroll = false
		oX, oY := v.Origin()
		if err := v.SetOrigin(oX, oY+dy); err != nil {
			return err
		}
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

// ######################################################################
//	All func below this are used to handle the commands
//  The 'fetchFOO' func's are the handlers
//  The 'getFOO' func' are the func used as goroutines as background task
//	The 'printFOO' prints the type of FOO specified
// ######################################################################

//func to get Rants
func getRants(resp chan []goRant.Rant) {
	rs, err := devRant.Rants(rantSetting.sort, rantSetting.limit, rantSetting.skip)
	if err != nil {
		output(false, "Error occoured!")
		return
	}
	resp <- rs
}

//func to get Rant
func getRant(resp chan goRant.Rant, com chan []goRant.Comment, ID int) {
	r, comms, err := devRant.GetRant(ID)
	if err != nil {
		output(false, "Error occoured!")
	}
	resp <- r
	com <- comms
}

//func to fetch single Rant
func fetchRant(num int) {
	output(true, "Fetching Rant....")
	if len(rants) <= 0 {
		output(false, "Fetch rants before opening them")
		return
	}
	if num >= len(rants) {
		output(false, "Only 20 rants fetched at a time\nIndex out of range :p ")
		return
	}
	res := make(chan goRant.Rant)
	com := make(chan []goRant.Comment)
	go getRant(res, com, rants[num].ID)
	select {
	case comms := <-com:
		output(false, "Done!!")
		r := <-res
		comments = comms
		printRant(r, comms)
	case <-time.After(time.Second * 5):
		output(false, "Timeout...")
	default:
		r, comms, err := devRant.GetRant(rants[num].ID)
		if err != nil {
			output(false, "Error occoured!")
		}
		comments = comms
		printRant(r, comms)
	}

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
	fmt.Fprintf(view, "%s\nScore: %d\n\nTags: %s\n\nComments:%d\n", rant, r.Score, tags, r.NumComments)
}

func printComment() {
	view, err := cui.View("comment")
	if err != nil {
		output(false, err.Error())
		return
	}
	view.Clear()
	if len(comments) > 0 {
		cm := comments[currCom]
		body := fmt.Sprintf("%s     %d+\n\n%s\nScore: %d", cm.Username, cm.UserScore, cm.Body, cm.Score)
		fmt.Fprintf(view, "%d] \n\n%s\n\n", currCom, body)
		/*
			for i, cm := range coms {
				body := fmt.Sprintf("%d>>\n%s     %d+\n%s\nScore: %d", i, cm.Username, cm.UserScore, cm.Body, cm.Score)
				fmt.Fprintf(view, "%s\n\n", body)
			}
		*/
	} else {
		fmt.Fprintf(view, "%s", "No comments to show!")
	}
}

//func rants command calls
func fetchRants() {
	output(true, "Fetching rants....")
	lastLim = 0
	rCh := make(chan []goRant.Rant)
	go getRants(rCh)
	select {
	case rs := <-rCh:
		rants = rs
		current = "rants"
		printRants(rs, 0)
		output(false, "Done!!")
	case <-time.After(time.Second * 5):
		output(false, "Timeout...")
	}
}

//func to print []Rant
func printRants(rs []goRant.Rant, start int) {
	rantOpen = false
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

func getProfile(resp chan goRant.User, name string) {
	p, err := devRant.Profile(name)
	if err != nil {
		output(false, "Error occoured!")
		return
	}
	resp <- p
}

func fetchProfile(name string) {
	output(false, "Fetching profile....")
	user := make(chan goRant.User)
	go getProfile(user, name)
	select {
	case p := <-user:
		output(false, "Done!!")
		printProfile(p)
	case <-time.After(time.Second * 5):
		output(false, "Timeout...")
	}
}

//func to print profile
func printProfile(user goRant.User) {
	rantOpen = false
	v := getMain()
	v.Clear()

	info := fmt.Sprintf("Username: %s     +%d\nAbout: %s\nSkills: %s\nLocation: %s\nGithub: %s\n", user.Username, user.Score, user.About, user.Skills, user.Location, user.Github)
	counts := user.Content.Counts
	fmt.Fprintf(v, "%s\n\nRants      : %d\n++'s       : %d\nComments   : %d\nFavourites : %d", info, counts.Rants, counts.Upvotes, counts.Comments, counts.Favourites)
}

func getRes(resp chan []goRant.Rant, term string) {
	rs, err := devRant.Search(term)
	if err != nil {
		output(false, "Error occoured!")
		return
	}
	resp <- rs
}

//func to print search result
func printRes(term string) {
	output(false, "Searching....")
	lastLim = 0
	listEnd = false
	res := make(chan []goRant.Rant)
	go getRes(res, term)
	select {
	case rs := <-res:
		output(false, "Done!!")
		rants = rs
		current = "search"
		printRants(rs, 0)
	case <-time.After(time.Second * 5):
		output(false, "Timeout...")
	}

}

func getStories(resp chan []goRant.Rant) {
	ss, err := devRant.Stories()
	if err != nil {
		output(false, "Error occoured!")
		return
	}
	resp <- ss
}

func fetchStories() {
	output(false, "Fetching Stories....")
	lastLim = 0
	listEnd = false
	res := make(chan []goRant.Rant, 1)
	go getStories(res)
	select {
	case ss := <-res:
		output(false, "Done!!")
		rants = ss
		current = "stories"
		printRants(ss, 0)
	case <-time.After(time.Second * 5):
		output(false, "Timeout...")
	}
}

func getWeeklyRants(resp chan []goRant.Rant) {
	wrs, err := devRant.WeeklyRant()
	if err != nil {
		output(false, "Error occoured!")
		return
	}
	resp <- wrs
}

func fetchWeeklyRants() {
	output(false, "Fetching weekly rants....")
	lastLim = 0
	listEnd = false
	res := make(chan []goRant.Rant, 1)
	go getWeeklyRants(res)
	select {
	case wrs := <-res:
		output(false, "Done!!")
		rants = wrs
		current = "weekly"
		printRants(wrs, 0)
	case <-time.After(time.Second * 5):
		output(false, "Timeout...")
	}
}

func getSurprise(resp chan goRant.Rant, com chan []goRant.Comment) {
	r, coms, err := devRant.Surprise()
	if err != nil {
		output(false, "Error occoured!")
		return
	}
	resp <- r
	com <- coms
}

func fetchSurprise() {
	output(false, "Fetching random rant....")
	res := make(chan goRant.Rant)
	com := make(chan []goRant.Comment)
	go getSurprise(res, com)
	select {
	case coms := <-com:
		output(false, "Done!!")
		r := <-res
		comments = coms
		printRant(r, coms)
	case <-time.After(time.Second * 10):
		output(false, "Timeout...")
	}
}

func getCollabs(resp chan []goRant.Rant) {
	rs, err := devRant.Collabs()
	if err != nil {
		output(false, "Error occoured!")
		return
	}
	resp <- rs
}

func fetchCollabs() {
	lastLim = 0
	listEnd = false
	output(false, "Fetching collabs....")
	res := make(chan []goRant.Rant)
	go getCollabs(res)
	select {
	case rs := <-res:
		output(false, "Done!!")
		rants = rs
		current = "collabs"
		printRants(rs, 0)
	case <-time.After(time.Second * 5):
		output(false, "Timeout...")
	}
}

//duh! func to print help
func printHelp() {
	rantOpen = false
	v := getMain()
	v.Clear()
	output(true, "")
	help := `Command: Output

help : help info
rants : Gets rants  (at once 20 rants are fetched) 
next or :n : Print next rants (after 20, new 'rants' are fetched automatically)
rant {number} : View the single rant
comment or :c : Open comment box to view comments of a rant
:m : To switch to main window to scroll for long rant
cd.. or back : Takes you back to rants view
profile {username} : Prints user info of given username
search {term}: Prints result of search of given term
weekly : Prints weekly rants
surprise or :!! : Prints a random rant
collabs : Prints Collabs
stories : Prints stories
exit or :q : Exit
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
	//output(true, "Ctrl + C to Quit")
	cui.Close()
	os.Exit(0)
}

//Clean input by removing the leading ":/ >" and trailing extra chars
func cleanComm(com string) string {
	ind := strings.IndexRune(com, '>')
	cleanInp := strings.ToLower(com[ind+1 : len(com)])
	return stripCtlAndExtFromBytes(cleanInp)
}

// Getting the clean Input and switching to match commands
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
			return
		}
		if strings.Compare(parts[0], "rant") == 0 {
			num, err := strconv.Atoi(parts[1])
			if err != nil {
				output(false, "Invalid rant number")
				return
			}
			rantOpen = true
			opened = num
			fetchRant(num)
			return
		}
	}

	//comment Command
	if strings.Compare(cleanInp, ":c") == 0 || strings.Compare(cleanInp, "comment") == 0 {
		if rantOpen {
			if err := UI.CommentView(cui); err != nil {
				output(false, err.Error())
			}
			printComment()
		} else {
			output(false, "Open rant to load comments")
		}
		return
	}

	//profile Command
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

	//search Command
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
		output(true, "")
		fetchStories()
		return
	}

	//weekly Command
	if strings.Compare(cleanInp, "weekly") == 0 || strings.Compare(cleanInp, "weekly rants") == 0 {
		output(true, "")
		fetchWeeklyRants()
		return
	}

	//surprise Command OR :!!
	if strings.Compare(cleanInp, "surprise") == 0 || strings.Compare(cleanInp, ":!!") == 0 {
		output(true, "")
		fetchSurprise()
		return
	}

	//collab(s) command
	if strings.Compare(cleanInp, "collab") == 0 || strings.Compare(cleanInp, "collabs") == 0 {
		output(true, "")
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
					fetchRants()
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
		default:
			{
				output(false, "No rants fetched")
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
			output(true, "")
			printRants(rants, lastLim)
		} else {
			output(false, "Cannot go back")
		}
		return
	}

	//mainView command
	if strings.Compare(cleanInp, ":m") == 0 {
		if _, err := cui.SetCurrentView("main"); err != nil {
			output(false, err.Error())
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
	if limit > 0 && limit <= 20 {
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

// Needed to clean the console Input!
// Due to gocui Line func, there is a trailing /x00 added to end of string
// This function is called while cleaning the input to remove trailing /x00
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
