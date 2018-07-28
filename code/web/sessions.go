package web

import (
	"math/rand"
	"SelfGrade/code/security"
)

type Session struct {
	Id string
	user security.User
	attributes map[string]interface{}
}

const SessionCookieName = "GSESSION"

var sessions = make(map[string]Session)

func (session *Session) Authenticated() bool {
	return session.user.Authenticated
}

func NewSession() Session {
	sessionID := generateSessionId()
	session := Session{Id: sessionID, attributes: make(map[string]interface{})}
	sessions[sessionID] = session
	return session
}

func SaveSession(session Session) {
	sessions[session.Id] = session
}

func RemoveSession(id string) bool {
	
	_, found := FindSession(id)

	if found {
		delete(sessions, id)
	}
	
	return found
}

func FindSession(sessionID string) (Session, bool) {
	session, ok := sessions[sessionID]
	return session, ok
}

func (session *Session) AddAttribute(attrName string, attr interface{}) {
	session.attributes[attrName] = attr
}

func (session *Session) GetAttribute(attr string) interface{} {
	return session.attributes[attr]
}

func generateSessionId() string {
	const POOL = "ABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	var id string

	for i := 0; i < 16; i++ {
		index := rand.Intn(len(POOL))
		id += POOL[index:index+1]
	}

	return id
}
