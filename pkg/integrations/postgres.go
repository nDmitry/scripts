package integrations

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

var container = os.Getenv("CONTAINER")
var user = os.Getenv("USER")
var database = os.Getenv("DATABASE")
var outfile = os.Getenv("OUTFILE")

type PostgresDumper struct{}

func (pd *PostgresDumper) DumpDocker() error {
	dockerExe, err := exec.LookPath("docker")

	if err != nil {
		return fmt.Errorf("could not lookup the docker location: %w", err)
	}

	cmd := exec.Command(
		"sh", "-c",
		strings.Join([]string{
			dockerExe, "exec", "--user", user, container,
			"bash", "-lc",
			fmt.Sprintf("\"pg_dump --format custom %s\"", database),
			">", outfile,
		}, " "),
	)

	cmd.Env = os.Environ()

	log.Println("Running the command:", cmd.String())

	out, err := cmd.CombinedOutput()

	log.Println("Dump command output:", string(out))

	if err != nil {
		return fmt.Errorf("could not run the docker command: %w", err)
	}

	return nil
}
