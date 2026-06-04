const { chromium } = require("playwright-extra");
const StealthPlugin = require("puppeteer-extra-plugin-stealth");

chromium.use(StealthPlugin());

(async () => {
  try {
    const browser = await chromium.launch({
      headless: false,
      channel: "chrome",
    });

    browser.on("disconnected", () => {
      console.log("BROWSER DISCONNECTED");
    });

    const page = await browser.newPage();

    page.on("close", () => {
      console.log("PAGE CLOSED");
    });

    page.on("crash", () => {
      console.log("PAGE CRASH");
    });

    await page.goto(
      "https://www.upwork.com/nx/search/jobs/?q=workflow%20engine",
      { waitUntil: "domcontentloaded" },
    );

    console.log("OPENED");

    await new Promise(() => {});
  } catch (e) {
    console.error(e);
  }
})();
