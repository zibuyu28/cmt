package main

import (
	"context"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"os"
	"os/exec"
	"reflect"
	"regexp"
	"strings"
	"text/template"
	"time"
)

var typeMap = map[string]string{
	"feat":     "A new feature",
	"fix":      "A bug fix",
	"docs":     "Documentation only changes",
	"style":    "Changes that do not affect the meaning of the code (white-space, formatting, missing semi-colons, etc)",
	"refactor": "A code change that neither fixes a bug nor adds a feature",
	"perf":     "A code change that improves performance",
	"test":     "Adding missing tests",
	"chore":    "Changes to the build process or auxiliary tools and libraries such as documentation generation",
	"revert":   "Revert to a commit",
	"WIP":      "Work in progress",
}

var typeMapZh = map[string]string{
	"feat":     "新功能",
	"fix":      "修复bug",
	"docs":     "修改文档",
	"style":    "改动的代码不影响代码含义（空格、格式和换行等）",
	"refactor": "既不修复错误也不添加功能的代码更改",
	"perf":     "提高性能的代码更改",
	"test":     "添加单测",
	"chore":    "对构建过程或辅助工具和库（例如文档生成）的更改",
	"revert":   "恢复至某次提交",
	"WIP":      "工作正在进行中",
}

type Answer struct {
	Type    string `survey:"type"`
	Scope   string `survey:"scope"`
	Message string `survey:"message"`
	Version string `survey:"version"`
	TBNum   string `survey:"tb_num"`
}

func main() {
	// the answers will be written to this struct
	ans := &Answer{}
	// get last commit message
	lc, _ := lastCommit()
	defaultConstruct(ans, lc)
	qs := []*survey.Question{tpq, scopeq, messageq, versionq, tbnumq}
	err := survey.Ask(qs, ans)
	if err != nil {
		return
	}
	commit(ans)
}

type simpleWriter struct {
	cmdReturn string
}

func (s *simpleWriter) Write(p []byte) (n int, err error) {
	s.cmdReturn = s.cmdReturn + string(p)
	return len(p), nil
}

const (
	commitTemplate = "{{.Type}}({{.Scope}}): {{.Message}}({{.Version}}-{{.TBNum}})"
)

func commit(ans *Answer) {
	parse, err := template.New("commit").Parse(commitTemplate)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	writer := &simpleWriter{}
	err = parse.Execute(writer, ans)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	commitmsg := writer.cmdReturn
	fmt.Printf("commit message : %s\n", commitmsg)
	commandContext := exec.Command("git", "commit", "-m", fmt.Sprintf("%s", commitmsg))
	commandContext.Stdout = os.Stdout
	commandContext.Stderr = os.Stdout
	err = commandContext.Run()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

var commitReg = regexp.MustCompile("^([^(]*)\\(([^)]*)\\):\\s([^(]*)\\(([v\\.0-9]{6,15})-([^)]*)\\)$")

func defaultConstruct(ans *Answer, lastCommit string) {
	submatch := commitReg.FindAllStringSubmatch(lastCommit, -1)
	if len(submatch) == 1 && len(submatch[0]) == 6 {
		ans.Type = submatch[0][1]
		tpq.Prompt.(*survey.Select).Default = ans.Type
		ans.Scope = submatch[0][2]
		scopeq.Prompt.(*survey.Input).Default = ans.Scope
		ans.Message = submatch[0][3]
		messageq.Prompt.(*survey.Input).Default = ans.Message
		ans.Version = submatch[0][4]
		versionq.Prompt.(*survey.Input).Default = ans.Version
		ans.TBNum = submatch[0][5]
		tbnumq.Prompt.(*survey.Input).Default = ans.TBNum
	}
}

func lastCommit() (string, error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*10)
	defer cancelFunc()
	sw := simpleWriter{}
	commandContext := exec.CommandContext(ctx, "git", "log", "--pretty=format:\"%s\"", "--no-merges", "-1")
	commandContext.Stdout = &sw
	err := commandContext.Run()
	if err != nil {
		return "", err
	}
	resp := strings.TrimSpace(sw.cmdReturn)
	resp = strings.TrimSuffix(resp, "\n")
	resp = strings.TrimPrefix(resp, "\"")
	resp = strings.TrimSuffix(resp, "\"")
	return resp, nil
}

var zh string
var tpq = &survey.Question{
	Name: "type",
	Prompt: &survey.Select{
		Message: "Select the type of change you're committing",
		//Options: []string{"feat", "fix", "docs", "style", "refactor", "perf", "test", "chore", "revert", "WIP"},
		Options: []string{"feat", "fix", "docs", "style", "perf", "test", "chore"},
		Description: func(value string, index int) string {
			if zh != "" {
				return typeMapZh[value]
			}
			return typeMap[value]
		},
	},
	Validate: survey.Required,
}

var scopeReg = regexp.MustCompile("^[a-z0-9]{2,20}$")
var scopeq = &survey.Question{
	Name: "scope",
	Prompt: &survey.Input{
		Message: "Input the scope of impact of your current commit",
		Help:    "one word to describe the scope of influence, composed of lowercase letters and numbers, the length range is 2-20 characters",
	},
	Validate: func(ans interface{}) error {
		if ans == nil {
			return fmt.Errorf("input is empty")
		}
		if str, ok := ans.(string); ok {
			if ans == "" {
				return fmt.Errorf("input is empty")
			}
			if !scopeReg.MatchString(str) {
				return fmt.Errorf("this input [%s] not in correct format, validate regular '[a-z0-9]{2,20}'", str)
			}
			// the input is fine
			return nil
		} else {
			// otherwise we cannot convert the value into a string and cannot enforce length
			return fmt.Errorf("cannot enforce length on response of type %v", reflect.TypeOf(ans).Name())
		}
	},
}

var messageq = &survey.Question{
	Name: "message",
	Prompt: &survey.Input{
		Message: "Input message of your current commit",
	},
	Validate: survey.Required,
}

var versionq = &survey.Question{
	Name: "version",
	Prompt: &survey.Input{
		Message: "Input repository version",
	},
}

var tbnumq = &survey.Question{
	Name: "tb_num",
	Prompt: &survey.Input{
		Message: "Input teambition number",
	},
}
