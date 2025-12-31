package chromex

import (
	"context"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/r3dpixel/toolkit/stringsx"
)

const (
	// DefaultTimeout is the default timeout for Chrome operations.
	DefaultTimeout time.Duration = 60 * time.Second
)

// Extractor is a function that extracts data from a Chrome context.
// It receives the chromedp context and should return the extracted data or an error.
type Extractor[T any] func(ctx context.Context) (T, error)

// Options holds configuration for Chrome execution.
type Options struct {
	// Path to the Chrome executable. If empty, chromedp will try to find Chrome automatically.
	Path string
	// Timeout for the entire Chrome operation. Defaults to 120 seconds if not set.
	Timeout time.Duration
	// Custom chromedp flags. If empty, DefaultFlags() will be used.
	Flags []chromedp.ExecAllocatorOption
}

// DefaultFlags returns the default chromedp flags used for automation.
// These flags help avoid detection and ensure stable operation.
func DefaultFlags() []chromedp.ExecAllocatorOption {
	return []chromedp.ExecAllocatorOption{
		// Disable first-run and default browser checks
		chromedp.NoFirstRun,
		chromedp.NoDefaultBrowserCheck,
		// Disable sandboxing
		chromedp.NoSandbox,
		// Run Chrome in incognito mode
		chromedp.Flag("incognito", true),
		// Disable Blink features that could interfere with automation
		chromedp.Flag("disable-blink-features", "AutomationControlled"),
	}
}

// RunChrome launches Chrome with the given configuration and executes the extractor function.
// It handles the Chrome lifecycle (launch, context creation, cleanup) and returns the extracted data.
//
// Parameters:
//   - config: Chrome configuration (path, timeout, flags)
//   - extractor: Function that performs the actual extraction using the chromedp context
//
// Returns the extracted data of type T and any error encountered.
func RunChrome[T any](config Options, extractor Extractor[T]) (T, error) {
	// Use default flags if none provided
	if len(config.Flags) == 0 {
		config.Flags = DefaultFlags()
	}

	// Add the custom Chrome path if specified (prepend to avoid mutating the original slice)
	flags := config.Flags
	if stringsx.IsNotBlank(config.Path) {
		newFlags := make([]chromedp.ExecAllocatorOption, len(flags)+1)
		newFlags[0] = chromedp.ExecPath(config.Path)
		copy(newFlags[1:], flags)
		flags = newFlags
	}

	// Create allocator context
	allocCtx, cancelAlloc := chromedp.NewExecAllocator(context.Background(), flags...)
	defer cancelAlloc()

	// Create chromedp context
	ctx, cancelCtx := chromedp.NewContext(allocCtx)
	defer cancelCtx()

	// Set timeout (default 60 seconds)
	timeout := config.Timeout
	if timeout == 0 {
		timeout = DefaultTimeout
	}
	ctx, cancelTimeout := context.WithTimeout(ctx, timeout)
	defer cancelTimeout()

	// Execute the extractor function
	return extractor(ctx)
}
