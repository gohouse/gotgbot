package router

import (
	"fmt"
	"github.com/gohouse/gotgbot/bot"
	"regexp"
	"strings"
)

type RouteGroup struct {
	module        bot.ModuleType
	routes        []*Route
	children      []*RouteGroup
	Middleware    []bot.HandleFunc
	CallbackQuery []bot.HandleFunc
	Wait          []bot.HandleFunc
	HandBot       []bot.HandleBot
}

func NewRouteGroup() (r *RouteGroup) {
	return &RouteGroup{}
}

func newRouteGroup(module ...bot.ModuleType) *RouteGroup {
	var mod = bot.MOTNil
	if len(module) > 0 {
		mod = module[0]
	}
	return &RouteGroup{module: mod}
}

func (rg *RouteGroup) ListenBot() *RouteGroup {
	if rg.module != bot.MOTNil {
		panic("please use root RouteGroup")
	}
	if rg.checkModuleExists(rg.module) {
		log.Panicf("module already exists: %s", rg.module.String())
	}
	listener := newRouteGroup(bot.MOTPrivate)
	rg.children = append(rg.children, listener)
	return listener
}
func (rg *RouteGroup) ListenChannel() *RouteGroup {
	if rg.module != bot.MOTNil {
		panic("please use root RouteGroup")
	}
	if rg.checkModuleExists(rg.module) {
		log.Panicf("module already exists: %s", rg.module.String())
	}
	listener := newRouteGroup(bot.MOTChannel)
	rg.children = append(rg.children, listener)
	return listener
}
func (rg *RouteGroup) ListenGroup() *RouteGroup {
	if rg.module != bot.MOTNil {
		panic("please use root RouteGroup")
	}
	if rg.checkModuleExists(rg.module) {
		log.Panicf("module already exists: %s", rg.module.String())
	}
	listener := newRouteGroup(bot.MOTSuperGroup)
	rg.children = append(rg.children, listener)
	return listener
}
func (rg *RouteGroup) ListenSuperGroup() *RouteGroup {
	if rg.module != bot.MOTNil {
		panic("please use root RouteGroup")
	}
	if rg.checkModuleExists(rg.module) {
		log.Panicf("module already exists: %s", rg.module.String())
	}
	listener := newRouteGroup(bot.MOTSuperGroup)
	rg.children = append(rg.children, listener)
	return listener
}

func (rg *RouteGroup) Use(h ...bot.HandleFunc) *RouteGroup {
	rg.Middleware = append(rg.Middleware, h...)
	return rg
}

func (rg *RouteGroup) GET(route string, h ...bot.HandleFunc) *Route {
	if rg.checkRouteExists(rg, strings.TrimLeft(route, "/"), rg.module, true) {
		log.Panicf("route already exists: %s", route)
	}

	routeSub := NewRoute(route, h...)
	rg.routes = append(rg.routes, routeSub)
	return routeSub
}

func (rg *RouteGroup) BuildHelpList(mot bot.ModuleType) (res []string) {
	var list = make(map[string]string)
	for _, v := range rg.routes {
		if v.Visible {
			list[v.Route] = v.Title
		}
	}
	for _, v2 := range rg.children {
		if v2.module == mot {
			for _, v := range v2.routes {
				if v.Visible {
					list[v.Route] = v.Title
				}
			}
			break
		}
	}
	for k, v := range list {
		res = append(res, fmt.Sprintf("%s %s", k, v))
	}
	return
}

func (rg *RouteGroup) checkModuleExists(mot bot.ModuleType) bool {
	for _, v := range rg.children {
		if v.module == mot {
			return true
		}
	}
	return false
}

