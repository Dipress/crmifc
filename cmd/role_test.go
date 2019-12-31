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

func TestRoleCreate(t *testing.T) {
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
			Username:     "Dipress",
			Email:        "dipress@example.com",
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

		s := setupServer(lis.Addr().String(), nil, services, authenticator)
		go s.Serve(lis)
		defer s.Close()

		t.Log("\ttest:0\tshould create a role.")
		{
			roleStr := `{"name": "Member"}`
			req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("http://%s/roles", s.Addr), strings.NewReader(roleStr))
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

func TestFindRole(t *testing.T) {
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
			Username:     "Dipress55",
			Email:        "dipress55@example.com",
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

		s := setupServer(lis.Addr().String(), nil, services, authenticator)
		go s.Serve(lis)
		defer s.Close()

		t.Log("\ttest:0\tshould find a role.")
		{
			req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://%s/roles/%d", s.Addr, rl.ID), nil)
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

func TestUpdateRole(t *testing.T) {
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
			Username:     "Dipress59",
			Email:        "dipress59@example.com",
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

		s := setupServer(lis.Addr().String(), nil, services, authenticator)
		go s.Serve(lis)
		defer s.Close()

		t.Log("\ttest:0\tshould update a role.")
		{
			roleStr := `{"name": "Manager"}`
			req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("http://%s/roles/%d", s.Addr, rl.ID), strings.NewReader(roleStr))
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

func TestDeleteRole(t *testing.T) {
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
			Username:     "Dipress60",
			Email:        "dipress60@example.com",
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

		s := setupServer(lis.Addr().String(), nil, services, authenticator)
		go s.Serve(lis)
		defer s.Close()

		t.Log("\ttest:0\tshould delete a role.")
		{
			req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("http://%s/roles/%d", s.Addr, rl.ID), nil)
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

func TestListRole(t *testing.T) {
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
			Username:     "Dipress59",
			Email:        "dipress59@example.com",
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

		s := setupServer(lis.Addr().String(), nil, services, authenticator)
		go s.Serve(lis)
		defer s.Close()

		t.Log("\ttest:0\tshould show all roles.")
		{
			req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://%s/roles", s.Addr), nil)
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
