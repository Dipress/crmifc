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
)

func TestRoleCreate(t *testing.T) {
	t.Log("with prepared server")
	{
		db, teardown := postgresDB(t)
		defer teardown()

		lis, err := net.Listen("tcp", "127.0.0.1:0")
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}

		s := setupServer(lis.Addr().String(), db, authenticator)
		go s.Serve(lis)
		defer s.Close()

		t.Log("\ttest:0\tshould create a role.")
		{
			roleStr := `{"name": "admin"}`
			req, err := http.NewRequest("POST", fmt.Sprintf("http://%s/roles", s.Addr), strings.NewReader(roleStr))
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

func TestFindRole(t *testing.T) {
	t.Log("with prepared server")
	{
		db, teardown := postgresDB(t)
		defer teardown()

		ctx, cancel := context.WithTimeout(context.Background(), caseTimeout)
		defer cancel()

		repo := postgres.NewRepository(db)

		nr := role.NewRole{
			Name: "Ingeneer",
		}

		var r role.Role
		if err := repo.CreateRole(ctx, &nr, &r); err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		lis, err := net.Listen("tcp", "127.0.0.1:0")
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}

		s := setupServer(lis.Addr().String(), db, authenticator)
		go s.Serve(lis)
		defer s.Close()

		t.Log("\ttest:0\tshould find a role.")
		{
			req, err := http.NewRequest("GET", fmt.Sprintf("http://%s/roles/%d", s.Addr, r.ID), nil)
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

		repo := postgres.NewRepository(db)

		nr := role.NewRole{
			Name: "Admin",
		}
		var r role.Role
		err := repo.CreateRole(ctx, &nr, &r)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		lis, err := net.Listen("tcp", "127.0.0.1:0")
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}

		s := setupServer(lis.Addr().String(), db, authenticator)
		go s.Serve(lis)
		defer s.Close()

		t.Log("\ttest:0\tshould update a role.")
		{
			roleStr := `{"name": "manager"}`
			req, err := http.NewRequest("PUT", fmt.Sprintf("http://%s/roles/%d", s.Addr, r.ID), strings.NewReader(roleStr))
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

func TestDeleteRole(t *testing.T) {
	t.Log("with prepared server")
	{
		db, teardown := postgresDB(t)
		defer teardown()

		ctx, cancel := context.WithTimeout(context.Background(), caseTimeout)
		defer cancel()

		repo := postgres.NewRepository(db)

		nr := role.NewRole{
			Name: "Idol",
		}
		var rl role.Role
		if err := repo.CreateRole(ctx, &nr, &rl); err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		lis, err := net.Listen("tcp", "127.0.0.1:0")
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}

		s := setupServer(lis.Addr().String(), db, authenticator)
		go s.Serve(lis)
		defer s.Close()

		t.Log("\ttest:0\tshould delete a role.")
		{
			req, err := http.NewRequest("DELETE", fmt.Sprintf("http://%s/roles/%d", s.Addr, rl.ID), nil)
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

func TestListRole(t *testing.T) {
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

		lis, err := net.Listen("tcp", "127.0.0.1:0")
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		s := setupServer(lis.Addr().String(), db, authenticator)
		go s.Serve(lis)
		defer s.Close()

		t.Log("\ttest:0\tshould show all roles.")
		{
			req, err := http.NewRequest("GET", fmt.Sprintf("http://%s/roles", s.Addr), nil)
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
