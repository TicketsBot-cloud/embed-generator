package bot

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/merlinfuchs/discordgo"
	"github.com/merlinfuchs/embed-generator/embedg-server/actions"
	"github.com/merlinfuchs/embed-generator/embedg-server/db/postgres/pgmodel"
	"github.com/merlinfuchs/embed-generator/embedg-server/util"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

var MessageIDRe = regexp.MustCompile("https?://(?:canary\\.|ptb\\.)?discord\\.com/channels/[0-9]+/([0-9]+)/([0-9]+)")

func (b *Bot) RegisterCommand() error {
	_, err := b.Session.ApplicationCommandBulkOverwrite(viper.GetString("discord.client_id"), "", []*discordgo.ApplicationCommand{
		{
			Name:        "help",
			Description: "Show help",
			IntegrationTypes: &[]discordgo.ApplicationIntegrationType{
				discordgo.ApplicationIntegrationGuildInstall,
				discordgo.ApplicationIntegrationUserInstall,
			},
		},
		{
			Name:        "invite",
			Description: "Invite the Embed Generator bot to your server",
			IntegrationTypes: &[]discordgo.ApplicationIntegrationType{
				discordgo.ApplicationIntegrationGuildInstall,
			},
		},
		{
			Name:        "website",
			Description: "Open the Embed Generator website",
			IntegrationTypes: &[]discordgo.ApplicationIntegrationType{
				discordgo.ApplicationIntegrationGuildInstall,
				discordgo.ApplicationIntegrationUserInstall,
			},
		},
		{
			Name:        "format",
			Description: "Get the API format for mentions, channels, roles, & custom emojis",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "text",
					Description: "Get the API format for a text with multiple mentions, channels, & custom emojis",
					Options: []*discordgo.ApplicationCommandOption{
						{
							Type:        discordgo.ApplicationCommandOptionString,
							Name:        "text",
							Description: "The text that you want to format (usually containing mentions or custom emojis)",
							Required:    true,
						},
					},
				},
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "user",
					Description: "Get the API format for mentioning a user",
					Options: []*discordgo.ApplicationCommandOption{
						{
							Type:        discordgo.ApplicationCommandOptionUser,
							Name:        "user",
							Description: "The user you want to mention",
							Required:    true,
						},
					},
				},
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "channel",
					Description: "Get the API format for mentioning a channel",
					Options: []*discordgo.ApplicationCommandOption{
						{
							Type:        discordgo.ApplicationCommandOptionChannel,
							Name:        "channel",
							Description: "The channel you want to mention",
							Required:    true,
						},
					},
				},
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "role",
					Description: "Get the API format for mentioning a role",
					Options: []*discordgo.ApplicationCommandOption{
						{
							Type:        discordgo.ApplicationCommandOptionRole,
							Name:        "role",
							Description: "The role you want to mention",
							Required:    true,
						},
					},
				},
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "emoji",
					Description: "Get the API format for a standard or custom emoji",
					Options: []*discordgo.ApplicationCommandOption{
						{
							Type:        discordgo.ApplicationCommandOptionString,
							Name:        "emoji",
							Description: "The standard or custom emoji you want to use",
							Required:    true,
						},
					},
				},
			},
			IntegrationTypes: &[]discordgo.ApplicationIntegrationType{
				discordgo.ApplicationIntegrationGuildInstall,
				discordgo.ApplicationIntegrationUserInstall,
			},
		},
		{
			Name:         "image",
			Description:  "Get the image URL for different entities",
			DMPermission: util.Ptr(false),
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "avatar",
					Description: "Get the avatar URL for a user",
					Options: []*discordgo.ApplicationCommandOption{
						{
							Type:        discordgo.ApplicationCommandOptionUser,
							Name:        "user",
							Description: "The user you want to get the avatar for",
							Required:    true,
						},
						{
							Type:        discordgo.ApplicationCommandOptionBoolean,
							Name:        "static",
							Description: "Whether animated avatars should be converted to static images",
						},
					},
				},
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "icon",
					Description: "Get the icon URL for this server",
					Options: []*discordgo.ApplicationCommandOption{
						{
							Type:        discordgo.ApplicationCommandOptionBoolean,
							Name:        "static",
							Description: "Whether animated icons should be converted to static images",
						},
					},
				},
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "emoji",
					Description: "Get the image URL for a custom or standard emoji",
					Options: []*discordgo.ApplicationCommandOption{
						{
							Type:        discordgo.ApplicationCommandOptionString,
							Name:        "emoji",
							Description: "The custom emoji you want the image URL for",
							Required:    true,
						},
						{
							Type:        discordgo.ApplicationCommandOptionBoolean,
							Name:        "static",
							Description: "Whether animated emojis should be converted to static images",
						},
					},
				},
			},
			IntegrationTypes: &[]discordgo.ApplicationIntegrationType{
				discordgo.ApplicationIntegrationGuildInstall,
			},
		},
		{

			Name:        "message",
			Description: "Get JSON for or restore a message on Embed Generator",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "restore",
					Description: "Restore a message on Embed Generator",
					Options: []*discordgo.ApplicationCommandOption{
						{
							Type:        discordgo.ApplicationCommandOptionString,
							Name:        "message_id_or_url",
							Description: "ID or URL of the message you want to restore",
							Required:    true,
						},
					},
				},
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "dump",
					Description: "Get the JSON code for a message",
					Options: []*discordgo.ApplicationCommandOption{
						{
							Type:        discordgo.ApplicationCommandOptionString,
							Name:        "message_id_or_url",
							Description: "ID or URL of the message you want to restore",
							Required:    true,
						},
					},
				},
			},
			IntegrationTypes: &[]discordgo.ApplicationIntegrationType{
				discordgo.ApplicationIntegrationGuildInstall,
			},
		},
		{
			Type: discordgo.MessageApplicationCommand,
			Name: "Restore Message",
			IntegrationTypes: &[]discordgo.ApplicationIntegrationType{
				discordgo.ApplicationIntegrationGuildInstall,
				discordgo.ApplicationIntegrationUserInstall,
			},
		},
		{
			Type: discordgo.MessageApplicationCommand,
			Name: "Dump Message",
			IntegrationTypes: &[]discordgo.ApplicationIntegrationType{
				discordgo.ApplicationIntegrationGuildInstall,
				discordgo.ApplicationIntegrationUserInstall,
			},
		},
		{
			Type: discordgo.UserApplicationCommand,
			Name: "Avatar Url",
			IntegrationTypes: &[]discordgo.ApplicationIntegrationType{
				discordgo.ApplicationIntegrationGuildInstall,
				discordgo.ApplicationIntegrationUserInstall,
			},
		},
		{
			Name:                     "embed",
			Description:              "Create an embed message",
			DMPermission:             util.Ptr(false),
			DefaultMemberPermissions: util.Ptr(int64(discordgo.PermissionManageWebhooks)),
			IntegrationTypes: &[]discordgo.ApplicationIntegrationType{
				discordgo.ApplicationIntegrationGuildInstall,
			},
		},
	})
	return err
}

