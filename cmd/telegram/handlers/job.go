package handlers

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
	"mgufrone.dev/job-alerts/internal/domain/job"
	job2 "mgufrone.dev/job-alerts/internal/usecases/job"
	"mgufrone.dev/job-alerts/pkg/payload"

	"regexp"
	"strings"
)

var (
	recognizablePattern = regexp.MustCompile(`(list|get|latest|search\s(?P<keyword>[\w\s-]+)).*jobs?(\s(of|about)\s(?P<skills>[\w,\s-]+))?(\s(in|at)\s(?P<location>[\w,\s]+))?`)
	searchPattern       = regexp.MustCompile(`search (?P<keyword>[\w-]+) jobs`)
)

type jobHandler struct {
	uc     *job2.UseCase
	logger logrus.FieldLogger
	bot    *tgbotapi.BotAPI
}

func newJobHandler(query *job2.UseCase, logger logrus.FieldLogger, bot *tgbotapi.BotAPI) *jobHandler {
	return &jobHandler{uc: query, logger: logger, bot: bot}
}

func (j *jobHandler) Name() string {
	return "job_handler"
}

func (j *jobHandler) Run(ctx context.Context, py payload.Payload) error {
	var msg tgbotapi.Message
	err := py.As("", &msg)
	if err != nil {
		return err
	}
	matches := recognizablePattern.FindStringSubmatch(msg.Text)
	skillsIdx := recognizablePattern.SubexpIndex("skills")
	keywordIdx := recognizablePattern.SubexpIndex("keyword")
	var (
		skills  []string
		keyword string
	)
	if len(matches) > 0 {
		if matches[skillsIdx] != "" {
			skills = strings.Split(matches[skillsIdx], ",")
			for k, v := range skills {
				skills[k] = strings.TrimSpace(v)
			}
		}
		if matches[keywordIdx] != "" {
			keyword = strings.TrimSpace(matches[keywordIdx])
		}
	}
	jbs, total, err := j.uc.List(ctx, job2.ListInput{
		Skills:     skills,
		Keyword:    strings.ToLower(keyword),
		User:       nil, // can be forced to list recommendation from user
		Pagination: job2.Pagination{Page: 1, PerPage: 3},
		Query: job2.Query{
			Fields: []string{"description", "id", "role", "updated_at", "tags", "company_name"},
		},
	})
	if err != nil {
		return err
	}
	text := "Here are some recommendations for you"
	if total == 0 {
		text = "I can't recommendation any job at the moment for you. Try again later."
	}
	reply := tgbotapi.NewMessage(msg.Chat.ID, text)
	reply.ParseMode = "markdown"
	j.bot.Send(reply)
	for _, jb := range jbs {
		jbMsg := tgbotapi.NewMessage(msg.Chat.ID, j.renderJob(jb))
		jbMsg.ParseMode = "markdown"
		j.bot.Send(jbMsg)
	}
	return nil
}

func (j *jobHandler) renderJob(entity *job.Entity) string {
	infoTag := []string{
		fmt.Sprintf("%s", entity.UpdatedAt().Format("02 Jan 2006")),
	}
	if len(entity.Salary()) > 0 && entity.SalaryCurrency() != "" {
		baseSalary := entity.Salary()[0]
		endSalary := baseSalary
		if len(entity.Salary()) > 1 {
			endSalary = entity.Salary()[1]
		}
		currency := entity.SalaryCurrency()
		salary := fmt.Sprintf("%s%.f - %s%.f", currency, baseSalary, currency, endSalary)
		if baseSalary == endSalary {
			salary = fmt.Sprintf("%s%.f", currency, baseSalary)
		}
		infoTag = append(infoTag, salary)
	} else {
		infoTag = append(infoTag, "no salary info")
	}
	if entity.Location() != "" {
		infoTag = append(infoTag, entity.Location())
	}
	if entity.IsRemote() {
		infoTag = append(infoTag, "remote work")
	}
	tags := entity.Tags()
	return fmt.Sprintf(fmt.Sprintf(`[%s](%s)
%s
%s
`, entity.Role()+" at "+entity.CompanyName(), entity.JobURL(), strings.Join(infoTag, " | "), "skills/tags: "+strings.Join(tags, ", ")))
}

func (j *jobHandler) When(py payload.Payload) bool {
	var msg tgbotapi.Message
	err := py.As("", &msg)
	return err == nil && (recognizablePattern.MatchString(msg.Text) ||
		searchPattern.MatchString(msg.Text))
}

func (j *jobHandler) SubscribeAt() []string {
	return []string{"chat"}
}

func (j *jobHandler) freelingParse() []string {
	return []string{}
	//posTagger := nlp.NewPosTagger()
	//
	//// Tokenize a sentence into words
	//sentence := "The quick brown fox jumps over the lazy dog."
	//words := nlp.WordTokenizer(sentence)
	//
	//// Tag each word with its POS
	//posTags := posTagger.Tag(words)
	//
	//// Print the tagged words
	//for i, word := range words {
	//	fmt.Printf("%s -> %s\n", word, posTags[i])
	//}
}
