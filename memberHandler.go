package main

import (
	"fmt"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func addLoginCookie(memberData Member, c *gin.Context) bool {
	email := memberData.Email
	result, err := findMemberByEmail(email)
	if err != nil {
		fmt.Printf("[NOTE] new member: %v", err)
		insertResult, err := insertMember(memberData)
		if err != nil {
			return false
		} else {
			session := sessions.Default(c)
			session.Set("id", insertResult.Id)
			session.Save()
			return true
		}
	} else {
		session := sessions.Default(c)
		fmt.Println("id:", result.Id)
		session.Set("id", result.Id)
		session.Save()
		return true
	}
}
