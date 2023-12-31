# chadpole

Discord bot using the discordgo library that I am using to experiment with golang. This bot is not meant to be used in a serious capacity. Features are implemented because I think they are cool and/or useful to learn from.

## Usage

```
$ DC_TOKEN=<discord bot token> ./chadpole
```

## Contents

### Main package
| File | Purpose |
|------|---------|
| main.go | Creates the connection between the bot and Discord. Invokes setup. |
| setup.go | Adds all commands and handlers at run time. |

### Util package
| File | Purpose |
|------|---------|
| msg.go | Simple message listener/handler. Not intended to be used to listen for legacy prefix commands. |
| frog_api.go | Simple API that serves a JSON response of 5 frogs |
| ocr.go | Functions that allow the bot to scan an image against Google's Tesseract OCR |

### Widgets package
This package is designed to contain portable code to use in other projects(maybe)
| File | Purpose |
|------|---------|
| paginator.go | Contains all structs and functions to create a fully functional pagination message. See `ribbit-pagination.go` for example usage |
### Commands package
This bot uses slash commands exclusively.
| File | Purpose |
|------|---------|
| odesli.go | Takes any link and uses the Odesli API to return a song.link |
| ribbit.go | Frog themed ping-pong impl. |
| ribbit-embed.go | Ping-pong impl. with an embed |
| ribbit-button.go | Ping-pong impl. with some buttons |
| ribbit-btn-edit.go | Message with buttons that when pressed edit the message or do something |
| ribbit-pagination.go | Example impl. of a paginated embed using discord buttons. To my knowledge this is the only public example of pagination in discordgo with buttons. This one also handles multiple pagination sessions (in a slightly ham-fisted manner). Pulls data from frog API. |

## Documentation Used
- https://0x2142.com/how-to-discordgo-bot/
- [Odesli API Docs](https://linktree.notion.site/API-d0ebe08a5e304a55928405eb682f6741)
- [Discordgo Docs](https://pkg.go.dev/github.com/bwmarrin/discordgo)
- [Discordgo Examples](https://github.com/bwmarrin/discordgo/tree/master/examples)