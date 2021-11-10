package router

import (
	"container/list"
	"strings"
	"sync"
)

type (
	SubscriptionPattern struct {
		p string
		t *list.List
	}

	// SubscriptionTable -
	SubscriptionTable struct {
		sub []ISubscriptionPattern
	}
)

// NewSubscribePattern -
func NewSubscribePattern(pattern string) (ISubscriptionPattern, error) {
	s := &SubscriptionPattern{
		p: pattern,
		t: list.New(),
	}
	if e := s.create(pattern); e != nil {
		return nil, e
	}
	return s, nil
}

//
// create - создает структуру для паттерна подписки
func (sp *SubscriptionPattern) create(s string) error {

	sp.t.Init()
	isTail := false

	for _, v := range strings.Split(strings.Trim(s, charDelimiter), charDelimiter) {

		if nodeValue, e := createNodeValue(v); e == nil {
			switch nodeValue.Type() {
			case nodeTypeSuffix, nodeTypePrefix:
				return errInvalidSubscribeTopic

			case nodeTypeTail:
				if isTail {
					return errInvalidSubscribeTopic
				}
				isTail = true
			}
			sp.t.PushBack(nodeValue)

		} else {
			return e
		}
	}
	return nil
}

func (sp *SubscriptionPattern) Pattern() string {
	return sp.p
}

// Match - сравнивает топик с паттерном подписки
func (sp *SubscriptionPattern) Match(topic string) (result bool) {

	tn := strings.Split(strings.Trim(topic, charDelimiter), charDelimiter)
	if len(tn) < 1 {
		return false
	}

	ptr := sp.t.Front()

	for _, v := range tn {
		if ptr == nil {
			return false
		}

		switch ptr.Value.(INode).Type() {
		case nodeTypeString:
			if !ptr.Value.(INode).Validate(v) {
				return false
			}
		case nodeTypePlus:
			if !ptr.Value.(INode).Validate(v) {
				return false
			}
		case nodeTypeTail:
			return ptr.Value.(INode).Validate(v)

		default:
			return false
		}
		ptr = ptr.Next()
	}
	return ptr.Value == nil
	return
}

func NewSubscriptionTable() ISubscriptionTable {
	return &SubscriptionTable{
		sub: nil,
	}
}

//
func (st *SubscriptionTable) Match(topic string) bool {
	for _, v := range st.sub {
		if v.Match(topic) {
			return true
		}
	}
	return false
}

//
func (st *SubscriptionTable) Add(pattern string) (e error) {
	if p, e := NewSubscribePattern(pattern); e == nil {
		for _, v := range st.sub {
			if v.Pattern() == pattern {
				return errPatternIsPresenr
			}
		}
		st.sub = append(st.sub, p)
	}
	return
}

//
func (st *SubscriptionTable) Remove(pattern string) {
	if len(st.sub) > 0 {
		var tmp []ISubscriptionPattern
		mu := sync.Mutex{}
		mu.Lock()
		for _, v := range st.sub {
			if v.Pattern() != pattern {
				tmp = append(tmp, v)
			}
		}
		st.sub = tmp
		mu.Unlock()
	}
}
