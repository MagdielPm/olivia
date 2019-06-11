package analysis

import (
	"encoding/json"
	"fmt"
	"github.com/olivia-ai/olivia/modules"
	"github.com/olivia-ai/olivia/util"
	"io/ioutil"
	"sort"
)

type Intent struct {
	Tag       string   `json:"tag"`
	Patterns  []string `json:"patterns"`
	Responses []string `json:"responses"`
	Context   string   `json:"context"`
}

type Document struct {
	Sentence Sentence
	Tag      string
}

func ReadIntents() []byte {
	bytes, err := ioutil.ReadFile("res/intents.json")
	if err != nil {
		fmt.Println(err)
	}

	return bytes
}

func SerializeIntents() []Intent {
	var intents []Intent

	err := json.Unmarshal(ReadIntents(), &intents)
	if err != nil {
		panic(err)
	}

	return intents
}

// SerializeModulesIntents retrieves all the registered modules and returns an array of Intents
func SerializeModulesIntents() []Intent {
	registeredModules := modules.GetModules()
	intents := make([]Intent, len(registeredModules))

	for _, module := range registeredModules {
		intents = append(intents, Intent{
			Tag:       module.Tag,
			Patterns:  module.Patterns,
			Responses: module.Responses,
			Context:   "",
		})
	}

	return intents
}

// Organize intents with an array of all words, an array with a representative word of each tag
// and an array of Documents which contains a word list associated with a tag
func Organize() (words, classes []string, documents []Document) {
	// Append the modules intents to the intents from res/intents.json
	intents := append(SerializeIntents(), SerializeModulesIntents()...)

	for _, intent := range intents {
		for _, pattern := range intent.Patterns {
			// Tokenize the pattern's sentence
			patternSentence := Sentence{pattern}

			// Add each word to response
			for _, word := range patternSentence.Tokenize() {

				if !util.Contains(words, word) {
					words = append(words, word)
				}
			}

			// Add a new document
			documents = append(documents, Document{
				patternSentence,
				intent.Tag,
			})

			// Add the intent tag to class if it doesn't exists
			if !util.Contains(classes, intent.Tag) {
				classes = append(classes, intent.Tag)
			}
		}
	}

	sort.Strings(words)
	sort.Strings(classes)

	return words, classes, documents
}
