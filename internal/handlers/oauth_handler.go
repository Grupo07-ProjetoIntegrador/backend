package handlers

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/Grupo07-ProjetoIntegrador/backend/internal/database"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func oauthConfig() (*oauth2.Config, error) {
	clientID := os.Getenv("GOOGLE_CLIENT_ID")
	clientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	redirectURL := os.Getenv("GOOGLE_OAUTH_REDIRECT_URL")
	if clientID == "" || clientSecret == "" || redirectURL == "" {
		return nil, fmt.Errorf("variaveis GOOGLE_CLIENT_ID, GOOGLE_CLIENT_SECRET e GOOGLE_OAUTH_REDIRECT_URL sao obrigatorias")
	}

	return &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Scopes: []string{
			"https://www.googleapis.com/auth/forms.body",
			"https://www.googleapis.com/auth/drive",
			"https://www.googleapis.com/auth/gmail.send",
		},
		Endpoint: google.Endpoint,
	}, nil
}

func generateState() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func tokenHasRequiredScopes(accessToken string, requiredScopes []string) (bool, error) {
	if accessToken == "" {
		return false, fmt.Errorf("access token vazio")
	}

	endpoint := fmt.Sprintf("https://oauth2.googleapis.com/tokeninfo?access_token=%s", url.QueryEscape(accessToken))
	resp, err := http.Get(endpoint)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	if resp.StatusCode >= 400 {
		return false, fmt.Errorf("tokeninfo retornou status %d: %s", resp.StatusCode, string(body))
	}

	var payload struct {
		Scope string `json:"scope"`
	}
	if err := json.Unmarshal(body, &payload); err != nil {
		return false, err
	}

	hasAll := true
	for _, required := range requiredScopes {
		if !strings.Contains(payload.Scope, required) {
			hasAll = false
			break
		}
	}

	return hasAll, nil
}

// GoogleOAuthStartHandler inicia o fluxo OAuth para um usuario
func GoogleOAuthStartHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, "Método não permitido. Use GET.", http.StatusMethodNotAllowed)
		return
	}

	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		http.Error(w, "O user_id é obrigatório", http.StatusBadRequest)
		return
	}

	config, err := oauthConfig()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	state, err := generateState()
	if err != nil {
		http.Error(w, "Erro ao gerar state", http.StatusInternalServerError)
		return
	}

	expiresAt := time.Now().Add(10 * time.Minute)
	_, err = database.DB.Exec(
		`INSERT INTO oauth_states (state, user_id, expires_at) VALUES ($1, $2, $3)`,
		state,
		userID,
		expiresAt,
	)
	if err != nil {
		http.Error(w, "Erro ao salvar state", http.StatusInternalServerError)
		return
	}

	url := config.AuthCodeURL(
		state,
		oauth2.AccessTypeOffline,
		oauth2.ApprovalForce,
		oauth2.SetAuthURLParam("prompt", "consent"),
	)
	http.Redirect(w, r, url, http.StatusFound)
}

