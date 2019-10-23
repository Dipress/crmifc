package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/dipress/crmifc/internal/kit/auth"
)

func TestArticleCreate(t *testing.T) {
	t.Log("with prepared server")
	{
		db, teardown := postgresDB(t)
		defer teardown()

		ctx, cancel := context.WithTimeout(context.Background(), caseTimeout)
		defer cancel()

		claims := auth.NewClaims("admin@example.com", time.Now(), time.Hour)

		token, err := authenticator.GenerateToken(ctx, claims)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		token = "Bearer " + token

		lis, err := net.Listen("tcp", "127.0.0.1:0")
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}

		s := setupServer(lis.Addr().String(), db, authenticator)
		go s.Serve(lis)
		defer s.Close()

		t.Log("\ttest:0\tshould create a article.")
		{
			articleStr := `{"category_id":1,"title":"my title","body":"my body"}`
			req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("http://%s/articles", s.Addr), strings.NewReader(articleStr))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Add("Authorization", token)

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if resp.StatusCode != http.StatusOK {
				t.Errorf("unexpected status code: %d expected: %d", resp.StatusCode, http.StatusOK)
			}
		}
	}
}
