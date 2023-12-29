package mj

import (
	"chatplus/core/types"
	logger2 "chatplus/logger"
	"chatplus/utils"
	discordgo "github.com/bg5t/mydiscordgo"
	"github.com/gorilla/websocket"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

// MidJourney 机器人

var logger = logger2.GetLogger()

type Bot struct {
	config  types.MidJourneyConfig
	bot     *discordgo.Session
	name    string
	service *Service
}

func NewBot(name string, proxy string, config types.MidJourneyConfig, service *Service) (*Bot, error) {
	bot, err := discordgo.New("Bot " + config.BotToken)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	// use CDN reverse proxy
	if config.UseCDN {
		discordgo.SetEndpointDiscord(config.DiscordAPI)
		discordgo.SetEndpointCDN(config.DiscordCDN)
		discordgo.SetEndpointStatus(config.DiscordAPI + "/api/v2/")
		bot.MjGateway = config.DiscordGateway + "/"
	} else { // use proxy
		discordgo.SetEndpointDiscord("https://discord.com")
		discordgo.SetEndpointCDN("https://cdn.discordapp.com")
		discordgo.SetEndpointStatus("https://discord.com/api/v2/")
		bot.MjGateway = "wss://gateway.discord.gg"

		if proxy != "" {
			proxy, _ := url.Parse(proxy)
			bot.Client = &http.Client{
				Transport: &http.Transport{
					Proxy: http.ProxyURL(proxy),
				},
			}
			bot.Dialer = &websocket.Dialer{
				Proxy: http.ProxyURL(proxy),
			}
		}

	}

	return &Bot{
		config:  config,
		bot:     bot,
		name:    name,
		service: service,
	}, nil
}

func (b *Bot) Run() error {
	b.bot.Identify.Intents = discordgo.IntentsAllWithoutPrivileged | discordgo.IntentsGuildMessages | discordgo.IntentMessageContent
	b.bot.AddHandler(b.messageCreate)
	b.bot.AddHandler(b.messageUpdate)

	logger.Infof("Starting MidJourney %s", b.name)
	err := b.bot.Open()
	if err != nil {
		logger.Errorf("Error opening Discord connection for %s, error: %v", b.name, err)
		return err
	}
	logger.Infof("Starting MidJourney %s successfully!", b.name)
	return nil
}

type TaskStatus string

const (
	Start    = TaskStatus("Started")
	Running  = TaskStatus("Running")
	Stopped  = TaskStatus("Stopped")
	Finished = TaskStatus("Finished")
)

type Image struct {
	URL      string `json:"url"`
	ProxyURL string `json:"proxy_url"`
	Filename string `json:"filename"`
	Width    int    `json:"width"`
	Height   int    `json:"height"`
	Size     int    `json:"size"`
	Hash     string `json:"hash"`
}

func (b *Bot) messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// ignore messages for other channels
	if m.GuildID != b.config.GuildId || m.ChannelID != b.config.ChanelId {
		return
	}
	// ignore messages for self
	if m.Author == nil || m.Author.ID == s.State.User.ID {
		return
	}

	logger.Debugf("CREATE: %s", utils.JsonEncode(m))
	var referenceId = ""
	if m.ReferencedMessage != nil {
		referenceId = m.ReferencedMessage.ID
	}
	if strings.Contains(m.Content, "(Waiting to start)") && !strings.Contains(m.Content, "Rerolling **") {
		// parse content
		req := CBReq{
			ChannelId:   m.ChannelID,
			MessageId:   m.ID,
			ReferenceId: referenceId,
			Prompt:      extractPrompt(m.Content),
			Content:     m.Content,
			Progress:    0,
			Status:      Start}
		b.service.Notify(req)
		return
	}

	b.addAttachment(m.ChannelID, m.ID, referenceId, m.Content, m.Attachments)
}

func (b *Bot) messageUpdate(s *discordgo.Session, m *discordgo.MessageUpdate) {
	// ignore messages for other channels
	if m.GuildID != b.config.GuildId || m.ChannelID != b.config.ChanelId {
		return
	}
	// ignore messages for self
	if m.Author == nil || m.Author.ID == s.State.User.ID {
		return
	}

	logger.Debugf("UPDATE: %s", utils.JsonEncode(m))

	var referenceId = ""
	if m.ReferencedMessage != nil {
		referenceId = m.ReferencedMessage.ID
	}
	if strings.Contains(m.Content, "(Stopped)") {
		req := CBReq{
			ChannelId:   m.ChannelID,
			MessageId:   m.ID,
			ReferenceId: referenceId,
			Prompt:      extractPrompt(m.Content),
			Content:     m.Content,
			Progress:    extractProgress(m.Content),
			Status:      Stopped}
		b.service.Notify(req)
		return
	}

	b.addAttachment(m.ChannelID, m.ID, referenceId, m.Content, m.Attachments)

}

func (b *Bot) addAttachment(channelId string, messageId string, referenceId string, content string, attachments []*discordgo.MessageAttachment) {
	progress := extractProgress(content)
	var status TaskStatus
	if progress == 100 {
		status = Finished
	} else {
		status = Running
	}
	for _, attachment := range attachments {
		if attachment.Width == 0 || attachment.Height == 0 {
			continue
		}
		image := Image{
			URL:      attachment.URL,
			Height:   attachment.Height,
			ProxyURL: attachment.ProxyURL,
			Width:    attachment.Width,
			Size:     attachment.Size,
			Filename: attachment.Filename,
			Hash:     extractHashFromFilename(attachment.Filename),
		}
		req := CBReq{
			ChannelId:   channelId,
			MessageId:   messageId,
			ReferenceId: referenceId,
			Image:       image,
			Prompt:      extractPrompt(content),
			Content:     content,
			Progress:    progress,
			Status:      status,
		}
		b.service.Notify(req)
		break // only get one image
	}
}

// extract prompt from string
func extractPrompt(input string) string {
	pattern := `\*\*(.*?)\*\*`
	re := regexp.MustCompile(pattern)
	matches := re.FindStringSubmatch(input)
	if len(matches) > 1 {
		return strings.TrimSpace(matches[1])
	}
	return ""
}

func extractProgress(input string) int {
	pattern := `\((\d+)\%\)`
	re := regexp.MustCompile(pattern)
	matches := re.FindStringSubmatch(input)
	if len(matches) > 1 {
		return utils.IntValue(matches[1], 0)
	}
	return 100
}

func extractHashFromFilename(filename string) string {
	if !strings.HasSuffix(filename, ".png") {
		return ""
	}

	index := strings.LastIndex(filename, "_")
	if index != -1 {
		return filename[index+1 : len(filename)-4]
	}
	return ""
}
