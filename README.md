# devRant_CUI
A [devRant.io](https://www.devrant.io/) Console User Interface  

![Screenshot](https://raw.github.com/Jay9596/devRant_CUI/master/docs/images/devRant_CUI.png)

# Tabel of Content
1. [Installation](#installation)  
2. [Usage](#usage)  
3. [Todo](#todo)  

## Installation
 ### Dependancy
 - [Go](https://golang.org/)
 - [goRant](https://www.github.com/Jay9596/goRant)
 - [gocui](https://www.github.com/jroimartin/gocui)
 
 ### To complie and run  
  * Clone this repo  
    ` git clone https://github.com/Jay9596/devRant_CUI.git `  
    ` cd devRant_CUI `
  * Install Dependancies  
    ` go install github.com/Jay9596/goRant `  
    ` go install github.com/jroimartin/gocui `
  * Run main.go  
    ` go run main.go `
  * Build and Run  
    ` go build `  
    ` ./devRant_cui `
  * Install  
    ` go install github.com/Jay9596/devRant_CUI `
 
 ### Download
  [Download](https://github.com/Jay9596/devRant_CUI/releases/tag/v0.7)
 ### RUN
  This can run on most terminals ,(not tried many) ,but has some limitations on Windows command prompt when resizing.
  It can run on command prompt ,but for Windows use [Cmder](http://cmder.net/),[Hyper](https://github.com/zeit/hyper) ,etc.; for a better experience.  
## Usage
  #### UI 
  UI consists of 6 Views, 3 of which are editable.  
  The main view on left, logo, and output display are not interactive.  
  The __Input__ window is the main interactive View.  
  The __Sort__ window is used to set the sort method for fetching rants.  
  The __Limit__ window is used to set number of rants to display.  
  
  ### Views
  - Sort  
     It can take only three inputs: "algo", "recent", "top"  
     Any other input will not be accepted  
     "algo" is the default sort method  
       
     * To change, just type and press Enter  
  - Limit  
     It can take any int input  
     Default value is 5  
       
     * To change ,just type and press Enter  
  - Input  
     It is a dummy terminal with commands to run the CUI  
  ###### Note:   
       It will only accept input after ":/ >" ,sometimes it can break. Please try to keep inputs after ":/ >" until i can fix this.  
  
  ### Commands  
  - **help**  
   This command displays the help ,very useful for first timers.  
  - **rants**  
   This command fetched rants using the Sort and prints them.  
   It fetches 20 at a time, but prints only the Limit specified.  
  - **next** or **:n**  
   This command is used to go to next ,i.e next page of list.  
   This command prints the next rants ,but prints only the limited number.  
   This command will auto fetch rants ,or specifies if end of section to fetch something else.  
   Can be used after 'rants', 'stories' ,'weekly' ,'search' ,and 'collabs'  
  - **rant {int}**  
   This command is used to open a rant.  
   It need a number, the number is shown before a rant, example ">>0".  
  - **back** or **cd ..**  
   This command is used to go beack to the rants view after opening a rant.  
  - **profile {username}**  
   This command will fetch the user info and print it.  
   If there is an error, it will be shown in the output view above the input window  .
  - **search {term}**  
   This command will search for the term on devrant.io, fetch, and print the search result,i.e the rants.  
  - **weekly**  
   This command will fetch the rants posted today, and tagged weekly topic,i.e 'wk-00'.  
  - **collabs**  
   This commande will fetch the collanfs and print them.  
  - **stories**  
   This command will fetch and print stories.  
   For this, keep the limit small as stories are quiet long, and not many can fit.  
  - **exit** or **:q**  
   This command can be used to call exit, but exit won't close the program.  
   For the vim lovers of devRant :q is also available.  
   To close _"Ctrl + C"_ is used.  
  - **clear**  
   This will clear the input and output views.  
  - **clean**  
   This will clean the main view.  
  
## Todo
- [ ] Refactor Code
- [ ] Improve the input terminal  
- [ ] Add option to view content of user  

## Contact
Any and all criticism are appreciated.  

If you try this and find any bug, suggestion, or feature.  
Feel free to contact me.  
For appreciation via ++ := username: BYTE  
