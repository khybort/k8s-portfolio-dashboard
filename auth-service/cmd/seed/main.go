package main

import (
	"context"
	"log"
	"github.com/portfolio/auth-service/internal/config"
	"github.com/portfolio/auth-service/internal/model"
	"github.com/portfolio/auth-service/internal/repository"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, err := gorm.Open(postgres.Open(cfg.Database.DSN()), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	ctx := context.Background()
	userRepo := repository.NewUserRepository(db)

	// Seed admin user
	if err := seedAdminUser(ctx, userRepo, cfg); err != nil {
		log.Printf("Error seeding admin user: %v", err)
	}

	log.Println("Auth service seeding completed!")
}

func seedAdminUser(ctx context.Context, repo repository.UserRepository, cfg *config.Config) error {
	email := cfg.Seeder.AdminEmail
	password := cfg.Seeder.AdminPassword
	name := cfg.Seeder.AdminName
	
	// Check if user exists
	existing, err := repo.GetByEmail(ctx, email)
	if err == nil && existing != nil {
		log.Printf("Admin user already exists: %s", email)
		return nil
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	adminUser := &model.User{
		Email:        email,
		PasswordHash: string(hashedPassword),
		Name:         name,
		Role:         "admin",
		EmailVerified: true,
	}

	return repo.Create(ctx, adminUser)
}