func (b *Bot) handleCommandInteraction(s *discordgo.Session, i *discordgo.Interaction, data discordgo.ApplicationCommandInteractionData) error {
	switch data.Name {
	case "invite":
		return b.handleInviteCommand(s, i, data)
	case "website":
		return b.handleWebsiteCommand(s, i, data)
	case "help":
		return b.handleHelpCommand(s, i, data)
	case "format":
		return b.handleFormatCommand(s, i, data)
	case "image":
		return b.handleImageCommand(s, i, data)
	case "message":
		return b.handleMessageCommand(s, i, data)
	case "Restore Message":
		return b.handleRestoreContextCommand(s, i, data)
	case "Dump Message":
		return b.handleJSONContextCommand(s, i, data)
	case "Avatar Url":
		return b.handleAvatarUrlContextCommand(s, i, data)
	case "embed":
		return b.handleEmbedCommand(s, i, data)
	}
	return nil
}

func (b *Bot) handleHelpCommand(s *discordgo.Session, i *discordgo.Interaction, data discordgo.ApplicationCommandInteractionData) error {
	return helpResponse(s, i)
}

func (b *Bot) handleInviteCommand(s *discordgo.Session, i *discordgo.Interaction, data discordgo.ApplicationCommandInteractionData) error {
	return helpResponse(s, i)
}

