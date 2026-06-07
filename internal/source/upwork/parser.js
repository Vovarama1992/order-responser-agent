const { chromium } = require("playwright-extra");
const StealthPlugin = require("puppeteer-extra-plugin-stealth");

chromium.use(StealthPlugin());

(async () => {
  const context = await chromium.launchPersistentContext(
    "/Users/volodzya13/Library/Application Support/Google/Chrome/Default",
    {
      headless: true,
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

  const cards = page.locator("[data-test='JobTile']");
  const count = await cards.count();

  const orders = [];

  for (let i = 0; i < Math.min(count, 10); i++) {
    const card = cards.nth(i);

    const title =
      (await card.locator("a[href*='/jobs/']").first().textContent()) || "";

    const href =
      (await card.locator("a[href*='/jobs/']").first().getAttribute("href")) ||
      "";

    const description = (await card.textContent()) || "";

    const match = href.match(/~0?(\d+)/);

    orders.push({
      source: "upwork",
      id: match ? match[1] : href,
      url: `https://www.upwork.com${href}`,
      title: title.trim(),
      description: description.trim(),
      budget: "",
    });
  }

  console.log(JSON.stringify(orders));

  await context.close();
})();
