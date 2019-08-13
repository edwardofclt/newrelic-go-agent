package nrmicro

import (
	"context"

	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/metadata"
	newrelic "github.com/newrelic/go-agent"
)

type nrWrapper struct {
	client.Client
}

func startExternal(ctx context.Context, procedure, host string) (context.Context, newrelic.ExternalSegment) {
	var seg newrelic.ExternalSegment
	if txn := newrelic.FromContext(ctx); nil != txn {
		seg = newrelic.ExternalSegment{
			StartTime: newrelic.StartSegmentNow(txn),
			Procedure: procedure,
			Library:   "Micro",
			Host:      host,
		}
		payload := txn.CreateDistributedTracePayload()
		if txt := payload.Text(); "" != txt {
			md, _ := metadata.FromContext(ctx)
			md = metadata.Copy(md)
			md[newrelic.DistributedTracePayloadHeader] = txt
			ctx = metadata.NewContext(ctx, md)
		}
	}
	return ctx, seg
}

func (n *nrWrapper) Publish(ctx context.Context, msg client.Message, opts ...client.PublishOption) error {
	ctx, seg := startExternal(ctx, "Publish", n.Options().Broker.Address())
	defer seg.End()
	return n.Client.Publish(ctx, msg, opts...)
}

func (n *nrWrapper) Call(ctx context.Context, req client.Request, rsp interface{}, opts ...client.CallOption) error {
	ctx, seg := startExternal(ctx, req.Endpoint(), req.Service())
	defer seg.End()
	return n.Client.Call(ctx, req, rsp, opts...)
}

func ClientWrapper(c client.Client) client.Client {
	return &nrWrapper{c}
}
