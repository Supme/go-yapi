package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	yapi "go-yapi"
	"golang.org/x/oauth2"
	"log"
	"math/rand"
	"os"
	"time"
)

// https://oauth.yandex.ru/
func main() {
	var clientID, clientSecret, tokFile string
	var orgID int

	clientID = os.Getenv("CLIENT_ID")
	clientSecret = os.Getenv("CLIENT_SECRET")

	flag.StringVar(&tokFile, "f", ".token", "Token file")
	flag.IntVar(&orgID, "o", 0, "Organization id")
	flag.Parse()

	if orgID == 0 {
		fmt.Println("Organization ID required")
		os.Exit(1)
	}

	if clientID == "" {
		fmt.Println("Client ID required")
		os.Exit(1)
	}

	if clientSecret == "" {
		fmt.Println("Client secret required")
		os.Exit(1)
	}

	ctx := context.Background()

	conf := yapi.NewOauth2Config(clientID, clientSecret, []string{"directory:write_users"})

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

	client := conf.Client(ctx, tok)

	directory := yapi.NewDirectory(client)

	yapi.Debug = true

	newUser := yapi.DirectoryUser{
		DepartmentID: 1,
		Password:     "test100500",
		Nickname:     "testapi",
		Name: &yapi.DirectoryUserName{
			First: "Test",
			Last:  "Api",
		},
	}
	err = directory.CreateUser(orgID, &newUser)
	if err != nil {
		log.Print("Create user ", err)
		os.Exit(1)
	}
	pretty("Create user", newUser)

	err = directory.AddAliasUser(orgID, newUser.ID, "test-api")
	if err != nil {
		log.Print("Add alias user ", err)
	}
	pretty("Add alias", `{"name":"test-api"}`)

	//orgs, err := directory.GetOrganizations(yapi.DirectoryOrganizationAllParameters)
	//if err != nil {
	//	log.Fatal("organizations ", err)
	//}
	//pretty("Organizations", orgs)

	rand.Seed(time.Now().Unix())

	users, err := directory.GetUsers(orgID, yapi.Parameters{"fields": []string{"id", "name", "gender", "nickname", "birthday"}, "per_page": []string{"3"}})
	if err != nil {
		log.Print("users ", err)
	} else {
		pretty("Users", users)
	}

	user, err := directory.GetUser(orgID, users.Result[rand.Intn(len(users.Result))].ID, yapi.DirectoryUserAllParameters)
	if err != nil {
		log.Print("user ", err)
	} else {
		pretty("User", user)
	}

	groups, err := directory.GetGroups(orgID, yapi.DirectoryGroupAllParameters)
	if err != nil {
		log.Print("groups ", err)
	} else {
		pretty("Groups", groups)
	}

	group, err := directory.GetGroup(orgID, groups.Result[rand.Intn(len(groups.Result))].ID, yapi.DirectoryGroupAllParameters)
	if err != nil {
		log.Print("group ", err)
	} else {
		pretty("Group", group)
	}

	newDep, err := directory.CreateDepartment(orgID, yapi.DirectoryNewDepartment{ParentID: 1, HeadID: newUser.ID, Name: "New test department"})
	if err != nil {
		log.Print("dep ", err)
	} else {
		pretty("New department", newDep)
	}

	modDep, err := directory.ModifyDepartment(orgID, newDep.ID, yapi.DirectoryNewDepartment{Name: "Modify test department"})
	if err != nil {
		log.Print("dep ", err)
	} else {
		pretty("Modify department", modDep)
	}

	deps, err := directory.GetDepartments(orgID, yapi.DirectoryDepartmentAllParameters)
	if err != nil {
		log.Print("deps ", err)
	} else {
		pretty("Departments", deps)
	}

	dep, err := directory.GetDepartment(orgID, deps.Result[rand.Intn(len(deps.Result))].ID, yapi.DirectoryDepartmentAllParameters)
	if err != nil {
		log.Print("dep ", err)
	} else {
		pretty("Department", dep)
	}

	err = directory.DeleteDepartment(orgID, newDep.ID)
	if err != nil {
		log.Print("dep ", err)
	}
	pretty("Delete department", "Success")

	modifyUser := yapi.DirectoryUser{IsDismissed: true}
	err = directory.ModifyUser(orgID, newUser.ID, &modifyUser)
	if err != nil {
		log.Print("Modify user ", err)
		os.Exit(1)
	}
	pretty("Modify user", modifyUser)

	domains, err := directory.GetDomains(orgID, yapi.DirectoryDomainAllParameters)
	if err != nil {
		log.Print("domains ", err)
	} else {
		pretty("Domains", domains)
	}
}

func pretty(h string, v interface{}) {
	b, _ := json.MarshalIndent(v, "", "  ")
	fmt.Println("--------------", h)
	fmt.Println(string(b))
	fmt.Println()
}
