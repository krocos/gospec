package gospec_test

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/krocos/gospec"
)

// Define composite specifications of any complexity

type Document struct {
	Title   string
	Content string
	Date    time.Time
}

type TitleContainsWordSpec struct {
	gospec.Spec
	word string
}

func NewTitleContainsWordSpec(word string) *TitleContainsWordSpec {
	s := &TitleContainsWordSpec{word: word}
	s.Spec = gospec.New(s)
	return s
}

func (s *TitleContainsWordSpec) Describe() string {
	return fmt.Sprintf("doc title must contain '%s'", s.word)
}

func (s *TitleContainsWordSpec) IsSatisfiedBy(_ context.Context, candidate any) (bool, error) {
	return strings.Contains(candidate.(*Document).Title, s.word), nil
}

type ContentContainsWordSpec struct {
	gospec.Spec
	word string
}

func NewContentContainsWordSpec(word string) *ContentContainsWordSpec {
	s := &ContentContainsWordSpec{word: word}
	s.Spec = gospec.New(s)
	return s
}

func (s *ContentContainsWordSpec) Describe() string {
	return fmt.Sprintf("doc content must contain '%s'", s.word)
}

func (s *ContentContainsWordSpec) IsSatisfiedBy(_ context.Context, candidate any) (bool, error) {
	return strings.Contains(candidate.(*Document).Content, s.word), nil
}

type DateLowerSpec struct {
	gospec.Spec
	date time.Time
}

func NewDateLowerSpec(date time.Time) *DateLowerSpec {
	s := &DateLowerSpec{date: date}
	s.Spec = gospec.New(s)
	return s
}

func (s *DateLowerSpec) Describe() string {
	return fmt.Sprintf("doc date must be lower than '%s'", s.date.Format(time.RFC3339))
}

func (s *DateLowerSpec) IsSatisfiedBy(_ context.Context, candidate any) (bool, error) {
	return s.date.UnixNano() > candidate.(*Document).Date.UnixNano(), nil
}

func TestCompositeSpec(t *testing.T) {
	// Composite specification example

	date := time.Date(2023, 6, 27, 20, 56, 0, 0, time.UTC)

	// Declaring our docs
	docs := []*Document{
		{Title: "First title", Content: "First doc content", Date: date},
		{Title: "Second title", Content: "Second doc content", Date: date.Add(time.Hour)},
		{Title: "Third title", Content: "Third doc content", Date: date.Add(2 * time.Hour)},
		{Title: "Fourth title", Content: "Fourth doc content", Date: date.Add(3 * time.Hour)},
	}

	// Declaring our specs
	titleContainsFirst := NewTitleContainsWordSpec("First")
	titleContainsThird := NewTitleContainsWordSpec("Third")

	contentContainsFirst := NewContentContainsWordSpec("First")
	contentContainsThird := NewContentContainsWordSpec("Third")

	dateLowerNowSpec := NewDateLowerSpec(time.Date(2023, 6, 27, 23, 0, 0, 0, time.UTC))

	// Composing spec
	contentSpec := titleContainsFirst.And(contentContainsFirst).Or(titleContainsThird.And(contentContainsThird))
	spec := dateLowerNowSpec.And(contentSpec)

	ctx := context.Background()
	for _, doc := range docs {
		satisfies, _ := spec.IsSatisfiedBy(ctx, doc)
		t.Log(satisfies)
	}

	// Describing our spec
	t.Log(spec.Describe())

	// Output:
	// true
	// false
	// true
	// false
	// (doc date must be lower than '2023-06-27T23:00:00Z' AND ((doc title
	// must contain 'First' AND doc content must contain 'First') OR (doc
	// title must contain 'Third' AND doc content must contain 'Third')))
}

func TestSetOperators(t *testing.T) {
	gospec.SetOperators("&&", "||", "!=", "!")
	t.Cleanup(func() {
		gospec.SetOperators("AND", "OR", "XOR", "NOT")
	})

	// Declaring our specs
	titleContainsFirst := NewTitleContainsWordSpec("First")
	titleContainsThird := NewTitleContainsWordSpec("Third")

	contentContainsFirst := NewContentContainsWordSpec("First")
	contentContainsThird := NewContentContainsWordSpec("Third")

	dateLowerNowSpec := NewDateLowerSpec(time.Date(2023, 6, 27, 23, 0, 0, 0, time.UTC))

	// Composing spec
	contentSpec := titleContainsFirst.And(contentContainsFirst).Or(titleContainsThird.And(contentContainsThird))
	spec := dateLowerNowSpec.And(contentSpec)

	t.Log(spec.Describe())

	// (doc date must be lower than '2023-06-27T23:00:00Z' && ((doc title
	// must contain 'First' && doc content must contain 'First') || (doc
	// title must contain 'Third' && doc content must contain 'Third')))
}
