package main

import (
	"chadpole/commands"
	"chadpole/util"
	"log"
	"sync"

	"github.com/bwmarrin/discordgo"
)

/*
Declare a list of commands and command handlers here. The functionality
of the handlers will be in the 'commands' package
*/
var (
	commandsList = []*discordgo.ApplicationCommand{
		{
			Name:        "ribbit",
			Description: "Ribbit",
		},
		{
			Name:        "ribbit-embed",
			Description: "Ribbit, but embeded",
		},
		{
			Name:        "ribbit-button",
			Description: "Ribbit with buttons",
		},
		{
			Name:        "ribbit-btn-edit",
			Description: "Editable Ribbits",
		},
		{
			Name:        "odesli",
			Description: "Convert any music link into a song.link",
			Options:     commands.OdesliOptions,
		},
		{
			Name:        "ribbit-pagination",
			Description: "Example pagination",
		},
	}
	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"ribbit":            commands.RibbitHandler,
		"ribbit-embed":      commands.RibbitEmbedHandler,
		"ribbit-button":     commands.RibbitButtonHandler,
		"ribbit-btn-edit":   commands.RibbitBtnEditHandler,
		"odesli":            commands.OdesliHandler,
		"ribbit-pagination": commands.RibbitPaginationHandler,
	}
	componentHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"primary_test": commands.PrimaryTestBtnHandler,
		"rp_prev":      commands.RPPrevBtnHandler,
		"rp_next":      commands.RPNextBtnHandler,
	}
)

// Register all commands for the bot
func RegisterAllCommands(s *discordgo.Session) {
	registeredCommands := make([]*discordgo.ApplicationCommand, len(commandsList))
	for i, v := range commandsList {
		cmd, err := s.ApplicationCommandCreate(s.State.User.ID, s.State.Application.GuildID, v)
		if err != nil {
			log.Fatal(err)
		}
		registeredCommands[i] = cmd
	}
	log.Println("[INFO] All commands registered to Chadpole")
}

// Attach all of the handlers required for functionality
func SetupAllHandlers(s *discordgo.Session) {

	// Attach type-specific handlers
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		switch i.Type {
		// Attach command handlers for slash commands
		case discordgo.InteractionApplicationCommand:
			if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
				h(s, i)
			}
		// Attach component handlers, such as handlers for buttons
		case discordgo.InteractionMessageComponent:
			if h, ok := componentHandlers[i.MessageComponentData().CustomID]; ok {
				h(s, i)
			}
		}
	})

	// Allow the created message handler to monitor messages
	s.AddHandler(util.MessageCreateHandler)

	log.Println("[INFO] All handlers attached to Chadpole")
}

func SetupStatus(s *discordgo.Session) {
	s.UpdateStatusComplex(discordgo.UpdateStatusData{
		AFK: false,
		Activities: []*discordgo.Activity{
			{
				Name: "Frogging 🐸",
				Type: discordgo.ActivityTypeCompeting,
			},
		},
	})
}

func MainSetup(s *discordgo.Session) {
	/*
		RegisterAllCommands(chadpole)
		SetupAllHandlers(chadpole)
		SetupStatus(chadpole)

		Setup does not need to be done with goroutines at all here
		mainly exists to have an example of several functions running off of goroutines with waitgroups and waiting
	*/

	var setup_wg sync.WaitGroup

	setup_wg.Add(3)

	go func() {
		defer setup_wg.Done()
		RegisterAllCommands(s)
	}()

	go func() {
		defer setup_wg.Done()
		SetupAllHandlers(s)
	}()

	go func() {
		defer setup_wg.Done()
		SetupStatus(s)
	}()

	setup_wg.Wait()

	// Start frog API on its own goroutine
	go util.StartFrogAPI()
}
