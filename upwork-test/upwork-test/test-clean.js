const { chromium } = require("playwright");

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

  console.log("OPENED");

  await new Promise(() => {});
})();
