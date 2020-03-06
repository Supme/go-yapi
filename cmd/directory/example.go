package main

import (
	"context"
	"encoding/json"
	"fmt"
	yapi "go-yapi"
	"golang.org/x/oauth2"
	"log"
	"math/rand"
	"os"
	"time"
)

const tokFile = ".token"

// https://oauth.yandex.ru/
func main() {
	var clientID, clientSecret string

	clientID = os.Getenv("CLIENT_ID")
	clientSecret = os.Getenv("CLIENT_SECRET")
	//flag.StringVar(&clientID, "id", "", "Client ID")
	//flag.StringVar(&clientSecret, "secret", "", "Client secret")
	//flag.Parse()

	if clientID == "" {
		fmt.Println("Client ID required")
		os.Exit(1)
	}

	if clientSecret == "" {
		fmt.Println("Client secret required")
		os.Exit(1)
	}

	ctx := context.Background()

	conf := yapi.NewOauth2Config(clientID, clientSecret)

	var (
		tok *oauth2.Token
		err error
	)

	tok, err = yapi.TokenFromFile(tokFile)
	if err != nil {
		url := conf.AuthCodeURL("state", oauth2.AccessTypeOffline)
		fmt.Printf("Visit the URL for the auth dialog: %v\n", url)

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
	//fmt.Printf("Token: %#v\n", tok)

	client := conf.Client(ctx, tok)

	directory := yapi.NewDirectory(client)

	orgs, err := directory.GetOrganizations(yapi.DirectoryOrganizationAllParameters)
	if err != nil {
		log.Fatal("organizations ", err)
	}
	pretty("Organizations", *orgs)

	if len(orgs.Result) > 0 {
		rand.Seed(time.Now().Unix())

		users, err := directory.GetUsers(orgs.Result[0].ID, yapi.Parameters{"fields": []string{"id", "name", "gender", "nickname", "birthday"}, "per_page": []string{"3"}})
		if err != nil {
			log.Print("users ", err)
		} else {
			pretty("Users", *users)
		}

		user, err := directory.GetUser(orgs.Result[0].ID, users.Result[rand.Intn(len(users.Result))].ID, yapi.DirectoryUserAllParameters)
		if err != nil {
			log.Print("user ", err)
		} else {
			pretty("User", *user)
		}

		groups, err := directory.GetGroups(orgs.Result[0].ID, yapi.DirectoryGroupAllParameters)
		if err != nil {
			log.Print("groups ", err)
		} else {
			pretty("Groups", *groups)
		}

		group, err := directory.GetGroup(orgs.Result[0].ID, groups.Result[rand.Intn(len(groups.Result))].ID, yapi.DirectoryGroupAllParameters)
		if err != nil {
			log.Print("group ", err)
		} else {
			pretty("Group", *group)
		}

		deps, err := directory.GetDepartments(orgs.Result[0].ID, yapi.DirectoryDepartmentAllParameters)
		if err != nil {
			log.Print("deps ", err)
		} else {
			pretty("Departments", *deps)
		}

		dep, err := directory.GetDepartment(orgs.Result[0].ID, deps.Result[rand.Intn(len(deps.Result))].ID, yapi.DirectoryDepartmentAllParameters)
		if err != nil {
			log.Print("dep ", err)
		} else {
			pretty("Department", *dep)
		}

		domains, err := directory.GetDomains(orgs.Result[0].ID, yapi.DirectoryDomainAllParameters)
		if err != nil {
			log.Print("domains ", err)
		} else {
			pretty("Domains", domains)
		}
	}
}

func pretty(h string, v interface{}) {
	b, _ := json.MarshalIndent(v, "", "  ")
	fmt.Println("--------------", h)
	fmt.Println(string(b))
	fmt.Println()
}