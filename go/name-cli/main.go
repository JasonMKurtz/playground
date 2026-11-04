// Command name-cli is a small CLI for name.com DNS records and domain lookups.
//
// Credentials come from the environment:
//
//	NAMECOM_USERNAME  your name.com username (append "-test" for the sandbox)
//	NAMECOM_TOKEN     your API token
//	NAMECOM_API_URL   optional, default https://api.name.com
//	                  (sandbox: https://api.dev.name.com)
package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"
)

const usage = `Usage:
  name-cli records list   <domain>
  name-cli records get    <domain> <id>
  name-cli records create <domain> -type CNAME -host www -answer target.example.com [-ttl 300] [-priority 10]
  name-cli records update <domain> <id> [-type T] [-host H] [-answer A] [-ttl N] [-priority N]
           (only the flags given change; the rest is kept from the current record)
  name-cli records delete <domain> <id> [<id>...]
  name-cli domains list
  name-cli domains get    <domain>
  name-cli domains check  <domain> [<domain>...]

Environment:
  NAMECOM_USERNAME  name.com username (sandbox: append "-test")
  NAMECOM_TOKEN     API token
  NAMECOM_API_URL   default https://api.name.com (sandbox: https://api.dev.name.com)
`

var errUsage = errors.New("usage")

func main() {
	if err := run(context.Background(), os.Args[1:], os.Stdout); err != nil {
		if errors.Is(err, errUsage) {
			fmt.Fprint(os.Stderr, usage)
			os.Exit(2)
		}
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, args []string, out io.Writer) error {
	if len(args) < 2 {
		return errUsage
	}

	username := strings.TrimSpace(os.Getenv("NAMECOM_USERNAME"))
	token := strings.TrimSpace(os.Getenv("NAMECOM_TOKEN"))
	if username == "" || token == "" {
		return errors.New("NAMECOM_USERNAME and NAMECOM_TOKEN must be set")
	}
	baseURL := os.Getenv("NAMECOM_API_URL")
	if baseURL == "" {
		baseURL = "https://api.name.com"
	}
	c := NewClient(baseURL, username, token)

	switch args[0] + " " + args[1] {
	case "records list":
		return recordsList(ctx, c, args[2:], out)
	case "records get":
		return recordsGet(ctx, c, args[2:], out)
	case "records create":
		return recordsCreate(ctx, c, args[2:], out)
	case "records update":
		return recordsUpdate(ctx, c, args[2:], out)
	case "records delete":
		return recordsDelete(ctx, c, args[2:], out)
	case "domains list":
		return domainsList(ctx, c, args[2:], out)
	case "domains get":
		return domainsGet(ctx, c, args[2:], out)
	case "domains check":
		return domainsCheck(ctx, c, args[2:], out)
	}
	return errUsage
}

// parseInterspersed lets flags appear before or after positional args
// (the stdlib flag package stops at the first positional).
func parseInterspersed(fs *flag.FlagSet, args []string) ([]string, error) {
	var pos []string
	for {
		if err := fs.Parse(args); err != nil {
			return nil, err
		}
		if fs.NArg() == 0 {
			return pos, nil
		}
		pos = append(pos, fs.Arg(0))
		args = fs.Args()[1:]
	}
}

func recordsList(ctx context.Context, c *Client, args []string, out io.Writer) error {
	if len(args) != 1 {
		return errUsage
	}
	recs, err := c.ListRecords(ctx, args[0])
	if err != nil {
		return err
	}
	printRecords(out, recs)
	return nil
}

func recordsGet(ctx context.Context, c *Client, args []string, out io.Writer) error {
	if len(args) != 2 {
		return errUsage
	}
	id, err := strconv.Atoi(args[1])
	if err != nil {
		return fmt.Errorf("record id %q is not a number", args[1])
	}
	r, err := c.GetRecord(ctx, args[0], id)
	if err != nil {
		return err
	}
	printRecords(out, []Record{r})
	return nil
}

func recordsCreate(ctx context.Context, c *Client, args []string, out io.Writer) error {
	fs := flag.NewFlagSet("records create", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	typ := fs.String("type", "", "record type: A, AAAA, ANAME, CNAME, MX, NS, SRV, TXT")
	host := fs.String("host", "", `host relative to the zone ("" or "@" for apex)`)
	answer := fs.String("answer", "", "IP, target hostname or TXT value")
	ttl := fs.Int("ttl", 300, "TTL in seconds (minimum 300)")
	priority := fs.Int("priority", 0, "priority (MX and SRV only)")

	pos, err := parseInterspersed(fs, args)
	if err != nil {
		return fmt.Errorf("%w: %v", errUsage, err)
	}
	if len(pos) != 1 {
		return errUsage
	}
	if *typ == "" || *answer == "" {
		return errors.New("-type and -answer are required")
	}

	r, err := c.CreateRecord(ctx, pos[0], Record{
		Type:     strings.ToUpper(*typ),
		Host:     *host,
		Answer:   *answer,
		TTL:      *ttl,
		Priority: *priority,
	})
	if err != nil {
		return err
	}
	printRecords(out, []Record{r})
	return nil
}

// recordsUpdate changes only the fields whose flags were given. The API's
// update is a full replace, so it reads the record first and sends it back
// with those fields changed; otherwise e.g. `-answer x` would reset the TTL.
func recordsUpdate(ctx context.Context, c *Client, args []string, out io.Writer) error {
	fs := flag.NewFlagSet("records update", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	typ := fs.String("type", "", "record type: A, AAAA, ANAME, CNAME, MX, NS, SRV, TXT")
	host := fs.String("host", "", `host relative to the zone ("" or "@" for apex)`)
	answer := fs.String("answer", "", "IP, target hostname or TXT value")
	ttl := fs.Int("ttl", 0, "TTL in seconds (minimum 300)")
	priority := fs.Int("priority", 0, "priority (MX and SRV only)")

	pos, err := parseInterspersed(fs, args)
	if err != nil {
		return fmt.Errorf("%w: %v", errUsage, err)
	}
	if len(pos) != 2 {
		return errUsage
	}
	domain := pos[0]
	id, err := strconv.Atoi(pos[1])
	if err != nil {
		return fmt.Errorf("record id %q is not a number", pos[1])
	}

	// Which flags were actually passed (fs.Visit skips unset ones), so an
	// explicit -host "" (apex) still counts.
	set := map[string]bool{}
	fs.Visit(func(f *flag.Flag) { set[f.Name] = true })
	if len(set) == 0 {
		return errors.New("nothing to update: give at least one of -type, -host, -answer, -ttl, -priority")
	}

	r, err := c.GetRecord(ctx, domain, id)
	if err != nil {
		return err
	}
	if set["type"] {
		r.Type = strings.ToUpper(*typ)
	}
	if set["host"] {
		r.Host = *host
	}
	if set["answer"] {
		r.Answer = *answer
	}
	if set["ttl"] {
		r.TTL = *ttl
	}
	if set["priority"] {
		r.Priority = *priority
	}

	updated, err := c.UpdateRecord(ctx, domain, id, r)
	if err != nil {
		return err
	}
	printRecords(out, []Record{updated})
	return nil
}

// recordsDelete deletes each record ID in turn, looking each one up first so
// the output shows what was removed. Every ID must be a number before
// anything is deleted. It stops at the first API error, still listing what
// it had already deleted.
//
// name.com has been seen to answer a DELETE with success and keep the record
// (2026-11-04: deleting a host's A then AAAA back to back left the AAAA). So
// each delete is checked with a GET, retried once, and reported as an error
// if the record is still there.
func recordsDelete(ctx context.Context, c *Client, args []string, out io.Writer) error {
	if len(args) < 2 {
		return errUsage
	}
	domain := args[0]

	ids := make([]int, 0, len(args)-1)
	for _, a := range args[1:] {
		id, err := strconv.Atoi(a)
		if err != nil {
			return fmt.Errorf("record id %q is not a number", a)
		}
		ids = append(ids, id)
	}

	var deleted []Record
	var err error
	for _, id := range ids {
		var r Record
		if r, err = c.GetRecord(ctx, domain, id); err != nil {
			break
		}
		if err = deleteAndVerify(ctx, c, domain, id); err != nil {
			break
		}
		deleted = append(deleted, r)
	}

	if len(deleted) > 0 {
		fmt.Fprintln(out, "Deleted:")
		printRecords(out, deleted)
	}
	return err
}

// deleteAndVerify deletes record id and confirms it's gone (GET returns 404),
// deleting once more if it isn't.
func deleteAndVerify(ctx context.Context, c *Client, domain string, id int) error {
	for attempt := 1; attempt <= 2; attempt++ {
		if err := c.DeleteRecord(ctx, domain, id); err != nil {
			return err
		}
		_, err := c.GetRecord(ctx, domain, id)
		if IsNotFound(err) {
			return nil
		}
		if err != nil {
			return err
		}
	}
	return fmt.Errorf("record %d still exists after deleting it twice", id)
}

func domainsList(ctx context.Context, c *Client, args []string, out io.Writer) error {
	if len(args) != 0 {
		return errUsage
	}
	ds, err := c.ListDomains(ctx)
	if err != nil {
		return err
	}
	tw := tabwriter.NewWriter(out, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "DOMAIN\tEXPIRES\tAUTORENEW\tLOCKED")
	for _, d := range ds {
		fmt.Fprintf(tw, "%s\t%s\t%t\t%t\n", d.DomainName, d.ExpireDate, d.AutorenewEnabled, d.Locked)
	}
	return tw.Flush()
}

func domainsGet(ctx context.Context, c *Client, args []string, out io.Writer) error {
	if len(args) != 1 {
		return errUsage
	}
	d, err := c.GetDomain(ctx, args[0])
	if err != nil {
		return err
	}
	tw := tabwriter.NewWriter(out, 0, 0, 2, ' ', 0)
	fmt.Fprintf(tw, "Domain:\t%s\n", d.DomainName)
	fmt.Fprintf(tw, "Created:\t%s\n", d.CreateDate)
	fmt.Fprintf(tw, "Expires:\t%s\n", d.ExpireDate)
	fmt.Fprintf(tw, "Autorenew:\t%t\n", d.AutorenewEnabled)
	fmt.Fprintf(tw, "Locked:\t%t\n", d.Locked)
	fmt.Fprintf(tw, "Privacy:\t%t\n", d.PrivacyEnabled)
	fmt.Fprintf(tw, "Nameservers:\t%s\n", strings.Join(d.Nameservers, ", "))
	return tw.Flush()
}

func domainsCheck(ctx context.Context, c *Client, args []string, out io.Writer) error {
	if len(args) == 0 || len(args) > 50 {
		return errUsage
	}
	results, err := c.CheckAvailability(ctx, args)
	if err != nil {
		return err
	}
	tw := tabwriter.NewWriter(out, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "DOMAIN\tAVAILABLE\tPRICE\tRENEWAL\tNOTE")
	for _, r := range results {
		price, renewal := "-", "-"
		if r.Purchasable {
			price = fmt.Sprintf("$%.2f", r.PurchasePrice)
			renewal = fmt.Sprintf("$%.2f", r.RenewalPrice)
		}
		note := r.Reason
		if r.Premium {
			note = strings.TrimSpace("premium " + note)
		}
		fmt.Fprintf(tw, "%s\t%t\t%s\t%s\t%s\n", r.DomainName, r.Purchasable, price, renewal, note)
	}
	return tw.Flush()
}

func printRecords(out io.Writer, recs []Record) {
	tw := tabwriter.NewWriter(out, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "ID\tTYPE\tFQDN\tANSWER\tTTL\tPRIORITY")
	for _, r := range recs {
		prio := "-"
		if r.Priority != 0 {
			prio = strconv.Itoa(r.Priority)
		}
		fmt.Fprintf(tw, "%d\t%s\t%s\t%s\t%d\t%s\n", r.ID, r.Type, r.FQDN, r.Answer, r.TTL, prio)
	}
	tw.Flush()
}
