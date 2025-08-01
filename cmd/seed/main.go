package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kirigaikabuto/saas-microcourses-management/internal/db"
)

type SeedCompany struct {
	Name             string
	SubscriptionPlan string
}

var seedCompanies = []SeedCompany{
	{"Acme Corporation", "enterprise"},
	{"TechStart Inc", "premium"},
	{"Digital Solutions Ltd", "basic"},
	{"Innovation Labs", "premium"},
	{"Global Enterprises", "enterprise"},
	{"StartupCo", "basic"},
	{"Enterprise Solutions", "enterprise"},
	{"Cloud Services Ltd", "premium"},
	{"Data Analytics Inc", "basic"},
	{"AI Innovations", "premium"},
	{"Future Tech", "basic"},
	{"Quantum Computing Co", "enterprise"},
	{"Mobile First Ltd", "premium"},
	{"Blockchain Ventures", "enterprise"},
	{"Green Energy Solutions", "basic"},
}

func main() {
	// Get database URL from environment
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL environment variable is required")
	}

	// Connect to database
	pool, err := pgxpool.New(context.Background(), databaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer pool.Close()

	// Test connection
	if err := pool.Ping(context.Background()); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	// Create queries instance
	queries := db.New(pool)
	ctx := context.Background()

	// Check command line arguments
	if len(os.Args) > 1 && os.Args[1] == "--clear" {
		fmt.Println("Clearing existing seed data...")
		if err := clearSeedData(ctx, queries); err != nil {
			log.Fatalf("Failed to clear seed data: %v", err)
		}
		fmt.Println("Seed data cleared successfully")
		return
	}

	// Seed companies
	fmt.Println("Seeding companies...")
	created := 0
	skipped := 0

	for _, seedCompany := range seedCompanies {
		// Check if company already exists
		companies, err := queries.ListCompanies(ctx, db.ListCompaniesParams{
			Limit:  1000, // Get all companies to check names
			Offset: 0,
		})
		if err != nil {
			log.Fatalf("Failed to list companies: %v", err)
		}

		// Check if company name already exists
		exists := false
		for _, company := range companies {
			if company.Name.String == seedCompany.Name {
				exists = true
				break
			}
		}

		if exists {
			fmt.Printf("  ✓ Company '%s' already exists, skipping\n", seedCompany.Name)
			skipped++
			continue
		}

		// Create company
		company, err := queries.CreateCompany(ctx, db.CreateCompanyParams{
			Name:             pgtype.Text{String: seedCompany.Name, Valid: true},
			SubscriptionPlan: pgtype.Text{String: seedCompany.SubscriptionPlan, Valid: true},
		})
		if err != nil {
			log.Fatalf("Failed to create company '%s': %v", seedCompany.Name, err)
		}

		fmt.Printf("  ✓ Created company '%s' with plan '%s' (ID: %s)\n", 
			company.Name.String, 
			company.SubscriptionPlan.String,
			uuid.UUID(company.ID.Bytes).String())
		created++
	}

	fmt.Printf("\nSeeding completed: %d created, %d skipped\n", created, skipped)
	
	// Show total count
	total, err := queries.CountCompanies(ctx)
	if err != nil {
		log.Printf("Warning: Failed to count companies: %v", err)
	} else {
		fmt.Printf("Total companies in database: %d\n", total)
	}
}

func clearSeedData(ctx context.Context, queries *db.Queries) error {
	// Get all companies
	companies, err := queries.ListCompanies(ctx, db.ListCompaniesParams{
		Limit:  1000,
		Offset: 0,
	})
	if err != nil {
		return fmt.Errorf("failed to list companies: %w", err)
	}

	// Create a map of seed company names for fast lookup
	seedNames := make(map[string]bool)
	for _, seedCompany := range seedCompanies {
		seedNames[seedCompany.Name] = true
	}

	// Delete seed companies
	deleted := 0
	for _, company := range companies {
		if seedNames[company.Name.String] {
			if err := queries.DeleteCompany(ctx, company.ID); err != nil {
				return fmt.Errorf("failed to delete company '%s': %w", company.Name.String, err)
			}
			fmt.Printf("  ✓ Deleted company '%s'\n", company.Name.String)
			deleted++
		}
	}

	fmt.Printf("Deleted %d seed companies\n", deleted)
	return nil
}