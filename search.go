package twitter

type SearchResult interface {
  GetCreatedAt() string;
  GetFromUser() string;
  GetToUserId() int64;
  GetText() string;
  GetId() int64;
  GetFromUserId() int64;
  GetGeo() string;
  GetIsoLanguageCode() string;
  GetSource() string;
}

type tTwitterSearch struct {
  Results []tTwitterSearchResult;
}

type tTwitterSearchResult struct {
  Profile_image_url string;
  Created_at string;
  From_user string;
  To_user_id int64;
  Text string;
  Id int64;
  From_user_id int64;
  Geo string;
  Iso_language_code string;
  Source string;
  Error string;
}

func (self *tTwitterSearchResult) GetError() string {
  return self.Error;
}

func (self *tTwitterSearchResult) GetCreatedAt() string {
  return self.Created_at;
}

func (self *tTwitterSearchResult) GetFromUser() string {
  return self.From_user;
}

func (self *tTwitterSearchResult) GetToUserId() int64 {
  return self.To_user_id;
}

func (self *tTwitterSearchResult) GetText() string {
  return self.Text;
}

func (self *tTwitterSearchResult) GetId() int64 {
  return self.Id;
}

func (self *tTwitterSearchResult) GetFromUserId() int64 {
  return self.From_user_id;
}

func (self *tTwitterSearchResult) GetGeo() string {
  return self.Geo;
}

func (self *tTwitterSearchResult) GetIsoLanguageCode() string {
  return self.Geo;
}

func (self *tTwitterSearchResult) GetSource() string {
  return self.Source;
}
