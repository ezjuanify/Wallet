package integration

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

const (
	TEST_CONTAINER_NAME = "wallet-test-db"
	PG_USER             = "test_db_wallet_app"
	PG_PASS             = "test_db_wallet_app"
	PG_DB               = "test_db_wallet_app"
	PG_HOST             = "localhost"
	PG_PORT             = 5433
	PG_SSL              = "disable"
)

func startTestDB() {
	cmd := exec.Command("docker", "run",
		"--rm",
		"--name", TEST_CONTAINER_NAME,
		"-e", fmt.Sprintf("POSTGRES_DB=%s", PG_DB),
		"-e", fmt.Sprintf("POSTGRES_USER=%s", PG_USER),
		"-e", fmt.Sprintf("POSTGRES_PASSWORD=%s", PG_PASS),
		"-p", fmt.Sprintf("%d:5432", PG_PORT),
		"-v", fmt.Sprintf("%s:/docker-entrypoint-initdb.d/init.sql", findInitSQL()),
		"-d", "postgres:17.5-alpine",
	)

	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("Failed to start test DB container: %v\nOutput: %s", err, out)
	}
}

func stopTestDB() {
	cmd := exec.Command("docker", "stop", TEST_CONTAINER_NAME)
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Failed to stop test DB container: %v\nOutput: %s", err, out)
	}
}

func findInitSQL() string {
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		log.Fatalf("Unable to determine current file path")
	}

	currentDir := filepath.Dir(thisFile)
	for {
		initSQL := filepath.Join(currentDir, "db", "init.sql")
		if _, err := os.Stat(initSQL); err == nil {
			absPath, err := filepath.Abs(initSQL)
			if err != nil {
				log.Fatalf("Failed to get absolute path: %v", err)
			}
			return absPath
		}

		parent := filepath.Dir(currentDir)
		if parent == currentDir {
			break
		}
		currentDir = parent
	}
	log.Fatalf("Could not find init.sql")
	return ""
}
