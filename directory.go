// https://yandex.ru/dev/connect/directory/
package go_yapi

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

const (
	directoryAPI = "https://api.directory.yandex.net/"
	directoryURL = directoryAPI + VersionAPI
)

type Directory struct {
	client *http.Client
}

func NewDirectory(client *http.Client) *Directory {
	return &Directory{client: client}
}

//     ____ ___
//    |    |   \______ ___________  ______
//    |    |   /  ___// __ \_  __ \/  ___/
//    |    |  /\___ \\  ___/|  | \/\___ \
//    |______//____  >\___  >__|  /____  >
//                 \/     \/           \/

type DirectoryUser struct {
	IsRobot      bool                      `json:"is_robot,omitempty"`
	ExternalID   interface{}               `json:"external_id,omitempty"`
	Position     string                    `json:"position,omitempty"`
	Departments  []DirectoryUserDepartment `json:"departments,omitempty"`
	OrgID        int                       `json:"org_id,omitempty"`
	Gender       string                    `json:"gender,omitempty"`
	Created      string                    `json:"created,omitempty"`
	Name         *DirectoryUserName        `json:"name,omitempty"`
	About        string                    `json:"about,omitempty"`
	Nickname     string                    `json:"nickname,omitempty"`
	Groups       []DirectoryUserGroup      `json:"groups,omitempty"`
	IsAdmin      bool                      `json:"is_admin,omitempty"`
	Birthday     string                    `json:"birthday,omitempty"`
	DepartmentID int                       `json:"department_id,omitempty"`
	Email        string                    `json:"email,omitempty"`
	Department   *DirectoryUserDepartment  `json:"department,omitempty"`
	Contacts     []DirectoryUserContact    `json:"contacts,omitempty"`
	Aliases      []string                  `json:"aliases,omitempty"`
	ID           int                       `json:"id,omitempty"`
	IsDismissed  bool                      `json:"is_dismissed,omitempty"`
	Password     string                    `json:"password,omitempty"`
}

type directoryUserSliceID struct {
	ID int `json:"id"`
}
type DirectoryUserDepartment directoryUserSliceID
type DirectoryUserGroup directoryUserSliceID

type DirectoryUserName struct {
	First  string `json:"first,omitempty"`
	Last   string `json:"last,omitempty"`
	Middle string `json:"middle,omitempty"`
}

type DirectoryUserContact struct {
	Value     string `json:"value,omitempty"`
	Type      string `json:"type,omitempty"`
	Main      bool   `json:"main,omitempty"`
	Alias     bool   `json:"alias,omitempty"`
	Synthetic bool   `json:"synthetic,omitempty"`
}

type DirectoryUsers struct {
	Page    int             `json:"page"`
	Total   int             `json:"total"`
	PerPage int             `json:"per_page"`
	Result  []DirectoryUser `json:"result"`
	Pages   int             `json:"pages"`
	Links   struct {
		Next  string `json:"next"`
		Prev  string `json:"prev"`
		Last  string `json:"last"`
		First string `json:"first"`
	} `json:"links"`
}

var DirectoryUserAllParameters = Parameters{
	"fields": []string{
		"is_robot",
		"external_id",
		"departments",
		"org_id",
		"gender",
		"created",
		"name",
		"about",
		"nickname",
		"groups",
		"is_admin",
		"birthday",
		"department_id",
		"email",
		"department",
		"contacts",
		"aliases",
		"id",
		"is_dismissed",
	},
}

// GetUsers ...
func (d Directory) GetUsers(orgID int, params Parameters) (DirectoryUsers, error) {
	var users DirectoryUsers
	err := Get(
		d.client,
		directoryURL+"/users/",
		params,
		headerOrgID(orgID),
		&users,
	)
	return users, err
}

// GetUser ...
func (d Directory) GetUser(orgID, userID int, params Parameters) (DirectoryUser, error) {
	var user DirectoryUser
	err := Get(
		d.client,
		directoryURL+"/users/"+strconv.Itoa(userID)+"/",
		params,
		headerOrgID(orgID),
		&user,
	)
	return user, err
}

func (d Directory) CreateUser(orgID int, user *DirectoryUser) error {
	j, err := json.Marshal(user)
	if err != nil {
		return err
	}

	err = Post(
		d.client,
		directoryURL+"/users/",
		nil,
		headerOrgID(orgID),
		bytes.NewReader(j),
		&user,
	)

	return err
}

func (d Directory) ModifyUser(orgID, userID int, user *DirectoryUser) error {
	j, err := json.Marshal(user)
	if err != nil {
		return err
	}

	err = Patch(
		d.client,
		directoryURL+"/users/"+strconv.Itoa(userID)+"/",
		nil,
		headerOrgID(orgID),
		bytes.NewReader(j),
		&user,
	)

	return err
}

