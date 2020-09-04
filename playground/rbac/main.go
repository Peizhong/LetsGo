package main

import (
	"log"
	"net/http"

	defaultrolemanager "github.com/casbin/casbin/rbac/default-role-manager"
	"github.com/casbin/casbin/util"
	"github.com/casbin/casbin/v2"
	"github.com/gin-contrib/authz"
	"github.com/gin-gonic/gin"
)

func main() {
	e, _ := casbin.NewEnforcer("rbac_model.conf", "rbac_policy.csv")
	if res, _ := e.Enforce("alice", "data1", "read"); res {
		log.Println("ok")
	} else {
		log.Println("not ok")
	}
	alicep, _ := e.GetImplicitPermissionsForUser("alice")
	log.Println(alicep)
	alicer, _ := e.GetImplicitRolesForUser("alice")
	log.Println(alicer)
	r := gin.New()
	auth := r.Group("/admin", gin.BasicAuth(gin.Accounts{
		"foo": "bar",
	}))                                                                                 
	auth.Use(authz.NewAuthorizer(e))
	auth.GET("/hi", func(g *gin.Context) {
		user := g.MustGet(gin.AuthUserKey).(string)
		if user != "" {
			g.JSON(http.StatusOK, gin.H{"wpz": "你好" + user})
		} else {
			g.JSON(http.StatusOK, gin.H{"wpz": "你是谁"})
		}
	})
	r.Run()
}

func NewRoleManager() {
	rm := defaultrolemanager.NewRoleManager(10)
	rm.(*defaultrolemanager.RoleManager).AddMatchingFunc("", util.KeyMatch2)
}
