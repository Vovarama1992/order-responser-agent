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

  await page.waitForSelector("[data-test='JobTile']", {
    timeout: 30000,
  });

  const jobs = page.locator("[data-test='JobTile']");

  const count = await jobs.count();

  console.log("JOBS:", count);

  const limit = Math.min(count, 3);

  for (let i = 0; i < limit; i++) {
    const job = jobs.nth(i);

    let title = "";
    let url = "";
    let description = "";

    try {
      title = (await job.locator("h2, h3, a").first().textContent())?.trim();
    } catch {}

    try {
      url = await job.locator("a").first().getAttribute("href");
    } catch {}

    try {
      description = (await job.textContent())?.trim();
    } catch {}

    console.log("\n====================");
    console.log("JOB:", i + 1);
    console.log("TITLE:", title);
    console.log("URL:", url);
    console.log("DESCRIPTION:");
    console.log(description?.slice(0, 1000));
  }

  await new Promise(() => {});
})();
