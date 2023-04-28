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
	wrapper "mgufrone.dev/job-alerts/pkg/worker_wrapper"
	"net/http"
	"strings"
)

type Handler struct {
	client worker2.IHTTPClient
}

func (h *Handler) Fetch(ctx context.Context) ([]*job.Entity, error) {
	feedURL := "https://weworkremotely.com/remote-jobs.rss"
	fp := gofeed.NewParser()
	fp.Client = h.client.ToHTTPClient()
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
	content, err := h.client.Do(req)
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

	//fmt.Println("found tags", contentTags)
	doc.Find(".listing-header .listing-header-container > a").Each(func(i int, selection *goquery.Selection) {
		if selection.Find(".listing-tag").Length() <= 0 {
			return
		}
		tag := strings.TrimSpace(selection.Text())
		tag = strings.ToLower(tag)
		if strings.Contains(tag, "View") {
			return
		}
		if tag != "" {
			if strings.Contains(tag, "remote") {
				isRemote = true
				return
			}
			if strings.Contains(tag, strings.ToLower(job2.JobType())) {
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

func NewHandler(client worker2.IHTTPClient) worker2.IWorker {
	return &Handler{client: client}
}
func Default() worker2.IWorker {
	return NewHandler(wrapper.NewHTTPClient(http.DefaultClient))
}
