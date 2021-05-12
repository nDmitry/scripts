package integrations

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

type PostgresDumper struct {
	Container string
	User      string
	Database  string
	Outfile   string
}

func (pd *PostgresDumper) DumpDocker() error {
	dockerExe, err := exec.LookPath("docker")

	if err != nil {
		return fmt.Errorf("could not lookup the docker location: %w", err)
	}

	cmd := exec.Command(
		"sh", "-c",
		strings.Join([]string{
			dockerExe, "exec", "--user", pd.User, pd.Container,
			"bash", "-lc",
			fmt.Sprintf("\"pg_dump --format custom %s\"", pd.Database),
			">", pd.Outfile,
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
