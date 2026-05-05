package transformer_utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"path"
	"reflect"
	"strconv"
	"strings"
	"unicode"

	"github.com/bbalet/stopwords"

	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/PuerkitoBio/goquery"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/parser"
	"github.com/k3a/html2text"
)

func ConvertHTMLToMarkdown(html string) (string, error) {
	converter := md.NewConverter("", true, nil)

	markdown, err := converter.ConvertString(html)
	if err != nil {
		return "", err
	}

	return markdown, nil
}

func LowercaseStringArray(arrayToUse []string) []string {
	var arrayToReturn []string
	for _, stringToUse := range arrayToUse {
		arrayToReturn = append(arrayToReturn, strings.ToLower(stringToUse))
	}
	return arrayToReturn
}
func StringArrayContainsStrings(arrayToCheck []string, valuesToCheckFor []string) bool {

	arrayToCheck = DedupArray(arrayToCheck)
	valuesToCheckFor = DedupArray(valuesToCheckFor)

	set := make(map[string]int)
	for _, value := range arrayToCheck {
		set[value] += 1
	}

	for _, value := range valuesToCheckFor {
		if count, found := set[value]; !found {
			return false
		} else if count < 1 {
			return false
		} else {
			set[value] = count - 1
		}
	}
	return true
}

func ConvertMarkdownToHTML(markdownString string) string {
	md := []byte(markdownString)
	extensions := parser.CommonExtensions | parser.HardLineBreak
	parser := parser.NewWithExtensions(extensions)
	output := markdown.ToHTML(md, parser, nil)
	return string(output)
}

func ConvertMarkdownToPlainText(markdownString string) string {
	return ConvertHTMLToPlainText(ConvertMarkdownToHTML(markdownString))
}

// Abs returns the absolute value of x.
func AbsoluteValue(x float32) float32 {
	if x < 0 {
		return -x
	}
	return x
}

func TruncateString(text string, max int) string {
	lastSpaceIx := -1
	length := 0
	for i, r := range text {
		if unicode.IsSpace(r) {
			lastSpaceIx = i
		}
		length++
		if length >= max {
			if lastSpaceIx != -1 {
				return text[:lastSpaceIx] + "..."
			}
		}
	}
	if len(text) < max {
		return text
	}
	return text[:max]
}

func SelectFirstCharacters(text string, charCount int) string {
	i := 0
	for j := range text {
		if i == charCount {
			return text[:j]
		}
		i++
	}
	return text
}

func ConvertHTMLToPlainText(html string) string {
	return html2text.HTML2Text(html)
}

func RemoveSpaces(value string) string {
	// if strings.Contains(value, "\t\t") {
	// 	return RemoveSpaces(strings.ReplaceAll(value, "\t\t", ""))
	// }
	return strings.Join(strings.Fields(strings.TrimSpace(value)), " ")
	// return value
}

func ExtractWordsFromString(query string) []string {
	querySplitted := strings.Split(query, " ")
	queryKeywords := []string{}

	for _, keyword := range querySplitted {
		if keyword != "" {
			queryKeywords = append(queryKeywords, keyword)
		}
	}

	return DedupArray(queryKeywords)
}

func DedupArray(arrayToDedup []string) []string {
	mapForElements := make(map[string]int)
	for _, value := range arrayToDedup {
		mapForElements[value] = 1
	}
	dedupedArray := []string{}
	for key := range mapForElements {
		dedupedArray = append(dedupedArray, key)
	}
	return dedupedArray
}

func GetStringFromInterface(val interface{}) string {
	if val != nil {
		valueToReturn, ok := val.(string)
		if !ok {
			return ""
		}
		return valueToReturn
	}
	return ""
}

func GetBooleanFromInterface(val interface{}) bool {
	if val != nil {
		valueToReturn, ok := val.(bool)
		if !ok {
			return false
		}
		return valueToReturn
	}
	return false
}

func GetBooleanFromInterfaceArray(val interface{}) bool {
	stringArray := GetStringArrayFromInterface(val)
	if len(stringArray) > 0 {
		return stringArray[0] == "true"
	}
	return false
}
func GetStringArrayFromInterfaceArray(val []interface{}) []string {
	valueToReturn := []string{}
	for _, value := range val {
		valueToReturn = append(valueToReturn, value.(string))
	}
	return valueToReturn

}

func GetStringArrayFromInterface(val interface{}) []string {
	valueToReturn, ok := val.([]string)
	if !ok {
		interfaceArray, ok := val.([]interface{})
		if !ok {
			// is the interface a string, then return it as an array
			stringValue, ok := val.(string)
			if ok {
				return []string{stringValue}
			}
			return []string{}
		}

		valueToReturn = []string{}
		for _, val := range interfaceArray {
			stringValue, ok := val.(string)
			if ok {
				valueToReturn = append(valueToReturn, stringValue)
			}
		}
		return valueToReturn
	}
	return valueToReturn
}

func GetInterfaceArrayFromInterface(val interface{}) []interface{} {
	if val == nil {
		return []interface{}{}
	}
	interfaceArray, ok := val.([]interface{})
	if !ok {
		return []interface{}{}
	}
	return interfaceArray
}

