/*
Exmaple pagination implementation with buttons.
Inspired from https://github.com/Necroforger/dgwidgets/blob/master/paginator.go
As well as from my own personal work at: https://github.com/CooldudePUGS/Spectre/blob/90463d95839caf6a8551cf6fa91ac2dc952101b5/cogs/ModSearch.py

A lot of variables here start with "TS" which stands for thunderstore, as this was originally
supposed to be an implementation of thunderstore mod searching. Might save that for another module
and get extra confusing lol.

KNOWN ISSUES:
  - Logging information will repeat several times depending on how many times pagination has been called.
    Actual functionality of this is not affected since we make sure there are not multiple paginations running at once
    However this is pretty far from ideal
  - The above has been solved by generating a random 8 character prefix for the handler CustomIDs.
    This will ensure that different handlers are assigned to the pagination every time.
    The drawback of this is that there are now essentially "dead handlers" attached to the bot. Maybe
    the garbage collector deals with it at some point. The fuck do I know. At least it works now
*/
package commands

import (
	"crypto/rand"
	"fmt"
	"log"
	"sync"

	"github.com/bwmarrin/discordgo"
)

/*
Need to create a class/struct for the pagination view.
Add in any variables you will need in here. Some notable things
would be variables to hold JSON values if you're displaying a search
from an API for example
*/
type PaginationView struct {
	sync.Mutex
	index           int
	embedtitle      string
	embeddesc       string
	handerPrefix    string
	pageBtnHandlers map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate)
	exampleData     []FrogExample // Create a slice of the struct where our data will be stored.
	currentPage     int
}

type FrogExample struct {
	frogName string
	frogDesc string
}

/*
General setup housekeeping should go here
*/
func (p *PaginationView) Setup(s *discordgo.Session, i *discordgo.InteractionCreate) {
	p.Lock()
	p.handerPrefix = p.GenPrefix() // Generate prefix to uniquely identify paginator controls
	p.Unlock()
	// Assign handlers to their respective CustomIDs
	p.pageBtnHandlers[p.handerPrefix+"pg_next"] = p.TSNextBtnHandler
	p.pageBtnHandlers[p.handerPrefix+"pg_prev"] = p.TSPrevBtnHandler
	p.pageBtnHandlers[p.handerPrefix+"pg_stop"] = p.TSStopBtnHandler
	p.pageBtnHandlers[p.handerPrefix+"pg_done"] = p.TSDoneBtnHandler
	p.pageBtnHandlers[p.handerPrefix+"pg_first"] = p.TSFirstBtnHandler
	p.pageBtnHandlers[p.handerPrefix+"pg_last"] = p.TSLastBtnHandler
	// Add the handlers to the bot
	p.TSAddHandlers(s, i)
	p.currentPage = p.index + 1
}

/*
Generate a unique random prefix to identify a unique paginator's button CustomIDs
Unfortunately this is how I decided to work around being unable to cleanly add
buttons to a specific instance rather than globally
*/
func (p *PaginationView) GenPrefix() string {
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	length := 8
	prefix := make([]byte, length)
	_, err := rand.Read(prefix)
	if err != nil {
		log.Println("[WARN] Could not generate random prefix. Things might get weird.")
		return ""
	}
	for i, b := range prefix {
		prefix[i] = charset[b%byte(len(charset))]
	}
	log.Printf("[INFO] Using prefix %s", prefix)
	return string(prefix)
}

/*
Creates the embed. This function is called every time there is an update to the message
*/
func (p *PaginationView) CreateEmbed() []*discordgo.MessageEmbed {
	// Pull data from Data slice based on the index which is modified with the buttons
	p.embedtitle = fmt.Sprintf("%s", p.exampleData[p.index].frogName)
	p.embeddesc = fmt.Sprintf("%s", p.exampleData[p.index].frogDesc)

	embed := []*discordgo.MessageEmbed{
		{
			Title:       p.embedtitle,
			Description: p.embeddesc,
			Footer: &discordgo.MessageEmbedFooter{
				Text: fmt.Sprintf("Page: %d/%d", p.currentPage, len(p.exampleData)),
			},
		},
	}

	return embed
}

/*
Since we need to create these buttons multiple times in the code, throw them in a function to improve readability
*/
func (p *PaginationView) CreateBtns() []discordgo.MessageComponent {
	component := []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.Button{
					Label:    "<<",
					Style:    discordgo.PrimaryButton,
					CustomID: p.handerPrefix + "pg_first",
				},
				discordgo.Button{
					Label:    "<",
					Style:    discordgo.PrimaryButton,
					CustomID: p.handerPrefix + "pg_prev",
				},
				discordgo.Button{
					Label:    ">",
					Style:    discordgo.PrimaryButton,
					CustomID: p.handerPrefix + "pg_next",
				},
				discordgo.Button{
					Label:    ">>",
					Style:    discordgo.PrimaryButton,
					CustomID: p.handerPrefix + "pg_last",
				},
			},
		},
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.Button{
					Label:    "Done",
					Style:    discordgo.SuccessButton,
					CustomID: p.handerPrefix + "pg_done",
				},
				discordgo.Button{
					Label:    "Stop",
					Style:    discordgo.DangerButton,
					CustomID: p.handerPrefix + "pg_stop",
				},
			},
		},
	}

	return component
}

