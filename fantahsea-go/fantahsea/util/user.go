package util

import (
	"strings"

	"fantahsea/err"

	"github.com/gin-gonic/gin"
)

type User struct {
	UserId   string
	UserNo   string
	Username string
	Role     string
	Services []string
}

/* Extract User from request headers */
func ExtractUser(c *gin.Context) (*User, error) {
	id := c.GetHeader("id")
	if id == "" {
		return nil, err.NewWebErr("Please sign up first")
	}

	var services []string
	servicesStr := c.GetHeader("services")
	if servicesStr == "" {
		services = make([]string, 0)
	} else {
		services = strings.Split(servicesStr, ",")
	}

	return &User{
		UserId:   id,
		Username: c.GetHeader("username"),
		UserNo:   c.GetHeader("userno"),
		Role:     c.GetHeader("role"),
		Services: services,
	}, nil
}

// Check if the user is a guest
func IsGuest(user *User) bool {
	if user == nil {
		panic("user == nil")
	}

	return user.Role == "guest"
}
