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
	"github.com/dipress/crmifc/internal/role"
	"github.com/dipress/crmifc/internal/storage/postgres"
	"github.com/dipress/crmifc/internal/user"
)

func TestUserCreate(t *testing.T) {
	t.Log("with prepared server")
	{
		db, teardown := postgresDB(t)
		defer teardown()

		ctx, cancel := context.WithTimeout(context.Background(), caseTimeout)
		defer cancel()

		userRepo := postgres.NewUserRepository(db)
		roleRepo := postgres.NewRoleRepository(db)

		nr := role.NewRole{
			Name: "Admin",
		}

		var rl role.Role
		if err := roleRepo.Create(ctx, &nr, &rl); err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		nu := user.NewUser{
			RoleID:       rl.ID,
			Username:     "Dipress91",
			Email:        "dipress91@example.com",
			PasswordHash: "$2y$12$e4.VBLqKAanAZs10dRL65O8.b0kHBC34pcGCN1HdJIchCi9im40Ei",
		}

		var u user.User
		if err := userRepo.Create(ctx, &nu, &u); err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		claims := auth.NewClaims(u.Email, time.Now(), time.Hour)
		authenticator := authenticatorSetup(db)

		token, err := authenticator.GenerateToken(ctx, claims)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		token = "Bearer " + token

		lis, err := net.Listen("tcp", "127.0.0.1:0")
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}

		services := setupServices(db, authenticator)

		s := setupServer(lis.Addr().String(), services, authenticator)
		go s.Serve(lis)
		defer s.Close()

		t.Log("\ttest:0\tshould create a user.")
		{
			userStr := `{
				"username": "dmitry",
				"email": "dmitry@example.com",
				"password": "password123",
				"role_id": 1
			}`
			req, err := http.NewRequest(
				http.MethodPost,
				fmt.Sprintf("http://%s/users", s.Addr),
				strings.NewReader(userStr),
			)
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

func TestFindUser(t *testing.T) {
	t.Log("with prepared server")
	{
		db, teardown := postgresDB(t)
		defer teardown()

		ctx, cancel := context.WithTimeout(context.Background(), caseTimeout)
		defer cancel()

		userRepo := postgres.NewUserRepository(db)
		roleRepo := postgres.NewRoleRepository(db)

		nr := role.NewRole{
			Name: "Admin",
		}

		var rl role.Role
		if err := roleRepo.Create(ctx, &nr, &rl); err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		nu := user.NewUser{
			RoleID:       rl.ID,
			Username:     "Dipress92",
			Email:        "dipress92@example.com",
			PasswordHash: "$2y$12$e4.VBLqKAanAZs10dRL65O8.b0kHBC34pcGCN1HdJIchCi9im40Ei",
		}

		var u user.User
		if err := userRepo.Create(ctx, &nu, &u); err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		claims := auth.NewClaims(u.Email, time.Now(), time.Hour)
		authenticator := authenticatorSetup(db)

		token, err := authenticator.GenerateToken(ctx, claims)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		token = "Bearer " + token

		lis, err := net.Listen("tcp", "127.0.0.1:0")
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}

		services := setupServices(db, authenticator)

		s := setupServer(lis.Addr().String(), services, authenticator)
		go s.Serve(lis)
		defer s.Close()

		t.Log("\ttest:0\tshould find a user.")
		{
			req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://%s/users/%d", s.Addr, u.ID), nil)
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

func TestUpdateUser(t *testing.T) {
	t.Log("with prepared server")
	{
		db, teardown := postgresDB(t)
		defer teardown()

		ctx, cancel := context.WithTimeout(context.Background(), caseTimeout)
		defer cancel()

		userRepo := postgres.NewUserRepository(db)
		roleRepo := postgres.NewRoleRepository(db)

		nr := role.NewRole{
			Name: "Admin",
		}

		var rl role.Role
		if err := roleRepo.Create(ctx, &nr, &rl); err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		nu := user.NewUser{
			RoleID:       rl.ID,
			Username:     "Dipress93",
			Email:        "dipress93@example.com",
			PasswordHash: "$2y$12$e4.VBLqKAanAZs10dRL65O8.b0kHBC34pcGCN1HdJIchCi9im40Ei",
		}

		var u user.User
		if err := userRepo.Create(ctx, &nu, &u); err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		claims := auth.NewClaims(u.Email, time.Now(), time.Hour)
		authenticator := authenticatorSetup(db)

		token, err := authenticator.GenerateToken(ctx, claims)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		token = "Bearer " + token

		lis, err := net.Listen("tcp", "127.0.0.1:0")
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}

		services := setupServices(db, authenticator)

		s := setupServer(lis.Addr().String(), services, authenticator)
		go s.Serve(lis)
		defer s.Close()

		t.Log("\ttest:0\tshould update a user.")
		{
			userStr := `{
				"username": "roman",
				"email": "roman@example.com",
				"password": "password1234",
				"role_id": 1
			}`
			req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("http://%s/users/%d", s.Addr, u.ID), strings.NewReader(userStr))
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

func TestDeleteUser(t *testing.T) {
	t.Log("with prepared server")
	{
		db, teardown := postgresDB(t)
		defer teardown()

		ctx, cancel := context.WithTimeout(context.Background(), caseTimeout)
		defer cancel()

		userRepo := postgres.NewUserRepository(db)
		roleRepo := postgres.NewRoleRepository(db)

		nr := role.NewRole{
			Name: "Admin",
		}

		var rl role.Role
		if err := roleRepo.Create(ctx, &nr, &rl); err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		nu := user.NewUser{
			RoleID:       rl.ID,
			Username:     "Dipress94",
			Email:        "dipress94@example.com",
			PasswordHash: "$2y$12$e4.VBLqKAanAZs10dRL65O8.b0kHBC34pcGCN1HdJIchCi9im40Ei",
		}

		var u user.User
		if err := userRepo.Create(ctx, &nu, &u); err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		claims := auth.NewClaims(u.Email, time.Now(), time.Hour)
		authenticator := authenticatorSetup(db)

		token, err := authenticator.GenerateToken(ctx, claims)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		token = "Bearer " + token

		lis, err := net.Listen("tcp", "127.0.0.1:0")
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}

		services := setupServices(db, authenticator)

		s := setupServer(lis.Addr().String(), services, authenticator)
		go s.Serve(lis)
		defer s.Close()

		t.Log("\ttest:0\tshould delete a user.")
		{
			req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("http://%s/users/%d", s.Addr, u.ID), nil)
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

func TestListUser(t *testing.T) {
	t.Log("with prepared server")
	{
		db, teardown := postgresDB(t)
		defer teardown()

		ctx, cancel := context.WithTimeout(context.Background(), caseTimeout)
		defer cancel()

		userRepo := postgres.NewUserRepository(db)
		roleRepo := postgres.NewRoleRepository(db)

		nr := role.NewRole{
			Name: "Admin",
		}

		var rl role.Role
		if err := roleRepo.Create(ctx, &nr, &rl); err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		nu := user.NewUser{
			RoleID:       rl.ID,
			Username:     "Dipress93",
			Email:        "dipress93@example.com",
			PasswordHash: "$2y$12$e4.VBLqKAanAZs10dRL65O8.b0kHBC34pcGCN1HdJIchCi9im40Ei",
		}

		var u user.User
		if err := userRepo.Create(ctx, &nu, &u); err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		claims := auth.NewClaims(u.Email, time.Now(), time.Hour)
		authenticator := authenticatorSetup(db)

		token, err := authenticator.GenerateToken(ctx, claims)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		token = "Bearer " + token

		lis, err := net.Listen("tcp", "127.0.0.1:0")
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}

		services := setupServices(db, authenticator)

		s := setupServer(lis.Addr().String(), services, authenticator)
		go s.Serve(lis)
		defer s.Close()

		t.Log("\ttest:0\tshould show all users.")
		{
			req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://%s/users", s.Addr), nil)
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
