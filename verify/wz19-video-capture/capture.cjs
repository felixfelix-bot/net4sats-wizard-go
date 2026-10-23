// WZ1.9 E2E capture — drives the wizard UI through the real happy path.
// Records one continuous take at 1920x1080. No cuts, no mocks.
//
// Project root for playwright is the npx cache; run via capture.sh.
const { chromium } = require('/home/c08r4d0r/.npm/_npx/c61c9351a0dbcfa7/node_modules/playwright');

const WIZARD = process.env.WIZARD_URL || 'http://127.0.0.1:8099/';
const OUT = process.env.OUT || '/tmp/wz19-take';
const ROUTER_IP = process.env.ROUTER_IP || '192.168.2.30';
const ROUTER_PW = process.env.ROUTER_PASSWORD || '';
const LNURL = process.env.LNURL || 'demo@getalby.com';

(async () => {
  const browser = await chromium.launch({
    headless: false,
    executablePath: '/home/c08r4d0r/.cache/ms-playwright/chromium-1248/chrome-linux64/chrome',
    args: ['--no-sandbox', '--disable-dev-shm-usage', '--window-size=1920,1080'],
  });
  const ctx = await browser.newContext({
    viewport: { width: 1920, height: 1080 },
    recordVideo: { dir: OUT, size: { width: 1920, height: 1080 } },
    deviceScaleFactor: 1,
  });
  const page = await ctx.newPage();
  const stepLog = [];
  const stamp = (m) => { const l = `${new Date().toISOString()} ${m}`; stepLog.push(l); console.log(l); };

  try {
    // ── 1. START ────────────────────────────────────────────────
    stamp('STEP 1 start — navigate to wizard root');
    await page.goto(WIZARD, { waitUntil: 'domcontentloaded', timeout: 30000 });
    await page.waitForTimeout(1500);

    // ── 2. DISCOVERY ────────────────────────────────────────────
    stamp('STEP 2 discovery — waiting for /api/scan to populate router select');
    await page.waitForFunction(
      () => {
        const s = document.querySelector('#router-select');
        return s && s.options.length > 0 && !/scanning|searching/i.test(s.options[0].textContent || '');
      },
      null,
      { timeout: 90000 },
    );
    const opts = await page.$$eval('#router-select option', (els) =>
      els.map((e) => ({ value: e.value, text: e.textContent.trim() })));
    stamp('STEP 2 discovery — routers offered: ' + JSON.stringify(opts));

    // ── 3. SELECT ROUTER + PASSWORD ─────────────────────────────
    stamp('STEP 3 config — select router ' + ROUTER_IP);
    await page.selectOption('#router-select', ROUTER_IP);
    await page.waitForTimeout(800);
    stamp('STEP 3 config — enter router password');
    await page.fill('#password', ROUTER_PW);
    await page.waitForTimeout(500);

    // ── 4. UPSTREAM MODE (default WAN) ──────────────────────────
    stamp('STEP 4 config — upstream mode WAN (default, visible)');
    await page.waitForTimeout(700);

    // ── 5. LIGHTNING ADDRESS ────────────────────────────────────
    stamp('STEP 5 config — enter Lightning payout address');
    await page.fill('#lnurl', LNURL);
    await page.waitForTimeout(700);

    // ── 6. ADVANCED (optional, opened to prove it is not required) ──
    stamp('STEP 6 advanced — open <details> to show optional settings');
    await page.click('details.advanced summary').catch(() => {});
    await page.waitForTimeout(1200);
    await page.click('details.advanced summary').catch(() => {});
    await page.waitForTimeout(600);

    const btnState = await page.$eval('#deploy-btn', (b) => ({
      disabled: b.disabled, text: b.textContent.trim(),
    }));
    stamp('STEP 6 ready-state — deploy button ' + JSON.stringify(btnState));

    // ── 7. DEPLOY (the real thing) ──────────────────────────────
    stamp('STEP 7 deploy — clicking deploy button');
    await page.click('#deploy-btn', { timeout: 15000 });
    await page.waitForTimeout(2000);
    await page.screenshot({ path: OUT + '/shot-deploy-start.png' });

    // Watch the deploy for up to 150s, capturing the step list as it advances.
    let lastSteps = '';
    let sawSuccess = false;
    let sawError = false;
    for (let i = 0; i < 75; i++) {
      await page.waitForTimeout(2000);
      const state = await page.evaluate(() => {
        const vis = (id) => {
          const el = document.querySelector(id);
          return el && !el.classList.contains('hidden');
        };
        const steps = document.querySelector('#steps-list');
        const logEl = document.querySelector('#deploy-log');
        return {
          success: vis('#success-view'),
          error: vis('#error-view'),
          errorDetail: (document.querySelector('#error-detail') || {}).textContent || '',
          steps: steps ? steps.innerText : '',
          log: logEl ? logEl.innerText.slice(-1200) : '',
        };
      });
      if (state.steps !== lastSteps) {
        lastSteps = state.steps;
        stamp('STEP 7 deploy-progress — ' + JSON.stringify(state.steps.split('\n').filter(Boolean)));
      }
      if (state.error) {
        sawError = true;
        await page.screenshot({ path: OUT + '/shot-error.png' });
        stamp('STEP 7 deploy-result — ERROR VIEW SHOWN. detail=' + JSON.stringify(state.errorDetail));
        stamp('STEP 7 deploy-log tail — ' + JSON.stringify(state.log));
        break;
      }
      if (state.success) {
        sawSuccess = true;
        await page.waitForTimeout(3500);
        await page.screenshot({ path: OUT + '/shot-success.png' });
        stamp('STEP 8 success — success view visible');
        break;
      }
    }
    if (!sawSuccess && !sawError) {
      await page.screenshot({ path: OUT + '/shot-timeout.png' });
      stamp('STEP 7 deploy-result — no terminal state within watch window (still running)');
      const st = await page.evaluate(() => {
        const el = document.querySelector('#deploy-log');
        return { steps: (document.querySelector('#steps-list') || {}).innerText || '',
                 log: el ? el.innerText.slice(-2000) : '' };
      });
      stamp('STEP 7 deploy-log tail — ' + JSON.stringify(st));
    }

    // hold the final frame so the outcome is unambiguously readable
    await page.waitForTimeout(4000);
    stamp('END take — sawSuccess=' + sawSuccess + ' sawError=' + sawError);
  } catch (e) {
    stamp('EXCEPTION — ' + e.message);
    try { await page.screenshot({ path: OUT + '/shot-exception.png' }); } catch (_) {}
    await page.waitForTimeout(3000);
  } finally {
    require('fs').writeFileSync(OUT + '/step-log.txt', stepLog.join('\n') + '\n');
    await ctx.close();
    await browser.close();
  }
})();
