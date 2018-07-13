package web

import (
	"math/rand"
	"SelfGrade/code/security"
	"SelfGrade/code/utils"
)

type Session struct {
	Id string
	User security.User
	Authenticated bool
	attributes map[string]interface{}
}

const SessionCookieName = "GSESSION"

var sessions = make(map[string]Session)

func NewSession() Session {
	sessionID := generateSessionId()
	session := Session{Id: sessionID, attributes: make(map[string]interface{})}
	sessions[sessionID] = session
	return session
}

func SaveSession(session Session) {
	sessions[session.Id] = session
}

func FindSession(sessionID string) (Session, error) {
	session, ok := sessions[sessionID]

	if !ok {
		return session, utils.Error{Message:"Session with ID " + sessionID + " not found.", Cause:nil}
	}

	return session, nil
}

func (session Session) AddAttribute(attrName string, attr interface{}) {
	session.attributes[attrName] = attr
}


func (session Session) GetAttribute(attr string) (interface{}, error) {
	attribute, ok := session.attributes[attr]

	if !ok {
		return attribute, utils.Error{Message:"Session attribute with ID " + attr + " not found.", Cause:nil}
	}

	return attribute, nil
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
