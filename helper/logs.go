package helper

import (
	"context"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
	"path/filepath"
	"time"
)

// FileHook is a custom hook for logging to a file with a different formatter
type FileHook struct {
	Writer    io.Writer
	Formatter logrus.Formatter
	LevelsVal []logrus.Level
}

func NewFileHook(levels []logrus.Level, writer io.Writer, formatter logrus.Formatter) *FileHook {
	return &FileHook{
		Writer:    writer,
		Formatter: formatter,
		LevelsVal: levels,
	}
}

func (hook *FileHook) Levels() []logrus.Level {
	return hook.LevelsVal
}

func (hook *FileHook) Fire(entry *logrus.Entry) error {
	if os.Getenv("GRPC_PORT") == "50052" {
		entry.Data["Environment"] = "Development"
	} else {
		entry.Data["Environment"] = "Production"
	}

	line, err := hook.Formatter.Format(entry)
	if err != nil {
		logrus.Errorf("Error formatting log entry for file: %v", err)
		return err
	}

	// Write the formatted entry to the writer
	_, err = hook.Writer.Write(line)
	if err != nil {
		logrus.Errorf("Error writing log entry to file: %v", err)
		return err
	}
	return nil
}

// SetupLogger initializes the logger with both terminal and file logging
func SetupLogger() {
	logDir := "logs"
	if err := os.MkdirAll(logDir, 0755); err != nil {
		logrus.Fatalf("Failed to create log directory: %v", err)
	}

	logFilePath := filepath.Join(logDir, "app.log")

	fileLogger := &lumberjack.Logger{
		Filename:   logFilePath,
		MaxSize:    10, // Megabytes
		MaxBackups: 3,
		MaxAge:     28, // Days
		Compress:   true,
	}

	// Set up terminal logger (this will be the default output)
	logrus.SetOutput(os.Stdout)
	logrus.SetFormatter(&logrus.TextFormatter{
		TimestampFormat:        "2006-01-02 15:04:05",
		FullTimestamp:          true,
		ForceColors:            true,  // Enable colors for terminal output
		DisableColors:          false, // Keep colors in terminal
		QuoteEmptyFields:       true,
		DisableQuote:           true,
		DisableLevelTruncation: true,
		PadLevelText:           false,
	})

	// Set log level
	logrus.SetLevel(logrus.InfoLevel)

	// Add custom hook for file logging with a different formatter (no colors)
	logrus.AddHook(NewFileHook(logrus.AllLevels, fileLogger, &logrus.TextFormatter{
		TimestampFormat:  "2006-01-02 15:04:05",
		FullTimestamp:    true,
		ForceColors:      false, // Disable colors for file output
		DisableColors:    true,
		QuoteEmptyFields: true,
	}))
}

// LogrusLoggerUnaryInterceptor is a gRPC unary interceptor that logs the request and response
func LogrusLoggerUnaryInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	// Save start time
	start := time.Now()

	// Call the handler to complete the normal execution of RPC
	resp, err := handler(ctx, req)

	// Measure latency
	latency := time.Since(start)

	// Fetch gRPC metadata (acts as a substitute for headers, query params, etc.)
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		md = metadata.New(nil)
	}

	// Serialize the request message to calculate its length
	reqBytes, err := proto.Marshal(req.(proto.Message)) // Casting to proto.Message and Marshalling
	if err != nil {
		logrus.WithError(err).Warn("Failed to marshal request message")
	}

	// Prepare log entry with relevant gRPC data
	logrus.WithFields(logrus.Fields{
		"method":         info.FullMethod, // gRPC method name
		"latency":        latency,
		"user_agent":     getMetadataValue(md, "user-agent"),
		"request_id":     getMetadataValue(md, "x-request-id"),
		"grpc_status":    grpc.Code(err).String(),
		"content_length": len(reqBytes),
		"error":          err,
	}).Info("gRPC request")

	return resp, err
}

// Helper function to get metadata value
func getMetadataValue(md metadata.MD, key string) string {
	if val, ok := md[key]; ok && len(val) > 0 {
		return val[0]
	}
	return ""
}
