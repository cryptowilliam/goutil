package telegram

import (
	"encoding/json"
	"fmt"
	"github.com/cryptowilliam/goutil/basic/gerrors"
	"github.com/cryptowilliam/goutil/basic/glog"
	"github.com/cryptowilliam/goutil/container/gstring"
	"github.com/cryptowilliam/goutil/net/ghttp"
	"github.com/cryptowilliam/goutil/sys/gtime"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"net/http"
	"time"
)

// Privacy mode与接收消息
// 如果需要机器人只接收/***的command消息，则必须在BotFather那里将其Privacy mode打开，默认情况就是打开的
// 如果需要机器人接收全部消息，则必须将其关闭
// 而且，修改设置前加过的群是不生效的，需要将机器人踢了重新加入才能生效

type Bot struct {
	key        string
	proxy      string
	httpClient *http.Client
	api        *tgbotapi.BotAPI
	logger     glog.Interface
}

type User tgbotapi.User

type Msg tgbotapi.Message

type MsgType int

const (
	MsgTyepJoin MsgType = iota + 0
	MsgTypeInvite
	MsgTypeText
	MsgTypeImage
	MsgTypeAudio
	MsgTypeVideo
	MsgTypeFile
	MsgTypeLocation
	MsgTypeUnknown
)

type ChatType int

const (
	ChatTypeSuperGroup ChatType = iota + 0
	ChatTypeChannel
	ChatTypePrivate // 一对一私聊
	ChatTypeUnknown
)

type Group tgbotapi.Chat

func (m *Msg) MsgType() MsgType {
	if m.NewChatMembers != nil {
		if m.From.ID == (*m.NewChatMembers)[0].ID {
			return MsgTyepJoin
		} else {
			return MsgTypeInvite
		}
	}

	if len(m.Text) > 0 {
		return MsgTypeText
	}

	if m.Photo != nil && len(*m.Photo) > 0 {
		return MsgTypeImage
	}

	if m.Voice != nil && len(m.Voice.FileID) > 0 {
		return MsgTypeAudio
	}

	if m.Document != nil && len(m.Document.FileID) > 0 && m.Document.Thumbnail != nil && len(m.Document.Thumbnail.FileID) > 0 {
		return MsgTypeVideo
	}

	if m.Document != nil && len(m.Document.FileID) > 0 {
		return MsgTypeFile
	}

	if m.Location != nil {
		return MsgTypeLocation
	}

	return MsgTypeUnknown
}

