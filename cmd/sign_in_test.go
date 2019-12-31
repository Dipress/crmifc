package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
	"testing"

	"github.com/dipress/crmifc/internal/role"
	"github.com/dipress/crmifc/internal/storage/postgres"
	"github.com/dipress/crmifc/internal/user"
)

func TestSignIn(t *testing.T) {
	t.Log("with prepared server")
	{
		db, teardown := postgresDB(t)
		defer teardown()

		ctx, cancel := context.WithTimeout(context.Background(), caseTimeout)
		defer cancel()

		userRepo := postgres.NewUserRepository(db)
		roleRepo := postgres.NewRoleRepository(db)

		nr := role.NewRole{
			Name: "Manager",
		}

		var rl role.Role
		if err := roleRepo.Create(ctx, &nr, &rl); err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		nu := user.NewUser{
			RoleID:       rl.ID,
			Username:     "username6",
			Email:        "username6@example.com",
			PasswordHash: "$2y$12$e4.VBLqKAanAZs10dRL65O8.b0kHBC34pcGCN1HdJIchCi9im40Ei",
		}

		var u user.User
		if err := userRepo.Create(ctx, &nu, &u); err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		lis, err := net.Listen("tcp", "127.0.0.1:0")
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}

		authenticator := authenticatorSetup(db)

		services := setupServices(db, authenticator)

		s := setupServer(lis.Addr().String(), nil, services, authenticator)
		go s.Serve(lis)
		defer s.Close()

		t.Log("\ttest:0\tshould authenticate a user.")
		{
			authStr := `{"email": "username6@example.com", "password": "password123"}`
			req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("http://%s/signin", s.Addr), strings.NewReader(authStr))
			req.Header.Set("Content-Type", "application/json")

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
