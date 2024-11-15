package cloudwatchwriter

import (
	"context"
	"io"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	logTypes "github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
)

type SyncWriter struct {
	ctx               context.Context
	client            CWClient
	logGroupName      string
	logStreamName     string
	nextSequenceToken *string

	streamCreateErr error
	streamExists    bool
}

var _ io.Writer = (*SyncWriter)(nil)

func NewSync(ctx context.Context, client CWClient, groupName, streamName string) *SyncWriter {
	return &SyncWriter{
		ctx:           ctx,
		client:        client,
		logGroupName:  groupName,
		logStreamName: streamName,
	}
}

func (w *SyncWriter) Write(data []byte) (int, error) {

	timeNow := time.Now()

	if err := w.ensureStream(); err != nil {
		return 0, err
	}

	out, err := w.client.PutLogEvents(w.ctx, &cloudwatchlogs.PutLogEventsInput{
		LogGroupName:  &w.logGroupName,
		LogStreamName: &w.logStreamName,
		LogEvents: []logTypes.InputLogEvent{
			{
				Message:   aws.String(string(data)),
				Timestamp: aws.Int64(timeNow.UnixMilli()),
			},
		},
		SequenceToken: w.nextSequenceToken,
	})
	if err != nil {
		return 0, err
	}

	w.nextSequenceToken = out.NextSequenceToken

	return len(data), nil
}

func (w *SyncWriter) ensureStream() error {
	if w.streamExists {
		return nil
	}

	if w.streamCreateErr != nil {
		return w.streamCreateErr
	}

	resp, err := w.client.DescribeLogStreams(w.ctx, &cloudwatchlogs.DescribeLogStreamsInput{
		Limit:               aws.Int32(1),
		LogGroupName:        aws.String(w.logGroupName),
		LogStreamNamePrefix: aws.String(w.logStreamName),
	})
	if err != nil {
		w.streamCreateErr = err
		return err
	}

	if len(resp.LogStreams) > 0 {
		w.streamExists = true
		w.nextSequenceToken = resp.LogStreams[0].UploadSequenceToken
		return nil
	}

	_, err = w.client.CreateLogStream(w.ctx, &cloudwatchlogs.CreateLogStreamInput{
		LogGroupName:  aws.String(w.logGroupName),
		LogStreamName: aws.String(w.logStreamName),
	})
	if err != nil {
		w.streamCreateErr = err
		return err
	}

	w.streamExists = true
	return nil

}
