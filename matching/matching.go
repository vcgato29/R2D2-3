package matching

import "regexp"

// Video content
type Content struct {
	Name string
	Year string
	Film bool
	Tv bool
	Season string
	Number string
}

func MatchContent(name string) Content {
	regexes := map[string]string {
		"TNSE": `(?P<title>[A-z]+) - (?P<season>[0-9]+)x(?P<number>[0-9]+) - (?P<episode>([A-z](\s)*|[0-9](\s)*)+)`,
	}

	for _, value := range regexes {
		regex := regexp.MustCompile(value)
		
		if regex.MatchString(name) {
			matches := regex.FindStringSubmatch(name)
			groups := getMatchGroups(matches, regex)

			return Content{groups["title"], groups["year"], false, false, groups["season"], groups["number"]}
		}
	}

	return Content{"NA", "NA", false, false, "NA", "NA"}	
}

func getMatchGroups(matches []string, exp *regexp.Regexp) map[string]string {
	result := make(map[string]string)

	for i, name := range exp.SubexpNames() {
		if i != 0 { result[name] = matches[i] }
	}
	
	return result
}