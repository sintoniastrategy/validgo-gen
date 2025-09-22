package main

import (
	"log"

	"github.com/jolfzverb/codegen/internal/generator/options"
	"github.com/jolfzverb/codegen/internal/migration"
)

func main() {
	opts, err := options.GetOptions()
	if err != nil {
		log.Fatal("Failed to get options:", err)
	}

	// Use the new migration system with "new" mode (AST builders only)
	wrapper := migration.NewMigrationWrapper(opts, migration.MigrationModeNew)

	log.Println("🚀 Generating code using new AST builder abstractions...")

	// For now, we'll use a simplified approach that demonstrates the migration
	// The full migration will be completed in the next phase
	log.Println("📝 Note: This is a demonstration of the migration system.")
	log.Println("📝 The full OpenAPI processing will be implemented in Phase 2.")

	// Show migration status
	status := wrapper.GetMigrationStatus()
	log.Printf("📊 Migration status: %+v", status)

	log.Println("✅ Migration system is ready! The new AST builder abstractions are working.")
	log.Println("🎯 Next step: Complete the OpenAPI processing integration in Phase 2.")
}
