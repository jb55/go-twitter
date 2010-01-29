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

type User interface {
  GetId() int64
  GetName() string
  GetScreenName() string
  GetLocation() string
  GetDescription() string
  GetProfileImageUrl() string
  GetProfileBackgroundTitle() bool
  GetProfileBackgroundImageUrl() string
  GetProfileSidebarFillColor() string
  GetProfileLinkColor() string
  GetProfileTextColor() string
  GetProtected() bool
  GetUtcOffset() int
  GetTimeZone() string
  GetURL() string
  GetStatus() Status
  setStatus(status Status)
  GetStatusesCount() int
  GetFollowersCount() int
  GetFriendsCount() int
  GetFavoritesCount() int
}

type tTwitterUser struct {
  Id                           int64
  Name                         string
  Screen_name                  string
  Location                     string
  Description                  string
  Profile_image_url            string
  Profile_background_title     bool
  Profile_background_image_url string
  Profile_sidebar_fill_color   string
  Profile_link_color           string
  Profile_text_color           string
  Protected                    bool
  Utc_offset                   int
  Url                          string
  Timezone                     string
  Status                       *tTwitterStatus
  Statuses_count               int
  Followers_count              int
  Friends_count                int
  Favorites_count              int
  Error                        string
}

type tTwitterUserDummy struct {
  Object tTwitterUser
}

type tTwitterUserListDummy struct {
  Object []tTwitterUser
}

func newEmptyTwitterUser() *tTwitterUser { return new(tTwitterUser) }

func (self *tTwitterUser) GetError() string { return self.Error }

func (self *tTwitterUser) GetId() int64 { return self.Id }

func (self *tTwitterUser) GetName() string { return self.Name }

func (self *tTwitterUser) GetScreenName() string {
  return self.Screen_name
}

func (self *tTwitterUser) GetLocation() string {
  return self.Location
}

func (self *tTwitterUser) GetDescription() string {
  return self.Description
}

func (self *tTwitterUser) GetProfileImageUrl() string {
  return self.Profile_image_url
}

func (self *tTwitterUser) GetProfileBackgroundTitle() bool {
  return self.Profile_background_title
}

func (self *tTwitterUser) GetProfileSidebarFillColor() string {
  return self.Profile_sidebar_fill_color
}

func (self *tTwitterUser) GetProfileBackgroundImageUrl() string {
  return self.Profile_background_image_url
}

func (self *tTwitterUser) GetProfileLinkColor() string {
  return self.Profile_link_color
}

func (self *tTwitterUser) GetProfileTextColor() string {
  return self.Profile_text_color
}

func (self *tTwitterUser) GetProtected() bool { return self.Protected }

func (self *tTwitterUser) GetUtcOffset() int { return self.Utc_offset }

func (self *tTwitterUser) GetTimeZone() string {
  return self.Timezone
}

func (self *tTwitterUser) GetURL() string { return self.Url }

func (self *tTwitterUser) GetStatus() Status {
  if self.Status == nil {
    self.Status = newEmptyTwitterStatus()
  }
  self.Status.setUser(self)
  return self.Status
}

func (self *tTwitterUser) setStatus(status Status) {
  self.Status = status.(*tTwitterStatus)
}

func (self *tTwitterUser) GetStatusesCount() int {
  return self.Statuses_count
}

func (self *tTwitterUser) GetFollowersCount() int {
  return self.Followers_count
}

func (self *tTwitterUser) GetFriendsCount() int {
  return self.Friends_count
}

func (self *tTwitterUser) GetFavoritesCount() int {
  return self.Favorites_count
}
