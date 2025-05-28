package middleware

import (
    "net/http"
    "strings"
    "github.com/juancruzestevez/go-gorm-restapi/utils"
    "context"
)

type key string

const UserIDKey key = "userID"

func AuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        authHeader := r.Header.Get("Authorization")
        if authHeader == "" {
            http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
            return
        }

        parts := strings.Split(authHeader, " ")
        if len(parts) != 2 || parts[0] != "Bearer" {
            http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
            return
        }

        tokenStr := parts[1]

        userID, err := utils.ParseJWT(tokenStr)
        if err != nil {
            http.Error(w, "Invalid token: "+err.Error(), http.StatusUnauthorized)
            return
        }

        // Lo guardamos en el contexto para usarlo en handlers
        ctx := context.WithValue(r.Context(), UserIDKey, userID)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

// Función helper para obtener el userID del contexto
func GetUserIDFromContext(ctx context.Context) (uint, bool) {
    userID, ok := ctx.Value(UserIDKey).(uint)
    return userID, ok
}
