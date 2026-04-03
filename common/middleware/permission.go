package middleware

import (
	"net/http"

	"github.com/casbin/casbin/v2/util"
	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/response"
)

/*
## How it works in this repo
`AuthCheckRole()` does this order:
1. Read JWT payload
2. If admin role -> allow
3. If path+method is in `CasbinExclude` -> allow
4. Else call `e.Enforce(rolekey, path, method)` (Casbin)
5. If denied -> `403`

```text
Request
  |
  |- JWT middleware (identity/auth)
  |
  \- AuthCheckRole
       |- admin? yes -> pass
       |- in CasbinExclude? yes -> pass
       \- Casbin Enforce(role, path, method)
             |- true -> pass
             \- false -> 403
```

So the exclusion list prevents valid users from getting blocked on endpoints that are intentionally "shared" (profile/info/tree/options, etc.) or public-ish infra endpoints.

## Why this exists (practical reasons)

- Avoid creating/maintaining tons of Casbin policies for endpoints everyone should use.
- Prevent bootstrap deadlocks (logged in, but cannot load base UI info/menu trees).
- Keep backward compatibility with older route policy expectations.
*/
func AuthCheckRole() gin.HandlerFunc{
	return func(c *gin.Context)  {
		log := api.GetRequestLogger(c)
		data, _ := c.Get(jwtauth.JwtPayloadKey)
		v := data.(jwtauth.MapClaims)
		e := sdk.Runtime.GetCasbinByTenant(c.Request.Host)

		var res, casbinExclude bool
		var err error

				if v["rolekey"] == "admin" {
			res = true
			c.Next()
			return
		}
		for _, i := range CasbinExclude {
			if util.KeyMatch2(c.Request.URL.Path, i.Url) && c.Request.Method == i.Method {
				casbinExclude = true
				break
			}
		}
		if casbinExclude {
			log.Infof("Casbin exclusion, no validation method:%s path:%s", c.Request.Method, c.Request.URL.Path)
			c.Next()
			return
		}
		res, err = e.Enforce(v["rolekey"], c.Request.URL.Path, c.Request.Method)
		if err != nil {
			log.Errorf("AuthCheckRole error:%s method:%s path:%s", err, c.Request.Method, c.Request.URL.Path)
			response.Error(c, 500, err, "")
			return
		}

		if res {
			log.Infof("isTrue: %v role: %s method: %s path: %s", res, v["rolekey"], c.Request.Method, c.Request.URL.Path)
			c.Next()
		} else {
			log.Warnf("isTrue: %v role: %s method: %s path: %s message: %s", res, v["rolekey"], c.Request.Method, c.Request.URL.Path, "current request has no permission, please contact the administrator")
			c.JSON(http.StatusOK, gin.H{
				"code": 403,
				"msg":  "Sorry, you do not have permission to access this endpoint. Please contact the administrator.",
			})
			c.Abort()
			return
		}
	}
}