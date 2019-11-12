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

	"github.com/dipress/crmifc/internal/article"
	"github.com/dipress/crmifc/internal/kit/auth"
	"github.com/dipress/crmifc/internal/role"
	"github.com/dipress/crmifc/internal/storage/postgres"
	"github.com/dipress/crmifc/internal/user"
)

func TestArticleCreate(t *testing.T) {
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

		var rol role.Role
		err := roleRepo.Create(ctx, &nr, &rol)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		nu := user.NewUser{
			Username:     "username21",
			Email:        "username21@example.com",
			PasswordHash: "$2y$12$e4.VBLqKAanAZs10dRL65O8.b0kHBC34pcGCN1HdJIchCi9im40Ei",
			RoleID:       rol.ID,
		}

		var u user.User
		if err := userRepo.Create(ctx, &nu, &u); err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		claims := auth.NewClaims(u.Email, time.Now(), time.Hour)

		authenticator := authenticatorSetup(db)

		token, err := authenticator.GenerateToken(ctx, claims.StandardClaims)
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

func TestFindArticle(t *testing.T) {
	t.Log("with prepared server")
	{
		db, teardown := postgresDB(t)
		defer teardown()

		ctx, cancel := context.WithTimeout(context.Background(), caseTimeout)
		defer cancel()

		userRepo := postgres.NewUserRepository(db)
		roleRepo := postgres.NewRoleRepository(db)
		articleRepo := postgres.NewArticleRepository(db)

		nr := role.NewRole{
			Name: "Admin",
		}

		var rol role.Role
		err := roleRepo.Create(ctx, &nr, &rol)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		nu := user.NewUser{
			Username:     "username21",
			Email:        "username21@example.com",
			PasswordHash: "$2y$12$e4.VBLqKAanAZs10dRL65O8.b0kHBC34pcGCN1HdJIchCi9im40Ei",
			RoleID:       rol.ID,
		}

		var u user.User
		if err := userRepo.Create(ctx, &nu, &u); err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		claims := auth.NewClaims(u.Email, time.Now(), time.Hour)

		authenticator := authenticatorSetup(db)

		na := article.NewArticle{
			UserID:     u.ID,
			CategoryID: 10,
			Title:      "my title",
			Body:       "my body",
		}

		var art article.Article
		err = articleRepo.Create(ctx, &na, &art)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

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

		t.Log("\ttest:0\tshould find a article.")
		{
			req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://%s/articles/%d", s.Addr, art.ID), nil)
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

func TestUpdateArticle(t *testing.T) {
	t.Log("with prepared server")
	{
		db, teardown := postgresDB(t)
		defer teardown()

		ctx, cancel := context.WithTimeout(context.Background(), caseTimeout)
		defer cancel()

		userRepo := postgres.NewUserRepository(db)
		roleRepo := postgres.NewRoleRepository(db)
		articleRepo := postgres.NewArticleRepository(db)

		nr := role.NewRole{
			Name: "Admin",
		}

		var rol role.Role
		err := roleRepo.Create(ctx, &nr, &rol)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		nu := user.NewUser{
			Username:     "username21",
			Email:        "username21@example.com",
			PasswordHash: "$2y$12$e4.VBLqKAanAZs10dRL65O8.b0kHBC34pcGCN1HdJIchCi9im40Ei",
			RoleID:       rol.ID,
		}

		var u user.User
		if err := userRepo.Create(ctx, &nu, &u); err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		claims := auth.NewClaims(u.Email, time.Now(), time.Hour)

		authenticator := authenticatorSetup(db)

		na := article.NewArticle{
			UserID:     u.ID,
			CategoryID: 10,
			Title:      "my title",
			Body:       "my body",
		}

		var art article.Article
		err = articleRepo.Create(ctx, &na, &art)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

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

		t.Log("\ttest:0\tshould update a post.")
		{
			articleStr := `{"category_id":22, "title":"my awesome title", "body":"my awesome body"}`
			req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("http://%s/articles/%d", s.Addr, art.ID), strings.NewReader(articleStr))
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

func TestDeleteArticle(t *testing.T) {
	t.Log("with prepared server")
	{
		db, teardown := postgresDB(t)
		defer teardown()

		ctx, cancel := context.WithTimeout(context.Background(), caseTimeout)
		defer cancel()

		userRepo := postgres.NewUserRepository(db)
		roleRepo := postgres.NewRoleRepository(db)
		articleRepo := postgres.NewArticleRepository(db)

		nr := role.NewRole{
			Name: "Admin",
		}

		var rol role.Role
		err := roleRepo.Create(ctx, &nr, &rol)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		nu := user.NewUser{
			Username:     "username21",
			Email:        "username21@example.com",
			PasswordHash: "$2y$12$e4.VBLqKAanAZs10dRL65O8.b0kHBC34pcGCN1HdJIchCi9im40Ei",
			RoleID:       rol.ID,
		}

		var u user.User
		if err := userRepo.Create(ctx, &nu, &u); err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		claims := auth.NewClaims(u.Email, time.Now(), time.Hour)

		authenticator := authenticatorSetup(db)

		na := article.NewArticle{
			UserID:     u.ID,
			CategoryID: 10,
			Title:      "my title",
			Body:       "my body",
		}

		var art article.Article
		err = articleRepo.Create(ctx, &na, &art)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

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

		t.Log("\ttest:0\tshould delete a article.")
		{
			req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("http://%s/articles/%d", s.Addr, art.ID), nil)
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
