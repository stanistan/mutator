package container

import (
	corev1 "k8s.io/api/core/v1"

	"github.com/stanistan/mutator/internal/lens"
)

// Mutator mutates a *corev1.Container or error.
//
// There is guarantee that a container _was not_ mutated
// if there was a failure here.
type Mutator interface {
	MutateContainer(Container) error
}

type Mutators []Mutator

func (ms Mutators) Mutate(c Container) error {
	for _, m := range ms {
		if err := m.MutateContainer(c); err != nil {
			return err
		}
	}

	return nil
}

type MutatorFunc func(Container) error

func (f MutatorFunc) MutateContainer(c Container) error {
	return f(c)
}

type Filter func(Container) bool

func Filtered(f Filter, m Mutator) Mutator {
	return MutatorFunc(func(c Container) error {
		if f != nil && !f(c) {
			return nil
		}

		return m.MutateContainer(c)
	})
}

func containerLens[T any](
	get func(*corev1.Container) T,
	set func(*corev1.Container, T),
) lens.Lens[T, Container, *corev1.Container, MutatorFunc] {
	return lens.Lens[T, Container, *corev1.Container, MutatorFunc]{Get: get, Set: set}
}

var (
	securityContextLens = containerLens(
		func(c *corev1.Container) *corev1.SecurityContext { return c.SecurityContext },
		func(c *corev1.Container, val *corev1.SecurityContext) { c.SecurityContext = val },
	)

	UpdateSecurityContext = securityContextLens.Mutator
	SetSecurityContext    = securityContextLens.InfallibleMutator
)

var (
	resourceLens = containerLens(
		func(c *corev1.Container) corev1.ResourceRequirements { return c.Resources },
		func(c *corev1.Container, val corev1.ResourceRequirements) { c.Resources = val },
	)

	UpdateResources = resourceLens.Mutator
	SetResources    = resourceLens.InfallibleMutator
)
