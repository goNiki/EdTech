package auth

import "net/http"

func (m *middleware) RequireRole(roles []string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			role := m.GetUserRole(r.Context())
			if role == "" {
				http.Error(w, "missing role", http.StatusForbidden)
				return
			}

			for _, allowed := range roles {
				if allowed == role {
					next.ServeHTTP(w, r)
					return
				}

			}
			http.Error(w, "forbidden", http.StatusForbidden)
		})
	}
}
