package gospec

import (
	"context"
	"fmt"
	"strings"
)

type Satisfiable interface {
	Describe() string
	IsSatisfiedBy(ctx context.Context, candidate any) (bool, error)
}

type Spec interface {
	Satisfiable
	And(condition Satisfiable, others ...Satisfiable) Spec
	Or(condition Satisfiable, others ...Satisfiable) Spec
	Xor(condition Satisfiable) Spec
	Not() Spec
}

func New(initial Satisfiable) Spec {
	return &compositeSpec{Satisfiable: initial}
}

type compositeSpec struct {
	Satisfiable
}

func (s *compositeSpec) And(condition Satisfiable, others ...Satisfiable) Spec {
	return newAndSpec(append([]Satisfiable{s.Satisfiable, condition}, others...)...)
}

func (s *compositeSpec) Or(condition Satisfiable, others ...Satisfiable) Spec {
	return newOrSpec(append([]Satisfiable{s.Satisfiable, condition}, others...)...)
}

func (s *compositeSpec) Xor(condition Satisfiable) Spec {
	return newXorSpec(s.Satisfiable, condition)
}

func (s *compositeSpec) Not() Spec {
	return newNotSpec(s.Satisfiable)
}

type andSpec struct {
	Spec
	conditions []Satisfiable
}

func newAndSpec(conditions ...Satisfiable) *andSpec {
	s := &andSpec{conditions: conditions}
	s.Spec = New(s)
	return s
}

func (s *andSpec) IsSatisfiedBy(ctx context.Context, candidate any) (bool, error) {
	for _, condition := range s.conditions {
		satisfied, err := condition.IsSatisfiedBy(ctx, candidate)
		if err != nil {
			return false, err
		}

		if !satisfied {
			return false, nil
		}
	}

	return true, nil
}

func (s *andSpec) Describe() string {
	var descriptions []string
	for _, condition := range s.conditions {
		descriptions = append(descriptions, condition.Describe())
	}

	return fmt.Sprintf("(%s)", strings.Join(descriptions, " AND "))
}

type orSpec struct {
	Spec
	conditions []Satisfiable
}

func newOrSpec(conditions ...Satisfiable) *orSpec {
	s := &orSpec{conditions: conditions}
	s.Spec = New(s)
	return s
}

func (s *orSpec) IsSatisfiedBy(ctx context.Context, candidate any) (bool, error) {
	for _, condition := range s.conditions {
		satisfied, err := condition.IsSatisfiedBy(ctx, candidate)
		if err != nil {
			return false, err
		}

		if satisfied {
			return true, nil
		}
	}

	return false, nil
}

func (s *orSpec) Describe() string {
	var descriptions []string
	for _, condition := range s.conditions {
		descriptions = append(descriptions, condition.Describe())
	}

	return fmt.Sprintf("(%s)", strings.Join(descriptions, " OR "))
}

type xorSpec struct {
	Spec
	left  Satisfiable
	right Satisfiable
}

func newXorSpec(left Satisfiable, right Satisfiable) *xorSpec {
	s := &xorSpec{left: left, right: right}
	s.Spec = New(s)
	return s
}

func (s *xorSpec) IsSatisfiedBy(ctx context.Context, candidate any) (bool, error) {
	l, err := s.left.IsSatisfiedBy(ctx, candidate)
	if err != nil {
		return false, err
	}

	r, err := s.right.IsSatisfiedBy(ctx, candidate)
	if err != nil {
		return false, err
	}

	return l != r, nil
}

func (s *xorSpec) Describe() string {
	return fmt.Sprintf("(%s XOR %s)", s.left.Describe(), s.right.Describe())
}

type notSpec struct {
	Spec
	condition Satisfiable
}

func newNotSpec(condition Satisfiable) *notSpec {
	s := &notSpec{condition: condition}
	s.Spec = New(s)
	return s
}

func (s *notSpec) IsSatisfiedBy(ctx context.Context, candidate any) (bool, error) {
	b, err := s.condition.IsSatisfiedBy(ctx, candidate)
	if err != nil {
		return false, err
	}

	return !b, nil
}

func (s *notSpec) Describe() string {
	return fmt.Sprintf("NOT(%s)", s.condition.Describe())
}