//func (rg *RouteGroup) CheckRouteExists(rgs *RouteGroup, route string, mot bot.ModuleType) bool {
//	if rg.checkRouteExists(rgs, route, mot, true) {
//		return true
//	}
//	return rg.checkRouteExists(rgs, route, mot, false)
//}
func (*RouteGroup) checkRouteExists(rgs *RouteGroup, route string, mot bot.ModuleType, mustMatch bool) bool {
	for _, v := range rgs.routes {
		if mot != rgs.module {
			break
		}
		if mustMatch {
			if v.Route == route {
				return true
			}
		} else {
			all := strings.ReplaceAll(route, ":id", `(\w+)`)
			exp := regexp.MustCompile(all)
			if exp.Match([]byte(v.Route)) {
				return true
			}
		}
	}
	if len(rgs.children) > 0 {
		for _, v2 := range rgs.children {
			if mot != v2.module {
				break
			}
			for _, v := range v2.routes {
				if mustMatch {
					if v.Route == route {
						return true
					}
				} else {
					all := strings.ReplaceAll(route, ":id", `(\w+)`)
					exp := regexp.MustCompile(all)
					if exp.Match([]byte(v.Route)) {
						return true
					}
				}
			}
		}
	}
	return false
}

func (rg *RouteGroup) ExtractRoute(route string, mot bot.ModuleType) (m []bot.HandleFunc, h []bot.HandleFunc) {
	route = fmt.Sprintf("/%s", (strings.Split(route, "@"))[0])
	//// 打印一下
	//log.Info(rg.module.String())
	//for _,v2 := range rg.routes {
	//	log.Info(v2.Route, v2.Title)
	//}
	//for _,v := range rg.children {
	//	log.Info(v.module.String())
	//	for _,v2 := range v.routes {
	//		log.Info(v2.Route, v2.Title)
	//	}
	//}
	m2, h2 := rg.extractRoute(route, mot, true)
	if len(h2) > 0 {
		return m2, h2
	}
	return rg.extractRoute(route, mot, false)
}
func (rg *RouteGroup) extractRoute(route string, mot bot.ModuleType, mustMatch bool) (m []bot.HandleFunc, h []bot.HandleFunc) {
	var hand []bot.HandleFunc
	var cors = rg.Middleware // 合并中间件
	for _, v := range rg.routes {
		if mustMatch {
			if v.Route == route {
				hand = v.handlers[:]
				break
			}
		} else {
			all := strings.ReplaceAll(v.Route, ":id", `(\w+)`)
			exp := regexp.MustCompile(all)
			if exp.Match([]byte(route)) {
				hand = v.handlers[:]
				break
			}
		}
	}
	if len(rg.children) > 0 {
	LOOP:
		for _, v2 := range rg.children {
			if v2.module != mot {
				continue
			}
			// 合并共用中间件
			cors = append(cors, v2.Middleware...)
			for _, v := range v2.routes {
				if mustMatch {
					if v.Route == route {
						hand = v.handlers[:]
						// 合并当前路由中间件
						cors = append(cors, v.middleware...)
						break LOOP
					}
				} else {
					all := strings.ReplaceAll(v.Route, ":id", `(\w+)`)
					exp := regexp.MustCompile(all)
					if exp.Match([]byte(route)) {
						hand = v.handlers[:]
						// 合并当前路由中间件
						cors = append(cors, v.middleware...)
						break
					}
				}
			}
		}
	}
	return cors, hand
}

func (rg *RouteGroup) OnCallbackQuery(h ...bot.HandleFunc) *RouteGroup {
	rg.CallbackQuery = append(rg.CallbackQuery, h...)
	return rg
}
func (rg *RouteGroup) OnWait(h ...bot.HandleFunc) *RouteGroup {
	rg.Wait = append(rg.Wait, h...)
	return rg
}
func (rg *RouteGroup) OnInlineQuery(h ...bot.HandleFunc) *RouteGroup {
	rg.Wait = append(rg.Wait, h...)
	return rg
}

func (rg *RouteGroup) Do(h ...bot.HandleBot) *RouteGroup {
	rg.HandBot = append(rg.HandBot, h...)
	return rg
}
