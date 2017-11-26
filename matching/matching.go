package matching

import(
	"regexp"
	"strconv"
	"fmt"
)

// Video content
type Content struct {
	Name string
	Episode string
	Year int
	Film bool
	Tv bool
	Season int
	Number int
}

func MatchContent(name string) Content {
	regexes := map[string]string {
		"TNSE": `(?P<title>[A-z]+) - (?P<season>[0-9]+)x(?P<number>[0-9]+) - (?P<episode>([A-z](\s)*|[0-9](\s)*)+)`,
		"E1": `(?P<title>([A-z]|[0-9]|\.)+)(\.|\s)S(?P<season>[0-9]{2})E(?P<number>[0-9]{2})`,
		"NSE": `(?P<title>([A-z]|[0-9]|\s|\.)+) E(?P<number>[0-9]{1,2}) - (?P<episode>([A-z]|[0-9]|\s|\.)+)`,
	}

	for key, value := range regexes {
		regex := regexp.MustCompile(value)
		
		if regex.MatchString(name) {
			fmt.Println(fmt.Sprintf("[r2d2] Matched regex rule %s", key))
			matches := regex.FindStringSubmatch(name)
			groups := getMatchGroups(matches, regex)

			fmt.Println(groups)
			
			intYear, intEpisode, intSeason := -1, -1, -1
			
			intYear, errYear := strconv.Atoi(groups["year"])
			intEpisode, errEpisode := strconv.Atoi(groups["number"])
			intSeason, errSeason := strconv.Atoi(groups["season"])

			if errYear != nil || errEpisode != nil || errSeason != nil {
				//fmt.Println(fmt.Sprintf("[ERR] Year: %s, Episode: %s, Season: %s", errYear, errEpisode, errSeason))
			}

			// No season captured, assume SE01 - sometimes SE0 is used for specials, which is why
			// we're using SE-1 as the default
			if intSeason == -1 {
				fmt.Println("[r2d2] No season matched, defaulting to 1")
				intSeason = 1
			}

			return Content{
				groups["title"], 
				groups["episode"], 
				intYear, 
				false, 
				true, 
				intSeason, 
				intEpisode,
			}
		}
	}

	return Content{"NA", "NA",0, false, false, 0, 0}	
}

func getMatchGroups(matches []string, exp *regexp.Regexp) map[string]string {
	result := make(map[string]string)

	for i, name := range exp.SubexpNames() {
		if i != 0 { result[name] = matches[i] }
	}
	
	return result
}