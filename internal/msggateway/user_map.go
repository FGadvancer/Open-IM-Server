// Copyright © 2023 OpenIM. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package msggateway

import (
	"context"
	"hash/fnv"
	"sync"

	"github.com/OpenIMSDK/tools/log"
	"github.com/OpenIMSDK/tools/utils"
)

type UserMap struct {
	maps []sync.Map
}

func newUserMap(numMaps int) *UserMap {
	um := &UserMap{
		maps: make([]sync.Map, numMaps),
	}
	return um
}

func (u *UserMap) getMap(key string) *sync.Map {
	hash := fnv.New32a()
	hash.Write([]byte(key))
	index := hash.Sum32() % uint32(len(u.maps))
	return &u.maps[index]
}

func (u *UserMap) GetAll(key string) ([]*Client, bool) {
	sm := u.getMap(key)
	allClients, ok := sm.Load(key)
	if ok {
		return allClients.([]*Client), ok
	}
	return nil, ok
}

func (u *UserMap) Get(key string, platformID int) ([]*Client, bool, bool) {
	sm := u.getMap(key)
	allClients, userExisted := sm.Load(key)
	if userExisted {
		var clients []*Client
		for _, client := range allClients.([]*Client) {
			if client.PlatformID == platformID {
				clients = append(clients, client)
			}
		}
		if len(clients) > 0 {
			return clients, userExisted, true
		}
		return clients, userExisted, false
	}
	return nil, userExisted, false
}

func (u *UserMap) Set(key string, v *Client) {
	sm := u.getMap(key)
	allClients, existed := sm.Load(key)
	if existed {
		log.ZDebug(context.Background(), "Set existed", "user_id", key, "client", *v)
		oldClients := allClients.([]*Client)
		oldClients = append(oldClients, v)
		sm.Store(key, oldClients)
	} else {
		log.ZDebug(context.Background(), "Set not existed", "user_id", key, "client", *v)
		var clients []*Client
		clients = append(clients, v)
		sm.Store(key, clients)
	}
}

func (u *UserMap) delete(key string, connRemoteAddr string) (isDeleteUser bool) {
	sm := u.getMap(key)
	allClients, existed := sm.Load(key)
	if existed {
		oldClients := allClients.([]*Client)
		var a []*Client
		for _, client := range oldClients {
			if client.ctx.GetRemoteAddr() != connRemoteAddr {
				a = append(a, client)
			}
		}
		if len(a) == 0 {
			sm.Delete(key)
			return true
		} else {
			sm.Store(key, a)
			return false
		}
	}
	return existed
}

func (u *UserMap) deleteClients(key string, clients []*Client) (isDeleteUser bool) {
	sm := u.getMap(key)
	m := utils.SliceToMapAny(clients, func(c *Client) (string, struct{}) {
		return c.ctx.GetRemoteAddr(), struct{}{}
	})
	allClients, existed := sm.Load(key)
	if existed {
		oldClients := allClients.([]*Client)
		var a []*Client
		for _, client := range oldClients {
			if _, ok := m[client.ctx.GetRemoteAddr()]; !ok {
				a = append(a, client)
			}
		}
		if len(a) == 0 {
			sm.Delete(key)
			return true
		} else {
			sm.Store(key, a)
			return false
		}
	}
	return existed
}

func (u *UserMap) DeleteAll(key string) {
	sm := u.getMap(key)
	sm.Delete(key)
}
