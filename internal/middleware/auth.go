package middleware

import(
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

func Auth(secret string) func(http.Handler) http.Handler{
	return func(next http.Handler) http.Handler{
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
			authorizationHeader := r.Header.Get("Authorization")
			if authorizationHeader == ""{
				http.Error(w,"missing authorization header",http.StatusUnauthorized)
				return
			}
			tokenString := strings.TrimPrefix(authorizationHeader,"Bearer ")
			token,err := jwt.Parse(tokenString, func(_ *jwt.Token)(interface{},error){
				return []byte(secret),nil
			})
			if err != nil || !token.Valid{
				http.Error(w,"invalid or expired token",http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w,r)
		})
	}
}