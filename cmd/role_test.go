package main

// func TestRoleCreate(t *testing.T) {
// 	t.Log("with prepared server")
// 	{
// 		db, teardown := postgresDB(t)
// 		defer teardown()

// 		ctx, cancel := context.WithTimeout(context.Background(), caseTimeout)
// 		defer cancel()

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

// 		t.Log("\ttest:0\tshould create a role.")
// 		{
// 			roleStr := `{"name": "admin"}`
// 			req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("http://%s/roles", s.Addr), strings.NewReader(roleStr))
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

// func TestFindRole(t *testing.T) {
// 	t.Log("with prepared server")
// 	{
// 		db, teardown := postgresDB(t)
// 		defer teardown()

// 		ctx, cancel := context.WithTimeout(context.Background(), caseTimeout)
// 		defer cancel()

// 		repo := postgres.NewRepository(db)

// 		nr := role.NewRole{
// 			Name: "Ingeneer",
// 		}

// 		var r role.Role
// 		if err := repo.CreateRole(ctx, &nr, &r); err != nil {
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

// 		t.Log("\ttest:0\tshould find a role.")
// 		{
// 			req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://%s/roles/%d", s.Addr, r.ID), nil)
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

// func TestUpdateRole(t *testing.T) {
// 	t.Log("with prepared server")
// 	{
// 		db, teardown := postgresDB(t)
// 		defer teardown()

// 		ctx, cancel := context.WithTimeout(context.Background(), caseTimeout)
// 		defer cancel()

// 		repo := postgres.NewRepository(db)

// 		nr := role.NewRole{
// 			Name: "Admin",
// 		}

// 		var r role.Role
// 		err := repo.CreateRole(ctx, &nr, &r)
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

// 		t.Log("\ttest:0\tshould update a role.")
// 		{
// 			roleStr := `{"name": "manager"}`
// 			req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("http://%s/roles/%d", s.Addr, r.ID), strings.NewReader(roleStr))
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

// func TestDeleteRole(t *testing.T) {
// 	t.Log("with prepared server")
// 	{
// 		db, teardown := postgresDB(t)
// 		defer teardown()

// 		ctx, cancel := context.WithTimeout(context.Background(), caseTimeout)
// 		defer cancel()

// 		repo := postgres.NewRepository(db)

// 		nr := role.NewRole{
// 			Name: "Idol",
// 		}

// 		var rl role.Role
// 		if err := repo.CreateRole(ctx, &nr, &rl); err != nil {
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

// 		t.Log("\ttest:0\tshould delete a role.")
// 		{
// 			req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("http://%s/roles/%d", s.Addr, rl.ID), nil)
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

// func TestListRole(t *testing.T) {
// 	t.Log("with prepared server")
// 	{
// 		db, teardown := postgresDB(t)
// 		defer teardown()

// 		ctx, cancel := context.WithTimeout(context.Background(), caseTimeout)
// 		defer cancel()

// 		repo := postgres.NewRepository(db)

// 		nr := role.NewRole{
// 			Name: "Member",
// 		}

// 		var rl role.Role
// 		if err := repo.CreateRole(ctx, &nr, &rl); err != nil {
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

// 		t.Log("\ttest:0\tshould show all roles.")
// 		{
// 			req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://%s/roles", s.Addr), nil)
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
