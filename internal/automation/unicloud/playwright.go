package playwright

import (
	"fmt"
	"log"

	"github.com/playwright-community/playwright-go"
)

func UpdateCert(
	spaceId string,
	username string,
	password string,
	cert string,
	key string,
) error {
	// Use default paths for playwright
	log.Println("[unicloud] using default paths for playwright")

	// Try to run playwright
	runOptions := &playwright.RunOptions{
		SkipInstallBrowsers: true,
	}
	pw, err := playwright.Run(runOptions)
	if err != nil {
		return err
	}
	defer pw.Stop()

	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(true),
	})
	if err != nil {
		return err
	}
	defer browser.Close()

	context, err := browser.NewContext(playwright.BrowserNewContextOptions{
		Locale:     playwright.String("zh-CN"),
		TimezoneId: playwright.String("Asia/Shanghai"),
		Viewport: &playwright.Size{
			Width:  1366,
			Height: 768,
		},
	})
	if err != nil {
		return err
	}

	page, err := context.NewPage()
	if err != nil {
		return err
	}

	// Set user agent and other headers
	page.SetExtraHTTPHeaders(map[string]string{
		"User-Agent":      "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/145.0.0.0 Safari/537.36",
		"Accept":          "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7",
		"Accept-Language": "zh-CN,zh;q=0.9",
		"Accept-Encoding": "gzip, deflate, br, zstd",
		"Cache-Control":   "no-cache",
		"Pragma":          "no-cache",
		"Priority":        "u=0, i",
		"Sec-CH-UA":       "\"Not:A-Brand\";v=\"99\", \"Google Chrome\";v=\"145\", \"Chromium\";v=\"145\"",
		"Sec-CH-UA-Mobile": "?0",
		"Sec-CH-UA-Platform": "\"Windows\"",
		"Sec-Fetch-Dest":   "document",
		"Sec-Fetch-Mode":   "navigate",
		"Sec-Fetch-Site":   "same-origin",
		"Sec-Fetch-User":   "?1",
		"Upgrade-Insecure-Requests": "1",
	})

	// 访问登录页面
	log.Println("[unicloud] navigating to login page...")

	_, err = page.Goto(
		"https://unicloud.dcloud.net.cn/pages/login/login",
		playwright.PageGotoOptions{
			WaitUntil: playwright.WaitUntilStateLoad,
			Timeout:   playwright.Float(60000),
		},
	)
	if err != nil {
		log.Println("[unicloud] failed to navigate to login page:", err)
		return err
	}

	// iframe 登录
	frame := page.FrameLocator("iframe")
	frame.Locator(`input[type="text"]`).Fill(username)
	frame.Locator(`input[type="password"]`).Fill(password)
	frame.Locator(`uni-button[type="primary"]`).Click()

	page.WaitForLoadState(playwright.PageWaitForLoadStateOptions{
		State: playwright.LoadStateNetworkidle,
	})

	// 进入空间
	space := page.Locator(
		fmt.Sprintf(`uni-view.list-title.left:has-text("%s")`, spaceId),
	)
	space.ScrollIntoViewIfNeeded()
	// page.WaitForTimeout(1000)

	space.
		Locator("xpath=../..").
		Locator("uni-text.title-text.to-link").
		Click()

	// 等待页面导航完成
	page.WaitForLoadState(playwright.PageWaitForLoadStateOptions{
		State: playwright.LoadStateNetworkidle,
	})

	// 前端网页托管
	// 等待下一步要操作的元素出现
	page.Locator(".title span:has-text('前端网页托管')").First().WaitFor(playwright.LocatorWaitForOptions{
		State: playwright.WaitForSelectorStateVisible,
	})
	page.Locator(".title span:has-text('前端网页托管')").First().Click()
	// page.WaitForTimeout(1000)
	// 进入参数配置页签
	// 等待下一步要操作的元素出现
	page.Locator(".el-tabs__item:has-text('参数配置')").First().WaitFor(playwright.LocatorWaitForOptions{
		State: playwright.WaitForSelectorStateVisible,
	})

	page.Locator(".el-tabs__item:has-text('参数配置')").Click()
	// page.WaitForTimeout(1000)

	// 等待更新证书按钮出现
	page.Locator("uni-button:has-text('更新证书')").WaitFor(playwright.LocatorWaitForOptions{
		State: playwright.WaitForSelectorStateVisible,
	})
	page.Locator("uni-button:has-text('更新证书')").Click()
	// page.WaitForTimeout(1000)

	// 填写证书
	page.Locator("textarea.uni-textarea-textarea").Nth(0).Fill(cert)
	page.Locator("textarea.uni-textarea-textarea").Nth(1).Fill(key)

	page.Locator("uni-button:has-text('确定')").Click()
	// page.WaitForTimeout(3000)

	page.WaitForURL("**/unicloud/**", playwright.PageWaitForURLOptions{
		WaitUntil: playwright.WaitUntilStateNetworkidle,
	})

	log.Println("[unicloud] cert updated success")

	return nil
}
