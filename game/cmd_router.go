package game

import (
	"errors"

)

type HandlerFunc func(*CmdContext)

type CmdContext struct {
	// The game object, for fetching data about other players
	Game *Game

	// The current player
	Client *Client

	//The text of the command
	Command string

	// Any arguments passed to the Command
	Args string
}

type Route struct {
	// Slice of subroutes
	Routes []*Route

	Name string
	Description string

	// Route parent
	Parent *Route

	Handler HandlerFunc
	Matcher func(string) bool
}

// Sets the Description, for command chaining
func (r *Route) Desc(description string) *Route {
	r.Description = description

	return r
}

// Creates a new Context
func NewContext(g *Game, c *Client, cmd string, args string) *CmdContext {
	return &CmdContext {
		Game: g,
			Client: c,
			Command: cmd,
			Args: args,
	}
}

// Creates a new Router
func NewRouter() *Route {
	return &Route {
		Routes: []*Route{},
	}
}

// Registers a route with the provided Handler function
// name : Name of the command, what players type to run it
// handler : Function to be run on match
func (r *Route) On(name string, handler HandlerFunc) *Route {
	rt := r.OnMatch(name, nil, handler)
	rt.Matcher = NewNameMatcher(rt)
	return rt
}

// Finds a route that matches the command and executes it if it exists
// Creates a context and passes it to the HandlerFunc
//
// returns false if no route is found
//
// g : Game object, passed to Context
// p : Player object, passed to Context
// cmd : Command used by player
// args : Arguments passed to the command by the player
func (r *Route) FindAndExecute(g *Game, c *Client, cmd string, args string) bool {
	if rt := r.Find(cmd); rt != nil {
		rt.Handler(NewContext(g, c, cmd, args))
		return true
	}

	return false
}

func (r *Route) OnMatch(name string, matcher func(string) bool, handler HandlerFunc) *Route {
	if rt := r.Find(name); rt != nil {
		return rt
	}

	rt := &Route {
		Name: name,
			Handler: handler,
			Matcher: matcher,
	}
	r.AddRoute(rt)
	return rt
}

func (r *Route) AddRoute(route *Route) error {
	if rt := r.Find(route.Name); rt != nil {
		return errors.New("Route already exists")
	}

	route.Parent = r
	r.Routes = append(r.Routes, route)
	return nil
}

// Finds a route with the given name
// Returns nil if nothing is found
func (r *Route) Find(name string) *Route {
	for _, v := range r.Routes {
		if v.Matcher(name) {
			return v
		}
	}
	return nil
}
