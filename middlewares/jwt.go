package middlewares

import (
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	//"github.com/hexops/valast"
)

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization format must be Bearer {token}"})
			c.Abort()
			return
		}
		tokenString := parts[1]
		token, claims, err := ParseToken(tokenString)
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}
		// Set user ID to context
		//fmt.Println("JWTAuthMiddleware userID=", claims.UserID)
		c.Set("user_id", claims.UserID)
		c.Next()
	}
}

// neung with sha256 hex
var jwtSecret = []byte("9e21758d56efc1bde0694e859a9f350c305f6901063a8c07b0384ff13f76b051")

type Claims struct {
	UserID uint `json:"user_id"`
	jwt.StandardClaims
}

// GenerateToken generates a JWT for a given user ID
func GenerateToken(userID uint) (string, error) {
	claims := Claims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "gotestbackend",
		},
	}
	//fmt.Println("GenerateToken")
	//fmt.Println(valast.String(claims))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// ParseToken parses and validates a JWT
func ParseToken(tokenString string) (*jwt.Token, *Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	return token, claims, err
}

/*
Encoded

eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJleHAiOjE3MTk0NDY1OTIsImlhdCI6MTcxOTM2MDE5MiwiaXNzIjoiZ290ZXN0YmFja2VuZCJ9.B5Cg8CupC1oWkvqazd4gVKDd375NUbiMUb6aUjriLN8

HEADER
{
  "alg": "HS256",
  "typ": "JWT"
}

PAYLOAD
{
  "user_id": 1,
  "exp": 1719446592,
  "iat": 1719360192,
  "iss": "gotestbackend"
}

*/
