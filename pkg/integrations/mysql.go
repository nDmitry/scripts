package integrations

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

type MysqlDumper struct {
	Container string
	User      string
	Password  string
	Database  string
	Outfile   string
}

func (md *MysqlDumper) Dump() error {
	dockerExe, err := exec.LookPath("docker")

	if err != nil {
		return fmt.Errorf("could not lookup the docker location: %w", err)
	}

	cmd := exec.Command(
		"sh", "-c",
		strings.Join([]string{
			dockerExe, "exec", md.Container,
			"bash", "-lc",
			fmt.Sprintf(
				"\"mysqldump -u %s -p%s %s\"",
				md.User, md.Password, md.Database,
			),
			"|", "gzip", ">", md.Outfile,
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
