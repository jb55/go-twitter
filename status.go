//
// Copyright 2009 Bill Casarin <billcasarin@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
package twitter

type Status interface {
  GetCreatedAt() string
  GetCreatedAtInSeconds() int64
  GetFavorited() bool
  GetId() int64
  GetText() string
  GetInReplyToScreenName() string
  GetInReplyToStatusId() int64
  GetInReplyToUserId() int64
  GetNow() int
  GetUser() User
  setUser(user User)
}

type errorSource interface {
  GetError() string
}

// Our internal status struct
// the naming is odd so that
// json.Unmarshal can do its thing properly
type tTwitterStatus struct {
  Text                    string
  Created_at              string
  Favorited               bool
  Id                      int64
  In_reply_to_screen_name string
  In_reply_to_status_id   int64
  In_reply_to_user_id     int64
  Error                   string
  User                    *tTwitterUser
  now                     int
  createdAtSeconds        int64
}

func newEmptyTwitterStatus() *tTwitterStatus { return new(tTwitterStatus) }

func (self *tTwitterStatus) GetError() string { return self.Error }

func (self *tTwitterStatus) GetCreatedAt() string {
  return self.Created_at
}

func (self *tTwitterStatus) GetUser() User {
  if self.User == nil {
    self.User = newEmptyTwitterUser()
  }
  self.User.setStatus(self)
  return self.User
}

func (self *tTwitterStatus) setUser(user User) {
  self.User = user.(*tTwitterUser)
}

func (self *tTwitterStatus) GetCreatedAtInSeconds() int64 {
  if self.createdAtSeconds == 0 {
    self.createdAtSeconds = parseTwitterDate(self.Created_at).Seconds()
  }
  return self.createdAtSeconds;
}

func (self *tTwitterStatus) GetFavorited() bool {
  return self.Favorited
}

func (self *tTwitterStatus) GetId() int64 { return self.Id }

func (self *tTwitterStatus) GetInReplyToScreenName() string {
  return self.In_reply_to_screen_name
}

func (self *tTwitterStatus) GetText() string { return self.Text }

func (self *tTwitterStatus) GetInReplyToStatusId() int64 {
  return self.In_reply_to_status_id
}

func (self *tTwitterStatus) GetInReplyToUserId() int64 {
  return self.In_reply_to_user_id
}

func (self *tTwitterStatus) GetNow() int { return self.now }
