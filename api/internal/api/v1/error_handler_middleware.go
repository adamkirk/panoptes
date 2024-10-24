package v1

import (
	"context"
)

func ErrorHandler[Req any, Resp any](debugErrors bool, handler func(context.Context, *Req) (*Resp, error)) (func (ctx context.Context, req *Req) (*Resp, error)) {
	return func (ctx context.Context, req *Req) (*Resp, error) {
		resp, err :=  handler(ctx, req)

		if err == nil {
			return resp, nil
		}

		// TODO: Add error hanlding logic
		return resp, err
	}
}
