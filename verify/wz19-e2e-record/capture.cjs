#!/usr/bin/env node
/**
 * WZ1.9 capture driver — real UI drive of the wizard happy path.
 *
 * Drives the SERVED wizard page through real clicks/fills (no API calls, no
 * mocks, no stubbing) and records the whole take at 1920x1080 with Playwright.
 * A frame-timestamped step log is written alongside so every claim in the
 * report can be checked against the footage.
 *
 * Env:
 *   WIZARD_URL      default http://127.0.0.1:8199
 *   ROUTER_IP       default 192.168.2.1
 *   ROUTER_PASSWORD default "" (empty)
 *   LNURL           default demo@getalby.com
 *   OUT_DIR         default /tmp/wz19-rec/raw
 *   CHROME          explicit chromium binary path
 */
const fs = require('fs');
const path = require('path');
const { chromium } = require('@playwright/test');

const WIZARD_URL = process.env.WIZARD_URL || 'http://127.0.0.1:8199';
const ROUTER_IP = process.env.ROUTER_IP || '192.168.2.1';
const ROUTER_PASSWORD = process.env.ROUTER_PASSWORD || '';
const LNURL = process.env.LNURL || 'demo@getalby.com';
const OUT_DIR = process.env.OUT_DIR || '/tmp/wz19-rec/raw';
const CHROME = process.env.CHROME ||
  path.join(process.env.HOME, '.cache/ms-playwright/chromium-1217/chrome-linux64/chrome');

fs.mkdirSync(OUT_DIR, { recursive: true });
const t0 = Date.now();
const stepLog = [];
function log(step, detail) {
  const line = `t=+${((Date.now() - t0) / 1000).toFixed(2)}s\t${step}\t${detail || ''}`;
  stepLog.push(line);
  console.log(line);
}
const sleep = (ms) => new Promise((r) => setTimeout(r, ms));

