package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-ldap/ldap"
	"github.com/go-playground/validator"
	"github.com/labstack/echo"
	mw "github.com/labstack/echo/middleware"
)

var (
	LDAP_SERVER     string   = "10.255.1.254"
	LDAP_DOMAIN     string   = "touzhijia.net"
	LDAP_PORT       uint16   = 389
	LDAP_BASE_DN    string   = "OU=TouZhiJia,DC=touzhijia,DC=net"
	LDAP_FILTER     string   = "(sAMAccountName=%s)"
	LDAP_ATTRIBUTES []string = []string{"name", "sAMAccountName", "mail", "memberOf"}
)

type (
	User struct {
		Username string `json:"username" validate:"required"`
		Password string `json:"password" validate:"required"`
	}

	LAResponse struct {
		Email        string `json:"email,omitempty"`
		Status       bool   `json:"status"`
		Message      string `json:"message,omitempty"`
		Username     string `json:"username,omitempty"`
		Nickname     string `json:"nickname,omitempty"`
		Department   string `json:"department,omitempty"`
		Organization string `json:"organization,omitempty"`
	}

	DataValidator struct {
		validator *validator.Validate
	}
)

func ldapAuth(username string, password string) (int, interface{}) {
	ldap.DefaultTimeout = 2 * time.Second

	// Conn Ldap Server
	_ldap, err := ldap.Dial(
		"tcp", fmt.Sprintf("%s:%d", LDAP_SERVER, LDAP_PORT),
	)

	if err != nil {
		return http.StatusInternalServerError, LAResponse{
			Status:  false,
			Message: "Network Error!",
		}
	}

	defer _ldap.Close()
	// l.Debug = true

	err = _ldap.Bind(fmt.Sprintf("%s@%s", username, LDAP_DOMAIN), password)
	if err != nil {
		return http.StatusUnauthorized, LAResponse{
			Status:  false,
			Message: "Username or Password Error!",
		}
	}

	_search := ldap.NewSearchRequest(
		LDAP_BASE_DN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf(LDAP_FILTER, username),
		LDAP_ATTRIBUTES,
		nil,
	)

	_searchResult, err := _ldap.Search(_search)
	if err != nil {
		return http.StatusInternalServerError, LAResponse{
			Status:  false,
			Message: "Server Error!",
		}
	}

	var res LAResponse
	for _, entry := range _searchResult.Entries {
		var (
			attrMail     string = entry.GetAttributeValue("mail")
			attrUsername string = entry.GetAttributeValue("sAMAccountName")
			attrNickname string = entry.GetAttributeValue("name")
			attrMemberOf string = entry.GetAttributeValue("memberOf")
		)
		if attrMail == "" || attrUsername == "" || attrNickname == "" || attrMemberOf == "" {
			continue
		}

		res = LAResponse{
			Email:        attrMail,
			Nickname:     attrNickname,
			Username:     attrUsername,
			Organization: strings.Replace(strings.Split(attrMemberOf, ",")[1], "OU=", "", 1),
			Department:   strings.Replace(strings.Split(attrMemberOf, ",")[0], "CN=", "", 1),
			Status:       true,
		}

		break
	}

	return http.StatusOK, res
}

func ldapUserList() (int, interface{}) {
	ldap.DefaultTimeout = 2 * time.Second

	// Conn Ldap Server
	_ldap, err := ldap.Dial(
		"tcp", fmt.Sprintf("%s:%d", LDAP_SERVER, LDAP_PORT),
	)

	if err != nil {
		return http.StatusInternalServerError, LAResponse{
			Status:  false,
			Message: "Network Error!",
		}
	}

	defer _ldap.Close()
	// l.Debug = true

	err = _ldap.Bind("auth@touzhijia.net", "iw8u55a6Z@soerlt7")
	if err != nil {
		return http.StatusUnauthorized, LAResponse{
			Status:  false,
			Message: "Password Error!",
		}
	}

	_search := ldap.NewSearchRequest(
		LDAP_BASE_DN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf(LDAP_FILTER, "*"),
		LDAP_ATTRIBUTES,
		nil,
	)

	_searchResult, err := _ldap.Search(_search)
	if err != nil {
		return http.StatusInternalServerError, LAResponse{
			Status:  false,
			Message: "Server Error!",
		}
	}

	var users = make([]interface{}, 0)
	for _, entry := range _searchResult.Entries {
		var (
			attrMail     string = entry.GetAttributeValue("mail")
			attrUsername string = entry.GetAttributeValue("sAMAccountName")
			attrNickname string = entry.GetAttributeValue("name")
			attrMemberOf string = entry.GetAttributeValue("memberOf")
		)
		if attrMail == "" || attrUsername == "" || attrNickname == "" || attrMemberOf == "" {
			continue
		}

		users = append(users, LAResponse{
			Email:        attrMail,
			Nickname:     attrNickname,
			Username:     attrUsername,
			Organization: strings.Replace(strings.Split(attrMemberOf, ",")[1], "OU=", "", 1),
			Department:   strings.Replace(strings.Split(attrMemberOf, ",")[0], "CN=", "", 1),
			Status:       true,
		})
	}

	return http.StatusOK, users
}

func (dv *DataValidator) Validate(i interface{}) error {
	return dv.validator.Struct(i)
}

func main() {
	e := echo.New()
	e.Use(mw.Logger())
	e.Use(mw.Recover())
	e.Validator = &DataValidator{validator: validator.New()}

	e.GET("/", func(c echo.Context) error {
		return c.JSON(
			http.StatusNotFound,
			LAResponse{
				Status:  false,
				Message: "Not Found.",
			},
		)
	})

	e.GET("/api", func(c echo.Context) error {
		return c.JSON(
			http.StatusOK,
			LAResponse{
				Status:  true,
				Message: "LDAP Auth Service.",
			},
		)
	})

	e.POST("/api/auth", func(c echo.Context) (err error) {
		u := new(User)
		if err = c.Bind(u); err != nil {
			return c.String(http.StatusOK, err.Error())
		}

		if err = c.Validate(u); err != nil {
			e.Logger.Printf("Request Data invalid:", err.Error())
			return c.JSON(
				http.StatusBadRequest,
				LAResponse{
					Status:  false,
					Message: "Request Data invalid",
				},
			)
		}

		resCode, resData := ldapAuth(u.Username, u.Password)
		return c.JSON(resCode, resData)
	})

	e.GET("/api/users", func(c echo.Context) (err error) {
		resCode, resData := ldapUserList()

		return c.JSON(resCode, resData)
	})

	e.Logger.Fatal(e.Start(":8389"))
}
