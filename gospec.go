package gospec

import (
	"context"
	"fmt"
)

type Satisfiable interface {
	Describe() string
	IsSatisfiedBy(ctx context.Context, candidate any) (bool, error)
}

type Spec interface {
	Satisfiable
	And(condition Satisfiable) Spec
	Or(condition Satisfiable) Spec
	Xor(condition Satisfiable) Spec
	Not() Spec
}

func New(initial Satisfiable) Spec {
	return &compositeSpec{Satisfiable: initial}
}

type compositeSpec struct {
	Satisfiable
}

func (s *compositeSpec) And(condition Satisfiable) Spec {
	return newAndSpec(s.Satisfiable, condition)
}

func (s *compositeSpec) Or(condition Satisfiable) Spec {
	return newOrSpec(s.Satisfiable, condition)
}

func (s *compositeSpec) Xor(condition Satisfiable) Spec {
	return newXorSpec(s.Satisfiable, condition)
}

func (s *compositeSpec) Not() Spec {
	return newNotSpec(s.Satisfiable)
}

type andSpec struct {
	Spec
	left  Satisfiable
	right Satisfiable
}

func newAndSpec(left Satisfiable, right Satisfiable) *andSpec {
	s := &andSpec{left: left, right: right}
	s.Spec = New(s)
	return s
}

func (s *andSpec) IsSatisfiedBy(ctx context.Context, candidate any) (bool, error) {
	l, err := s.left.IsSatisfiedBy(ctx, candidate)
	if err != nil {
		return false, err
	}

	r, err := s.right.IsSatisfiedBy(ctx, candidate)
	if err != nil {
		return false, err
	}

	return l && r, nil
}

func (s *andSpec) Describe() string {
	return fmt.Sprintf("(%s AND %s)", s.left.Describe(), s.right.Describe())
}

type orSpec struct {
	Spec
	left  Satisfiable
	right Satisfiable
}

func newOrSpec(left Satisfiable, right Satisfiable) *orSpec {
	s := &orSpec{left: left, right: right}
	s.Spec = New(s)
	return s
}

func (s *orSpec) IsSatisfiedBy(ctx context.Context, candidate any) (bool, error) {
	l, err := s.left.IsSatisfiedBy(ctx, candidate)
	if err != nil {
		return false, err
	}

	r, err := s.right.IsSatisfiedBy(ctx, candidate)
	if err != nil {
		return false, err
	}

	return l || r, nil
}

func (s *orSpec) Describe() string {
	return fmt.Sprintf("(%s OR %s)", s.left.Describe(), s.right.Describe())
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

func (s *xorSpec) String() string {
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