func (b *Bot) handleWebsiteCommand(s *discordgo.Session, i *discordgo.Interaction, data discordgo.ApplicationCommandInteractionData) error {
	return helpResponse(s, i)
}

func helpResponse(s *discordgo.Session, i *discordgo.Interaction) error {
	return fancyResponse(s, i, "**The best way to generate rich embed messages for your Discord Server!**\n\nhttps://www.youtube.com/watch?v=DnFP0MRJPIg", []*discordgo.MessageEmbed{}, []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.Button{
					Style: discordgo.LinkButton,
					Label: "Website",
					URL:   "https://message.style",
				},
				discordgo.Button{
					Style: discordgo.LinkButton,
					Label: "Invite Bot",
					URL:   util.BotInviteURL(),
				},
				discordgo.Button{
					Style: discordgo.LinkButton,
					Label: "Discord Server",
					URL:   viper.GetString("links.discord"),
				},
			},
		},
	})
}

func (b *Bot) handleFormatCommand(s *discordgo.Session, i *discordgo.Interaction, data discordgo.ApplicationCommandInteractionData) error {
	subCMD := data.Options[0]

	switch subCMD.Name {
	case "text":
		value := subCMD.Options[0].StringValue()
		return textResponse(s, i, fmt.Sprintf("API format for the provided text: ```%s```", value))
	case "user":
		user := subCMD.Options[0].UserValue(nil)
		return textResponse(s, i, fmt.Sprintf("API format for <@%s>: ```<@%s>```", user.ID, user.ID))
	case "channel":
		channel := subCMD.Options[0].ChannelValue(nil)
		return textResponse(s, i, fmt.Sprintf("API format for <#%s>: ```<#%s>```", channel.ID, channel.ID))
	case "role":
		role := subCMD.Options[0].RoleValue(nil, i.GuildID)
		return textResponse(s, i, fmt.Sprintf("API format for <@&%s>: ```<@&%s>```", role.ID, role.ID))
	case "emoji":
		emoji := subCMD.Options[0].StringValue()
		// TODO
		return textResponse(s, i, fmt.Sprintf("API format for %s: ```%s```", emoji, emoji))
	}

	return nil
}

func (b *Bot) handleImageCommand(s *discordgo.Session, i *discordgo.Interaction, data discordgo.ApplicationCommandInteractionData) error {
	subCMD := data.Options[0]

	makeStatic := func(url string, option int) string {
		if len(subCMD.Options) <= option {
			return url
		}

		static := subCMD.Options[option].BoolValue()
		if static {
			return strings.Replace(url, ".gif", ".png", 1)
		}

		return url
	}

	switch subCMD.Name {
	case "avatar":
		userID := subCMD.Options[0].UserValue(nil).ID
		user := data.Resolved.Users[userID]

		avatarURL := makeStatic(user.AvatarURL("1024"), 1)
		return imageUrlResponse(s, i, avatarURL)
	case "icon":
		guild, err := b.State.Guild(i.GuildID)
		if err != nil {
			return err
		}
		if guild.Icon == "" {
			return textResponse(s, i, "This server has no icon.")
		}
		iconURL := makeStatic(guild.IconURL("1024"), 1)
		return imageUrlResponse(s, i, iconURL)
	case "emoji":
		// emoji := subCMD.Options[0].StringValue()
		// TODO: get emoji id from regex
	}

	return nil
}

