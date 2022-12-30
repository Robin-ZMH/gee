package gee

type routerGroup struct {
	// embedding router
	*router
	*Engine
	// new fields
	prefix      string
	parent      *routerGroup
}

func newRouterGroup() *routerGroup{
	return &routerGroup{router: newRouter()}
}


