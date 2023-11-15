package gospec

import (
	"context"
)

type Satisfiable[T any] interface {
	IsSatisfiedBy(ctx context.Context, candidate T) (bool, error)
}

type SatisfiableFn[T any] func(ctx context.Context, candidate T) (bool, error)

func (f SatisfiableFn[T]) IsSatisfiedBy(ctx context.Context, candidate T) (bool, error) {
	return f(ctx, candidate)
}

type Spec[T any] interface {
	Satisfiable[T]
	And(condition Satisfiable[T], others ...Satisfiable[T]) Spec[T]
	Or(condition Satisfiable[T], others ...Satisfiable[T]) Spec[T]
	Xor(condition Satisfiable[T]) Spec[T]
	Not() Spec[T]
}

func NewInline[T any](fn SatisfiableFn[T]) Spec[T] {
	return New[T](fn)
}

func New[T any](initial Satisfiable[T]) Spec[T] {
	return &compositeSpec[T]{Satisfiable: initial}
}

type compositeSpec[T any] struct {
	Satisfiable[T]
}

func (s *compositeSpec[T]) And(condition Satisfiable[T], others ...Satisfiable[T]) Spec[T] {
	return newAndSpec(append([]Satisfiable[T]{s.Satisfiable, condition}, others...)...)
}

func (s *compositeSpec[T]) Or(condition Satisfiable[T], others ...Satisfiable[T]) Spec[T] {
	return newOrSpec(append([]Satisfiable[T]{s.Satisfiable, condition}, others...)...)
}

func (s *compositeSpec[T]) Xor(condition Satisfiable[T]) Spec[T] {
	return newXorSpec(s.Satisfiable, condition)
}

func (s *compositeSpec[T]) Not() Spec[T] {
	return newNotSpec(s.Satisfiable)
}

type andSpec[T any] struct {
	Spec[T]
	conditions []Satisfiable[T]
}

func newAndSpec[T any](conditions ...Satisfiable[T]) *andSpec[T] {
	s := &andSpec[T]{conditions: conditions}
	s.Spec = New[T](s)
	return s
}

func (s *andSpec[T]) IsSatisfiedBy(ctx context.Context, candidate T) (bool, error) {
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

type orSpec[T any] struct {
	Spec[T]
	conditions []Satisfiable[T]
}

func newOrSpec[T any](conditions ...Satisfiable[T]) *orSpec[T] {
	s := &orSpec[T]{conditions: conditions}
	s.Spec = New[T](s)
	return s
}

func (s *orSpec[T]) IsSatisfiedBy(ctx context.Context, candidate T) (bool, error) {
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

type xorSpec[T any] struct {
	Spec[T]
	left  Satisfiable[T]
	right Satisfiable[T]
}

func newXorSpec[T any](left Satisfiable[T], right Satisfiable[T]) *xorSpec[T] {
	s := &xorSpec[T]{left: left, right: right}
	s.Spec = New[T](s)
	return s
}

func (s *xorSpec[T]) IsSatisfiedBy(ctx context.Context, candidate T) (bool, error) {
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

type notSpec[T any] struct {
	Spec[T]
	condition Satisfiable[T]
}

func newNotSpec[T any](condition Satisfiable[T]) *notSpec[T] {
	s := &notSpec[T]{condition: condition}
	s.Spec = New[T](s)
	return s
}

func (s *notSpec[T]) IsSatisfiedBy(ctx context.Context, candidate T) (bool, error) {
	b, err := s.condition.IsSatisfiedBy(ctx, candidate)
	if err != nil {
		return false, err
	}

	return !b, nil
}
