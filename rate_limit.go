package twitter

type RateLimit interface {
  GetRemainingHits() int
  GetHourlyLimit() int
  GetResetTimeInSeconds() int64
  GetResetTime() string
}

type tTwitterRateLimit struct {
  Remaining_hits int
  Hourly_limit int
  Reset_time_in_seconds int64
  Reset_time string
}

func (self *tTwitterRateLimit) GetRemainingHits() int {
  return self.Remaining_hits
}

func (self *tTwitterRateLimit) GetHourlyLimit() int {
  return self.Hourly_limit
}

func (self *tTwitterRateLimit) GetResetTimeInSeconds() int64 {
  return self.Reset_time_in_seconds
}

func (self *tTwitterRateLimit) GetResetTime() string {
  return self.Reset_time
}
