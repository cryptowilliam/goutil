package twitter

import "github.com/cryptowilliam/goutil/basic/gerrors"

func (a *Api) GetUserById(userId int64) (*SimpleUser, error) {
	users, err := a.inApi.GetUsersLookupByIds([]int64{userId}, nil)
	if err != nil {
		return nil, err
	}
	if len(users) == 0 {
		return nil, gerror.Errorf("Can't find user by Id %d", userId)
	}
	u := SimpleUser{}
	u.Id = userId
	u.UserName = users[0].ScreenName
	u.DisplayName = users[0].Name
	u.FollowersCount = users[0].FollowersCount
	return &u, nil
}

// username is @username, NOT nickname(showname)
func (a *Api) GetUserByUsername(username string) (*SimpleUser, error) {
	users, err := a.inApi.GetUsersLookup(username, nil)
	if err != nil {
		return nil, err
	}
	if len(users) == 0 {
		return nil, gerror.Errorf("Can't find user by unique name %s", username)
	}
	u := SimpleUser{}
	u.Id = users[0].Id
	u.UserName = users[0].Name
	u.DisplayName = users[0].ScreenName
	u.FollowersCount = users[0].FollowersCount
	return &u, nil
}