func (b *Bot) handleMessageCommand(s *discordgo.Session, i *discordgo.Interaction, data discordgo.ApplicationCommandInteractionData) error {
	subCMD := data.Options[0]

	messageID := subCMD.Options[0].StringValue()

	match := MessageIDRe.FindStringSubmatch(messageID)
	if match != nil {
		messageID = match[2]
	}

	message, err := s.ChannelMessage(i.ChannelID, messageID)
	if err != nil {
		if util.IsDiscordRestErrorCode(err, discordgo.ErrCodeUnknownMessage) {
			return textResponse(s, i, "Message not found.")
		}
		log.Error().Err(err).Msg("Failed to get message")
		return textResponse(s, i, "Failed to get message.")
	}

	components, err := b.ActionParser.UnparseMessageComponents(message.Components)
	if err != nil {
		log.Error().Err(err).Msg("Failed to unparse message components")
		return textResponse(s, i, "Failed to unparse message components.")
	}

	actionSets, err := b.ActionParser.RetrieveActionsForMessage(context.TODO(), messageID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to retrieve actions for message")
		return textResponse(s, i, "Failed to retrieve actions for message.")
	}

	messageDump, err := json.MarshalIndent(actions.MessageWithActions{
		Username:   message.Author.Username,
		AvatarURL:  message.Author.AvatarURL("1024"),
		Content:    message.Content,
		Embeds:     message.Embeds,
		Components: components,
		Actions:    actionSets,
	}, "", "  ")
	if err != nil {
		return err
	}

	switch subCMD.Name {
	case "restore":
		msg, err := b.pg.Q.InsertSharedMessage(context.TODO(), pgmodel.InsertSharedMessageParams{
			ID:        util.UniqueID(),
			CreatedAt: time.Now().UTC(),
			ExpiresAt: time.Now().UTC().Add(time.Hour * 24),
			Data:      messageDump,
		})
		if err != nil {
			log.Error().Err(err).Msg("Failed to insert shared message")
			return textResponse(s, i, "Failed to create shared message.")
		}

		url := fmt.Sprintf("%s/editor/share/%s", viper.GetString("app.public_url"), msg.ID)
		return textResponse(s, i, fmt.Sprintf("Click this link to restore the message: [message.style](<%s>)", url))
	case "dump":
		paste, err := util.CreateVaultBinPaste(string(messageDump), "json")
		if err != nil {
			log.Error().Err(err).Msg("Failed to create vaultb.in paste")
			return textResponse(s, i, "Failed to create vaultb.in paste.")
		}

		return textResponse(s, i, fmt.Sprintf("You can find the JSON code here: <%s>", paste.URL()))
	}

	return nil
}

