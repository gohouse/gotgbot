package router

import (
	"github.com/gohouse/gotgbot/bot"
)

type Route struct {
	Route      string
	Title      string
	//Available  bool
	Visible    bool
	handlers   []bot.HandleFunc
	middleware []bot.HandleFunc
}

func NewRoute(route string, h ...bot.HandleFunc) *Route {
	return &Route{Route: route, handlers: h, Visible: true}
}

func (r *Route) Use(h ...bot.HandleFunc) *Route {
	r.middleware = append(r.middleware, h...)
	return r
}

func (r *Route) SetTitle(title string) *Route {
	r.Title = title
	return r
}

func (r *Route) Show() *Route {
	r.Visible = true
	return r
}

func (r *Route) Hide() *Route {
	r.Visible = false
	return r
}

//func (r *Route) Enable() *Route {
//	r.Available = true
//	return r
//}
//
//func (r *Route) Disable() *Route {
//	r.Available = false
//	return r
//}
