package gberr

import (
	"errors"
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/forbot161602/pbc-golang-lib/source/core/base/gbmtmsg"
)

func Is(err, target error) bool {
	return errors.Is(err, target)
}

func As(err error, target any) bool {
	return errors.As(err, target)
}

func AsError(err error) (target *Error, ok bool) {
	return target, As(err, &target)
}

func AsWrapError(err error) (target *WrapError, ok bool) {
	return target, As(err, &target)
}

func AsGenericError(err error) (target GenericError, ok bool) {
	return target, As(err, &target)
}

func AsInternalError(err error) (target *InternalError, ok bool) {
	return target, As(err, &target)
}

func AsValidationError(err error) (target *ValidationError, ok bool) {
	return target, As(err, &target)
}

func AsUnexpectedError(err error) (target *UnexpectedError, ok bool) {
	return target, As(err, &target)
}

func Aggravate(err error) error {
	if verr, ok := AsValidationError(err); ok {
		err = Unexpected(verr.Message(), verr.Options(), verr.Unwrap()...)
	}
	return err
}

type Error struct {
	message string
}

func New(message string) *Error {
	return &Error{
		message: message,
	}
}

func Newf(message string, args []any) *Error {
	return New(fmt.Sprintf(message, args...))
}

func (err *Error) Error() string {
	return fmt.Sprintf("<Error| %s>", err.message)
}

type WrapError struct {
	message string
	errs    []error
}

func Wrap(message string, errs ...error) *WrapError {
	return &WrapError{
		message: message,
		errs:    errs,
	}
}

func Wrapf(message string, args []any, errs ...error) *WrapError {
	return Wrap(fmt.Sprintf(message, args...), errs...)
}

func (err *WrapError) Error() string {
	return fmt.Sprintf("<WrapError| %s> caused from: %v", err.message, err.errs)
}

func (err *WrapError) Unwrap() []error {
	return err.errs
}

type GenericError interface {
	Error() string
	Unwrap() []error
	Message() *Message
	Options() *Options
	OutText() string
	LogText() string
}

type InternalError struct {
	message   *Message
	options   *Options
	errs      []error
	outArgs   []any
	logArgs   []any
	logFields LogFields
}

type InternalErrorOptions struct {
	OutArgs   []any
	LogArgs   []any
	LogFields LogFields
}

type (
	Message   = gbmtmsg.MetaMessage
	Options   = InternalErrorOptions
	LogFields = logrus.Fields
)

func Internal(message *Message, options *Options, errs ...error) *InternalError {
	err := (&internalErrorBuilder{options: options}).
		initialize().
		setMessage(message).
		setOptions().
		setErrors(errs...).
		setOutArgs().
		setLogArgs().
		setLogFields().
		build()
	return err
}

func (err *InternalError) Error() string {
	return fmt.Sprintf("<InternalError| %s> caused from: %v", err.LogText(), err.errs)
}

func (err *InternalError) Unwrap() []error {
	return err.errs
}

func (err *InternalError) Message() *Message {
	return err.message
}

func (err *InternalError) Options() *Options {
	return err.options
}

func (err *InternalError) OutText() string {
	return err.message.GetOutText(err.outArgs...)
}

func (err *InternalError) LogText() string {
	return err.message.GetLogText(err.logArgs...)
}

type internalErrorBuilder struct {
	err     *InternalError
	options *Options
}

func (builder *internalErrorBuilder) build() *InternalError {
	return builder.err
}

func (builder *internalErrorBuilder) initialize() *internalErrorBuilder {
	builder.err = &InternalError{}
	if builder.options == nil {
		builder.options = &Options{}
	}
	return builder
}

func (builder *internalErrorBuilder) setMessage(message *Message) *internalErrorBuilder {
	builder.err.message = message
	return builder
}

func (builder *internalErrorBuilder) setOptions() *internalErrorBuilder {
	builder.err.options = builder.options
	return builder
}

func (builder *internalErrorBuilder) setErrors(errs ...error) *internalErrorBuilder {
	builder.err.errs = errs
	return builder
}

func (builder *internalErrorBuilder) setOutArgs() *internalErrorBuilder {
	builder.err.outArgs = builder.options.OutArgs
	return builder
}

func (builder *internalErrorBuilder) setLogArgs() *internalErrorBuilder {
	builder.err.logArgs = builder.options.LogArgs
	return builder
}

func (builder *internalErrorBuilder) setLogFields() *internalErrorBuilder {
	builder.err.logFields = builder.options.LogFields
	return builder
}

type ValidationError struct {
	*InternalError
}

func Validation(message *Message, options *Options, errs ...error) *ValidationError {
	return &ValidationError{
		InternalError: Internal(message, options, errs...),
	}
}

func (err *ValidationError) Error() string {
	return fmt.Sprintf("<ValidationError| %s> caused from: %v", err.LogText(), err.errs)
}

type UnexpectedError struct {
	*InternalError
}

func Unexpected(message *Message, options *Options, errs ...error) *UnexpectedError {
	return &UnexpectedError{
		InternalError: Internal(message, options, errs...),
	}
}

func (err *UnexpectedError) Error() string {
	return fmt.Sprintf("<UnexpectedError| %s> caused from: %v", err.LogText(), err.errs)
}
