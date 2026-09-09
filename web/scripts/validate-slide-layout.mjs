// スライドHTML(1920x1080pxのキャンバスに収まっているか)を実ブラウザで検証するスクリプト。
// budoux-jaの改行処理・実際のフォントメトリクスを反映した実測値で判定するため、
// puppeteer-core(ローカルのGoogle Chrome/Chromiumを操作、バンドルDLなし)で実描画する。
// 図表はHTML/CSSで直接組むのが基本方針のため非同期描画の待ち合わせは不要。
// 画像を<img>で使う場合も、通常の画像と同様に waitUntil: 'load' で読み込み完了を待てば十分。
// @font-face でWebフォントを同梱している場合はフォント読み込みが load イベントをブロックしない
// ことがあるため、document.fonts.ready を待ってから計測する。
// 見出し(h2)と直後の要素との間隔も全スライドで一致しているかを検証する(margin-bottomと
// 親要素のgapが二重に加算される等でスライドごとに間隔がずれた不具合が過去にあったため)。
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
const GAP_TOLERANCE_PX = 1; // 見出し-本文間隔の許容誤差(サブピクセルの丸め)

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
    await page.evaluate(() => document.fonts.ready);

    const slides = await page.evaluate((tolerance) => {
      return Array.from(document.querySelectorAll('.slide')).map((slide, i) => {
        const rect = slide.getBoundingClientRect();
        const overflowX = slide.scrollWidth - slide.clientWidth;
        const overflowY = slide.scrollHeight - slide.clientHeight;

        // 見出し(h2)と直後の要素との間隔(px)。h2が無いスライド(表紙・高橋メソッド風等)はnull
        let headingGap = null;
        const h2 = slide.querySelector('h2');
        const next = h2 ? h2.nextElementSibling : null;
        if (h2 && next) {
          const h2Rect = h2.getBoundingClientRect();
          const nextRect = next.getBoundingClientRect();
          headingGap = Math.round((nextRect.top - h2Rect.bottom) * 100) / 100;
        }

        return {
          index: i + 1,
          renderedWidth: rect.width,
          renderedHeight: rect.height,
          overflowX,
          overflowY,
          overflows: overflowX > tolerance || overflowY > tolerance,
          headingGap,
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

    // 見出し(h2)と本文の間隔が全スライドで統一されているかを検証する
    const gaps = slides
      .filter((s) => s.headingGap !== null)
      .map((s) => ({ index: s.index, gap: s.headingGap }));

    if (gaps.length > 0) {
      const baseline = gaps[0].gap;
      const inconsistent = gaps.filter((g) => Math.abs(g.gap - baseline) > GAP_TOLERANCE_PX);

      if (inconsistent.length > 0) {
        console.error(`NG: 見出し(h2)と本文の間隔がスライドごとに異なります(基準 ${baseline}px):`);
        for (const g of inconsistent) {
          console.error(`  - ${g.index}枚目: ${g.gap}px`);
        }
        hasError = true;
      } else {
        console.log(`OK: 見出しと本文の間隔は全スライドで統一されています(${baseline}px)`);
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
