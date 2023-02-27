package upwork

import (
	"context"
	"github.com/google/uuid"
	"github.com/mmcdole/gofeed"
	log "github.com/sirupsen/logrus"
	"html"
	"mgufrone.dev/job-alerts/internal/domain/job"
	worker2 "mgufrone.dev/job-alerts/pkg/worker"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

type Handler struct {
}

func (h *Handler) Fetch(ctx context.Context) ([]*job.Entity, error) {
	feedURL := "https://www.upwork.com/ab/feed/jobs/rss"
	fp := gofeed.NewParser()
	res, err := fp.ParseURLWithContext(feedURL, ctx)
	if err != nil {
		return nil, err
	}
	if len(res.Items) == 0 {
		return nil, nil
	}
	jbs := make([]*job.Entity, 0, len(res.Items))
	for _, feed := range res.Items {
		var (
			skills  []string
			country string
		)
		hourlyRegExp := regexp.MustCompile(`(?i)hourly range.*: (?P<range>(?P<currencyStart>.)(?P<hourlyStart>([\d.,]+))-(?P<currencyEnd>.)(?P<hourlyEnd>([\d.,]+)))`)
		//postDateExp := regexp.MustCompile(`(?i)posted on.*: (?P<postDate>([\w,: ]+))<br.*`)
		roleExp := regexp.MustCompile(`(?i)category.*: (?P<role>([\w- ]+))<br`)
		skillsExp := regexp.MustCompile(`(?i)skills.*: (?P<skills>.*)<br.*country`)
		countryExp := regexp.MustCompile(`(?i)country.*: (?P<country>.*)`)
		skillIdx := skillsExp.SubexpIndex("skills")
		countryIdx := countryExp.SubexpIndex("country")
		skillsRange := skillsExp.FindStringSubmatch(feed.Description)
		hourlyRange := hourlyRegExp.FindStringSubmatch(feed.Description)
		countryRange := countryExp.FindStringSubmatch(feed.Description)
		isFixed := false
		if len(hourlyRange) == 0 {
			hourlyRegExp = regexp.MustCompile(`(?i)budget.*: (?P<range>(?P<currencyStart>.)(?P<hourlyStart>([\d.,]+)))`)
			hourlyRange = hourlyRegExp.FindStringSubmatch(feed.Description)
			isFixed = true
		}
		if len(skillsRange) > 0 {
			sks := strings.Split(skillsRange[skillIdx], ",")
			for _, v := range sks {
				val, _ := url.QueryUnescape(v)
				skills = append(skills, strings.TrimSpace(html.UnescapeString(val)))
			}
		}
		if len(countryRange) > 0 {
			country = countryRange[countryIdx]
		}
		//rangeIdx := hourlyRegExp.SubexpIndex("range")
		startIdx := hourlyRegExp.SubexpIndex("hourlyStart")

		endIdx := hourlyRegExp.SubexpIndex("hourlyEnd")
		jobType := "unknown"
		if isFixed {
			endIdx = startIdx
			jobType = "fixed-price"
		} else {
			jobType = "hourly"
		}
		currencyIdx := hourlyRegExp.SubexpIndex("currencyStart")
		// post date
		//postDateRange := postDateExp.FindStringSubmatch(feed.Description)
		//postDateIdx := postDateExp.SubexpIndex("postDate")
		var (
			start, end = 0.01, 0.01
			currency   = "$"
		)
		if len(hourlyRange) > 0 {
			currency = hourlyRange[currencyIdx]
			start, _ = strconv.ParseFloat(strings.ReplaceAll(hourlyRange[startIdx], ",", ""), 10)
			end, _ = strconv.ParseFloat(strings.ReplaceAll(hourlyRange[endIdx], ",", ""), 10)
		}
		// role
		roleIdx := roleExp.SubexpIndex("role")
		roles := roleExp.FindStringSubmatch(feed.Description)

		//roleSplit := strings.Split(roleString, ":")
		//companyName, role := strings.TrimSpace(roleSplit[0]), strings.TrimSpace(roleSplit[1])
		companyName := "-"
		var role string
		if len(roles) > 0 {
			role = roles[roleIdx]
		}
		companyURL := "-"
		jb, err1 := job.NewJob(
			uuid.New(),
			role,
			companyName, companyURL,
			feed.Description, feed.Link, WorkerName, country, []string{})
		if err1 != nil {
			log.Error("cannot parse correctly", err1)
			continue
		}
		jb.SetJobType(jobType)
		jb.SetCreatedAt(*feed.PublishedParsed)
		jb.SetSalary([]float64{start, end})
		jb.SetSalaryCurrency(currency)
		jb.SetTags(skills)
		jbs = append(jbs, jb)
	}
	return jbs, nil
}

func (h *Handler) FetchJob(ctx context.Context, job2 *job.Entity) (*job.Entity, error) {
	return nil, nil
	//req, err := http.NewRequestWithContext(ctx, "GET", job2.JobURL(), nil)
	//if err != nil {
	//	return nil, err
	//}
	//content, err := h.agent.Do(req)
	//if err != nil {
	//	return nil, err
	//}
	//if content.StatusCode != http.StatusOK {
	//	return nil, fmt.Errorf("something wrong when retrieving data: %v", job2.JobURL())
	//}
	//defer content.Body.Close()
	//doc, _ := goquery.NewDocumentFromReader(content.Body)
	//doc.Find(".company-card").Each(func(i int, selection *goquery.Selection) {
	//	job2.SetLocation(strings.TrimSpace(strings.ReplaceAll(selection.Find("h3:nth-of-type(1)").Text(), "–", "")))
	//})
	//var tags []string
	//isRemote := false
	//doc.Find(".listing-header > .listing-header-container > a").Each(func(i int, selection *goquery.Selection) {
	//	tag := strings.TrimSpace(selection.Text())
	//	if strings.Contains(tag, "View") {
	//		return
	//	}
	//	if tag != "" {
	//		if strings.Contains(strings.ToLower(tag), "remote") {
	//			isRemote = true
	//			return
	//		}
	//		if strings.Contains(strings.ToLower(tag), strings.ToLower(job2.JobType())) {
	//			return
	//		}
	//		tags = append(tags, tag)
	//	}
	//})
	//err = job2.SetTags(tags)
	//if err != nil {
	//	return nil, err
	//}
	//job2.SetIsRemote(isRemote)
	//return job2, nil
}

func NewHandler() worker2.IWorker {
	return &Handler{}
}
