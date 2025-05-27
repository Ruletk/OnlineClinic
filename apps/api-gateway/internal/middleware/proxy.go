package middleware

import (
	"fmt"
	"github.com/Ruletk/OnlineClinic/pkg/logging"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type ServiceConfig struct {
	Name string
	Host string
	Port string
}

func ReverseProxy(service *ServiceConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		remote, err := url.Parse(fmt.Sprintf("http://%s:%s", service.Host, service.Port))
		if err != nil {
			logging.Logger.Errorf("Failed to parse remote URL: %v", err)
			c.AbortWithStatusJSON(500, gin.H{"error": "Internal Server Error"})
			return
		}

		proxy := httputil.NewSingleHostReverseProxy(remote)
		proxy.Director = func(req *http.Request) {
			req.Header = c.Request.Header
			req.Host = remote.Host
			req.URL.Scheme = remote.Scheme
			req.URL.Host = remote.Host
			req.URL.Path = c.Param("proxyPath")

			if accessToken, exists := c.Get("access_token"); exists {
				req.Header.Set("X-Access-Token", accessToken.(string))
			}
		}

		proxy.ServeHTTP(c.Writer, c.Request)
	}
}
