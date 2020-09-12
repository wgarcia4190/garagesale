package main

import (
	"flag"
	"log"

	"github.com/wgarcia4190/garagesale/internal/platform/database"
	"github.com/wgarcia4190/garagesale/internal/schema"
)

func main() {
	// =========================================================================
	// Setup Dependencies

	db, err := database.Open()

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	flag.Parse()
	switch flag.Arg(0) {
	case "migrate":
		if err := schema.Migrate(db); err != nil {
			log.Fatal("applying migrations", err)
		}
		log.Println("Migrations complete")
		return

	case "seed":
		if err := schema.Seed(db); err != nil {
			log.Fatal("applying seed data", err)
		}
		log.Println("Seed data inserted")
		return
	}
}
