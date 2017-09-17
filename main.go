package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strings"
	"time"
)

func getPageLinks(link string, domain string, data map[string]bool) map[string]bool {
	// fmt.Println(link)
	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
	}
	client := &http.Client{Transport: tr}

	request, err := http.NewRequest("GET", link, nil)

	request.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.115 Safari/537.36")
	resp, err := client.Do(request)

	if err != nil {
		fmt.Println("Error get content")
	} else {
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusNotFound && resp.Header["Content-Type"][0][0:9] == "text/html" {
			fmt.Println(link)
			bodyBytes, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Println("Error read bytes")
			} else {
				bodyString := string(bodyBytes)
				//fmt.Println(bodyString)
				regQuery := `href="([^"]+)"`
				re := regexp.MustCompile(regQuery)
				cArray := re.FindAllString(bodyString, -1)
				if len(cArray) > 0 {
					for _, element := range cArray {
						rawlink := strings.ToLower(element[6 : len(element)-1])
						// fmt.Println(rawlink)
						if strings.Index(rawlink, "?") > -1 {
							rawlink = rawlink[0:strings.Index(rawlink, "?")]
						}
						if strings.Index(rawlink, "&") > -1 {
							rawlink = rawlink[0:strings.Index(rawlink, "&")]
						}
						if strings.Index(rawlink, "#") > -1 {
							rawlink = rawlink[0:strings.Index(rawlink, "#")]
						}
						// fmt.Println(rawlink)
						if rawlink != "" && rawlink != "#" && rawlink != "javascript:;" && rawlink != "javascript:void(0);" {
							if len(rawlink) < 2 || rawlink[0:2] != "//" {
								var cleanlink string
								if len(rawlink) > 4 && rawlink[0:4] == "http" {
									cleanlink = rawlink
								} else {
									if rawlink[0:1] == "/" {
										cleanlink = domain + rawlink[1:len(rawlink)]
									} else {
										cleanlink = domain + rawlink
									}
								}
								// fmt.Println(cleanlink)
								if len(cleanlink) >= len(domain) &&
									cleanlink[0:len(domain)] == domain &&
									cleanlink[len(cleanlink)-4:len(cleanlink)] != ".css" &&
									cleanlink[len(cleanlink)-3:len(cleanlink)] != ".js" &&
									cleanlink[len(cleanlink)-4:len(cleanlink)] != ".jpg" &&
									cleanlink[len(cleanlink)-5:len(cleanlink)] != ".jpeg" &&
									cleanlink[len(cleanlink)-4:len(cleanlink)] != ".png" &&
									cleanlink[len(cleanlink)-4:len(cleanlink)] != ".gif" &&
									cleanlink[len(cleanlink)-4:len(cleanlink)] != ".pdf" &&
									cleanlink[len(cleanlink)-4:len(cleanlink)] != ".zip" &&
									cleanlink[len(cleanlink)-4:len(cleanlink)] != ".doc" &&
									cleanlink[len(cleanlink)-4:len(cleanlink)] != ".rar" &&
									cleanlink[len(cleanlink)-4:len(cleanlink)] != ".xls" &&
									cleanlink[len(cleanlink)-5:len(cleanlink)] != ".docx" &&
									cleanlink[len(cleanlink)-5:len(cleanlink)] != ".xlsx" {
									_, ok := data[cleanlink]
									if ok == false {
										// fmt.Println(cleanlink)
										data[cleanlink] = true
										data = getPageLinks(cleanlink, domain, data)
									}
								}
							}
						}
					}
				}
			}
		} else {
			delete(data, link)
		}
	}
	return data
}

func main() {
	fmt.Println("Start", time.Now())

	if len(os.Args) > 1 {
		domain := os.Args[1]

		data := map[string]bool{domain: true}

		data = getPageLinks(domain, domain, data)

		fmt.Println("Finish", time.Now())

		sitemap := make([]string, len(data))
		i := 0
		for k, _ := range data {
			sitemap[i] = k
			i++
		}
		sort.Strings(sitemap)

		t := time.Now()
		t.Format("2006-01-02 15:04:05")

		f, _ := os.Create("sitemap." + t.Format("2006_01_02_15_04_05") + ".xml")
		f.WriteString(`<?xml version="1.0" encoding="UTF-8"?>` + "\n")
		f.WriteString(`<urlset xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:schemaLocation="http://www.sitemaps.org/schemas/sitemap/0.9 http://www.sitemaps.org/schemas/sitemap/0.9/sitemap.xsd" xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">` + "\n")
		for _, val := range sitemap {
			// fmt.Println(val)
			f.WriteString("\t<url>\n")
			f.WriteString("\t\t<loc>" + val + "</loc>\n")
			f.WriteString("\t\t<lastmod>" + t.Format("2006-01-02") + "</lastmod>\n")
			f.WriteString("\t\t<priority>0.5</priority>\n")
			f.WriteString("\t\t<changefreq>weekly</changefreq>\n")
			f.WriteString("\t</url>\n")
		}
		f.WriteString(`</urlset>`)
	}
}