/*
Function is called once at the start to create the first instance of the pagination.
Importantly, there is a check to make sure that there are no running paginations currently

Handers are added several times because currently they are set to add once to avoid duplication
*/
func (p *PaginationView) SendMessage(s *discordgo.Session, i *discordgo.InteractionCreate) {
	p.Setup(s, i)               // Call setup function first so the handler prefix can be generated and all handlers added
	p.currentPage = p.index + 1 // Display page number normally
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			// You can make searches only visible to the invoker with the following:
			// (Note: The stop button won't work)
			// Flags:      discordgo.MessageFlagsEphemeral,
			Embeds:     p.CreateEmbed(),
			Components: p.CreateBtns(),
		},
	})
	log.Printf("[INFO] Sent initial pagination view for %s", p.handerPrefix)
}

/*
Updates the message every time next or prev is pressed
*/
func (p *PaginationView) UpdateMessage(s *discordgo.Session, i *discordgo.InteractionCreate) {
	p.currentPage = p.index + 1 // Display page number normally
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			// You can make searches only visible to the invoker with the following:
			// (Note: The stop button won't work)
			// Flags:      discordgo.MessageFlagsEphemeral
			Embeds:     p.CreateEmbed(),
			Components: p.CreateBtns(),
		},
	})
	log.Printf("[INFO] Updated pagination %s", p.handerPrefix)
}

/*
Increments the index value. In a proper situation you would use this index value to get more data from
a data structure
*/
func (p *PaginationView) TSNextBtnHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	p.Lock()
	defer p.Unlock()
	p.index = (p.index + 1) % len(p.exampleData) // Circular pagination
	log.Printf("[INFO] Pagination %s data incremented", p.handerPrefix)
	p.UpdateMessage(s, i)
}

func (p *PaginationView) TSPrevBtnHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	p.Lock()
	defer p.Unlock()
	if p.index == 0 { // Prevent running out of bounds. Function as a "last" button if index is at 0
		p.index = len(p.exampleData) - 1
		p.UpdateMessage(s, i)
	} else {
		p.index = (p.index - 1) % len(p.exampleData) // Circular pagination
		log.Printf("[INFO] Pagination %s data decremented", p.handerPrefix)
		p.UpdateMessage(s, i)
	}
}

func (p *PaginationView) TSFirstBtnHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	log.Printf("[INFO] Pagination %s data set to first", p.handerPrefix)
	p.index = 0
	p.UpdateMessage(s, i)
}

func (p *PaginationView) TSLastBtnHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	log.Printf("[INFO] Pagination %s data set to last index", p.handerPrefix)
	p.index = len(p.exampleData) - 1
	p.UpdateMessage(s, i)
}

/*
Deletes the message if pressed, also resets running status
*/
func (p *PaginationView) TSStopBtnHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	/*
		Holy fucking shit
		this is how you delete your own fucking message
		why did it take so long to figure this out?
		fucks sake

		Do note: you can get the Channel ID of the Interaction with i.Message.ChannelID
		and it's own message id with i.Message.ID

	*/
	log.Printf("[INFO] Pagination %s stopped", p.handerPrefix)
	s.ChannelMessageDelete(i.Message.ChannelID, i.Message.ID)
	log.Printf("[INFO] Pagination %s removed from channel", p.handerPrefix)
}

/*
Resets running status but instead of deleting the message, it just gets rid of the buttons so the
embed is permanent in the channel
*/
func (p *PaginationView) TSDoneBtnHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Embeds: p.CreateEmbed(),
		},
	})
	log.Printf("[INFO] Pagination %s stopped, embed remains in channel", p.handerPrefix)
}

/*
Add all the handlers for the buttons. Notably this is a direct copy-paste implementation from setup.go
Definitely not the best solution. Should probably change this later.
I am fairly confident (lol) that cleaner code here would not help the issue with concurrent paginations.
*/
func (p *PaginationView) TSAddHandlers(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		switch i.Type {
		// Attach component handlers, such as handlers for buttons
		case discordgo.InteractionMessageComponent:
			if h, ok := p.pageBtnHandlers[i.MessageComponentData().CustomID]; ok {
				h(s, i)
			}
		}
	})

	log.Printf("[INFO] Added pagination handlers to %s", p.handerPrefix)
}

/*
Entrypoint for the pagination system
*/
func RibbitPaginationHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	new_pagination := PaginationView{
		index:           0,
		pageBtnHandlers: make(map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate)), // Must create the map for the handler CustomIDs
		exampleData: []FrogExample{ // Declare our data here outside of the struct
			{"Red-eyed Tree Frog", "Known for its vibrant green coloration, the red-eyed tree frog is native to the rainforests of Central and South America. It has striking red eyes and orange feet, and it spends most of its life in trees near water bodies."},
			{"Poison Dart Frog", "Poison dart frogs are small and brightly colored frogs found in Central and South America. They are known for their toxic skin secretions, which have been used by indigenous people to poison the tips of blowgun darts for hunting."},
			{"African Clawed Frog", "Native to sub-Saharan Africa, the African clawed frog is an aquatic species that lacks tongue and teeth. Its distinctive \"claws\" on its hind feet are used for digging and defense."},
			{"Goliath Frog", "The Goliath frog is the largest frog species, found in Cameroon and Equatorial Guinea in Africa. It can grow to be over a foot long and has a unique appearance with robust body proportions."},
			{"Wood Frog", "Wood frogs are found in North America and are known for their remarkable adaptation to cold environments. They can survive freezing temperatures by entering a state of suspended animation and then thawing back to life when temperatures rise."},
		},
	}
	log.Println("[INFO] New pagination created")
	new_pagination.SendMessage(s, i) // Send the message, functions as the entrypoint for the pagination view

}
