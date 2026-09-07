// スライドHTML(1920x1080pxのキャンバスに収まっているか)を実ブラウザで検証するスクリプト。
// budoux-jaの改行処理・実際のフォントメトリクスを反映した実測値で判定するため、
// puppeteer-core(ローカルのGoogle Chrome/Chromiumを操作、バンドルDLなし)で実描画する。
//
// 使い方:
//   node web/scripts/validate-slide-layout.mjs [path/to/slide.html]
//   (引数省略時は web/static/content/templates/slide.html を検証する)
//
// 前提: ローカルにGoogle Chrome(またはChromium)がインストールされていること。
//   別の場所にある場合は PUPPETEER_EXECUTABLE_PATH で明示する。

import puppeteer from 'puppeteer-core';
import { existsSync } from 'node:fs';
import { resolve } from 'node:path';

const OVERFLOW_TOLERANCE_PX = 1; // サブピクセルの丸め誤差を吸収する

function findChromeExecutable() {
  const candidates = [
    process.env.PUPPETEER_EXECUTABLE_PATH,
    '/Applications/Google Chrome.app/Contents/MacOS/Google Chrome',
    '/Applications/Chromium.app/Contents/MacOS/Chromium',
    '/usr/bin/google-chrome',
    '/usr/bin/google-chrome-stable',
    '/usr/bin/chromium',
    '/usr/bin/chromium-browser',
  ].filter(Boolean);
  return candidates.find((path) => existsSync(path)) ?? null;
}

async function main() {
  const target = process.argv[2] ?? 'web/static/content/templates/slide.html';
  const absPath = resolve(target);

  if (!existsSync(absPath)) {
    console.error(`NG: ファイルが見つかりません: ${absPath}`);
    process.exitCode = 1;
    return;
  }

  const executablePath = findChromeExecutable();
  if (!executablePath) {
    console.error(
      'NG: ローカルにGoogle Chrome/Chromiumが見つかりません。' +
        'PUPPETEER_EXECUTABLE_PATH環境変数でパスを指定してください。'
    );
    process.exitCode = 1;
    return;
  }

  const browser = await puppeteer.launch({ executablePath, headless: true });
  try {
    const page = await browser.newPage();
    await page.setViewport({ width: 1920, height: 1080 });
    await page.goto(`file://${absPath}`, { waitUntil: 'load' });

    const slides = await page.evaluate((tolerance) => {
      return Array.from(document.querySelectorAll('.slide')).map((slide, i) => {
        const rect = slide.getBoundingClientRect();
        const overflowX = slide.scrollWidth - slide.clientWidth;
        const overflowY = slide.scrollHeight - slide.clientHeight;
        return {
          index: i + 1,
          renderedWidth: rect.width,
          renderedHeight: rect.height,
          overflowX,
          overflowY,
          overflows: overflowX > tolerance || overflowY > tolerance,
        };
      });
    }, OVERFLOW_TOLERANCE_PX);

    if (slides.length === 0) {
      console.error('NG: .slide 要素が1つも見つかりません');
      process.exitCode = 1;
      return;
    }

    let hasError = false;

    for (const s of slides) {
      const sizeOk = Math.abs(s.renderedWidth - 1920) <= OVERFLOW_TOLERANCE_PX
        && Math.abs(s.renderedHeight - 1080) <= OVERFLOW_TOLERANCE_PX;

      if (!sizeOk) {
        console.error(
          `NG: ${s.index}枚目 が1920x1080ではありません(実測 ${s.renderedWidth}x${s.renderedHeight}px)`
        );
        hasError = true;
        continue;
      }

      if (s.overflows) {
        console.error(
          `NG: ${s.index}枚目 の内容がキャンバスからはみ出しています` +
            `(縦の超過: ${s.overflowY}px, 横の超過: ${s.overflowX}px)`
        );
        hasError = true;
      } else {
        console.log(`OK: ${s.index}枚目`);
      }
    }

    if (hasError) {
      console.error(`NG: ${slides.length}枚中、レイアウトの制約に違反したスライドがあります`);
      process.exitCode = 1;
    } else {
      console.log(`OK: 全${slides.length}枚が1920x1080に収まっています`);
    }
  } finally {
    await browser.close();
  }
}

main().catch((err) => {
  console.error('NG: 検証中にエラーが発生しました');
  console.error(err);
  process.exitCode = 1;
});