func (m *Msg) ToJson() (string, error) {
	b, err := json.MarshalIndent(m, "", "\t")
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// 消息来自哪个群组
func (m *Msg) FromGroup() (isMsgFromGroup bool, group Group) {
	chatType := m.ChatType()
	if chatType == ChatTypeSuperGroup || chatType == ChatTypeChannel {
		return true, Group(*m.Chat)
	}
	return false, Group{}
}

func (m *Msg) FromUser() User {
	return User(*m.From)
}

func (m *Msg) ChatType() ChatType {
	if m.Chat.Type == "supergroup" {
		return ChatTypeSuperGroup
	}
	if m.Chat.Type == "channel" {
		return ChatTypeChannel
	}
	if m.Chat.Type == "private" {
		return ChatTypePrivate
	}
	fmt.Println(m.Chat.Type)
	return ChatTypeUnknown
}

type MsgChan <-chan Msg

func NewBot(key, proxy string, logger glog.Interface) (*Bot, error) {
	b := Bot{key: key, proxy: proxy}
	err := error(nil)

	if len(proxy) > 0 {
		b.httpClient = http.DefaultClient
		if err := ghttp.SetProxy(b.httpClient, proxy); err != nil {
			return nil, err
		}
		b.api, err = tgbotapi.NewBotAPIWithClient(key, tgbotapi.APIEndpoint, b.httpClient)
		if err != nil {
			return nil, err
		}
		return &b, nil
	} else {
		b.api, err = tgbotapi.NewBotAPI(key)
		if err != nil {
			return nil, err
		}
		return &b, nil
	}
}

func (b *Bot) GetMe() (User, error) {
	u, err := b.api.GetMe()
	return User(u), err
}

// 耗时较长
// 从消息查询消息出处的所有管理员
func (b *Bot) GetAdminsByMsg(m Msg) ([]User, error) {
	members, err := b.api.GetChatAdministrators(m.Chat.ChatConfig())
	if err != nil {
		return nil, err
	}
	admins := []User{}
	for _, v := range members {
		admins = append(admins, User(*v.User))
	}
	return admins, nil
}

// 消息是否是来自管理员
func (b *Bot) IsMsgFromAdmin(m Msg) (bool, error) {
	admins, err := b.GetAdminsByMsg(m)
	if err != nil {
		return false, err
	}
	for _, v := range admins {
		if m.From != nil && m.From.ID == v.ID {
			return true, nil
		}
	}
	return false, nil
}

func (b *Bot) GetNewMsgChan() (MsgChan, error) {
	ucfg := tgbotapi.NewUpdate(0)
	ucfg.Timeout = 60

	ch := make(chan Msg, b.api.Buffer)

	go func() {
		for {
			updates, err := b.api.GetUpdates(ucfg)
			if err != nil {
				b.logger.Erro(err)
				b.logger.Infof("Failed to get updates, retrying in 3 seconds...")
				time.Sleep(time.Second * 3)

				continue
			}

			for _, update := range updates {
				if update.UpdateID >= ucfg.Offset {
					ucfg.Offset = update.UpdateID + 1
					if update.Message != nil { // 偶尔可能出现空Message指针
						ch <- Msg(*update.Message)
					}
				}
			}
		}
	}()

	return ch, nil
}

func (b *Bot) Reply(m Msg, text string) (Msg, error) {
	msg := tgbotapi.NewMessage(m.Chat.ID, text)
	msg.ReplyToMessageID = m.MessageID
	res, err := b.api.Send(msg)
	return Msg(res), err
}

func (b *Bot) Answer(m Msg, text string) (Msg, error) {
	msg := tgbotapi.NewMessage(m.Chat.ID, text)
	res, err := b.api.Send(msg)
	return Msg(res), err
}

func (b *Bot) RemoveMsg(m Msg) error {
	if m.Chat != nil {
		return b.RemoveMsgById(m.Chat.ID, m.MessageID)
	} else {
		return b.RemoveMsgById(0, m.MessageID)
	}
}

func (b *Bot) RemoveMsgById(chatId int64, msgId int) error {
	config := tgbotapi.DeleteMessageConfig{MessageID: msgId}
	if chatId != 0 {
		config.ChatID = chatId
	}
	resp, err := b.api.DeleteMessage(config)
	if err != nil {
		return err
	}
	if !resp.Ok {
		return gerrors.Errorf("RemoveMsg error code %s, error message %s", resp.ErrorCode, resp.Description)
	}
	return nil
}

// Every private channel has a unique and solid chat id.
func (b *Bot) SendTextToPrivateChannel(chatIdOfChannel int64, text string) (Msg, error) {
	res, err := b.api.Send(tgbotapi.NewMessage(chatIdOfChannel, text))
	return Msg(res), err
}

// 只有加入了群，才能发送消息到该群
// channelUserName: for example, https://t.me/rich, the channel username is 'rich'.
func (b *Bot) SendTextToPublicChannel(channelUserName, text string) (Msg, error) {
	if !gstring.StartWith(channelUserName, "@") {
		channelUserName = "@" + channelUserName
	}
	res, err := b.api.Send(tgbotapi.NewMessageToChannel(channelUserName, text))
	return Msg(res), err
}

func (b *Bot) SendImageToPublicChannel(channelUserName, imageBuffer []byte) (Msg, error) {
	/*if !xstring.StartWith(channelUserName, "@") {
		channelUserName = "@" + channelUserName
	}

	res, err := b.api.Send(tgbotapi.NewPhotoUpload(channelUserName, text))
	return Msg(res), err*/

	return Msg{}, nil
}

// 加入黑名单并且永久地踢出去（连看消息的权限都没有）
// 经测试发现：是否填写SuperGroupUsername或者ChannelUsername根本没有关系，ChatId却是必须的，没有ChatId只有群组名称，函数是没有效果的，也不报错
// 当然，也可能是测试不严谨的误判
func (b *Bot) BlockGroupMember(chatId int64, userId int) error {
	resp, err := b.api.KickChatMember(tgbotapi.KickChatMemberConfig{ChatMemberConfig: tgbotapi.ChatMemberConfig{ChatID: chatId, UserID: userId}, UntilDate: time.Now().Add(gtime.Year365 * 100).Unix()})
	if err != nil {
		return err
	}
	if !resp.Ok {
		return gerrors.Errorf("BlockGroupMember error code %s, error message %s", resp.ErrorCode, resp.Description)
	}
	return nil
}

func (b *Bot) restrictGroupMember(chatId int64, channelUsername, superGroupUsername string, userId int, blockTime time.Duration, allowed bool) error {
	config := tgbotapi.RestrictChatMemberConfig{}
	config.ChatMemberConfig.ChatID = chatId
	config.ChatMemberConfig.ChannelUsername = channelUsername
	config.ChatMemberConfig.SuperGroupUsername = superGroupUsername
	config.ChatMemberConfig.UserID = userId
	config.UntilDate = time.Now().Add(blockTime).Unix()
	config.CanSendMessages = &allowed
	config.CanSendMediaMessages = &allowed
	config.CanSendOtherMessages = &allowed
	config.CanAddWebPagePreviews = &allowed
	resp, err := b.api.RestrictChatMember(config)
	if err != nil {
		return err
	}
	if !resp.Ok {
		return gerrors.Errorf("restrictGroupMember error code %s, error message %s", resp.ErrorCode, resp.Description)
	}
	return nil
}

// 加入黑名单并且禁言，但是仍然留在群里可以看不可以操作
func (b *Bot) BanGroupMember(chatId int64, userId int, blockTime time.Duration) error {
	return b.restrictGroupMember(chatId, "", "", userId, gtime.Year365*100, false)
}

func (b *Bot) UnBanGroupMember(chatId int64, channelUsername, superGroupUsername string, userId int, blockTime time.Duration) error {
	return b.restrictGroupMember(chatId, channelUsername, superGroupUsername, userId, gtime.Year365*100, true)
}

// 获取群组全部成员名单
func (b *Bot) GetGroupMembers() ([]User, error) {
	return nil, nil
}

// Does user have avatar or not
func (b *Bot) UserHasAvatar(userId int) (bool, error) {
	photo, err := b.api.GetUserProfilePhotos(tgbotapi.UserProfilePhotosConfig{UserID: userId})
	if err != nil {
		return false, err
	}
	return photo.TotalCount > 0, nil
}