func (d Directory) AddAliasUser(orgID, userID int, alias string) error {
	return Patch(
		d.client,
		directoryURL+"/users/"+strconv.Itoa(userID)+"/aliases/",
		nil,
		headerOrgID(orgID),
		strings.NewReader(`{"name": `+jsonParam(alias)+`}`),
		nil,
	)
}

//    ________                              __                         __
//    \______ \   ____ ___________ ________/  |_  _____   ____   _____/  |_  ______
//     |    |  \_/ __ \\____ \__  \\_  __ \   __\/     \_/ __ \ /    \   __\/  ___/
//     |    `   \  ___/|  |_> > __ \|  | \/|  | |  Y Y  \  ___/|   |  \  |  \___ \
//    /_______  /\___  >   __(____  /__|   |__| |__|_|  /\___  >___|  /__| /____  >
//            \/     \/|__|       \/                  \/     \/     \/          \/

type DirectoryDepartment struct {
	Name         string                      `json:"name,omitempty"`
	Email        string                      `json:"email,omitempty"`
	ExternalID   interface{}                 `json:"external_id,omitempty"`
	Removed      bool                        `json:"removed,omitempty"`
	ID           int                         `json:"id,omitempty"`
	Parents      []DirectoryDepartmentParent `json:"parents,omitempty"`
	Label        string                      `json:"label,omitempty"`
	Created      string                      `json:"created,omitempty"`
	Parent       DirectoryDepartmentParent   `json:"parent,omitempty"`
	Description  string                      `json:"description,omitempty"`
	MembersCount int                         `json:"members_count,omitempty"`
	Head         int                         `json:"head,omitempty"`
}

type DirectoryDepartmentParent struct {
	Name         string      `json:"name,omitempty"`
	Email        string      `json:"email,omitempty"`
	ExternalID   interface{} `json:"external_id,omitempty"`
	Removed      bool        `json:"removed,omitempty"`
	ID           int         `json:"id"`
	ParentID     int         `json:"parent_id,omitempty"`
	Label        string      `json:"label,omitempty"`
	Created      string      `json:"created,omitempty"`
	Description  string      `json:"description,omitempty"`
	MembersCount int         `json:"members_count,omitempty"`
}

type DirectoryDepartments struct {
	Page    int                   `json:"page"`
	Total   int                   `json:"total"`
	PerPage int                   `json:"per_page"`
	Result  []DirectoryDepartment `json:"result"`
	Pages   int                   `json:"pages"`
	Links   struct {
		Next  string `json:"next"`
		Prev  string `json:"prev"`
		Last  string `json:"last"`
		First string `json:"first"`
	} `json:"links"`
}

var DirectoryDepartmentAllParameters = Parameters{
	"fields": []string{
		"name",
		"email",
		"external_id",
		"removed",
		"id",
		"parents",
		"label",
		"created",
		"parent",
		"description",
		"members_count",
		"head",
	},
}

// GetDepartments ...
func (d Directory) GetDepartments(orgID int, params Parameters) (DirectoryDepartments, error) {
	var departments DirectoryDepartments
	err := Get(
		d.client,
		directoryURL+"/departments/",
		params,
		headerOrgID(orgID),
		&departments,
	)
	return departments, err
}

// GetDepartment ...
func (d Directory) GetDepartment(orgID, depID int, params Parameters) (DirectoryDepartment, error) {
	var department DirectoryDepartment
	err := Get(
		d.client,
		directoryURL+"/departments/"+strconv.Itoa(depID)+"/",
		params,
		headerOrgID(orgID),
		&department,
	)
	return department, err
}

// CreateDepartment ...
func (d Directory) CreateDepartment(orgID, parentID, headID int, name, label, description *string) (DirectoryDepartment, error) {
	var department DirectoryDepartment
	err := Post(
		d.client,
		directoryURL+"/departments/",
		nil,
		headerOrgID(orgID),
		strings.NewReader(
			`{`+
				`"parent_id": `+jsonParam(parentID)+
				`, "head_id": `+jsonParam(headID)+
				`, "label": `+jsonParam(label)+
				`, "name": `+jsonParam(name)+
				`, "description": `+jsonParam(description)+
				`}`,
		),
		&department,
	)
	return department, err
}

