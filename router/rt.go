package router

import (
	"errors"
)

// Errors
var (
	errInvalidNodeValue      = errors.New("Invalid node value. ")
	errInvalidSubscribeTopic = errors.New("Invalid subscribe topic. ")
	errPatternIsPresenr      = errors.New("Pattern already present in list. ")
	errInvalidRoutePattern   = errors.New("Invalid route pattern. ")
	errNotMatched            = errors.New("Not matched. ")
)

type (
	// IRoutePattern - запись содержащая маршрут
	IRoutePattern interface {
		Agent() Agent
		Pattern() string
		Match(subscribe string) error
		Subscribe(topic string) bool
	}

	// ISubscriptionTable -
	ISubscriptionTable interface {
		Add(pattern string) error
		Remove(pattern string)
		Match(topic string) bool
	}

	// ISubscriptionPattern -
	ISubscriptionPattern interface {
		Match(topic string) bool
		Pattern() string
	}

	// IRoute -
	IRoute interface {
		Broker() Agent
		Route() string
	}

	//
	// IRouteTable - таблица маршрутов
	//
	IRouteTable interface {
		// Remove - удаляет маршрут
		Remove(broker Agent, pattern string)
		// Set - добавляет маршрут в таблицу
		Set(broker Agent, pattern string) error
		// Match - возвращает список маршрутов для
		// топика
		Match(topic string) []IRoute
	}
	Agent string
)
