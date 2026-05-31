package logger

import "context"

type contextFieldsKey struct{}

func getFieldsFromContext(ctx context.Context) map[string]any {
	if ctx == nil {
		return nil
	}

	fields, ok := ctx.Value(contextFieldsKey{}).(map[string]any)
	if !ok {
		return nil
	}

	return fields
}

func contextWithFields(ctx context.Context, fields map[string]any) context.Context {
	return context.WithValue(ctx, contextFieldsKey{}, fields)
}
