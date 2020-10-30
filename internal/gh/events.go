package gh

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/google/go-github/v32/github"
	"golang.org/x/oauth2"

	"github.com/helix-collective/hxgh/internal/types"
)

type eventsCmd struct {
	rt *types.Root
}

type events2CsvCmd struct {
	rt                  *types.Root
	TimezoneHoursOffset int64
	Username            string `opts:"mode=arg"`
}

// NewEventsCmd events sub command
func NewEventsCmd(rt *types.Root) interface{} {
	cmd := eventsCmd{rt: rt}
	return &cmd
}

// NewCsvCmd outputs events as a csv sub command
func NewCsvCmd(rt *types.Root) interface{} {
	cmd := events2CsvCmd{rt: rt}
	return &cmd
}

type eventData struct {
	When    time.Time
	What    string
	Where   string
	Payload map[string]interface{}
}

func eventDataString(timezoneOffset time.Duration, ev eventData) string {
	dt := ev.When.Add(timezoneOffset)
	date := dt.Format("2006-01-02")
	time := dt.Format("15:04:05")
	return fmt.Sprintf("%s,%s,%s,%s", date, time, ev.What, ev.Where)
}
func (ev eventData) StringPayload(cols []string) string {
	b := &strings.Builder{}
	for _, c := range cols {
		if h, ex := ev.Payload[c]; ex {
			str := fmt.Sprintf("%v", h)
			str = strings.ReplaceAll(str, "\r\n", "")
			str = strings.ReplaceAll(str, "\n\r", "")
			str = strings.ReplaceAll(str, "\r", "")
			str = strings.ReplaceAll(str, "\n", "")
			str = strings.ReplaceAll(str, ",", " ")
			b.WriteString(str)
		}
		b.WriteString(",")
	}
	return b.String()
}

func (cmd *events2CsvCmd) Run() error {
	if cmd.rt.GhToken == "" {
		return fmt.Errorf(`GhToken string (env GITHUB_TOKEN) must be set`)
	}
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: cmd.rt.GhToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	data := make([]eventData, 0)
	page := 0
	header := make(map[string]int)
	for {
		fmt.Fprintf(os.Stderr, "page %v\n", page)
		events, resp, err := client.Activity.ListEventsPerformedByUser(ctx, cmd.Username, false, &github.ListOptions{Page: page, PerPage: 100})
		if err != nil {
			fmt.Fprintf(os.Stderr, "error on list events %v\n", err)
			break
		}
		for _, event := range events {
			ed := eventData{
				When:    event.GetCreatedAt(),
				What:    event.GetType(),
				Where:   event.GetRepo().GetName(),
				Payload: make(map[string]interface{}),
			}
			err = json.Unmarshal(event.GetRawPayload(), &ed.Payload)
			if err != nil {
				fmt.Fprintf(os.Stderr, "err %v %v\n", err, event.GetRawPayload())
			}
			for k := range ed.Payload {
				if _, ex := header[k]; !ex {
					header[k] = len(header)
				}
			}
			data = append(data, ed)
		}
		if page == resp.LastPage {
			fmt.Fprintf(os.Stderr, "last page\n")
			break
		}
		if resp.NextPage == 0 {
			fmt.Fprintf(os.Stderr, "next page is 0\n")
			break
		}
		page = resp.NextPage
	}
	cols := []string{}
	fmt.Printf("Date,Time,What,Where,")
	for k := range header {
		cols = append(cols, k)
	}
	sort.Strings(cols)
	for _, c := range cols {
		fmt.Printf("%s,", c)
	}
	fmt.Printf("\n")
	timezoneOffset := time.Hour * time.Duration(cmd.TimezoneHoursOffset)
	for _, ed := range data {
		fmt.Printf("%s,%s\n", eventDataString(timezoneOffset, ed), ed.StringPayload(cols))
	}
	return nil
}
