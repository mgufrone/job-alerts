package weworkremotely

import (
	"context"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/google/uuid"
	"github.com/mmcdole/gofeed"
	log "github.com/sirupsen/logrus"
	"mgufrone.dev/job-alerts/internal/domain/job"
	worker2 "mgufrone.dev/job-alerts/pkg/worker"
	"net/http"
	"strings"
)

type Handler struct {
}

func (h *Handler) Fetch(ctx context.Context) ([]*job.Entity, error) {
	feedURL := "https://weworkremotely.com/remote-jobs.rss"
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
		roleString := strings.Replace(feed.Title, "–", "-", -1)
		roleSplit := strings.Split(roleString, ":")
		companyName, role := strings.TrimSpace(roleSplit[0]), strings.TrimSpace(roleSplit[1])
		companyURL := fmt.Sprintf("%s/%s/%s", host, "company", strings.Replace(strings.ToLower(companyName), " ", "-", -1))
		jb, err1 := job.NewJob(
			uuid.New(),
			role,
			companyName, companyURL,
			feed.Description, feed.Link, WorkerName, "-", []string{})
		if err1 != nil {
			log.Error("cannot parse correctly", err)
			continue
		}
		jb.SetJobType(feed.Custom["type"])
		jb.SetCreatedAt(*feed.PublishedParsed)
		jbs = append(jbs, jb)
	}
	return jbs, nil
}

func (h *Handler) FetchJob(ctx context.Context, job2 *job.Entity) (*job.Entity, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", job2.JobURL(), nil)
	if err != nil {
		return nil, err
	}
	content, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	if content.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("something wrong when retrieving data: %v", job2.JobURL())
	}
	defer content.Body.Close()
	doc, _ := goquery.NewDocumentFromReader(content.Body)
	doc.Find(".company-card").Each(func(i int, selection *goquery.Selection) {
		job2.SetLocation(strings.TrimSpace(strings.ReplaceAll(selection.Find("h3:nth-of-type(1)").Text(), "–", "")))
	})
	var tags []string
	isRemote := false
	doc.Find(".listing-header > .listing-header-container > a").Each(func(i int, selection *goquery.Selection) {
		tag := strings.TrimSpace(selection.Text())
		if strings.Contains(tag, "View") {
			return
		}
		if tag != "" {
			if strings.Contains(strings.ToLower(tag), "remote") {
				isRemote = true
				return
			}
			if strings.Contains(strings.ToLower(tag), strings.ToLower(job2.JobType())) {
				return
			}
			tags = append(tags, tag)
		}
	})
	err = job2.SetTags(tags)
	if err != nil {
		return nil, err
	}
	job2.SetIsRemote(isRemote)
	return job2, nil
}

func NewHandler() worker2.IWorker {
	return &Handler{}
}
