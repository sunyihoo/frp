package server

type ControlManager struct {
	ctlsByRunID map[string]*Cont
}

type Control struct {
	rc *controller
}
