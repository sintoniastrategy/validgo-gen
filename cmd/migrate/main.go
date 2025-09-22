package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/jolfzverb/codegen/internal/generator/options"
	"github.com/jolfzverb/codegen/internal/migration"
)

func main() {
	var (
		migrationMode = flag.String("mode", "hybrid", "Migration mode: legacy, hybrid, new")
		validate      = flag.Bool("validate", true, "Validate migration output")
		compare       = flag.Bool("compare", false, "Compare legacy and migrated outputs")
		status        = flag.Bool("status", false, "Show migration status")
		plan          = flag.Bool("plan", false, "Show migration plan")
	)
	flag.Parse()

	// Get options
	opts, err := options.GetOptions()
	if err != nil {
		log.Fatal("Failed to get options:", err)
	}

	// Parse migration mode
	var mode migration.MigrationMode
	switch *migrationMode {
	case "legacy":
		mode = migration.MigrationModeLegacy
	case "hybrid":
		mode = migration.MigrationModeHybrid
	case "new":
		mode = migration.MigrationModeNew
	default:
		log.Fatal("Invalid migration mode. Use: legacy, hybrid, or new")
	}

	// Show migration plan if requested
	if *plan {
		showMigrationPlan()
		return
	}

	// Create migration wrapper
	wrapper := migration.NewMigrationWrapper(opts, mode)

	// Show status if requested
	if *status {
		showMigrationStatus(wrapper)
		return
	}

	// Generate code
	ctx := context.Background()
	err = wrapper.Generate(ctx)
	if err != nil {
		log.Fatal("Generation failed:", err)
	}

	// Validate migration if requested
	if *validate {
		err = wrapper.ValidateMigration()
		if err != nil {
			log.Printf("Migration validation failed: %v", err)
		} else {
			fmt.Println("✅ Migration validation passed")
		}
	}

	// Compare outputs if requested
	if *compare && mode == migration.MigrationModeHybrid {
		comparison, err := wrapper.CompareOutputs(ctx)
		if err != nil {
			log.Printf("Output comparison failed: %v", err)
		} else {
			showComparison(comparison)
		}
	}

	fmt.Printf("✅ Code generation completed using %s mode\n", *migrationMode)
}

func showMigrationPlan() {
	plan := migration.CreateMigrationPlan()
	fmt.Println("📋 Migration Plan:")
	fmt.Println("==================")

	for i, step := range plan.Steps {
		fmt.Printf("\n%d. %s (%s)\n", i+1, step.Phase, step.Duration)
		for _, task := range step.Tasks {
			fmt.Printf("   %s\n", task)
		}
	}
}

func showMigrationStatus(wrapper *migration.MigrationWrapper) {
	status := wrapper.GetMigrationStatus()
	fmt.Println("📊 Migration Status:")
	fmt.Println("===================")

	for key, value := range status {
		fmt.Printf("%s: %v\n", key, value)
	}
}

func showComparison(comparison *migration.MigrationComparison) {
	fmt.Println("🔍 Output Comparison:")
	fmt.Println("====================")

	if comparison.Success {
		fmt.Println("✅ Both generation methods succeeded")
	} else {
		if comparison.LegacyError != "" {
			fmt.Printf("❌ Legacy generation failed: %s\n", comparison.LegacyError)
		}
		if comparison.NewError != "" {
			fmt.Printf("❌ Migrated generation failed: %s\n", comparison.NewError)
		}
	}

	if len(comparison.Differences) > 0 {
		fmt.Println("\n📝 Differences:")
		for _, diff := range comparison.Differences {
			fmt.Printf("   %s\n", diff)
		}
	}
}
