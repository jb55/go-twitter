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

import "time"

const kExpireTime = 60 * 1

type IdProvider interface {
  GetId() int64
}

type tMemoryCacheEntry struct {
  data       IdProvider
  timeStored int64
}

type MemoryCache struct {
  mapIdToEntry map[int64]*tMemoryCacheEntry
}

type Cache interface {
  // Stores a value in the database, the key is determined
  // by the GetId() function in the IdProvider interface
  Store(data IdProvider)

  // Checks to see if the cache contains a given key
  HasId(id int64) bool

  // Gets a value from the cache
  Get(id int64) IdProvider

  // Gets the time a given key was stored
  GetTimeStored(id int64) int64
}

type CacheBackend struct {
  hit         int64
  store       int64
  userCache   Cache
  statusCache Cache
  expireTime  int64
}

func (self *MemoryCache) Store(data IdProvider) {
  id := data.GetId()
  var entry *tMemoryCacheEntry
  var notEmpty bool

  if entry, notEmpty = self.mapIdToEntry[id]; notEmpty {
    entry.data = data
  } else {
    entry = new(tMemoryCacheEntry)
    entry.data = data
    self.mapIdToEntry[id] = entry
  }

  entry.timeStored = time.Seconds()
}

func (self *MemoryCache) HasId(id int64) bool {
  _, hasId := self.mapIdToEntry[id]
  return hasId
}

func (self *MemoryCache) Get(id int64) IdProvider {
  if self.HasId(id) {
    return self.mapIdToEntry[id].data
  }

  return nil
}

func (self *MemoryCache) GetTimeStored(id int64) int64 {
  if self.HasId(id) {
    return self.mapIdToEntry[id].timeStored
  }

  return 0
}

func (self *CacheBackend) StoreUser(user User) {
  self.store++
  self.userCache.Store(user)
}

func (self *CacheBackend) StoreStatus(status Status) {
  self.store++
  self.statusCache.Store(status)
}

func (self *CacheBackend) GetUser(id int64) User {
  self.hit++
  return self.userCache.Get(id).(User)
}

func (self *CacheBackend) GetStatus(id int64) Status {
  self.hit++
  return self.statusCache.Get(id).(Status)
}

func (self *CacheBackend) HasUserExpired(id int64) bool {
  timeSinceStored := time.Seconds() - self.userCache.GetTimeStored(id)
  return timeSinceStored >= self.expireTime
}

func (self *CacheBackend) HasStatusExpired(id int64) bool {
  timeSinceStored := time.Seconds() - self.statusCache.GetTimeStored(id)
  return timeSinceStored >= self.expireTime
}

// Creates a custom cache backend
func NewCacheBackend(user Cache, status Cache, expireTime int64) *CacheBackend {
  backend := new(CacheBackend)

  backend.hit = 0
  backend.userCache = user
  backend.statusCache = status
  backend.expireTime = expireTime

  return backend
}

func NewMemoryCache() *MemoryCache {
  cache := new(MemoryCache)
  cache.mapIdToEntry = make(map[int64]*tMemoryCacheEntry)
  return cache
}
