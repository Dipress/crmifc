package main

// func TestCategoryCreate(t *testing.T) {
// 	t.Log("with prepared server")
// 	{
// 		db, teardown := postgresDB(t)
// 		defer teardown()

// 		ctx, cancel := context.WithTimeout(context.Background(), caseTimeout)
// 		defer cancel()

// 		repo := postgres.NewRepository(db)

// 		nu := user.NewUser{
// 			RoleID:       5,
// 			Username:     "Dipress",
// 			Email:        "dipress@example.com",
// 			PasswordHash: "$2y$12$e4.VBLqKAanAZs10dRL65O8.b0kHBC34pcGCN1HdJIchCi9im40Ei",
// 		}

// 		var u user.User
// 		if err := repo.CreateUser(ctx, &nu, &u); err != nil {
// 			t.Errorf("unexpected error: %v", err)
// 		}

// 		claims := auth.NewClaims(u.Email, time.Now(), time.Hour)

// 		token, err := authenticator.GenerateToken(ctx, claims)
// 		if err != nil {
// 			t.Errorf("unexpected error: %v", err)
// 		}

// 		token = "Bearer " + token

// 		lis, err := net.Listen("tcp", "127.0.0.1:0")
// 		if err != nil {
// 			log.Fatalf("failed to listen: %v", err)
// 		}

// 		s := setupServer(lis.Addr().String(), db, authenticator)
// 		go s.Serve(lis)
// 		defer s.Close()

// 		t.Log("\ttest:0\tshould create a category.")
// 		{
// 			categoryStr := `{"name": "Contacts"}`
// 			req, err := http.NewRequest(
// 				http.MethodPost,
// 				fmt.Sprintf("http://%s/categories", s.Addr),
// 				strings.NewReader(categoryStr),
// 			)
// 			req.Header.Set("Content-Type", "application/json")
// 			req.Header.Add("Authorization", token)

// 			if err != nil {
// 				t.Fatalf("unexpected error: %v", err)
// 			}

// 			resp, err := http.DefaultClient.Do(req)
// 			if err != nil {
// 				t.Errorf("unexpected error: %v", err)
// 			}

// 			if resp.StatusCode != http.StatusOK {
// 				t.Errorf("unexpected status code: %d expected: %d", resp.StatusCode, http.StatusOK)
// 			}
// 		}
// 	}

// }

// func TestFindCategory(t *testing.T) {
// 	t.Log("with prepared server")
// 	{
// 		db, teardown := postgresDB(t)
// 		defer teardown()

// 		ctx, cancel := context.WithTimeout(context.Background(), caseTimeout)
// 		defer cancel()

// 		repo := postgres.NewRepository(db)

// 		nc := category.NewCategory{
// 			Name: "Airmax",
// 		}

// 		var c category.Category
// 		if err := repo.CreateCategory(ctx, &nc, &c); err != nil {
// 			t.Errorf("unexpected error: %v", err)
// 		}

// 		claims := auth.NewClaims("admin@example.com", time.Now(), time.Hour)

// 		token, err := authenticator.GenerateToken(ctx, claims)
// 		if err != nil {
// 			t.Errorf("unexpected error: %v", err)
// 		}

// 		token = "Bearer " + token

// 		lis, err := net.Listen("tcp", "127.0.0.1:0")
// 		if err != nil {
// 			log.Fatalf("failed to listen: %v", err)
// 		}

// 		s := setupServer(lis.Addr().String(), db, authenticator)
// 		go s.Serve(lis)
// 		defer s.Close()

// 		t.Log("\ttest:0\tshould find a category.")
// 		{
// 			req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://%s/categories/%d", s.Addr, c.ID), nil)
// 			req.Header.Set("Content-Type", "application/json")
// 			req.Header.Add("Authorization", token)

// 			if err != nil {
// 				t.Fatalf("unexpected error: %v", err)
// 			}
// 			resp, err := http.DefaultClient.Do(req)
// 			if err != nil {
// 				t.Errorf("unexpected error: %v", err)
// 			}

// 			if resp.StatusCode != http.StatusOK {
// 				t.Errorf("unexpected status code: %d expected: %d", resp.StatusCode, http.StatusOK)
// 			}
// 		}
// 	}
// }

// func TestUpdateCategory(t *testing.T) {
// 	t.Log("with prepared server")
// 	{
// 		db, teardown := postgresDB(t)
// 		defer teardown()

// 		ctx, cancel := context.WithTimeout(context.Background(), caseTimeout)
// 		defer cancel()

// 		repo := postgres.NewRepository(db)

// 		nc := category.NewCategory{
// 			Name: "Contacts",
// 		}

// 		var cat category.Category
// 		err := repo.CreateCategory(ctx, &nc, &cat)
// 		if err != nil {
// 			t.Errorf("unexpected error: %v", err)
// 		}

// 		claims := auth.NewClaims("admin@example.com", time.Now(), time.Hour)

// 		token, err := authenticator.GenerateToken(ctx, claims)
// 		if err != nil {
// 			t.Errorf("unexpected error: %v", err)
// 		}

