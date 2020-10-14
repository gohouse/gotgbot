package router

import (
	"github.com/gohouse/gotgbot/bot"
)

type Route struct {
	Route      string
	Title      string
	handlers   []bot.HandleFunc
	middleware []bot.HandleFunc
}

func NewRoute(route string, title string, h ...bot.HandleFunc) *Route {
	return &Route{Route: route, Title: title, handlers: h}
}

func (r *Route) Use(h ...bot.HandleFunc) *Route {
	r.middleware = append(r.middleware, h...)
	return r
}

func (r *Route) SetTitle(title string) *Route {
	r.Title = title
	return r
}
