package pubsub

import (
	"errors"
	"time"
)

const (
	InvalidOptionAlias = iota
	PublisherOptionAlias
	PublishOptionAlias
	SubscriberOptionAlias

	PublishOptionsStructOptionType = -1
)

var (
	ErrInvalidOptionType        = errors.New("invalid option type")
	ErrInvalidPublishAtValue    = errors.New("invalid publish at value")
	ErrInvalidPublishAfterValue = errors.New("invalid publish after value")
	ErrInvalidRescheduleValue   = errors.New("invalid reschedule value")
)

// OptionType is the integer type of the option
// this is used to infer the option type.
type OptionType int

// Alias ...
type OptionAlias int

// Option defines the methods that an option must implement.
type Option interface {
	Type() OptionType
	Alias() OptionAlias
	Value() any
}

// PublisherOption is an alias for Option, that should only be passed to publishers.
type PublisherOption Option

// PublishOption is an alias for an Option that should only be passed to a Publisher.
type PublishOption Option

// SubscriberOption is an alias for Option, that should only be passed to subscribers.
type SubscriberOption Option

// ApplyPublisherOptions is a generic function that applies PublisherOptions to T.
type ApplyPublisherOptions[T any] func(pubOpts T, opts ...PublisherOption) (T, error)

// ApplySubscriberOptions is a generic function that applies SubscriberOptions to T.
type ApplySubscriberOptions[T any] func(subOpts T, opts ...SubscriberOption) (T, error)

// ApplyPublishOptions is a generic function that applies PublishOptions to T.
type ApplyPublishOptions[T any] func(pubOpts T, opts ...PublishOption) (T, error)

func PublishOptionsStructOption(value any) PublishOption {
	return OptionValue{OptionValue: value, OptionAlias: PublishOptionAlias, OptionType: -1}
}

func PublishOptionsPublishAt(publishAt time.Time) PublishOption {
	return OptionValue{OptionValue: publishAt, OptionAlias: PublishOptionAlias, OptionType: -2}
}

func PublishOptionsPublishAfter(publishAfter time.Duration) PublishOption {
	return OptionValue{OptionValue: publishAfter, OptionAlias: PublishOptionAlias, OptionType: -3}
}

func PublishOptionsWithReschedule(reschedule bool) PublishOption {
	return OptionValue{OptionValue: reschedule, OptionAlias: PublishOptionAlias, OptionType: -4}
}

type OptionValue struct {
	OptionValue any
	OptionAlias OptionAlias
	OptionType  OptionType
}

func (o OptionValue) Type() OptionType {
	return o.OptionType
}

func (o OptionValue) Value() any {
	return o.OptionValue
}

func (o OptionValue) Alias() OptionAlias {
	return o.OptionAlias
}
