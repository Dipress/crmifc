package postgres

import (
	"context"
	"testing"

	"github.com/dipress/crmifc/internal/role"

	"github.com/dipress/crmifc/internal/user"
	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
	t.Log("with initialize repository")
	{
		db, teardown := postgresDB(t)
		defer teardown()

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		roleRepo := NewRoleRepository(db)
		userRepo := NewUserRepository(db)

		nr := role.NewRole{
			Name: "Admin",
		}
		var rl role.Role
		err := roleRepo.Create(ctx, &nr, &rl)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		t.Log("\ttest:0\tshould insert a new user into the database")
		{
			nu := user.NewUser{
				Username:     "username",
				Email:        "username@example.com",
				PasswordHash: "$2y$12$gwoUXq7kCxNcucd.eFxOp.vJYYmo6917fSGuuEowfyNf3E8KySrWC",
				RoleID:       rl.ID,
			}

			var user user.User
			err := userRepo.Create(ctx, &nu, &user)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if user.ID == 0 {
				t.Error("expected to parse returned id")
			}
		}
	}
}

func TestUserFind(t *testing.T) {
	t.Log("with initialized repository")
	{
		db, teardown := postgresDB(t)
		defer teardown()

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		roleRepo := NewRoleRepository(db)
		userRepo := NewUserRepository(db)

		nr := role.NewRole{
			Name: "Admin",
		}
		var rl role.Role
		err := roleRepo.Create(ctx, &nr, &rl)
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
		err = userRepo.Create(ctx, &nu, &user)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		t.Log("\ttest:0\tshould find the user into the database")
		{
			_, err := userRepo.Find(ctx, user.ID)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		}
	}
}

func TestUniqueUsername(t *testing.T) {
	t.Log("with initialized repository")
	{
		db, teardown := postgresDB(t)
		defer teardown()

		userRepo := NewUserRepository(db)

		nu := user.NewUser{
			Username:     "username1",
			Email:        "username1@example.com",
			PasswordHash: "$2y$12$gwoUXq7kCxNcucd.eFxOp.vJYYmo6917fSGuuEowfyNf3E8KySrWC",
			RoleID:       1,
		}

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		var user user.User
		err := userRepo.Create(ctx, &nu, &user)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		t.Log("\ttest:0\tshould return error")
		{
			err := userRepo.UniqueUsername(ctx, "username1")
			assert.Error(t, err, "username already exists")
		}
		t.Log("\ttest:0\tshould return nil")
		{
			err := userRepo.UniqueUsername(ctx, "username2")
			assert.Nil(t, err)
		}
	}
}

func TestUpdateUser(t *testing.T) {
	t.Log("with initialized repository")
	{
		db, teardown := postgresDB(t)
		defer teardown()

		roleRepo := NewRoleRepository(db)
		userRepo := NewUserRepository(db)

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		nr := role.NewRole{
			Name: "Manager",
		}

		var rl role.Role
		err := roleRepo.Create(ctx, &nr, &rl)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		nu := user.NewUser{
			Username:     "username10",
			Email:        "username10@example.com",
			PasswordHash: "$2y$12$gwoUXq7kCxNcucd.eFxOp.vJYYmo6917fSGuuEowfyNf3E8KySrWC",
			RoleID:       rl.ID,
		}

		var user user.User
		err = userRepo.Create(ctx, &nu, &user)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		t.Log("\ttest:0\tshould update the role into the database")
		{

			user.Username = "Hacket"
			user.Email = "hacket@example.com"
			user.Role.ID = 3

			err := userRepo.Update(ctx, 1, &user)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		}
	}
}

func TestUniqueEmail(t *testing.T) {
	t.Log("with initialized repository")
	{
		db, teardown := postgresDB(t)
		defer teardown()
		userRepo := NewUserRepository(db)

		nu := user.NewUser{
			Username:     "username5",
			Email:        "username5@example.com",
			PasswordHash: "$2y$12$gwoUXq7kCxNcucd.eFxOp.vJYYmo6917fSGuuEowfyNf3E8KySrWC",
			RoleID:       1,
		}

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		var user user.User
		err := userRepo.Create(ctx, &nu, &user)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		t.Log("\ttest:0\tshould return error")
		{
			err := userRepo.UniqueEmail(ctx, "username5@example.com")
			assert.Error(t, err, "email already exists")
		}
		t.Log("\ttest:0\tshould return nil")
		{
			err := userRepo.UniqueEmail(ctx, "username6@example.com")
			assert.Nil(t, err)
		}
	}
}

func TestListUsers(t *testing.T) {
	t.Log("with initialized repository")
	{
		db, teardown := postgresDB(t)
		defer teardown()

		roleRepo := NewRoleRepository(db)
		userRepo := NewUserRepository(db)

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		nr := role.NewRole{
			Name: "Admin",
		}
		var rl role.Role
		err := roleRepo.Create(ctx, &nr, &rl)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		nu1 := user.NewUser{
			Username:     "username1",
			Email:        "username1@example.com",
			PasswordHash: "$2y$12$gwoUXq7kCxNcucd.eFxOp.vJYYmo6917fSGuuEowfyNf3E8KySrWC",
			RoleID:       rl.ID,
		}

		nu2 := user.NewUser{
			Username:     "username2",
			Email:        "username2@example.com",
			PasswordHash: "$2y$12$gwoUXq7kCxNcucd.eFxOp.vJYYmo6917fSGuuEowfyNf3E8KySrWC",
			RoleID:       rl.ID,
		}

		var usr1 user.User
		err = userRepo.Create(ctx, &nu1, &usr1)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		var usr2 user.User
		err = userRepo.Create(ctx, &nu2, &usr2)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		t.Log("\ttest:0\tshould show list of users")
		{
			var users user.Users
			err := userRepo.List(ctx, &users)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if len(users.Users) != 2 {
				t.Error("expected to slice of two users")
			}
		}
	}
}

func TestFindUserByEmail(t *testing.T) {
	t.Log("with initialized repository")
	{
		db, teardown := postgresDB(t)
		defer teardown()

		userRepo := NewUserRepository(db)
		roleRepo := NewRoleRepository(db)

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		nr := role.NewRole{
			Name: "Member",
		}
		var r role.Role
		err :=roleRepo.Create(ctx, &nr, &r)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		nu := user.NewUser{
			Username:     "username4",
			Email:        "username4@example.com",
			PasswordHash: "$2y$12$gwoUXq7kCxNcucd.eFxOp.vJYYmo6917fSGuuEowfyNf3E8KySrWC",
			RoleID: r.ID,
		}

		var usr user.User
		err = userRepo.Create(ctx, &nu, &usr)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		t.Log("\ttest:0\tshould find user by email")
		{
			_, err := userRepo.FindByEmail(ctx, usr.Email)
			assert.Nil(t, err)
		}
	}
}
