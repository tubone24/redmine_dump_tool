package redmine

import (
	"fmt"
	"github.com/goccy/go-json"
	"github.com/tubone24/redump/pkg/utils"
	"strconv"
)

var dat map[string]interface{}

type Project struct {
	Id int `json:"id"`
	Name string `json:"name"`
}

type Tracker struct {
	Id int `json:"id"`
	Name string `json:"name"`
}

type Status struct {
	Id int `json:"id"`
	Name string `json:"name"`
}

type Priority struct {
	Id int `json:"id"`
	Name string `json:"name"`
}

type Author struct {
	Id int `json:"id"`
	Name string `json:"name"`
}

type AssignedTo struct {
	Id int `json:"id"`
	Name string `json:"name"`
}

type CustomField struct {
	Id int `json:"id"`
	Name string `json:"name"`
	Multiple bool `json:"multiple"`
	Value interface{} `json:"value"`
}

type CustomFields []*CustomField

type Attachment struct {
	Id int `json:"id"`
	FileName string `json:"filename"`
	FileSize int64 `json:"filesize"`
	Description string `json:"description"`
	ContentUrl string `json:"content_url"`
	Author Author
	CreatedOn string `json:"created_on"`
}

type Attachments []*Attachment

type User struct {
	Id int `json:"id"`
	Name string `json:"name"`
}

type Detail struct {
	Property string `json:"property"`
	Name string `json:"name"`
	NewValue string `json:"new_value"`
	OldValue string `json:"old_value"`
}

type Details []*Detail

type Journal struct {
	Id int `json:"id"`
	User User `json:"user"`
	Notes string `json:"notes"`
	CreatedOn string `json:"created_on"`
	Details Details `json:"details"`
}

type Journals []*Journal

type Watcher struct {
	Id int
	Name string
}

type Watchers []*Watcher

type Issue struct{
	Id int `json:"id"`
	Project Project `json:"project"`
	Tracker Tracker `json:"tracker"`
	Status Status `json:"status"`
	Priority Priority `json:"priority"`
	Author Author `json:"author"`
	AssignedTo AssignedTo `json:"assigned_to"`
	Subject string `json:"subject"`
	Description string `json:"description"`
	StartDate string `json:"start_date"`
	DueDate string `json:"due_date"`
	DoneRatio int `json:"done_ratio"`
	CustomFields CustomFields `json:"custom_fields"`
	IsPrivate bool `json:"is_private"`
	EstimatedHours string `json:"estimated_hours"`
	CreatedOn string `json:"created_on"`
	UpdatedOn string `json:"updated_on"`
	ClosedOn string `json:"closed_on"`
	Attachments Attachments `json:"attachments"`
	Journals Journals `json:"journals"`
	Watchers Watchers `json:"watchers"`
}

type Issues []*Issue

type IssueParam struct {
	ProjectId     int          `json:"project_id,omitempty"`
	TrackerId     int          `json:"tracker_id,omitempty"`
	StatusId      int          `json:"status_id,omitempty"`
	PriorityId    int          `json:"priority_id,omitempty"`
	Subject       string       `json:"subject,omitempty"`
	Description   string       `json:"description,omitempty"`
	AssignedToId  int          `json:"assigned_to_id,omitempty"`
	ParentIssueId int          `json:"parent_issue_id,omitempty"`
	CustomFields  CustomFields `json:"custom_fields,omitempty"`
	Notes string `json:"notes,omitempty"`
}

type IssueParamJson struct {
	Issue IssueParam `json:"issue"`
}

var issuesResult struct {
	Issues Issues `json:"issues"`
	TotalCount int `json:"total_count"`
	Offset int `json:"offset"`
	Limit int `json:"limit"`
}

var issueResult struct {
	Issue Issue `json:"issue"`
}


func GetIssues(url, key string) (Issues, error){
	var issues Issues
	body, err := utils.Get(url + "/issues.json?key=" + key + "&limit=1&offset=0")
	if err != nil {
		return nil, err
	}
	err2 := json.Unmarshal(body, &issuesResult)
	if err2 != nil {
		return nil, err
	}
	fmt.Println(issuesResult.TotalCount)
	for offset := 0; offset < issuesResult.TotalCount; offset+=100 {
		body, err := utils.Get(url + "/issues.json?key=" + key + "&limit=100&offset=" + strconv.Itoa(offset))
		if err != nil {
			return nil, err
		}
		err2 := json.Unmarshal(body, &issuesResult)
		if err2 != nil {
			panic(err2)
		}
		issues = append(issues, issuesResult.Issues...)
	}
	return issues, nil
}

func GetIssue(url, key string, id int) (Issue, error) {
	var issue Issue
	body, err := utils.Get(url + "/issues/" + strconv.Itoa(id) + ".json?key=" + key + "&include=attachments,journals,watchers")
	if err != nil {
		return issue, err
	}
	err2 := json.Unmarshal(body, &issueResult)
	if err2 != nil {
		return issue, err
	}
	//fmt.Println(issueResult)
	return issueResult.Issue, nil
}

func CreateIssue(url, key string, issue IssueParam) (int, error) {
	issueJson, err := json.Marshal(IssueParamJson{Issue: issue})
	if err != nil {
		return 0, err
	}
	fmt.Println(string(issueJson))
	body, err := utils.Post(url + "/issues.json?key=" + key, "application/json", issueJson)
	if err != nil {
		return 0, err
	}
	err = json.Unmarshal(body, &issueResult)
	if err != nil {
		return 0, err
	}
	return issueResult.Issue.Id, nil
}

func UpdateIssueJournals(url, key string, id int, journals []string) error {
	issue := IssueParam{Notes: journals[0]}
	issueJson, err := json.Marshal(IssueParamJson{Issue: issue})
	if err != nil {
		return err
	}
	fmt.Println(string(issueJson))
	err = utils.Put(url + "/issues/" + strconv.Itoa(id) + ".json?key=" + key, "application/json", issueJson)
	if err != nil {
		return err
	}
	return nil
}
