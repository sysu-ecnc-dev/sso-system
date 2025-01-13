package utils

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sysu-ecnc-dev/sso-system/backend/internal/config"
	"github.com/sysu-ecnc-dev/sso-system/backend/internal/global"
	"github.com/sysu-ecnc-dev/sso-system/backend/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

// check the roles table is correctly initialized
func EnsureRolesTableInitialized(cfg *config.Config, q *repository.Queries, dbpool *pgxpool.Pool) error {
	// get all roles
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.Database.QueryTimeout)*time.Second)
	defer cancel()

	roles, err := q.GetAllRoles(ctx)
	if err != nil {
		return err
	}

	// check the roles table is correctly initialized
	isCorrect := true
	for key, value := range global.RoleLevelMap {
		// key is the name of the rol and value is the level of the role
		isFound := false

		for _, role := range roles {
			if role.Name == key && role.Level == value {
				isFound = true
				break
			}
		}

		if !isFound {
			isCorrect = false
			break
		}
	}

	// if the roles table is not correctly initialized, initialize it
	if !isCorrect {
		// begin transaction
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.Database.TransactionTimeout)*time.Second)
		defer cancel()

		tx, err := dbpool.Begin(ctx)
		if err != nil {
			return err
		}
		defer func() {
			_ = tx.Rollback(ctx)
		}()
		qtx := q.WithTx(tx)

		// delete all roles
		if err := qtx.DeleteAllRoles(ctx); err != nil {
			return err
		}

		// insert correct roles
		for key, value := range global.RoleLevelMap {
			_, err := qtx.CreateRole(ctx, repository.CreateRoleParams{
				Name:  key,
				Level: value,
			})
			if err != nil {
				return err
			}
		}

		// commit
		if err := tx.Commit(ctx); err != nil {
			return err
		}
	}

	return nil
}

// check if the initial administrator exists, if not, create one
func EnsureInitialAdminExists(cfg *config.Config, q *repository.Queries, dbpool *pgxpool.Pool) error {
	valid := true

	// check if the initial administrator exists
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.Database.QueryTimeout)*time.Second)
	defer cancel()

	initialAdmin, err := q.GetUserByUsername(ctx, cfg.InitialAdmin.Username)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			valid = false
		default:
			return err
		}
	}

	// check the role
	if valid {
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.Database.QueryTimeout)*time.Second)
		defer cancel()

		role, err := q.GetRoleById(ctx, initialAdmin.RoleID)
		if err != nil {
			return err
		}

		if role.Name != "black_core" {
			valid = false
		}
	}

	if valid {
		return nil
	}

	// delete the invalid initial administrator and create a new one
	ctx, cancel = context.WithTimeout(context.Background(), time.Duration(cfg.Database.TransactionTimeout)*time.Second)
	defer cancel()

	tx, err := dbpool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	qtx := q.WithTx(tx)
	if err := qtx.DeleteUserByUsername(ctx, cfg.InitialAdmin.Username); err != nil {
		return err
	}

	// hash the password
	password_hash, err := bcrypt.GenerateFromPassword([]byte(cfg.InitialAdmin.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// get the correct role id
	role, err := q.GetRoleByName(ctx, "black_core")
	if err != nil {
		return err
	}

	// insert the initial administrator
	_, err = qtx.CreateUser(ctx, repository.CreateUserParams{
		Username:     cfg.InitialAdmin.Username,
		PasswordHash: string(password_hash),
		FullName:     cfg.InitialAdmin.FullName,
		Email:        cfg.InitialAdmin.Email,
		RoleID:       role.ID,
	})
	if err != nil {
		return err
	}

	// commit
	if err := tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}
