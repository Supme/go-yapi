// https://yandex.ru/dev/connect/directory/
package go_yapi

import (
	"net/http"
	"strconv"
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
	IsRobot     bool        `json:"is_robot"`
	ExternalID  interface{} `json:"external_id"`
	Position    string      `json:"position"`
	Departments []struct {
		ID int `json:"id"`
	} `json:"departments"`
	OrgID   int    `json:"org_id"`
	Gender  string `json:"gender"`
	Created string `json:"created"`
	Name    struct {
		First  string `json:"first"`
		Last   string `json:"last"`
		Middle string `json:"middle"`
	} `json:"name"`
	About    string `json:"about"`
	Nickname string `json:"nickname"`
	Groups   []struct {
		ID int `json:"id"`
	} `json:"groups"`
	IsAdmin      bool   `json:"is_admin"`
	Birthday     string `json:"birthday"`
	DepartmentID int    `json:"department_id"`
	Email        string `json:"email"`
	Department   struct {
		ID int `json:"id"`
	} `json:"department"`
	Contacts []struct {
		Value     string `json:"value"`
		Type      string `json:"type"`
		Main      bool   `json:"main"`
		Alias     bool   `json:"alias"`
		Synthetic bool   `json:"synthetic"`
	} `json:"contacts"`
	Aliases     []string `json:"aliases"`
	ID          int      `json:"id"`
	IsDismissed bool     `json:"is_dismissed"`
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
func (d Directory) GetUsers(orgID int, params Parameters) (*DirectoryUsers, error) {
	var users DirectoryUsers
	err := Get(
		d.client,
		directoryURL+"/users/",
		params,
		map[string]string{"X-Org-ID": strconv.Itoa(orgID)},
		&users,
	)
	return &users, err
}

// GetUser ...
func (d Directory) GetUser(orgID, userID int, params Parameters) (*DirectoryUser, error) {
	var user DirectoryUser
	err := Get(
		d.client,
		directoryURL+"/users/"+strconv.Itoa(userID),
		params,
		map[string]string{"X-Org-ID": strconv.Itoa(orgID)},
		&user,
	)
	return &user, err
}

//    ________                              __                         __
//    \______ \   ____ ___________ ________/  |_  _____   ____   _____/  |_  ______
//     |    |  \_/ __ \\____ \__  \\_  __ \   __\/     \_/ __ \ /    \   __\/  ___/
//     |    `   \  ___/|  |_> > __ \|  | \/|  | |  Y Y  \  ___/|   |  \  |  \___ \
//    /_______  /\___  >   __(____  /__|   |__| |__|_|  /\___  >___|  /__| /____  >
//            \/     \/|__|       \/                  \/     \/     \/          \/

type DirectoryDepartment struct {
	Name       string      `json:"name"`
	Email      string      `json:"email"`
	ExternalID interface{} `json:"external_id"`
	Removed    bool        `json:"removed"`
	ID         int         `json:"id"`
	Parents    []struct {
		Name         string      `json:"name"`
		Email        string      `json:"email"`
		ExternalID   interface{} `json:"external_id"`
		Removed      bool        `json:"removed"`
		ID           int         `json:"id"`
		ParentID     int         `json:"parent_id"`
		Label        string      `json:"label"`
		Created      string      `json:"created"`
		Description  string      `json:"description"`
		MembersCount int         `json:"members_count"`
	} `json:"parents"`
	Label   string `json:"label"`
	Created string `json:"created"`
	Parent  struct {
		Name       string      `json:"name"`
		ID         int         `json:"id"`
		ExternalID interface{} `json:"external_id"`
		Removed    bool        `json:"removed"`
		ParentID   int         `json:"parent_id"`
	} `json:"parent"`
	Description  string `json:"description"`
	MembersCount int    `json:"members_count"`
	Head         int    `json:"head"`
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

func (d Directory) GetDepartments(orgID int, params Parameters) (*DirectoryDepartments, error) {
	var departments DirectoryDepartments
	err := Get(
		d.client,
		directoryURL+"/departments/",
		params,
		map[string]string{"X-Org-ID": strconv.Itoa(orgID)},
		&departments,
	)
	return &departments, err
}

func (d Directory) GetDepartment(orgID, depID int, params Parameters) (*DirectoryDepartment, error) {
	var department DirectoryDepartment
	err := Get(
		d.client,
		directoryURL+"/departments/"+strconv.Itoa(depID),
		params,
		map[string]string{"X-Org-ID": strconv.Itoa(orgID)},
		&department,
	)
	return &department, err
}

//      ________
//     /  _____/______  ____  __ ________  ______
//    /   \  __\_  __ \/  _ \|  |  \____ \/  ___/
//    \    \_\  \  | \(  <_> )  |  /  |_> >___ \
//     \______  /__|   \____/|____/|   __/____  >
//            \/                   |__|       \/

type DirectoryGroup struct {
	Name       string `json:"name"`
	Email      string `json:"email"`
	ExternalID string `json:"external_id"`
	ID         int    `json:"id"`
	Members    []struct {
		Type   string `json:"type"` // <user|group|department>
		Object struct {
			ID int `json:"id"`
		} `json:"object"`
	} `json:"members"`
	Label   string `json:"label"`
	Created string `json:"created"`
	Type    string `json:"type"`
	Admins  []struct {
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
	} `json:"admins"`
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
	} `json:"author"`
	Description  string `json:"description"`
	MembersCount int    `json:"members_count"`
	MemberOf     []int  `json:"member_of"`
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

func (d Directory) GetGroups(orgID int, params Parameters) (*DirectoryGroups, error) {
	var groups DirectoryGroups
	err := Get(
		d.client,
		directoryURL+"/groups/",
		params,
		map[string]string{"X-Org-ID": strconv.Itoa(orgID)},
		&groups,
	)
	return &groups, err
}

func (d Directory) GetGroup(orgID, groupID int, params Parameters) (*DirectoryGroup, error) {
	var group DirectoryGroup
	err := Get(
		d.client,
		directoryURL+"/groups/"+strconv.Itoa(groupID),
		params,
		map[string]string{"X-Org-ID": strconv.Itoa(orgID)},
		&group,
	)
	return &group, err
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
		map[string]string{"X-Org-ID": strconv.Itoa(orgID)},
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
func (d Directory) GetOrganizations(params Parameters) (*DirectoryOrganizations, error) {
	var organizations DirectoryOrganizations
	err := Get(
		d.client,
		directoryURL+"/organizations/",
		params,
		nil,
		&organizations,
	)
	return &organizations, err
}