func GetInterfaceArrayFromStringArray(array []string) []interface{} {
	interfaceArray := []interface{}{}
	if len(array) == 0 {
		return interfaceArray
	}

	for _, val := range array {
		interfaceArray = append(interfaceArray, val)
	}
	return interfaceArray
}

func GetInterfaceArrayFromCommaSeparatedString(val interface{}) []interface{} {
	if val == nil {
		return []interface{}{}
	}
	valStr, ok := val.(string)
	if !ok {
		return []interface{}{}
	}
	splittedVal := strings.Split(valStr, ",")
	interfaceArray := []interface{}{}
	for _, val := range splittedVal {
		interfaceArray = append(interfaceArray, val)
	}
	return interfaceArray
}

func GetMapArrayFromInterfaceArray(array interface{}) []map[string]interface{} {
	arrayOfMap, ok := array.([]map[string]interface{})
	if ok {
		return arrayOfMap
	}
	arrayOfInterface, ok := array.([]interface{})
	if !ok {
		return nil
	}
	arrayOfMaps := []map[string]interface{}{}
	for _, intf := range arrayOfInterface {
		mapToAdd, ok := intf.(map[string]interface{})
		if ok {
			arrayOfMaps = append(arrayOfMaps, mapToAdd)
		}
	}
	return arrayOfMaps
}

func GetPointerFromString(value string) *string {
	str := value
	return &str
}

func GetPointerFromInt(value int) *int {
	return &value
}

func GetIntFromPointer(value *int) int {
	if value == nil {
		return 0
	}
	return *value
}

func GetPointerFromFloat64(value float64) *float64 {
	return &value
}

func GetFloat64FromPointer(value *float64) float64 {
	if value == nil {
		return 0.0
	}
	return *value
}

func GetStringFromPointer(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

func GetStringFromMap(value map[string]interface{}) string {
	if value == nil {
		return "{}"
	}

	jsonValue, err := json.Marshal(value)
	if err != nil {
		return "{}"
	}
	return string(jsonValue)
}

func GetStringFromFloat64(value float64) string {
	return fmt.Sprintf("%v", value)
}

func GetIntFromString(value string) int {
	intValue, err := strconv.Atoi(value)
	if err != nil {
		return 0
	}
	return intValue
}
func ReplaceIframeWithLink(value string, linkColor string) (string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(value))
	if err != nil {
		return "", err
	}

	doc.Find("iframe").Each(func(i int, s *goquery.Selection) {
		src, _ := s.Attr("src")
		if strings.Contains(src, "youtube.com") {
			u, _ := url.Parse(src)
			u.RawQuery = ""
			youtube_id := path.Base(u.String())
			s.ReplaceWithHtml(`<a href='` + src + `' target='_blank' style="color:` + linkColor + `"> <img style="max-width:100%" src="https://i.ytimg.com/vi/` + youtube_id + `/maxresdefault.jpg"/><a/>`)
		} else {
			s.ReplaceWithHtml(`<a href='` + src + `' target='_blank' style="color:` + linkColor + `">` + src + `<a/>`)
		}
	})
	converted_html, _ := doc.Html()
	return converted_html, nil
}

func GetDomainFromURL(url string) string {
	domain := url
	if strings.HasPrefix(url, "http") {
		domain = strings.TrimPrefix(domain, "http://")
		domain = strings.TrimPrefix(domain, "https://")
	}
	return domain
}

func GetMapFromString(value string) map[string]interface{} {
	if value == "" {
		return nil
	}
	var valueMap map[string]interface{}
	err := json.Unmarshal([]byte(value), &valueMap)
	if err != nil {
		return nil
	}
	return valueMap
}

func GetByteArray(value interface{}) []byte {
	if value == nil {
		return nil
	}
	br, _ := json.Marshal(value)
	return br
}
func ConvertFloat64ArrayToFloat32(value []float64) []float32 {
	float32Array := []float32{}
	for _, val := range value {
		float32Array = append(float32Array, float32(val))
	}
	return float32Array
}

func CleanString(value string) string {
	return strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(value, "\n", ""), "\r", ""), "\u0026", ""), "  ", " ")
}

func CompareMapStringInterface(value1 map[string]interface{}, value2 map[string]interface{}) bool {
	// compare map string interface
	return reflect.DeepEqual(value1, value2)
}

// CompareMapStringInterfaceFields compares specific fields of two map[string]interface{} values, returns false if any field value not equal
func CompareMapStringInterfaceFields(value1, value2 map[string]interface{}, fields []string) bool {
	for _, field := range fields {
		if !reflect.DeepEqual(value1[field], value2[field]) {
			return false
		}
	}
	return true
}

func RemoveStopWordsFromString(text string) string {
	return stopwords.CleanString(text, "en", false)
}

func MergeMaps(dest, src map[string]interface{}) error {
	if dest == nil {
		return errors.New("destination map is nil")
	}
	if src == nil {
		return errors.New("source map is nil")
	}

	for key, value := range src {
		dest[key] = value
	}
	return nil
}

func DeduplicateInterfaceArray(array []interface{}) []interface{} {
	var deduplicatedArray = []interface{}{}
	for _, item := range array {
		itemExists := false
		for _, ditem := range deduplicatedArray {
			if ditem == item {
				itemExists = true
			}
		}
		if !itemExists {
			deduplicatedArray = append(deduplicatedArray, item)
		}
	}
	return deduplicatedArray
}