// GoogleOAuthCallbackHandler recebe o callback do Google
func GoogleOAuthCallbackHandler(w http.ResponseWriter, r *http.Request) {
	config, err := oauthConfig()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	state := r.URL.Query().Get("state")
	code := r.URL.Query().Get("code")
	if state == "" || code == "" {
		http.Error(w, "State ou code ausente", http.StatusBadRequest)
		return
	}

	var userID string
	var expiresAt time.Time
	row := database.DB.QueryRow(
		`SELECT user_id, expires_at FROM oauth_states WHERE state = $1`,
		state,
	)
	if err := row.Scan(&userID, &expiresAt); err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "State invalido", http.StatusBadRequest)
			return
		}
		http.Error(w, "Erro ao validar state", http.StatusInternalServerError)
		return
	}

	if time.Now().After(expiresAt) {
		_, _ = database.DB.Exec(`DELETE FROM oauth_states WHERE state = $1`, state)
		http.Error(w, "State expirado", http.StatusBadRequest)
		return
	}

	_, _ = database.DB.Exec(`DELETE FROM oauth_states WHERE state = $1`, state)

	token, err := config.Exchange(r.Context(), code)
	if err != nil {
		frontendURL := os.Getenv("FRONTEND_BASE_URL")
		if frontendURL == "" {
			frontendURL = "http://localhost:5173"
		}
		http.Redirect(w, r, fmt.Sprintf("%s/perfil?google=error", frontendURL), http.StatusFound)
		return
	}

	if ok, scopeErr := tokenHasRequiredScopes(token.AccessToken, config.Scopes); scopeErr != nil || !ok {
		_, _ = database.DB.Exec(`DELETE FROM google_oauth_tokens WHERE user_id = $1`, userID)
		frontendURL := os.Getenv("FRONTEND_BASE_URL")
		if frontendURL == "" {
			frontendURL = "http://localhost:5173"
		}
		errorParam := "missing_scope"
		if scopeErr != nil {
			errorParam = "missing_scope"
		}
		http.Redirect(w, r, fmt.Sprintf("%s/perfil?google=%s", frontendURL, errorParam), http.StatusFound)
		return
	}

	expires := token.Expiry
	scopes := strings.Join(config.Scopes, " ")
	_, err = database.DB.Exec(
		`INSERT INTO google_oauth_tokens (user_id, access_token, refresh_token, token_type, scope, expires_at)
		 VALUES ($1, $2, $3, $4, $5, $6)
		 ON CONFLICT (user_id) DO UPDATE SET
		 access_token = EXCLUDED.access_token,
		 refresh_token = COALESCE(NULLIF(EXCLUDED.refresh_token, ''), google_oauth_tokens.refresh_token),
		 token_type = EXCLUDED.token_type,
		 scope = EXCLUDED.scope,
		 expires_at = EXCLUDED.expires_at,
		 updated_at = NOW()`,
		userID,
		token.AccessToken,
		token.RefreshToken,
		token.TokenType,
		scopes,
		expires,
	)
	if err != nil {
		http.Error(w, "Erro ao salvar tokens", http.StatusInternalServerError)
		return
	}

	frontendURL := os.Getenv("FRONTEND_BASE_URL")
	if frontendURL == "" {
		frontendURL = "http://localhost:5173"
	}

	http.Redirect(w, r, fmt.Sprintf("%s/perfil?google=connected", frontendURL), http.StatusFound)
}

// GoogleOAuthStatusHandler retorna se o usuario esta conectado
func GoogleOAuthStatusHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, "Método não permitido. Use GET.", http.StatusMethodNotAllowed)
		return
	}

	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		http.Error(w, "O user_id é obrigatório", http.StatusBadRequest)
		return
	}

	var expiresAt time.Time
	row := database.DB.QueryRow(
		`SELECT expires_at FROM google_oauth_tokens WHERE user_id = $1`,
		userID,
	)

	connected := true
	if err := row.Scan(&expiresAt); err != nil {
		if err == sql.ErrNoRows {
			connected = false
		} else {
			http.Error(w, "Erro ao consultar tokens", http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"connected":  connected,
		"expires_at": expiresAt,
	})
}

// GoogleOAuthDisconnectHandler revoga o acesso do Google e remove os tokens armazenados.
func GoogleOAuthDisconnectHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodDelete {
		http.Error(w, "Método não permitido. Use DELETE.", http.StatusMethodNotAllowed)
		return
	}

	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		http.Error(w, "O user_id é obrigatório", http.StatusBadRequest)
		return
	}

	var accessToken string
	var refreshToken string
	row := database.DB.QueryRow(
		`SELECT access_token, refresh_token FROM google_oauth_tokens WHERE user_id = $1`,
		userID,
	)
	if err := row.Scan(&accessToken, &refreshToken); err != nil {
		if err == sql.ErrNoRows {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]any{"disconnected": true, "revoked": false})
			return
		}
		http.Error(w, "Erro ao consultar tokens", http.StatusInternalServerError)
		return
	}

	tokenToRevoke := refreshToken
	if tokenToRevoke == "" {
		tokenToRevoke = accessToken
	}

	if tokenToRevoke != "" {
		revokeURL := "https://oauth2.googleapis.com/revoke"
		form := url.Values{}
		form.Set("token", tokenToRevoke)

		resp, err := http.PostForm(revokeURL, form)
		if err != nil {
			http.Error(w, "Erro ao revogar acesso no Google", http.StatusBadGateway)
			return
		}
		defer resp.Body.Close()
		_, _ = io.ReadAll(resp.Body)

		if resp.StatusCode >= 400 {
			http.Error(w, "Google recusou a revogação do acesso", http.StatusBadGateway)
			return
		}
	}

	if _, err := database.DB.Exec(`DELETE FROM google_oauth_tokens WHERE user_id = $1`, userID); err != nil {
		http.Error(w, "Erro ao remover tokens locais", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"disconnected": true,
		"revoked":      tokenToRevoke != "",
	})
}
