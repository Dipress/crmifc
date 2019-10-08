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

		repo := postgres.NewRepository(db)

		nr := role.NewRole{
			Name: "Admin",
		}
		var rl role.Role
		err := repo.CreateRole(ctx, &nr, &rl)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		nu := user.NewUser{
			Username:     "username",
			Email:        "username@example.com",
			PasswordHash: "$2y$12$gwoUXq7kCxNcucd.eFxOp.vJYYmo6917fSGuuEowfyNf3E8KySrWC",
			RoleID:       rl.ID,
		}

		var user user.User
		err = repo.CreateUser(ctx, &nu, &user)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

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

		t.Log("\ttest:0\tshould find a user.")
		{
			req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://%s/users/%d", s.Addr, user.ID), nil)
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

		repo := postgres.NewRepository(db)

		nr := role.NewRole{
			Name: "Admin",
		}
		var rl role.Role
		err := repo.CreateRole(ctx, &nr, &rl)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		nu := user.NewUser{
			Username:     "username",
			Email:        "username@example.com",
			PasswordHash: "$2y$12$gwoUXq7kCxNcucd.eFxOp.vJYYmo6917fSGuuEowfyNf3E8KySrWC",
			RoleID:       rl.ID,
		}

		var user user.User
		err = repo.CreateUser(ctx, &nu, &user)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

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

		t.Log("\ttest:0\tshould update a user.")
		{
			userStr := `{
				"username": "dmitry2",
				"email": "dmitry2@example.com",
				"password": "password1234",
				"role_id": 1
			}`
			req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("http://%s/users/%d", s.Addr, user.ID), strings.NewReader(userStr))
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

		repo := postgres.NewRepository(db)

		nr := role.NewRole{
			Name: "Admin",
		}
		var rl role.Role
		err := repo.CreateRole(ctx, &nr, &rl)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		nu := user.NewUser{
			Username:     "username11",
			Email:        "username11@example.com",
			PasswordHash: "$2y$12$gwoUXq7kCxNcucd.eFxOp.vJYYmo6917fSGuuEowfyNf3E8KySrWC",
			RoleID:       rl.ID,
		}

		var user user.User
		err = repo.CreateUser(ctx, &nu, &user)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

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

		t.Log("\ttest:0\tshould delete a user.")
		{
			req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("http://%s/users/%d", s.Addr, user.ID), nil)
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

		repo := postgres.NewRepository(db)

		nr := role.NewRole{
			Name: "Member",
		}

		var rl role.Role
		if err := repo.CreateRole(ctx, &nr, &rl); err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		nu1 := user.NewUser{
			Username:     "username21",
			Email:        "username21@example.com",
			PasswordHash: "$2y$12$gwoUXq7kCxNcucd.eFxOp.vJYYmo6917fSGuuEowfyNf3E8KySrWC",
			RoleID:       rl.ID,
		}

		nu2 := user.NewUser{
			Username:     "username22",
			Email:        "username22@example.com",
			PasswordHash: "$2y$12$gwoUXq7kCxNcucd.eFxOp.vJYYmo6917fSGuuEowfyNf3E8KySrWC",
			RoleID:       rl.ID,
		}

		var user1 user.User
		if err := repo.CreateUser(ctx, &nu1, &user1); err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		var user2 user.User
		if err := repo.CreateUser(ctx, &nu2, &user2); err != nil {
			t.Errorf("unexpected error: %v", err)
		}

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
