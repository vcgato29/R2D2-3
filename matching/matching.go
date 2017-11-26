package matching

import(
	"regexp"
	"strconv"
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
	}

	for _, value := range regexes {
		regex := regexp.MustCompile(value)
		
		if regex.MatchString(name) {
			matches := regex.FindStringSubmatch(name)
			groups := getMatchGroups(matches, regex)

			intYear, intEpisode, intSeason := 0, 0, 0
			
			intYear, errYear := strconv.Atoi(groups["year"])
			intEpisode, errEpisode := strconv.Atoi(groups["number"])
			intSeason, errSeason := strconv.Atoi(groups["season"])

			if errYear != nil || errEpisode != nil || errSeason != nil {
				//fmt.Println(fmt.Sprintf("[ERR] Year: %s, Episode: %s, Season: %s", errYear, errEpisode, errSeason))
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