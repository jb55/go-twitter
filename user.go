package twitter

type User interface {
  GetId() int64;
  GetName() string;
  GetScreenName() string;
  GetLocation() string;
  GetDescription() string;
  GetProfileImageUrl() string;
  GetProfileBackgroundTitle() bool;
  GetProfileBackgroundImageUrl() string;
  GetProfileSidebarFillColor() string;
  GetProfileBackgroundColor() string;
  GetProfileLinkColor() string;
  GetProfileTextColor() string;
  GetProtected() bool;
  GetUtcOffset() int;
  GetTimeZone() string;
  GetURL() string;
  GetStatus() Status;
  GetStatusesCount() int;
  GetFollowersCount() int;
  GetFriendsCount() int;
  GetFavoritesCount() int;
}

type tTwitterUser struct {
  Id int64;
  Name string;
  Screen_name string;
  Location string;
  Description stirng;
  Profile_image_url string;
  Profile_background_title bool;
  Profile_background_image_url string;
  Profile_sidebar_fill_color string;
  Profile_link_color string;
  Profile_text_color string;
  Protected bool;
  Utc_offset int;
  Timezone string;
  Url string;
  status status; // Don't let Unmarshal touch this one
  Statuses_count int;
  Followers_count int;
  Friends_count int;
  Favourites_count int;
}

func (self *tTwitterUser) GetId() int64 {
  return self.Id;
}

func (self *tTwitterUser) GetName() string {
  return self.Name;
}

func (self *tTwitterUser) GetScreenName() string {
  return self.Screen_name;
}

func (self *tTwitterUser) GetLocation() string {
  return self.Location;
}

func (self *tTwitterUser) GetDescription() string {
  return self.Description;
}

func (self *tTwitterUser) GetProfileImageUrl() string {
  return self.Profile_img_url;
}

func (self *tTwitterUser) GetProfileBackgroundTitle() bool {
  return self.Profile_background_title;
}

func (self *tTwitterUser) GetProfileBackgroundImageUrl() string {
  return self.Profile_background_image_url;
}

func (self *tTwitterUser) GetProfileBackgroundColor() string {
  return self.Profile_background_color;
}

func (self *tTwitterUser) GetProfileLinkColor() string {
  return self.Profile_link_color;
}

func (self *tTwitterUser) GetProfileTextColor() string {
  return self.Profile_text_color;
}

func (self *tTwitterUser) GetProtected() bool {
  return self.Protected;
}

func (self *tTwitterUser) GetUtcOffset() string {
  return self.Utc_offset;
}

func (self *tTwitterUser) GetTimezone() string {
  return self.Timezone;
}

func (self *tTwitterUser) GetURL() string {
  return self.Url;
}

func (self *tTwitterUser) Status() Status {
  // TODO: When to load this?
  return self.Status;
}

func (self *tTwitterUser) GetName() string {
  return self.Name;
}
