package config

import (
	"fmt"
	"net/http"
)

const (
	// ParamBoardID is name of path param which represents board identifier
	ParamBoardID = "boardId"
	// ParamUserID is name of path param which represents user identifier
	ParamUserID = "userId"
	// ParamBoardMemberID is name of path param which represents board member identifier
	ParamBoardMemberID = "boardMemberId"
	// ParamTaskOrder is name of path param which represents order of task on board
	ParamTaskOrder = "taskOrder"
)

// URLPaths defines url paths which used by app router
type URLPaths struct {
	// BoardsHandler is url path to handlers.BoardHandler methods for working with multiple records
	BoardsHandler string
	// BoardsHandler is url path to handlers.BoardHandler methods for working with single record
	BoardHandler         string
	BoardMembersHandler  string
	BoardMemberHandler   string
	RefreshAccessHandler string
	RegistrationHandler  string
	LoginHandler         string
	LogoutHandler        string
	TasksHandler         string
	TaskHandler          string
	UsersHandler         string
	UserHandler          string
}

// AllowedHTTPMethods defines allowed http methods for handlers in URLPaths
type AllowedHTTPMethods struct {
	BoardsHandler        []string
	BoardHandler         []string
	BoardMembersHandler  []string
	BoardMemberHandler   []string
	LoginHandler         []string
	LogoutHandler        []string
	RefreshAccessHandler []string
	RegistrationHandler  []string
	UsersHandler         []string
	UserHandler          []string
	TasksHandler         []string
	TaskHandler          []string
}

// NewHTTPPaths returns config for working with http routing in app
func NewHTTPPaths() (*URLPaths, *AllowedHTTPMethods) {
	paths := &URLPaths{
		LoginHandler:         "/login",
		LogoutHandler:        "/logout",
		RefreshAccessHandler: "/refresh-access",
		RegistrationHandler:  "/registration",
		UsersHandler:         "/users",
		BoardsHandler:        "/boards",
		UserHandler:          fmt.Sprintf("/users/{%s}", ParamUserID),
		BoardHandler:         fmt.Sprintf("/boards/{%s}", ParamBoardID),
		BoardMembersHandler:  fmt.Sprintf("/boards/{%s}/members", ParamBoardID),
		BoardMemberHandler:   fmt.Sprintf("/boards/{%s}/members/{%s}", ParamBoardID, ParamBoardMemberID),
		TasksHandler:         fmt.Sprintf("/boards/{%s}/tasks", ParamBoardID),
		TaskHandler:          fmt.Sprintf("/boards/{%s}/tasks/{%s}", ParamBoardID, ParamTaskOrder),
	}
	allowedMethods := &AllowedHTTPMethods{
		BoardsHandler:        []string{http.MethodGet, http.MethodPost},
		BoardHandler:         []string{http.MethodGet, http.MethodPatch, http.MethodDelete},
		BoardMembersHandler:  []string{http.MethodGet, http.MethodPost},
		BoardMemberHandler:   []string{http.MethodGet, http.MethodPatch, http.MethodDelete},
		LoginHandler:         []string{http.MethodPost},
		LogoutHandler:        []string{http.MethodPost},
		RefreshAccessHandler: []string{http.MethodPost},
		RegistrationHandler:  []string{http.MethodPost},
		UsersHandler:         []string{http.MethodGet},
		UserHandler:          []string{http.MethodGet, http.MethodPatch, http.MethodDelete},
		TasksHandler:         []string{http.MethodGet, http.MethodPost, http.MethodDelete},
	}
	return paths, allowedMethods
}
