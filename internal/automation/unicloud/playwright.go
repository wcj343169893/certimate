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
	if err := playwright.Install(); err != nil {
		return err
	}

	pw, err := playwright.Run()
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
	// 访问登录页面
	log.Println("[unicloud] navigating to login page...")

	_, err = page.Goto(
		"https://unicloud.dcloud.net.cn/pages/login/login?change=1",
		playwright.PageGotoOptions{
			WaitUntil: playwright.WaitUntilStateNetworkidle,
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
