package gospec_test

import (
	"context"
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
	gospec.Spec[*Document]
	word string
}

func NewTitleContainsWordSpec(word string) *TitleContainsWordSpec {
	s := &TitleContainsWordSpec{word: word}
	s.Spec = gospec.New[*Document](s)
	return s
}

func (s *TitleContainsWordSpec) IsSatisfiedBy(_ context.Context, candidate *Document) (bool, error) {
	return strings.Contains(candidate.Title, s.word), nil
}

type ContentContainsWordSpec struct {
	gospec.Spec[*Document]
	word string
}

func NewContentContainsWordSpec(word string) *ContentContainsWordSpec {
	s := &ContentContainsWordSpec{word: word}
	s.Spec = gospec.New[*Document](s)
	return s
}

func (s *ContentContainsWordSpec) IsSatisfiedBy(_ context.Context, candidate *Document) (bool, error) {
	return strings.Contains(candidate.Content, s.word), nil
}

type DateLowerSpec struct {
	gospec.Spec[*Document]
	date time.Time
}

func NewDateLowerSpec(date time.Time) *DateLowerSpec {
	s := &DateLowerSpec{date: date}
	s.Spec = gospec.New[*Document](s)
	return s
}

func (s *DateLowerSpec) IsSatisfiedBy(_ context.Context, candidate *Document) (bool, error) {
	return s.date.UnixNano() > candidate.Date.UnixNano(), nil
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
	titleContainsFirst := gospec.NewInline[*Document](func(ctx context.Context, candidate *Document) (bool, error) {
		return strings.Contains(candidate.Title, "First"), nil
	})
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

	// Output:
	// true
	// false
	// true
	// false
}
