const { chromium } = require("playwright");

(async () => {
  const browser = await chromium.launch({ headless: false });
  const page = await browser.newPage();

  await page.goto("https://www.upwork.com/nx/search/jobs/?q=workflow%20engine");

  await page.waitForTimeout(600000);

  console.log(await page.title());

  await page.screenshot({ path: "screen.png" });

  await browser.close();
})();
