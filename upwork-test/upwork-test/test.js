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

  const cookies = await context.cookies();

  const auth = cookies.filter(c =>
    [
      "auth_session",
      "master_access_token",
      "master_refresh_token",
      "XSRF-TOKEN",
    ].includes(c.name),
  );

  console.log(JSON.stringify(auth, null, 2));

  await new Promise(() => {});
})();
