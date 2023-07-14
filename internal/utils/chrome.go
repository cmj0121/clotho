// Get the global utility for the internal packages.
package utils

import (
	"context"

	"github.com/chromedp/chromedp"
	"github.com/go-rod/stealth"
)

// The chrome wrapper with undetectable techniques.
type Chrome struct {
	Headless bool `group:"chrome" help:"Run the browser in headless mode." default:"true" negatable:""`

	// The ChromeDP parent context.
	exec_ctx    context.Context
	exec_cancel context.CancelFunc

	// The ChromeDP child context.
	ctx    context.Context
	cancel context.CancelFunc
}

// open the necessary resources.
func (c *Chrome) Prologue() {
	// The customized chrome options.
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		// general chrome options
		chromedp.DisableGPU,
		chromedp.NoFirstRun,
		chromedp.NoDefaultBrowserCheck,
		// make the browser window undetecteable for webdriver detection.
		chromedp.Flag("enable-automation", false),
		chromedp.Flag("disable-blink-features", "AutomationControlled"),
		// headless chrome options
		chromedp.Flag("headless", c.Headless),
	)

	c.exec_ctx, c.exec_cancel = chromedp.NewExecAllocator(context.Background(), opts...)
	c.ctx, c.cancel = chromedp.NewContext(c.exec_ctx)
}

// clean up the resources.
func (c *Chrome) Epilogue() {
	c.cancel()
	c.exec_cancel()
}

// Navigate to the given url.
func (c *Chrome) Navigate(url string) (err error) {
	err = chromedp.Run(
		c.ctx,
		// insert the stealth-js to the browser before the page is loaded
		chromedp.Evaluate(stealth.JS, nil),
		chromedp.Navigate(url),
	)
	return
}

// Run the sevearal actions.
func (c *Chrome) Run(actions ...chromedp.Action) (err error) {
	err = chromedp.Run(c.ctx, actions...)
	return
}
