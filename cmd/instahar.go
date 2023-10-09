package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	models "github.com/divya-rajput/instahar/models"
)

const target_request = "https://www.instagram.com/api/v1/clips/user/"

var (
	removeMPD     = regexp.MustCompile(`<MPD[^>]*>.*?</MPD>`)
	xmlRemover    = regexp.MustCompile(`\<.*?\>`)
	hashtagRegexp = regexp.MustCompile(`(^|\s)#([A-Za-z_][A-Za-z0-9_]*)`)
	IndiaTz, _    = time.LoadLocation("Asia/Kolkata")
)

func processHarEntry(entry models.HarEntry) []models.InstagramMediaEntry {
	entry.Response.Content.Text = strings.ReplaceAll(entry.Response.Content.Text, `\n`, ``)
	entry.Response.Content.Text = strings.ReplaceAll(entry.Response.Content.Text, `\u003c`, `<`)
	entry.Response.Content.Text = strings.ReplaceAll(entry.Response.Content.Text, `\u003e`, `>`)
	entry.Response.Content.Text = strings.ReplaceAll(entry.Response.Content.Text, `\u0026`, `&`)
	entry.Response.Content.Text = removeMPD.ReplaceAllString(entry.Response.Content.Text, ``)
	entry.Response.Content.Text = xmlRemover.ReplaceAllString(entry.Response.Content.Text, ``)

	os.WriteFile("cleanedpost.json", []byte(entry.Response.Content.Text), 0755)

	instagramMedia := models.InstagramGraphQLQueryResponse{}
	if err := json.Unmarshal([]byte(entry.Response.Content.Text), &instagramMedia); err != nil {
		fmt.Println("Error while processing the har entry", err.Error())
	}

	return instagramMedia.Items
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage:\ninstahar <filename>")
		os.Exit(0)
	}

	filename := os.Args[1]

	if matched, _ := regexp.MatchString(`^.*\.(har|HAR)$`, filename); !matched {
		fmt.Println("Please provide the HAR file as the argument. The filename must end with \"har\"")
		os.Exit(0)
	}

	harFileRawData, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}

	harFileData := models.HarFileData{}
	if err := json.Unmarshal(harFileRawData, &harFileData); err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}

	output := []models.InstagramMediaEntry{}

	for _, entry := range harFileData.Log.Entries {
		if entry.Request.Url != target_request {
			continue
		}
		output = append(output, processHarEntry(entry)...)
	}

	csvFile, err := os.Create("output.csv")
	if err != nil {
		log.Printf("failed creating file: %s\n", err)
	}
	defer csvFile.Close()

	fmt.Println("tags,like_count,play_count,comment_count,IST Time,Local Time,UTC Time")
	csvFile.WriteString("tags,like_count,play_count,comment_count,IST Time,Local Time,UTC Time\n")
	for _, entry := range output {
		tags := strings.Join(hashtagRegexp.FindAllString(entry.Media.Caption.Text, -1), " ")
		timestamp := time.Unix(int64(entry.Media.TakenAt), 0)
		csvEntry := strings.Join([]string{tags, strconv.FormatInt(entry.Media.LikeCount, 10), strconv.FormatInt(entry.Media.PlayCount, 10), strconv.FormatInt(entry.Media.CommentCount, 10), timestamp.In(IndiaTz).Format(time.RFC822), timestamp.Format(time.RFC822), timestamp.UTC().Format(time.RFC822)}, ",")
		fmt.Println(csvEntry)
		csvFile.WriteString(csvEntry + "\n")
	}
}
