const { chromium } = require("playwright-extra");
const StealthPlugin = require("puppeteer-extra-plugin-stealth");

chromium.use(StealthPlugin());

(async () => {
  const context = await chromium.launchPersistentContext(
    "/Users/volodzya13/Library/Application Support/Google/Chrome/Default",
    {
      headless: false,
      channel: "chrome",
    },
  );

  const page = context.pages()[0];

  await page.goto(
    "https://www.upwork.com/nx/search/jobs/?q=workflow%20engine",
    {
      waitUntil: "domcontentloaded",
    },
  );

  await page.waitForTimeout(5000);

  const jobs = await page.locator("[data-test='JobTile']").count();

  console.log("JOBS:", jobs);

  if (jobs > 0) {
    const first = page.locator("[data-test='JobTile']").first();

    console.log("FIRST JOB:");
    console.log(await first.textContent());
  }

  await new Promise(() => {});
})();
