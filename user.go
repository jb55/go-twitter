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
  Profile_background_img_url string;
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