// ChangeDepartment ...
func (d Directory) ChangeDepartment(orgID, depID, parentID, headID int, name, label, description *string) (DirectoryDepartment, error) {
	var department DirectoryDepartment
	err := Patch(
		d.client,
		directoryURL+"/departments/"+strconv.Itoa(depID)+"/",
		nil,
		headerOrgID(orgID),
		strings.NewReader(
			`{`+
				`"parent_id": `+jsonParam(parentID)+
				`, "head_id": `+jsonParam(headID)+
				`, "label": `+jsonParam(label)+
				`, "name": `+jsonParam(name)+
				`, "description": `+jsonParam(description)+
				`}`,
		),
		&department,
	)
	return department, err
}

// DeleteDepartment ...
func (d Directory) DeleteDepartment(orgID, depID int) error {
	return Delete(
		d.client,
		directoryURL+"/departments/"+strconv.Itoa(depID)+"/",
		nil,
		headerOrgID(orgID),
	)
}

//      ________
//     /  _____/______  ____  __ ________  ______
//    /   \  __\_  __ \/  _ \|  |  \____ \/  ___/
//    \    \_\  \  | \(  <_> )  |  /  |_> >___ \
//     \______  /__|   \____/|____/|   __/____  >
//            \/                   |__|       \/

type DirectoryGroup struct {
	Name       string `json:"name,omitempty"`
	Email      string `json:"email,omitempty"`
	ExternalID string `json:"external_id,omitempty"`
	ID         int    `json:"id,omitempty"`
	Members    []struct {
		Type   string `json:"type"` // <user|group|department>
		Object struct {
			ID int `json:"id"`
		} `json:"object"`
	} `json:"members,omitempty"`
	Label   string `json:"label,omitempty"`
	Created string `json:"created,omitempty"`
	Type    string `json:"type,omitempty"`
	Admins  []struct {
		Aliases      []string `json:"aliases"`
		ID           int      `json:"id"`
		Nickname     string   `json:"nickname"`
		DepartmentID int      `json:"department_id"`
		IsDismissed  bool     `json:"is_dismissed"`
		Position     string   `json:"position"`
		Groups       []struct {
			ID int `json:"id"`
		} `json:"groups,omitempty"`
		IsAdmin    bool   `json:"is_admin"`
		Birthday   string `json:"birthday"`
		Email      string `json:"email"`
		ExternalID string `json:"external_id"`
		Gender     string `json:"gender"`
		Contacts   []struct {
			Value     string `json:"value"`
			Type      string `json:"type"`
			Main      bool   `json:"main"`
			Alias     bool   `json:"alias"`
			Synthetic bool   `json:"synthetic"`
		} `json:"contacts"`
		Name struct {
			First  string `json:"first"`
			Last   string `json:"last"`
			Middle string `json:"middle"`
		} `json:"name"`
		About string `json:"about"`
	} `json:"admins,omitempty"`
	Author struct {
		Aliases      []string `json:"aliases"`
		ID           int      `json:"id"`
		Nickname     string   `json:"nickname"`
		DepartmentID int      `json:"department_id"`
		IsDismissed  bool     `json:"is_dismissed"`
		Position     string   `json:"position"`
		Groups       []struct {
			ID int `json:"id"`
		} `json:"groups"`
		IsAdmin    bool   `json:"is_admin"`
		Birthday   string `json:"birthday"`
		Email      string `json:"email"`
		ExternalID string `json:"external_id"`
		Gender     string `json:"gender"`
		Contacts   []struct {
			Value     string `json:"value"`
			Type      string `json:"type"`
			Main      bool   `json:"main"`
			Alias     bool   `json:"alias"`
			Synthetic bool   `json:"synthetic"`
		} `json:"contacts"`
		Name struct {
			First  string `json:"first"`
			Last   string `json:"last"`
			Middle string `json:"middle"`
		} `json:"name"`
		About string `json:"about"`
	} `json:"author,omitempty"`
	Description  string `json:"description,omitempty"`
	MembersCount int    `json:"members_count,omitempty"`
	MemberOf     []int  `json:"member_of,omitempty"`
}

type DirectoryGroups struct {
	Page    int              `json:"page"`
	Total   int              `json:"total"`
	PerPage int              `json:"per_page"`
	Result  []DirectoryGroup `json:"result"`
	Pages   int              `json:"pages"`
	Links   struct {
		Next  string `json:"next"`
		Prev  string `json:"prev"`
		Last  string `json:"last"`
		First string `json:"first"`
	} `json:"links"`
}

type DirectoryGroupMember struct {
	Type   string `json:"type"`
	Object struct {
		DepartmentID int    `json:"department_id"`
		ID           int    `json:"id"`
		Nickname     string `json:"nickname"`
		Email        string `json:"email"`
		Gender       string `json:"gender"`
		Name         struct {
			First  string `json:"first"`
			Last   string `json:"last"`
			Middle string `json:"middle"`
		} `json:"name"`
	} `json:"object"`
}

