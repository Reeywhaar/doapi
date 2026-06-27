package cli

import (
	"doctl/internal"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
)

func Run(args []string, client *internal.Client) {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "commands: dns")
		os.Exit(1)
	}
	switch args[0] {
	case "dns":
		dns(args[1:], client)
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n", args[0])
		fmt.Fprintln(os.Stderr, "commands: dns")
		os.Exit(1)
	}
}

func dns(args []string, client *internal.Client) {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "dns commands: list, create, delete")
		os.Exit(1)
	}
	switch args[0] {
	case "list":
		list(client)
	case "create":
		create(args[1:], client)
	case "delete":
		del(args[1:], client)
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n", args[0])
		fmt.Fprintln(os.Stderr, "dns commands: list, create, delete")
		os.Exit(1)
	}
}

func list(client *internal.Client) {
	result, err := client.ListDomains()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	out, _ := json.MarshalIndent(result, "", "  ")
	fmt.Println(string(out))
}

func create(args []string, client *internal.Client) {
	if len(args) < 4 {
		fmt.Fprintln(os.Stderr, "usage: doapi dns create <domain> <type> <name> <data> [ttl]")
		os.Exit(1)
	}
	ttl := 0
	if len(args) >= 5 {
		var err error
		ttl, err = strconv.Atoi(args[4])
		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid ttl: %s\n", args[4])
			os.Exit(1)
		}
	}
	data, status, err := client.CreateRecord(args[0], args[1], args[2], args[3], ttl)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if status != http.StatusCreated {
		fmt.Fprintf(os.Stderr, "DO API error %d: %s\n", status, data)
		os.Exit(1)
	}
	var pretty any
	json.Unmarshal(data, &pretty)
	out, _ := json.MarshalIndent(pretty, "", "  ")
	fmt.Println(string(out))
}

func del(args []string, client *internal.Client) {
	if len(args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: doapi dns delete <domain> <id>")
		os.Exit(1)
	}
	id, err := strconv.Atoi(args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "invalid id: %s\n", args[1])
		os.Exit(1)
	}
	status, err := client.DeleteRecord(args[0], id)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if status != http.StatusNoContent {
		fmt.Fprintf(os.Stderr, "DO API error %d\n", status)
		os.Exit(1)
	}
	fmt.Println("deleted")
}