func (b *Bot) handleRestoreContextCommand(s *discordgo.Session, i *discordgo.Interaction, data discordgo.ApplicationCommandInteractionData) error {
	messageID := data.TargetID
	message := data.Resolved.Messages[messageID]

	components, err := b.ActionParser.UnparseMessageComponents(message.Components)
	if err != nil {
		return fmt.Errorf("Failed to unparse message components: %w", err)
	}

	actionSets, err := b.ActionParser.RetrieveActionsForMessage(context.TODO(), messageID)
	if err != nil {
		return fmt.Errorf("Failed to retrieve actions for message: %w", err)
	}

	messageDump, err := json.MarshalIndent(actions.MessageWithActions{
		Username:   message.Author.Username,
		AvatarURL:  message.Author.AvatarURL("1024"),
		Content:    message.Content,
		Embeds:     message.Embeds,
		Components: components,
		Actions:    actionSets,
	}, "", "  ")
	if err != nil {
		return err
	}

	msg, err := b.pg.Q.InsertSharedMessage(context.TODO(), pgmodel.InsertSharedMessageParams{
		ID:        util.UniqueID(),
		CreatedAt: time.Now().UTC(),
		ExpiresAt: time.Now().UTC().Add(time.Hour * 24),
		Data:      messageDump,
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to insert shared message")
		return textResponse(s, i, "Failed to create shared message.")
	}

	url := fmt.Sprintf("%s/editor/share/%s", viper.GetString("app.public_url"), msg.ID)
	return textResponse(s, i, fmt.Sprintf("Click this link to restore the message: [message.style](<%s>)", url))
}

func (b *Bot) handleJSONContextCommand(s *discordgo.Session, i *discordgo.Interaction, data discordgo.ApplicationCommandInteractionData) error {
	messageID := data.TargetID
	message := data.Resolved.Messages[messageID]

	components, err := b.ActionParser.UnparseMessageComponents(message.Components)
	if err != nil {
		log.Error().Err(err).Msg("Failed to unparse message components")
		return textResponse(s, i, "Failed to unparse message components.")
	}

	actionSets, err := b.ActionParser.RetrieveActionsForMessage(context.TODO(), messageID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to retrieve actions for message")
		return textResponse(s, i, "Failed to retrieve actions for message.")
	}

	messageDump, err := json.MarshalIndent(actions.MessageWithActions{
		Username:   message.Author.Username,
		AvatarURL:  message.Author.AvatarURL("1024"),
		Content:    message.Content,
		Embeds:     message.Embeds,
		Components: components,
		Actions:    actionSets,
	}, "", "  ")
	if err != nil {
		return err
	}

	paste, err := util.CreateVaultBinPaste(string(messageDump), "json")
	if err != nil {
		log.Error().Err(err).Msg("Failed to create vaultb.in paste")
		return textResponse(s, i, "Failed to create vaultb.in paste.")
	}

	return textResponse(s, i, fmt.Sprintf("You can find the JSON code here: <%s>", paste.URL()))
}

func (b *Bot) handleAvatarUrlContextCommand(s *discordgo.Session, i *discordgo.Interaction, data discordgo.ApplicationCommandInteractionData) error {
	userId := data.TargetID
	user := data.Resolved.Users[userId]

	return imageUrlResponse(s, i, user.AvatarURL("1024"))
}

func (b *Bot) handleEmbedCommand(s *discordgo.Session, i *discordgo.Interaction, data discordgo.ApplicationCommandInteractionData) error {
	return fancyResponse(s, i, "If you want to have more options to customize your message go to [message.style](<https://message.style/app>)!", []*discordgo.MessageEmbed{}, embedEditComponent())
}

func textResponse(s *discordgo.Session, i *discordgo.Interaction, content string) error {
	return s.InteractionRespond(i, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
}

func imageUrlResponse(s *discordgo.Session, i *discordgo.Interaction, url string) error {
	return s.InteractionRespond(i, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags: discordgo.MessageFlagsEphemeral,
			Embeds: []*discordgo.MessageEmbed{
				{
					Description: url,
					Image: &discordgo.MessageEmbedImage{
						URL: url,
					},
				},
			},
		},
	})
}

func fancyResponse(s *discordgo.Session, i *discordgo.Interaction, content string, embeds []*discordgo.MessageEmbed, components []discordgo.MessageComponent) error {
	return s.InteractionRespond(i, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content:    content,
			Flags:      discordgo.MessageFlagsEphemeral,
			Embeds:     embeds,
			Components: components,
		},
	})
}

func embedEditComponent() []discordgo.MessageComponent {
	return []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.Button{
					Style:    discordgo.PrimaryButton,
					Label:    "Set Author",
					CustomID: "embed:author",
				},
				discordgo.Button{
					Style:    discordgo.PrimaryButton,
					Label:    "Set Title",
					CustomID: "embed:title",
				},
				discordgo.Button{
					Style:    discordgo.PrimaryButton,
					Label:    "Set Description",
					CustomID: "embed:description",
				},
				discordgo.Button{
					Style:    discordgo.PrimaryButton,
					Label:    "Set Color",
					CustomID: "embed:color",
				},
			},
		},
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.Button{
					Style:    discordgo.PrimaryButton,
					Label:    "Set Image",
					CustomID: "embed:image",
				},
				discordgo.Button{
					Style:    discordgo.PrimaryButton,
					Label:    "Set Thumbnail",
					CustomID: "embed:thumbnail",
				},
				discordgo.Button{
					Style:    discordgo.PrimaryButton,
					Label:    "Set Footer",
					CustomID: "embed:footer",
				},
			},
		},
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.Button{
					Style:    discordgo.DangerButton,
					Label:    "Cancel",
					CustomID: "embed:cancel",
				},
				discordgo.Button{
					Style:    discordgo.SuccessButton,
					Label:    "Submit",
					CustomID: "embed:submit",
				},
			},
		},
	}
}
