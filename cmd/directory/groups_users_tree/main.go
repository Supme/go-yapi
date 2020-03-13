package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	yapi "go-yapi"
	"golang.org/x/oauth2"
	"io"
	"log"
	"os"
	"strconv"
)

func main() {
	var clientID, clientSecret, tokFile string

	flag.StringVar(&clientID, "i", "", "Client ID")
	flag.StringVar(&clientSecret, "s", "", "Client secret")
	flag.StringVar(&tokFile, "f", ".token", "Token file")
	flag.Parse()

	ctx := context.Background()

	conf := yapi.NewOauth2Config(clientID, clientSecret, nil)

	var (
		tok *oauth2.Token
		err error
	)

	tok, err = yapi.TokenFromFile(tokFile)
	if err != nil {
		if clientID == "" {
			if clientID = os.Getenv("CLIENT_ID"); clientID == "" {
				fmt.Println("Client ID required")
				os.Exit(1)
			}
		}

		if clientSecret == "" {
			if clientSecret = os.Getenv("CLIENT_SECRET"); clientSecret == "" {
				fmt.Println("Client secret required")
				os.Exit(1)
			}
		}

		url := conf.AuthCodeURL("state", oauth2.AccessTypeOffline)
		fmt.Printf("Visit the URL for the auth dialog:\n%s\nAnd enter code: ", url)

		var code string
		if _, err := fmt.Scan(&code); err != nil {
			log.Fatal(err)
		}

		tok, err = conf.Exchange(ctx, code)
		if err != nil {
			log.Fatal(err)
		}
		tok.TokenType = "oauth"

		err = yapi.TokenToFile(tokFile, tok)
		if err != nil {
			log.Fatal(err)
		}
	}

	client := conf.Client(ctx, tok)

	buf := &bytes.Buffer{}
	err = printTable(yapi.NewDirectory(client), buf)
	if err != nil {
		log.Fatal(err)
	}
	_, err = io.Copy(os.Stdout, buf)
	if err != nil {
		log.Fatal(err)
	}
}

func printTable(directory *yapi.Directory, w io.Writer) error {
	orgs, err := directory.GetOrganizations(nil)
	if err != nil {
		log.Fatal("organizations ", err)
	}

	for n := range orgs.Result {
		_, err = fmt.Fprintln(w, "╒═══════════════════════════════════════════════════════════════════╕")
		if err != nil {
			return err
		}
		_, err = fmt.Fprintf(w, "%-15s %37s %15s\n", "│       ", "Organisation ID: "+strconv.Itoa(orgs.Result[n].ID), "       │")
		if err != nil {
			return err
		}
		groups, err := directory.GetGroups(orgs.Result[n].ID, yapi.Parameters{"fields": []string{"id", "name", "email"}, "per_page": []string{"1000"}})
		if err != nil {
			return fmt.Errorf("get groups %w", err)
		}

		for i := range groups.Result {
			if i == 0 {
				_, err = fmt.Fprintln(w, "├───────────────────────────────────────────────────────────────────┤")
				if err != nil {
					return err
				}
			}
			if groups.Result[i].Email == "" {
				continue
			}
			_, err = fmt.Fprintf(w, "│ %-65s │\n", groups.Result[i].Email+" ("+groups.Result[i].Name+")")
			if err != nil {
				return err
			}
			_, err = fmt.Fprintln(w, "├───────────────────────────────┬─────────────────┬─────────────────┤")
			if err != nil {
				return err
			}

			users, err := directory.GetUsers(orgs.Result[0].ID, yapi.Parameters{"fields": []string{"name", "email"}, "group_id": []string{strconv.Itoa(groups.Result[i].ID)}, "per_page": []string{"1000"}})
			if err != nil {
				return fmt.Errorf("get users %w", err)
			} else {
				for _, u := range users.Result {
					_, err = fmt.Fprintf(w, "│ %-29s │ %-15s │ %-15s │\n", u.Email, u.Name.Last, u.Name.First)
					if err != nil {
						return err
					}
				}
			}

			if i != len(groups.Result)-1 {
				_, err = fmt.Fprintln(w, "├───────────────────────────────┴─────────────────┴─────────────────┤")
				if err != nil {
					return err
				}
			}
		}

		_, err = fmt.Fprintln(w, "╘═══════════════════════════════╧═════════════════╧═════════════════╛")
		if err != nil {
			return err
		}

	}

	return nil
}
