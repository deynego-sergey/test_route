package router

import (
	"container/list"
	"strings"
	"sync"
)

type (
	// RoutePattern -
	RoutePattern struct {
		broker  Agent
		pattern string
		rt      *list.List
	}
	// Route -
	Route struct {
		broker Agent
		topic  string
	}
	// RouteTable -
	RouteTable struct {
		routes map[Agent][]IRoutePattern
	}
)

func NewRouteTable() IRouteTable {
	t := &RouteTable{routes: make(map[Agent][]IRoutePattern)}
	return t
}

// Set - добавляет в таблицу информацию о брокере и паттерн для роутинга
func (rt *RouteTable) Set(broker Agent, pattern string) (e error) {

	ptrn, e := NewRoutePattern(broker, pattern)

	if e != nil {
		return e
	}
	mu := sync.Mutex{}
	mu.Lock()
	rt.routes[broker] = append(rt.routes[broker], ptrn)
	mu.Unlock()
	return
}

// Remove - Удаляет паттерн роутинга
func (rt *RouteTable) Remove(broker Agent, pattern string) {
	mu := sync.Mutex{}
	mu.Lock()
	if _, ok := rt.routes[broker]; ok {
		var p []IRoutePattern
		for _, v := range rt.routes[broker] {
			if v.Pattern() == pattern {
				continue
			}
			p = append(p, v)
		}
		rt.routes[broker] = p
	}
	mu.Unlock()
}

// Match - находит в таблице маршруты для подключения к брокерам для заданного топика
// возвращает список брокеров для подключения с заданным топиком
func (rt *RouteTable) Match(topic string) (routes []IRoute) {
	mu := sync.Mutex{}
	mu.Lock()
	for _, v := range rt.routes {
		for _, p := range v {
			if p.Match(topic) {
				routes = append(routes, &Route{
					broker: p.Agent(),
					topic:  topic,
				})
			}
		}
	}
	mu.Unlock()
	return
}

func (r *Route) Broker() Agent {
	return r.broker
}

func (r *Route) Route() string {
	return r.topic
}

// NewRoutePattern - создает паттерн для подписки
//
func NewRoutePattern(broker Agent, pattern string) (IRoutePattern, error) {

	r := &RoutePattern{
		broker:  broker,
		pattern: pattern,
		rt:      list.New(),
	}
	if e := r.create(pattern); e != nil {
		return nil, e
	}
	return r, nil
}

// Agent -
func (rp *RoutePattern) Agent() Agent {
	return Agent(rp.broker)
}

// Pattern -
func (rp *RoutePattern) Pattern() string {
	return rp.pattern
}

// Match - проверяет возможность подписки для топика
func (rp *RoutePattern) Match(subscribe string) bool {

	//tp, e  := NewSubscribePattern(topic)
	if len(subscribe) > 1 {
		if rp.rt != nil {
			return rp.find(subscribe)
		}
	}
	return false
}

// Subscribe -
func (rp *RoutePattern) Subscribe(topic string) bool {
	if len(topic) > 1 {
		if rp.rt != nil {
			return rp.subscribe(topic)
		}
	}
	return false
}

// create -
func (rp *RoutePattern) create(route string) error {
	//ptr := list.New()
	rp.rt.Init()

	if strings.Count(route, charTail) > 1 {
		return errInvalidRoutePattern
	}
	isTail := false
	for _, v := range strings.Split(strings.Trim(route, charDelimiter), charDelimiter) {
		if nodeValue, e := createNodeValue(v); e != nil {
			return e
		} else {
			switch nodeValue.Type() {
			case nodeTypeTail:
				if isTail {
					return errInvalidRoutePattern
				}
				isTail = true
			default:

			}
			rp.rt.PushBack(nodeValue)
		}
	}
	return nil
}

// find - найти соответствие
func (rp *RoutePattern) find(topic string) bool {
	tn := strings.Split(strings.Trim(topic, charDelimiter), charDelimiter)
	if len(tn) < 1 {
		return false
	}

	ptr := rp.rt.Front()

	for _, v := range tn {

		if ptr == nil {
			return false
		}

		nv, e := createNodeValue(v)

		if e != nil {
			return false
		}

		switch ptr.Value.(INode).Type() {

		case nodeTypeString:
			switch nv.Type() {
			case nodeTypeTail:
				return true
			case nodeTypePlus:
				ptr = ptr.Next()
				continue
			case nodeTypeString:
				if !ptr.Value.(INode).Validate(v) {
					return false
				}
			default:
				return false
			}

		case nodeTypePlus:
			if !ptr.Value.(INode).Validate(v) {
				return false
			}

		case nodeTypePrefix:

			switch nv.Type() {
			case nodeTypeTail:
				return true
			case nodeTypePlus:
				ptr = ptr.Next()
				continue
			case nodeTypeString:
				if !ptr.Value.(INode).Validate(v) {
					return false
				}
			default:
				return false
			}

		case nodeTypeSuffix:
			switch nv.Type() {
			case nodeTypeTail:
				return true
			case nodeTypePlus:
				ptr = ptr.Next()
				continue
			case nodeTypeString:
				if !ptr.Value.(INode).Validate(v) {
					return false
				}
			default:
				return false
			}

		case nodeTypeTail:
			return ptr.Value.(INode).Validate(v)

		default:
			return false
		}
		ptr = ptr.Next()
	}
	//return ptr.Value == nil
	return true
}

// subscribe - проверить возможность подписки
func (rp *RoutePattern) subscribe(topic string) bool {
	ptr := rp.rt.Front()

	for _, v := range strings.Split(strings.Trim(topic, "/ "), "/") {
		switch v {
		case "+":
			ptr = ptr.Next()
		case "#":
			return true
		default:
			if !ptr.Value.(INode).Validate(v) {
				return false
			}
		}
		ptr = ptr.Next()

		if ptr == nil {
			return false
		}

	}
	return true
}
