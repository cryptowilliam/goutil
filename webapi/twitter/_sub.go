package twitter

import (
	"fmt"
	"github.com/ChimeraCoder/anaconda"
	"github.com/cryptowilliam/goutil/basic/glog"
	"github.com/cryptowilliam/goutil/safe/gchan"
	"net/url"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

type subUserTask struct {
	uid     int64
	sinceId int64
}

type SubChannel struct {
	stream       *anaconda.Stream
	subUserTasks []subUserTask
	runflag      atomic.Value
	exitMsg      chan struct{} // routine退出信号
	c            chan SimpleTweet
	wg           sync.WaitGroup
}

// 通过定时主动请求模拟订阅
// requestInterval: 发起API请求的时间间隔
func (a *Api) SubUsersSimulate(ioch chan SimpleTweet, requestInterval time.Duration, usernames ...string) (*SubChannel, error) {

	if ioch == nil {
		return nil, gerrors.New("SubUsersSimulate nil ioch param")
	}
	if len(usernames) == 0 {
		return nil, gerrors.New("SubUsersSimulate empty sub usernames")
	}
	sc := SubChannel{}
	sc.c = ioch
	for _, v := range usernames {
		time.Sleep(time.Second)
		u, err := a.GetUserByUsername(v)
		if err != nil {
			return nil, err
		}
		glog.Infof("Get uid %d of user %s", u.Id, v)
		sc.subUserTasks = append(sc.subUserTasks, subUserTask{uid: u.Id})
	}
	if len(sc.subUserTasks) == 0 {
		return nil, gerrors.New("No valid sub uid")
	}

	sc.wg.Add(1)
	go func() {
		defer sc.wg.Add(-1)

		// API请求的连续错误数
		continuousError := 0

		for len(sc.subUserTasks) > 0 {
			for i := range sc.subUserTasks {
				select {
				case <-sc.exitMsg:
					return
				default:
				}

				// 连续错误超过20，就认为这个session有问题，需要退出进程
				if continuousError > 20 {
					glog.Errof("Twitter api continue error more than 20")
					return
				}

				if sc.subUserTasks[i].sinceId == 0 {
					time.Sleep(time.Second)
				} else {
					time.Sleep(requestInterval)
				}
				if sc.subUserTasks[i].sinceId == 0 {
					sinceId, err := a.GetLatestPostId(sc.subUserTasks[i].uid)
					if err != nil {
						continuousError++
						glog.Erro(err)
					} else {
						sc.subUserTasks[i].sinceId = sinceId
						continuousError = 0
						glog.Infof("Get sinceId %d of uid %d", sc.subUserTasks[i].sinceId, sc.subUserTasks[i].uid)
					}
				} else {
					twts, err := a.GetTweetsSinceId(sc.subUserTasks[i].uid, &sc.subUserTasks[i].sinceId, nil)

					if err != nil {
						continuousError++
						glog.Erro(err)
					} else {
						continuousError = 0
						if len(twts) > 0 {
							sc.subUserTasks[i].sinceId = maxPostId(twts)
						}
						for _, v := range twts {
							// 这里加上exitMsg检测，保证主进程发出关闭信号后，尽快退出线程，而不会因为sc.c满阻塞导致线程卡住
							select {
							case <-sc.exitMsg:
								return
							case sc.c <- v:
							}
						}
					}
				}
			}
		}
	}()

	return &sc, nil
}

func readSubTweets(s *anaconda.Stream) *SimpleTweet {
	itf := <-s.C
	if itf == nil {
		return nil
	}
	switch itf.(type) {
	case anaconda.Tweet:
		tw := itf.(anaconda.Tweet) // try casting into a tweet
		return readSimpleTweet(&tw)
	case anaconda.StatusDeletionNotice: // 订阅用户且该用户删了某个帖子时触发
		fmt.Println("StatusDeletionNotice")
	case anaconda.DirectMessageDeletionNotice:
		fmt.Println("DirectMessageDeletionNotice")
	case anaconda.LocationDeletionNotice:
		fmt.Println("LocationDeletionNotice")
	case anaconda.LimitNotice:
		fmt.Println("LimitNotice")
	case anaconda.StatusWithheldNotice:
		fmt.Println("StatusWithheldNotice")
	case anaconda.UserWithheldNotice:
		fmt.Println("UserWithheldNotice")
	case anaconda.DisconnectMessage:
		fmt.Println("DisconnectMessage")
	case anaconda.StallWarning:
		fmt.Println("StallWarning")
	case anaconda.FriendsList:
		//_ := itf.(anaconda.FriendsList)
		fmt.Println("FriendsList")
	case anaconda.DirectMessage:
		fmt.Println("DirectMessage")
	case anaconda.EventTweet:
		fmt.Println("EventTweet")
	case anaconda.EventList:
		fmt.Println("EventList")
	case anaconda.Event:
		fmt.Println("Event")
	case anaconda.EventFollow:
		fmt.Println("EventFollow")
	}
	return nil
}

func (sc *SubChannel) readroutine() {
	defer sc.wg.Add(-1)
	for {
		select {
		case <-sc.exitMsg:
			return
		default:
		}

		twt := readSubTweets(sc.stream)
		if twt == nil {
			time.Sleep(time.Millisecond * 500)
			continue
		}

		// 这里加上exitMsg检测，保证主进程发出关闭信号后，尽快退出线程，而不会因为sc.c满阻塞导致线程卡住
		select {
		case <-sc.exitMsg:
			return
		case sc.c <- *twt:
		}
	}
}

func (sc *SubChannel) Wait() {
	sc.wg.Wait()
}

func (sc *SubChannel) Close() {
	gchan.SafeCloseChanStruct(sc.exitMsg)
	sc.wg.Wait()
	sc.stream.Stop()
	sc.stream = nil
}

// TODO: 这个接口是不是可以删了？工作很不稳定，随时可能获取不到实时反馈流
// 订阅的推特不仅仅是该发送者，还包括at他的内容
// username是指@后面的，注意与nickname区别开
func (a *Api) SubUsers(ioch chan SimpleTweet, usernames ...string) (*SubChannel, error) {
	if len(usernames) == 0 {
		return nil, gerrors.Errorf("Empty sub users")
	}

	uids := []string{}
	for _, v := range usernames {
		u, e := a.GetUserByUsername(v)
		if e != nil {
			glog.Erro(e)
		} else {
			uids = append(uids, strconv.FormatInt(u.Id, 10))
		}
	}
	if len(uids) == 0 {
		return nil, gerrors.Errorf("Can't get uids of usernames %s", usernames)
	}

	sc := SubChannel{}
	// 只支持uid，不支持username或者nickname
	val := url.Values{}
	for _, v := range uids {
		val.Add("follow", v)
	}
	sc.stream = a.inApi.PublicStreamFilter(val /*url.Values{"follow": uids}*/)
	sc.c = ioch
	sc.exitMsg = make(chan struct{})
	sc.wg.Add(1)
	go sc.readroutine()
	return &sc, nil
}

// This api works fine
// If you need AND, use "word1 word2".
// If you need OR, use "word1,word2", or put them in different kws array item.
func (a *Api) SubKeywords(ioch chan SimpleTweet, kws ...string) (*SubChannel, error) {
	if len(kws) == 0 {
		return nil, gerrors.Errorf("Empty tracks")
	}

	sc := SubChannel{}
	sc.stream = a.inApi.PublicStreamFilter(url.Values{"track": append([]string{}, kws...)})
	sc.c = ioch
	sc.exitMsg = make(chan struct{})
	sc.wg.Add(1)
	go sc.readroutine()
	return &sc, nil
}

func maxPostId(tws []SimpleTweet) int64 {
	if len(tws) == 0 {
		return 0
	}
	maxId := tws[0].Id
	for _, v := range tws {
		if v.Id > maxId {
			maxId = v.Id
		}
	}
	return maxId
}