(async () => {
  const browser = await chromium.launch({
    executablePath: CHROME,
    headless: true,
    args: ['--no-sandbox', '--force-device-scale-factor=1', '--hide-scrollbars'],
  });
  const context = await browser.newContext({
    viewport: { width: 1920, height: 1080 },
    recordVideo: { dir: path.join(OUT_DIR, 'video'), size: { width: 1920, height: 1080 } },
    deviceScaleFactor: 1,
  });
  const page = await context.newPage();
  const shots = path.join(OUT_DIR, 'shots');
  fs.mkdirSync(shots, { recursive: true });
  const shot = (n) => page.screenshot({ path: path.join(shots, `${n}.png`) });

  try {
    // ── Step 1: fresh load / Discover ───────────────────────────────
    log('START', `navigating to ${WIZARD_URL}`);
    await page.goto(WIZARD_URL, { waitUntil: 'domcontentloaded' });
    await page.waitForSelector('#scan-view', { state: 'visible', timeout: 15000 });
    await sleep(1200);
    await shot('01-discover-scanning');
    log('DISCOVER', 'scan view visible: "Scanning for routers on your network"');
    log('TITLE', JSON.stringify(await page.title()));
    log('H1', JSON.stringify((await page.textContent('h1')) || ''));

    // wait for the select view (the scan is sequential and slow)
    await page.waitForSelector('#select-view:not(.hidden)', { timeout: 180000 });
    await sleep(900);
    await shot('02-discover-results');
    const options = await page.$$eval('#router-select option', (os) =>
      os.map((o) => ({ value: o.value, label: o.textContent })));
    log('DISCOVER_DONE', `${options.length} routers: ` + JSON.stringify(options));

    // ── Step 2: Connect — pick router + password ────────────────────
    const target = options.find((o) => o.value === ROUTER_IP) || options[0];
    if (!target) throw new Error('no routers discovered — cannot continue');
    await page.selectOption('#router-select', target.value);
    await sleep(1100);
    await shot('03-connect-selected');
    log('CONNECT', `selected router ${target.value}`);
    await page.fill('#password', ROUTER_PASSWORD);
    await sleep(700);
    log('CONNECT_PW', `password field filled (len=${ROUTER_PASSWORD.length})`);

    // ── Step 3: Uplink — WAN ────────────────────────────────────────
    await page.click('#mode-wan');
    await sleep(900);
    await shot('04-uplink-wan');
    const wanSel = await page.getAttribute('#mode-wan', 'class');
    log('UPLINK', `WAN selected (class=${JSON.stringify(wanSel)})`);

    // ── Step 4: Payouts — Lightning address ─────────────────────────
    await page.fill('#lnurl', LNURL);
    await sleep(900);
    await shot('05-payouts');
    log('PAYOUTS', `lnurl=${LNURL}`);

    // ── Step 5: Advanced — open then close (optional, not required) ─
    const adv = await page.$('details.advanced');
    if (adv) {
      await page.click('details.advanced > summary');
      await sleep(900);
      await shot('06-advanced-open');
      const openProp = await page.$eval('details.advanced', (d) => d.open);
      log('ADVANCED_OPEN', `details.advanced open=${openProp}`);
      await page.click('details.advanced > summary');
      await sleep(900);
      const openProp2 = await page.$eval('details.advanced', (d) => d.open);
      log('ADVANCED_CLOSE', `details.advanced open=${openProp2}`);
    } else {
      log('ADVANCED_ABSENT', 'no details.advanced on the page');
    }

    // ── Gating: Deploy enabled only when requirements are met ──────
    const btn = await page.$eval('#deploy-btn', (b) => ({
      disabled: b.disabled, text: b.textContent.trim(),
    }));
    await shot('07-gating-ready');
    log('GATING', JSON.stringify(btn));
    if (btn.disabled) throw new Error('deploy button still disabled — cannot proceed');

    // ── Step 6: Deploy ─────────────────────────────────────────────
    await page.click('#deploy-btn');
    log('DEPLOY_CLICKED', 'deploy button clicked');
    await page.waitForSelector('#deploy-view:not(.hidden)', { timeout: 60000 });
    await sleep(800);
    await shot('08-deploy-steps');
    const stepList = await page.$$eval('#deploy-view li, #deploy-view .step', (ns) =>
      ns.map((n) => n.textContent.replace(/\s+/g, ' ').trim()));
    log('DEPLOY_STEPS', `${stepList.length} steps: ` + JSON.stringify(stepList));

    // ── Step 7: terminal state — success OR the real error ─────────
    const terminal = await page.waitForFunction(() => {
      const ok = document.querySelector('#success-view');
      const err = document.querySelector('#error-view');
      if (ok && !ok.classList.contains('hidden')) return 'success';
      if (err && !err.classList.contains('hidden')) return 'error';
      return null;
    }, null, { timeout: 240000 }).then((h) => h.jsonValue());
    log('TERMINAL', `terminal view = ${terminal}`);
    await sleep(1500);
    await shot(`09-terminal-${terminal}`);
    const terminalText = await page.evaluate(() => {
      const ok = document.querySelector('#success-view');
      const err = document.querySelector('#error-view');
      const el = ok && !ok.classList.contains('hidden') ? ok : err;
      return el ? el.innerText.replace(/\s+/g, ' ').trim() : '(none)';
    });
    log('TERMINAL_TEXT', JSON.stringify(terminalText));

    // final per-step statuses as rendered
    const statuses = await page.$$eval('#deploy-view li, #deploy-view .step', (ns) =>
      ns.map((n) => ({ text: n.textContent.replace(/\s+/g, ' ').trim(), cls: n.className })));
    log('DEPLOY_STATUSES', JSON.stringify(statuses));

    // hold the terminal frame so it is legible in the video
    await sleep(3000);
    log('HOLD_END', 'holding terminal frame for legibility');
    fs.writeFileSync(path.join(OUT_DIR, 'result.json'), JSON.stringify({
      wizard_url: WIZARD_URL, router: target.value, lnurl: LNURL,
      terminal, terminalText, stepList, statuses, options,
      title: await page.title(), h1: (await page.textContent('h1')) || '',
    }, null, 2));
  } catch (e) {
    log('EXCEPTION', e.message.split('\n')[0]);
    try { await shot('99-exception'); } catch (_) {}
    fs.writeFileSync(path.join(OUT_DIR, 'result.json'), JSON.stringify({
      wizard_url: WIZARD_URL, router: ROUTER_IP, lnurl: LNURL, exception: String(e),
    }, null, 2));
    process.exitCode = 2;
  } finally {
    await context.close();
    await browser.close();
    fs.writeFileSync(path.join(OUT_DIR, 'step-log.txt'), stepLog.join('\n') + '\n');
  }
})();
