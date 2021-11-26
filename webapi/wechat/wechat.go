package wechat

import (
	"fmt"
	"github.com/cryptowilliam/goutil/basic/gerrors"
	"github.com/cryptowilliam/goutil/basic/glog"
	"github.com/cryptowilliam/goutil/container/grand"
	"github.com/cryptowilliam/goutil/container/gstring"
	"github.com/cryptowilliam/goutil/net/gnet"
	"github.com/cryptowilliam/goutil/safe/gchan"
	"github.com/cryptowilliam/goutil/sys/gfs"
	"github.com/cryptowilliam/goutil/sys/gproc"
	"github.com/cryptowilliam/goutil/sys/gtime"
	"github.com/songtianyi/wechat-go/plugins/wxweb/cleaner"
	"github.com/songtianyi/wechat-go/plugins/wxweb/faceplusplus"
	"github.com/songtianyi/wechat-go/plugins/wxweb/forwarder"
	"github.com/songtianyi/wechat-go/plugins/wxweb/gifer"
	"github.com/songtianyi/wechat-go/plugins/wxweb/joker"
	"github.com/songtianyi/wechat-go/plugins/wxweb/replier"
	"github.com/songtianyi/wechat-go/plugins/wxweb/revoker"
	"github.com/songtianyi/wechat-go/plugins/wxweb/switcher"
	"github.com/songtianyi/wechat-go/plugins/wxweb/system"
	"github.com/songtianyi/wechat-go/wxweb"
	"net/http"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type Bot struct {
	sess           *wxweb.Session
	recvch         chan RecvMsg
	exitMsg        chan struct{}
	wg             sync.WaitGroup
	loginResult    interface{}
	loginResultMtx sync.RWMutex
	sendQueue      chan sendTask
	logger         glog.Interface

	// temp caches
	loginMode    LoginMode
	qrCodeUrlWAN string
	qrCodeUrlLAN string
}

// 个人和群组都属于User
type User struct {
	Uin         int    // User info, 只有自己的能获取，用于聊天记录缓存的加密解密，所以获取不到别人的
	DynUserName string // 动态的Username，是一串长度不定的加密字符串，过一段时间就会变，用它而不用微信号，是因为无法获取真实的唯一微信号
	NickName    string // 显示
	RemarkName  string // 备注名，群聊没有这个选项
	IsMale      bool
	Signature   string
	Province    string
	City        string
	IsGroup     bool
}

type MessageType int

const (
	MessageTypeText MessageType = iota + 0
	MessageTypeImage
	/*MessageTypeAudio
	MessageTypeVideo
	MessageTypeOtherFile*/
)

type LoginMode int

const (
	LoginModeTerminal LoginMode = iota + 0
	LoginModeWeb
)

type RecvMsg struct {
	MsgId       string
	SenderUid   string
	ReceiverUid string // 接收者，在个人聊天中就是自己，群聊消息中则不一定
	Content     []byte
	Type        MessageType
	IsGroup     bool // 是不是群组消息
}

// 发送任务，用于send queue，不要无时间间隙地连续发送消息，避免被封号
type sendTask struct {
	recvNameType RecvNameType
	recvName     string
	content      []byte
	msgType      MessageType
}

type RecvNameType int

const (
	RecvNameTypeNickName RecvNameType = iota + 0
	RecvNameTypeRemarkName
)

func wxRecvMsgToRecvMsg(sess wxweb.Session, msg wxweb.ReceivedMessage) (RecvMsg, error) {
	res := RecvMsg{}

	if msg.MsgType == wxweb.MSG_TEXT {
		res.Type = MessageTypeText
		res.Content = []byte(msg.Content)
	} else if msg.MsgType == wxweb.MSG_IMG {
		res.Type = MessageTypeImage
		b, err := sess.GetImg(msg.MsgId)
		if err != nil {
			return RecvMsg{}, err
		}
		res.Content = b
	} else {

	}

	res.ReceiverUid = msg.Who
	res.IsGroup = msg.IsGroup

	return res, nil
}

func wxUserToUser(src wxweb.User) User {
	return User{
		DynUserName: src.UserName,
		Uin:         src.Uin,
		NickName:    src.NickName,
		RemarkName:  src.RemarkName,
		IsMale:      src.Sex == 1,
		Signature:   src.Signature,
		Province:    src.Province,
		City:        src.City,
		IsGroup:     gstring.StartWith(src.UserName, "@@"),
	}
}

// 必须有的插件注册函数
// 指定session, 可以对不同用户注册不同插件
func addRecvChan(sess *wxweb.Session, name string, recvch chan RecvMsg) error {
	// 将插件注册到session并开启，插件其实就是消息处理函数
	// 第一个参数: 指定消息类型, 所有该类型的消息都会被转发到此插件
	// 第二个参数: 指定消息处理函数, 消息会进入此函数
	// 第三个参数: 自定义插件名，不能重名，switcher插件会用到此名称
	err := sess.HandlerRegister.Add(wxweb.MSG_TEXT,
		func(s *wxweb.Session, msg *wxweb.ReceivedMessage) {
			m, err := wxRecvMsgToRecvMsg(*sess, *msg)
			if err == nil {
				recvch <- m
			}
		},
		name)
	if err != nil {
		return err
	}
	if err := sess.HandlerRegister.EnableByName(name); err != nil {
		return err
	}
	return nil
}

// recvch可以为空，表示忽略收到的消息
// 此接口非阻塞
func NewBot(mode LoginMode, recvch chan RecvMsg, logger glog.Interface) (*Bot, error) {
	ip, err := gnet.GetPublicIPOL("")
	if err != nil {
		return nil, err
	}

	// Create session
	session := new(wxweb.Session)
	err = error(nil)
	webqrcodedir := ""
	switch mode {
	case LoginModeTerminal:
		session, err = wxweb.CreateSession(nil, nil, wxweb.TERMINAL_MODE)
	case LoginModeWeb:
		webqrcodedir, err = gproc.SelfDir()
		if err == nil {
			webqrcodedir = filepath.Join(webqrcodedir, "wechat-qrcode")
			gfs.MakeDir(webqrcodedir)
			gfs.CleanDir(webqrcodedir)
			session, err = wxweb.CreateWebSessionWithPath(nil, nil, webqrcodedir)
		}
	default:
		err = gerrors.Errorf("unknown LoginMode %d", mode)
	}
	if err != nil {
		return nil, gerrors.Wrap(err, "wechat NewBot()")
	}
	if session == nil {
		return nil, gerrors.New("wechat nil session")
	}

	// load plugins for this session
	faceplusplus.Register(session)
	replier.Register(session)
	switcher.Register(session)
	gifer.Register(session)
	cleaner.Register(session)
	joker.Register(session)
	revoker.Register(session)
	forwarder.Register(session)
	system.Register(session)
	//youdao.Register(session)

	// 不知道想干哈
	if mode == LoginModeTerminal {
		if err := session.HandlerRegister.EnableByType(wxweb.MSG_SYS); err != nil {
			return nil, err
		}
	}

	// 初始化
	c := Bot{sess: session, recvch: recvch, loginMode: mode, logger: logger}
	c.loginResult = false // 此状态表示登陆中
	c.exitMsg = make(chan struct{})
	c.sendQueue = make(chan sendTask, 1024)
	if recvch != nil {
		addRecvChan(c.sess, "myCustomPlugin", recvch)
	}
	c.sess.AfterLogin = func() error {
		c.loginResultMtx.Lock()
		c.loginResult = true
		c.loginResultMtx.Unlock()
		c.wg.Add(1)
		go c.popSendQueue() // 开始排空发送任务队列
		return nil
	}

	// 在routine中登陆并开启消息接收循环
	c.wg.Add(1)
	go func() {
		defer c.wg.Add(-1)

		// serve http
		if mode == LoginModeWeb {
			go http.ListenAndServe(":8080", http.FileServer(http.Dir(webqrcodedir)))
			c.qrCodeUrlWAN = fmt.Sprintf("http://%s:8080/%s", ip.String(), session.QrcodePath)
			c.qrCodeUrlLAN = fmt.Sprintf("http://127.0.0.1:8080/%s", session.QrcodePath)
			c.logger.Infof("please visit %s", c.qrCodeUrlWAN)
			c.logger.Infof("or visit %s", c.qrCodeUrlLAN)
		}

		err := c.sess.LoginAndServe(false)
		if err != nil {
			c.loginResultMtx.Lock()
			c.loginResult = gerrors.Wrap(err, "session LoginAndServe()")
			c.loginResultMtx.Unlock()
		}
	}()

	return &c, nil
}

func (b *Bot) GetSendQueueLen() int {
	return len(b.sendQueue)
}

// 按一定范围的时间间隔，发送任务队列中的发送消息任务
func (b *Bot) popSendQueue() {
	defer b.wg.Add(-1)

	for {
		millis := grand.RandomInt(2000, 7000) // 2000-7000毫秒
		time.Sleep(gtime.MulDuration(int64(millis), time.Millisecond))

		select {
		case task := <-b.sendQueue:
			b.logger.Infof("wechat to send: %s", string(task.content))
			err := b.sendMsg(task.recvNameType, task.recvName, task.content, task.msgType)
			if err != nil {
				b.logger.Erro(err, "b.sendMsg")
				groups := b.GetMyGroups()
				groupsstr := ""
				for _, v := range groups {
					groupsstr += "," + v.NickName
				}
				b.logger.Infof("my group contacts %s", groupsstr)
			}
		case <-b.exitMsg:
			return
		}
	}
}

// 等待登陆结果，如果用户选择中途Logout，则提前返回
func (b *Bot) WaitLogin() error {

	for {
		select {
		case <-b.exitMsg:
			break
		default:
		}

		b.loginResultMtx.RLock()
		rst := b.loginResult
		b.loginResultMtx.RUnlock()

		switch rst.(type) {
		case error: // 登陆失败
			err := rst.(error)
			return err
		case bool:
			if rst.(bool) { // 登陆成功
				return nil
			} else { // 登陆中
				time.Sleep(time.Second)
				continue
			}
		default:
			return gerrors.Errorf("NewBot result invalid value stored")
		}
	}

	return gerrors.Errorf("NewBot canceled by user")
}

// 获取自己的用户信息
func (b *Bot) GetMyself() User {
	return wxUserToUser(*b.sess.Bot)
}

// 从通讯录根据动态加密用户名获取个人和群组信息
// 无法获取没有保存到通讯录的群
func (b *Bot) GetContactByUserName(dynUsername string) (User, error) {
	u := b.sess.Cm.GetContactByUserName(dynUsername)
	if u == nil {
		return User{}, gerrors.Errorf("username %s not found", dynUsername)
	}
	return wxUserToUser(*u), nil
}

// 从通讯录根据昵称名获取个人和群组信息
// 无法获取没有保存到通讯录的群
func (b *Bot) GetContactsByNickName(nickName string, caseSensitive bool) []User {
	allusers := b.GetAllMyContacts()
	rst := []User{}
	for _, v := range allusers {
		if caseSensitive {
			if v.NickName == nickName {
				rst = append(rst, v)
			}
		} else {
			if strings.ToLower(v.NickName) == strings.ToLower(nickName) {
				rst = append(rst, v)
			}
		}
	}
	return rst
}

// 从通讯录根据备注名获取个人和群组信息
// 无法获取没有保存到通讯录的群
func (b *Bot) GetContactsByRemarkName(remarkName string, caseSensitive bool) []User {
	allusers := b.GetAllMyContacts()
	rst := []User{}
	for _, v := range allusers {
		if caseSensitive {
			if v.RemarkName == remarkName {
				rst = append(rst, v)
			}
		} else {
			if strings.ToLower(v.RemarkName) == strings.ToLower(remarkName) {
				rst = append(rst, v)
			}
		}
	}
	return rst
}

// 获取已保存到通讯录的全部个人和群组
// 无法获取没有保存到通讯录的群
func (b *Bot) GetAllMyContacts() []User {
	users := b.sess.Cm.GetAll()
	res := []User{}
	for _, v := range users {
		if v == nil {
			continue
		}
		item := wxUserToUser(*v)
		res = append(res, item)
	}
	return res
}

func (b *Bot) GetMyGroups() []User {
	users := b.sess.Cm.GetGroupContacts()
	res := []User{}
	for _, v := range users {
		if v == nil {
			continue
		}
		item := wxUserToUser(*v)
		res = append(res, item)
	}
	return res
}

func (b *Bot) sendMsg(recvNameType RecvNameType, recvName string, content []byte, msgType MessageType) error {
	errReceiverNotFound := gerrors.Errorf("receiver name (%s) not found", recvName)
	errReceiverUnUnique := gerrors.Errorf("receiver name (%s) is not unique", recvName)

	users := []User{}
	switch recvNameType {
	case RecvNameTypeNickName:
		users = b.GetContactsByNickName(recvName, true)
	case RecvNameTypeRemarkName:
		users = b.GetContactsByRemarkName(recvName, true)
	default:
		return gerrors.Errorf("unknown RecvNameType %d", recvNameType)
	}
	if len(users) == 0 {
		return errReceiverNotFound
	}
	if len(users) > 1 {
		return errReceiverUnUnique
	}

	err := error(nil)
	switch msgType {
	case MessageTypeText:
		_, _, err = b.sess.SendText(string(content), b.sess.Bot.UserName, users[0].DynUserName)
	case MessageTypeImage:
		b.sess.SendImgFromBytes(content, "", b.sess.Bot.UserName, users[0].DynUserName)
	default:
		return gerrors.Errorf("unknown MessageType %d", msgType)
	}
	return err
}

func (b *Bot) SendText(recvNameType RecvNameType, recvName string, content string) error {
	return b.sendMsg(recvNameType, recvName, []byte(content), MessageTypeText)
}

func (b *Bot) SendImage(recvNameType RecvNameType, recvName string, content []byte) error {
	return b.sendMsg(recvNameType, recvName, content, MessageTypeImage)
}

func (b *Bot) SendTextQueue(recvNameType RecvNameType, recvName string, content string) {
	b.sendQueue <- sendTask{recvNameType: recvNameType, recvName: recvName, content: []byte(content), msgType: MessageTypeText}
}

func (b *Bot) SendImageQueue(recvNameType RecvNameType, recvName string, content []byte) {
	b.sendQueue <- sendTask{recvNameType: recvNameType, recvName: recvName, content: content, msgType: MessageTypeImage}
}

func (b *Bot) Logout() error {
	gchan.SafeCloseChanStruct(b.exitMsg)
	b.loginResultMtx.Lock()
	b.loginResult = false
	b.loginResultMtx.Unlock()
	b.wg.Wait()
	return b.sess.Logout()
}
