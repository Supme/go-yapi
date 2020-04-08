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
	"net/http"
	_ "net/http/pprof"
	"os"
	"sort"
	"strconv"
	"sync"
	"time"
)

func main() {
	var clientID, clientSecret, tokFile, webPort string
	var startWeb bool
	var orgID int

	flag.StringVar(&clientID, "i", "", "Client ID")
	flag.StringVar(&clientSecret, "s", "", "Client secret")
	flag.IntVar(&orgID, "o", 0, "Organization ID (default all")
	flag.StringVar(&webPort, "p", "8080", "Listen port")
	flag.BoolVar(&startWeb, "w", false, "Start web server")
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

	if startWeb {
		fmt.Println("Start web server on port", webPort)
		log.Fatal(http.ListenAndServe(":"+webPort, handler(yapi.NewDirectory(client), orgID)))
	} else {
		buf := &bytes.Buffer{}
		err = printTable(yapi.NewDirectory(client), orgID, buf)
		if err != nil {
			log.Fatal(err)
		}
		_, err = io.Copy(os.Stdout, buf)
		if err != nil {
			log.Fatal(err)
		}
	}
}

var stor = storage{
	data: "Initialize storage",
}

type storage struct {
	data string
	sync.RWMutex
}

func (s *storage) update(directory *yapi.Directory, orgID int) error {
	buf := &bytes.Buffer{}
	err := printTable(directory, orgID, buf)
	if err != nil {
		return err
	}
	s.Lock()
	defer s.Unlock()
	s.data = buf.String()
	return nil
}

func handler(directory *yapi.Directory, orgID int) http.HandlerFunc {
	go func(directory *yapi.Directory) {
		for {
			log.Print("start update storage")
			err := stor.update(directory, orgID)
			if err != nil {
				log.Print(err)
				continue
			}
			log.Print("storage updated")
			time.Sleep(time.Minute * 10)
		}
	}(directory)

	return func(w http.ResponseWriter, r *http.Request) {
		stor.RLock()
		defer stor.RUnlock()
		w.Header().Add("Content-Type", "text/plain; charset=utf-8")
		w.Write([]byte(stor.data))
	}
}

type SortByEmail []yapi.DirectoryUser

func (e SortByEmail) Len() int           { return len(e) }
func (e SortByEmail) Swap(i, j int)      { e[i], e[j] = e[j], e[i] }
func (e SortByEmail) Less(i, j int) bool { return e[i].Email < e[j].Email }

func printTable(directory *yapi.Directory, orgID int, w io.Writer) error {
	orgs, err := directory.GetOrganizations(nil)
	if err != nil {
		return fmt.Errorf("get organizations %w", err)
	}

	if orgID != 0 {
		for n := range orgs.Result {
			if orgs.Result[n].ID == orgID {
				orgs.Result = orgs.Result[n : n+1]
				break
			}
		}
	}

	for n := range orgs.Result {
		_, err = fmt.Fprintln(w, "╒═══════════════════════════════════════════════════════════════════╕")
		if err != nil {
			return err
		}
		_, err = fmt.Fprintf(w, "%-15s %37s %15s\n", "│       ", "Organisation ID: "+strconv.Itoa(orgs.Result[n].ID)+" (name: \""+orgs.Result[n].Name+"\")", "       │")
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

			users, err := directory.GetUsers(orgs.Result[n].ID, yapi.Parameters{"fields": []string{"name", "email"}, "group_id": []string{strconv.Itoa(groups.Result[i].ID)}, "per_page": []string{"1000"}})
			if err != nil {
				return fmt.Errorf("get users %w", err)
			} else {
				sort.Sort(SortByEmail(users.Result))
				for _, u := range users.Result {
					_, err = fmt.Fprintf(w, "│ %-29s │ %-15s │ %-15s │\n", u.Email, u.Name.First, u.Name.Last)
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
