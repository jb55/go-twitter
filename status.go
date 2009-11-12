package twitter

type Status interface {
  AsJsonString() string;
  GetCreatedAt() string;
  GetCreatedAtInSeconds() int;
  GetFavorited() bool;
  GetId() int64;
  GetInReplyToScreenName() string;
  GetInReplyToStatusId() int64;
  GetInReplyToUserId() int64;
  GetNow() string;
}

type tTwitterStatus struct {
  jsonString string;
  CreatedAt string;
  CreatedAtSeconds int;
  Favorited bool;
  Id int64;
  InReplyToScreenName string;
  InReplyToStatusId int64;
  InReplyToUserId int64;
  Now string;
}

func (ts *tTwitterStatus) AsJsonString() string {
  return ts.jsonString;
}

func (ts *tTwitterStatus) GetCreatedAt() string {
  return ts.CreatedAt;
}

func (ts *tTwitterStatus) GetCreatedAtInSeconds() int {
  return ts.CreatedAtSeconds;
}

func (ts *tTwitterStatus) GetFavorited() bool {
  return ts.Favorited;
}

func (ts *tTwitterStatus) GetId() int {
  return 10;
}