// 		token = "Bearer " + token

// 		lis, err := net.Listen("tcp", "127.0.0.1:0")
// 		if err != nil {
// 			log.Fatalf("failed to listen: %v", err)
// 		}

// 		s := setupServer(lis.Addr().String(), db, authenticator)
// 		go s.Serve(lis)
// 		defer s.Close()

// 		t.Log("\ttest:0\tshould update a category.")
// 		{
// 			categoryStr := `{"name": "Partners"}`
// 			req, err := http.NewRequest(http.MethodPut,
// 				fmt.Sprintf("http://%s/categories/%d", s.Addr, cat.ID),
// 				strings.NewReader(categoryStr),
// 			)
// 			req.Header.Set("Content-Type", "application/json")
// 			req.Header.Add("Authorization", token)

// 			if err != nil {
// 				t.Fatalf("unexpected error: %v", err)
// 			}

// 			resp, err := http.DefaultClient.Do(req)
// 			if err != nil {
// 				t.Errorf("unexpected error: %v", err)
// 			}

// 			if resp.StatusCode != http.StatusOK {
// 				t.Errorf("unexpected status code: %d expected: %d", resp.StatusCode, http.StatusOK)
// 			}
// 		}
// 	}
// }

// func TestDeleteCategory(t *testing.T) {
// 	t.Log("with prepared server")
// 	{
// 		db, teardown := postgresDB(t)
// 		defer teardown()

// 		ctx, cancel := context.WithTimeout(context.Background(), caseTimeout)
// 		defer cancel()

// 		repo := postgres.NewRepository(db)

// 		nr := category.NewCategory{
// 			Name: "Flex",
// 		}

// 		var cat category.Category
// 		if err := repo.CreateCategory(ctx, &nr, &cat); err != nil {
// 			t.Errorf("unexpected error: %v", err)
// 		}

// 		claims := auth.NewClaims("admin@example.com", time.Now(), time.Hour)

// 		token, err := authenticator.GenerateToken(ctx, claims)
// 		if err != nil {
// 			t.Errorf("unexpected error: %v", err)
// 		}

// 		token = "Bearer " + token

// 		lis, err := net.Listen("tcp", "127.0.0.1:0")
// 		if err != nil {
// 			log.Fatalf("failed to listen: %v", err)
// 		}

// 		s := setupServer(lis.Addr().String(), db, authenticator)
// 		go s.Serve(lis)
// 		defer s.Close()

// 		t.Log("\ttest:0\tshould delete a category.")
// 		{
// 			req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("http://%s/categories/%d", s.Addr, cat.ID), nil)
// 			req.Header.Set("Content-Type", "application/json")
// 			req.Header.Add("Authorization", token)

// 			if err != nil {
// 				t.Fatalf("unexpected error: %v", err)
// 			}

// 			resp, err := http.DefaultClient.Do(req)
// 			if err != nil {
// 				t.Errorf("unexpected error: %v", err)
// 			}

// 			if resp.StatusCode != http.StatusOK {
// 				t.Errorf("unexpected status code: %d expected: %d", resp.StatusCode, http.StatusOK)
// 			}
// 		}
// 	}
// }

// func TestListCategory(t *testing.T) {
// 	t.Log("with prepared server")
// 	{
// 		db, teardown := postgresDB(t)
// 		defer teardown()

// 		ctx, cancel := context.WithTimeout(context.Background(), caseTimeout)
// 		defer cancel()

// 		repo := postgres.NewRepository(db)

// 		nc := category.NewCategory{
// 			Name: "Production",
// 		}

// 		var cat category.Category
// 		if err := repo.CreateCategory(ctx, &nc, &cat); err != nil {
// 			t.Errorf("unexpected error: %v", err)
// 		}

// 		claims := auth.NewClaims("admin@example.com", time.Now(), time.Hour)

// 		token, err := authenticator.GenerateToken(ctx, claims)
// 		if err != nil {
// 			t.Errorf("unexpected error: %v", err)
// 		}

// 		token = "Bearer " + token

// 		lis, err := net.Listen("tcp", "127.0.0.1:0")
// 		if err != nil {
// 			log.Fatalf("failed to listen: %v", err)
// 		}
// 		s := setupServer(lis.Addr().String(), db, authenticator)
// 		go s.Serve(lis)
// 		defer s.Close()

// 		t.Log("\ttest:0\tshould show all categories.")
// 		{
// 			req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://%s/categories", s.Addr), nil)
// 			req.Header.Set("Content-Type", "application/json")
// 			req.Header.Add("Authorization", token)

// 			if err != nil {
// 				t.Fatalf("unexpected error: %v", err)
// 			}
// 			resp, err := http.DefaultClient.Do(req)
// 			if err != nil {
// 				t.Errorf("unexpected error: %v", err)
// 			}

// 			if resp.StatusCode != http.StatusOK {
// 				t.Errorf("unexpected status code: %d expected: %d", resp.StatusCode, http.StatusOK)
// 			}
// 		}
// 	}
// }
