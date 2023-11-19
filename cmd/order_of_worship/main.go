package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
)

func main() {
	skipTopLinesRegex := false
	lineOneRegex := regexp.MustCompile(`Providence United Reformed Church in Strathroy`)
	lineTwoRegex := regexp.MustCompile(`I was glad when they said to me`)
	lordsDayRegex := regexp.MustCompile(`Lord’s Day [0-9]+ – [a-zA-Z]+ [0-9]+, [0-9]+`)

	serviceNineThirty := regexp.MustCompile(`9:30 AM`)
	serviceThree := regexp.MustCompile(`3:00 PM`)

	inList := false

	silentPrayerRegex := regexp.MustCompile(`^*Silent Prayer [^\s0-9].*`)
	scriptureReadingRegex := regexp.MustCompile(`^Scripture Reading ([^\s0-9].*)`)
	textReadingRegex := regexp.MustCompile(`^Text ([^\s0-9].*)`)
	sermonRegex := regexp.MustCompile(`^Sermon ([^\s0-9].*)`)
	psalmsRegex := regexp.MustCompile(`^Reading through the Psalms ([^\s0-9].*)`)
	whatWeBelieveRegex := regexp.MustCompile(`^What We Believe ([’a-zA-Z0-9\.()\s]+)`)

	songListItemRegex := regexp.MustCompile(`^([^\s0-9].*) (TPH [0-9A-Z,:\s]+)`)
	listItemRegex := regexp.MustCompile(`^[^\s0-9].*`)
	listItemInnerLinesRegex := regexp.MustCompile(`^[\s0-9].*`)

	nextSundayRegex := regexp.MustCompile(`Lord willing, we will gather next Sunday`)
	nextSundayAMServiceRegex := regexp.MustCompile(`AM [S|s]ervice: `)
	nextSundayPMServiceRegex := regexp.MustCompile(`PM [S|s]ervice: `)
	nextOfferingRegex := regexp.MustCompile(`Next Lord’s Day offering`)

	readFile, err := os.Open("in.txt")
	if err != nil {
		fmt.Println(err)
	}
	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)

	var outFileLines []string
	outFileLines = append(outFileLines,
		`<h2><img class="aligncenter wp-image-444 size-full" src="//wp-content/uploads/2019/06/1H7A1311_01-e1590750591648.png" alt="" width="3960" height="1633"></h2>`)

	for fileScanner.Scan() {
		line := fileScanner.Text()

		if !skipTopLinesRegex {
			if m := lineOneRegex.MatchString(line); m {
				continue
			}

			if m := lineTwoRegex.MatchString(line); m {
				continue
			}

			if m := lordsDayRegex.MatchString(line); m {
				outFileLines = append(outFileLines, fmt.Sprintf(`<h2 style="text-align: center;">%s</h2>`, line))
				outFileLines = append(outFileLines,
					`<blockquote><p><span class="indent-1"><span class="text Lam-3-25">I was glad when they said to me, "Let us go to the house of the Lord!" - <a href="https://www.biblegateway.com/passage/?search=Psalm+122%3A1&amp;version=ESV">Psalm 122:1</a></span></span></p></blockquote>`)
				skipTopLinesRegex = true
				continue
			}
		}

		if m := serviceNineThirty.MatchString(line); m {
			outFileLines = append(outFileLines, `<p style="text-align: center;"><strong>9:30 AM Service<em><br></em></strong></p>`)
			continue
		}

		if m := serviceThree.MatchString(line); m {
			outFileLines = append(outFileLines, `</li></ul><hr /><p style="text-align: center;"><strong>3:00 PM Service<em><br></em></strong></p>`)
			inList = false
			continue
		}

		if m := nextSundayRegex.MatchString(line); m {
			outFileLines = append(outFileLines, fmt.Sprintf(`</li></ul>%s`, line))
			continue
		}

		if m := nextSundayAMServiceRegex.MatchString(line); m {
			outFileLines = append(outFileLines, fmt.Sprintf(`<ul><li>%s</li>`, line))
			continue
		}

		if m := nextSundayPMServiceRegex.MatchString(line); m {
			outFileLines = append(outFileLines, fmt.Sprintf(`<li>%s</li></ul>`, line))
			continue
		}

		if m := nextOfferingRegex.MatchString(line); m {
			outFileLines = append(outFileLines, line)
			continue
		}

		if m := silentPrayerRegex.MatchString(line); m {
			if inList {
				outFileLines = append(outFileLines, fmt.Sprintf("</li><li>%s", line))
			} else {
				outFileLines = append(outFileLines, fmt.Sprintf("<ul><li>%s", line))
				inList = true
			}
			continue
		}

		if m := whatWeBelieveRegex.FindStringSubmatch(line); len(m) > 0 {
			if inList {
				outFileLines = append(outFileLines, fmt.Sprintf("</li><li>What We Believe: <b>%s</b>", m[1]))
			} else {
				outFileLines = append(outFileLines, fmt.Sprintf("<ul><li>What We Believe: <b>%s</b>", m[1]))
				inList = true
			}
			continue
		}

		if m := scriptureReadingRegex.FindStringSubmatch(line); len(m) > 0 {
			if inList {
				outFileLines = append(outFileLines, fmt.Sprintf("</li><li>Scripture Reading: <b>%s</b>", m[1]))
			} else {
				outFileLines = append(outFileLines, fmt.Sprintf("<ul><li>Scripture Reading: <b>%s</b>", m[1]))
				inList = true
			}
			continue
		}

		if m := textReadingRegex.FindStringSubmatch(line); len(m) > 0 {
			if inList {
				outFileLines = append(outFileLines, fmt.Sprintf("</li><li>Text: <b>%s</b>", m[1]))
			} else {
				outFileLines = append(outFileLines, fmt.Sprintf("<ul><li>Text: <b>%s</b>", m[1]))
				inList = true
			}
			continue
		}

		if m := sermonRegex.FindStringSubmatch(line); len(m) > 0 {
			if inList {
				outFileLines = append(outFileLines, fmt.Sprintf("</li><li>Sermon: <b>%s</b>", m[1]))
			} else {
				outFileLines = append(outFileLines, fmt.Sprintf("<ul><li>Sermon: <b>%s</b>", m[1]))
				inList = true
			}
			continue
		}

		if m := psalmsRegex.FindStringSubmatch(line); len(m) > 0 {
			if inList {
				outFileLines = append(outFileLines, fmt.Sprintf("</li><li>Reading through the Psalms: <b>%s</b>", m[1]))
			} else {
				outFileLines = append(outFileLines, fmt.Sprintf("<ul><li>Reading through the Psalms: <b>%s</b>", m[1]))
				inList = true
			}
			continue
		}

		if m := songListItemRegex.FindStringSubmatch(line); len(m) > 0 {
			if inList {
				outFileLines = append(outFileLines, fmt.Sprintf("</li><li>%s: <b>%s</b>", m[1], m[2]))
			} else {
				outFileLines = append(outFileLines, fmt.Sprintf("<ul><li>%s: <b>%s</b>", m[1], m[2]))
				inList = true
			}
			continue
		}

		if m := listItemRegex.MatchString(line); m {
			if inList {
				outFileLines = append(outFileLines, fmt.Sprintf("</li><li>%s", line))
			} else {
				outFileLines = append(outFileLines, fmt.Sprintf("<ul><li>%s", line))
				inList = true
			}
			continue
		}

		if m := listItemInnerLinesRegex.MatchString(line); m {
			outFileLines = append(outFileLines, fmt.Sprintf("<br /> %s", line))
			continue
		}
	}

	readFile.Close()

	file, err := os.OpenFile("out.html", os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}

	datawriter := bufio.NewWriter(file)

	for _, data := range outFileLines {
		_, _ = datawriter.WriteString(data)
		//_, _ = datawriter.WriteString("\n")
	}

	datawriter.Flush()
	file.Close()
}