var DirectoryGroupAllParameters = Parameters{
	"fields": []string{
		"name",
		"email",
		"external_id",
		"id",
		"members",
		"label",
		"created",
		"type",
		"admins",
		"author",
		"description",
		"members_count",
		"member_of",
	},
}

func (d Directory) GetGroups(orgID int, params Parameters) (DirectoryGroups, error) {
	var groups DirectoryGroups
	err := Get(
		d.client,
		directoryURL+"/groups/",
		params,
		headerOrgID(orgID),
		&groups,
	)
	return groups, err
}

func (d Directory) GetGroup(orgID, groupID int, params Parameters) (DirectoryGroup, error) {
	var group DirectoryGroup
	err := Get(
		d.client,
		directoryURL+"/groups/"+strconv.Itoa(groupID),
		params,
		headerOrgID(orgID),
		&group,
	)
	return group, err
}

//    ________                        .__
//    \______ \   ____   _____ _____  |__| ____   ______
//     |    |  \ /  _ \ /     \\__  \ |  |/    \ /  ___/
//     |    `   (  <_> )  Y Y  \/ __ \|  |   |  \\___ \
//    /_______  /\____/|__|_|  (____  /__|___|  /____  >
//            \/             \/     \/        \/     \/

type DirectoryDomain struct {
	Mx            bool   `json:"mx"`
	Delegated     bool   `json:"delegated"`
	Tech          bool   `json:"tech"`
	PopEnabled    bool   `json:"pop_enabled"`
	Master        bool   `json:"master"`
	PostmasterUID int    `json:"postmaster_uid"`
	Owned         bool   `json:"owned"`
	Country       string `json:"country"`
	Name          string `json:"name"`
	ImapEnabled   bool   `json:"imap_enabled"`
}

var DirectoryDomainAllParameters = Parameters{
	"fields": []string{
		"mx",
		"delegated",
		"tech",
		"pop_enabled",
		"master",
		"postmaster_uid",
		"owned",
		"country",
		"name",
		"imap_enabled",
	},
}

func (d Directory) GetDomains(orgID int, params Parameters) ([]DirectoryDomain, error) {
	var domains []DirectoryDomain
	err := Get(
		d.client,
		directoryURL+"/domains/",
		params,
		headerOrgID(orgID),
		&domains,
	)
	return domains, err
}

//
//    ________                            .__                __  .__
//    \_____  \_______  _________    ____ |__|____________ _/  |_|__| ____   ____   ______
//     /   |   \_  __ \/ ___\__  \  /    \|  \___   /\__  \\   __\  |/  _ \ /    \ /  ___/
//    /    |    \  | \/ /_/  > __ \|   |  \  |/    /  / __ \|  | |  (  <_> )   |  \\___ \
//    \_______  /__|  \___  (____  /___|  /__/_____ \(____  /__| |__|\____/|___|  /____  >
//            \/     /_____/     \/     \/         \/     \/                    \/     \/
//

type DirectoryOrganizations struct {
	Links  interface{} `json:"links"`
	Result []struct {
		Revision int    `json:"revision,omitempty"`
		ID       int    `json:"id,omitempty"`
		Label    string `json:"label,omitempty"`
		Domains  struct {
			Display string   `json:"display"`
			Master  string   `json:"master"`
			All     []string `json:"all"`
		} `json:"domains,omitempty"`
		AdminUID int    `json:"admin_uid,omitempty"`
		Email    string `json:"email,omitempty"`
		Services []struct {
			Slug  string `json:"slug"`
			Ready bool   `json:"ready"`
		} `json:"services,omitempty"`
		DiskLimit        int    `json:"disk_limit,omitempty"`
		SubscriptionPlan string `json:"subscription_plan,omitempty"`
		Country          string `json:"country,omitempty"`
		Language         string `json:"language,omitempty"`
		Name             string `json:"name,omitempty"`
		Fax              string `json:"fax,omitempty"`
		DiskUsage        int    `json:"disk_usage,omitempty"`
		PhoneNumber      string `json:"phone_number,omitempty"`
	} `json:"result"`
}

var DirectoryOrganizationAllParameters = Parameters{
	"fields": []string{
		"revision",
		"id",
		"label",
		"domains",
		"admin_uid",
		"email",
		"services",
		"disk_limit",
		"subscription_plan",
		"country",
		"language",
		"name",
		"fax",
		"disk_usage",
		"phone_number",
	},
}

// GetOrganizations ...
func (d Directory) GetOrganizations(params Parameters) (DirectoryOrganizations, error) {
	var organizations DirectoryOrganizations
	err := Get(
		d.client,
		directoryURL+"/organizations/",
		params,
		nil,
		&organizations,
	)
	return organizations, err
}
