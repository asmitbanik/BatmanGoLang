package internal

import (
	"encoding/gob"
	"os"
	"regexp"
	"strings"
	"sync"
	"errors"
	"math"
	"d:/Asmit/shazamCode/internal/data"
)

// NaiveBayesModel is a simple multinomial NB for framework detection
type NaiveBayesModel struct {
	Classes []string
	WordCounts map[string]map[string]int // class -> word -> count
	ClassCounts map[string]int
	Vocab map[string]struct{}
	mu sync.RWMutex
}

func NewNaiveBayesModel() *NaiveBayesModel {
	m := &NaiveBayesModel{
		Classes: []string{},
		WordCounts: map[string]map[string]int{},
		ClassCounts: map[string]int{},
		Vocab: map[string]struct{}{},
	}
	return m
}

func (m *NaiveBayesModel) Train(class, text string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.WordCounts[class]; !ok {
		m.Classes = append(m.Classes, class)
		m.WordCounts[class] = map[string]int{}
	}
	m.ClassCounts[class]++
	words := tokenize(text)
	for _, w := range words {
		m.WordCounts[class][w]++
		m.Vocab[w] = struct{}{}
	}
}

func (m *NaiveBayesModel) Predict(text string) (string, float64) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	words := tokenize(text)
	bestClass := "Unknown"
	bestScore := math.Inf(-1)
	for _, class := range m.Classes {
		score := math.Log(float64(m.ClassCounts[class]+1))
		total := 0
		for _, c := range m.WordCounts[class] {
			total += c
		}
		for _, w := range words {
			count := m.WordCounts[class][w]
			score += math.Log(float64(count+1)) - math.Log(float64(total+len(m.Vocab)))
		}
		if score > bestScore {
			bestScore = score
			bestClass = class
		}
	}
	return bestClass, bestScore
}

func (m *NaiveBayesModel) Save(path string) error {
	m.mu.RLock()
	defer m.mu.RUnlock()
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return gob.NewEncoder(f).Encode(m)
}

func (m *NaiveBayesModel) Load(path string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return gob.NewDecoder(f).Decode(m)
}

func tokenize(text string) []string {
	re := regexp.MustCompile(`[a-zA-Z_][a-zA-Z0-9_]+`)
	return re.FindAllString(strings.ToLower(text), -1)
}

var (
	model     *NaiveBayesModel
	modelOnce sync.Once
)

func getModel() *NaiveBayesModel {
	modelOnce.Do(func() {
		model = NewNaiveBayesModel()
		// Try to load from disk, else train from embedded data
		if err := model.Load("framework_model.gob"); err != nil {
			for fw, exs := range data.FrameworkExamples {
				for _, ex := range exs {
					model.Train(fw, ex)
				}
			}
			_ = model.Save("framework_model.gob")
		}
	})
	return model
}

// DetectFrameworkML attempts prediction using the trained model
func DetectFrameworkML(code string) (string, float64, error) {
	m := getModel()
	fw, score := m.Predict(code)
	if fw == "Unknown" {
		return fw, score, errors.New("no confident match")
	}
	return fw, score, nil
}
