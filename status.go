package twitter

import "json"
import "os"

type Status interface {
  AsJsonString() string;
  GetCreatedAt() string;
  GetCreatedAtInSeconds() int;
  GetFavorited() bool;
  GetId() int64;
  GetText() string;
  GetInReplyToScreenName() string;
  GetInReplyToStatusId() int64;
  GetInReplyToUserId() int;
  GetNow() int;
}

type tTwitterStatus struct {
  jsonString string;
  text string;
  createdAt string;
  createdAtSeconds int;
  favorited bool;
  id int64;
  inReplyToScreenName string;
  inReplyToStatusId int64;
  inReplyToUserId int;
  now int;
}

func jsonToStatus(raw string, j *json.Json, errors chan os.Error) Status {
  status := new(tTwitterStatus);

  status.jsonString = raw;
  status.createdAt = j.Get("created_at").String();
  status.createdAtSeconds = 0;
  status.favorited = j.Get("favorited").Bool();
  status.text = j.Get("text").String();
  status.id = int64(j.Get("id").Number());

  return status;
}

func (self *tTwitterStatus) AsJsonString() string {
  return self.jsonString;
}

func (self *tTwitterStatus) GetCreatedAt() string {
  return self.createdAt;
}

func (self *tTwitterStatus) GetCreatedAtInSeconds() int {
  return self.createdAtSeconds;
}

func (self *tTwitterStatus) GetFavorited() bool {
  return self.favorited;
}

func (self *tTwitterStatus) GetId() int64 {
  return self.id;
}

func (self *tTwitterStatus) GetInReplyToScreenName() string {
  return self.inReplyToScreenName;
}

func (self *tTwitterStatus) GetText() string {
  return self.text;
}

func (self *tTwitterStatus) GetInReplyToStatusId() int64 {
  return self.inReplyToStatusId;
}

func (self *tTwitterStatus) GetInReplyToUserId() int {
  return self.inReplyToUserId;
}

func (self *tTwitterStatus) GetNow() int {
  return self.now;
}
