// スライドHTMLを1920x1080pxのPDFとして書き出すスクリプト。
// puppeteer-core(ローカルのGoogle Chrome/Chromiumを操作、バンドルDLなし)でヘッドレス印刷する。
// web/static/content/templates/slide.html の @media print / @page { size: 1920px 1080px; margin: 0 }
// を前提にしており、preferCSSPageSize: true でそのページサイズをそのまま採用するため、
// 縮小・再エンコードによる画質劣化なしに1920x1080pxのPDFページを生成する
// (.slideごとに break-after: page で改ページされ、1スライド=1ページになる)。
//
// 注記: 生成されるPDFの/MediaBoxは[0 0 1440 810]になる(1920x1080ではない)。
// これは画質劣化ではなく、CSSのpx(1/96インチ)とPDFのpt(1/72インチ)の単位換算による
// 正しい値(1920px×72/96=1440pt、1080px×72/96=810pt、いずれも20x11.25インチ相当)。
// Chromiumの印刷パイプラインはこの96→72換算を必ず行うため、MediaBoxの数値を
// そのまま1920x1080にしようとして用紙サイズだけを拡大すると、スライド本体が
// ページ左上に小さく収まり余白ができる不具合が生じる(zoomやscaleオプションでの
// 拡大は印刷時のレイアウトに反映されないため補正できない)。数値上1920x1080に
// 揃えたい場合は、qpdf/pikepdf等の外部ツールでPDFのMediaBoxとページ内容の
// 変換行列を書き換える後処理が必要(このスクリプトでは行わない)。
//
// 使い方:
//   node web/scripts/export-slide-pdf.mjs <path/to/slide.html> [path/to/output.pdf]
//   (出力先省略時は入力と同じディレクトリに拡張子だけ.pdfへ変えたパスに書き出す。
//    web/static/content/slides/<id>.html → web/static/content/slides/<id>.pdf)
//
// 前提: ローカルにGoogle Chrome(またはChromium)がインストールされていること。
//   別の場所にある場合は PUPPETEER_EXECUTABLE_PATH で明示する。
// レイアウトのはみ出し検証は行わない(web/scripts/validate-slide-layout.mjsの責務)。
// 事前にそちらでOKになっていることを確認してから実行する。

import puppeteer from 'puppeteer-core';
import { existsSync } from 'node:fs';
import { resolve, dirname, basename, extname, join } from 'node:path';

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
  const input = process.argv[2];
  if (!input) {
    console.error('Usage: node web/scripts/export-slide-pdf.mjs <path/to/slide.html> [path/to/output.pdf]');
    process.exitCode = 2;
    return;
  }

  const absInput = resolve(input);
  if (!existsSync(absInput)) {
    console.error(`NG: ファイルが見つかりません: ${absInput}`);
    process.exitCode = 1;
    return;
  }

  const outputArg = process.argv[3];
  const absOutput = outputArg
    ? resolve(outputArg)
    : join(dirname(absInput), `${basename(absInput, extname(absInput))}.pdf`);

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
    await page.goto(`file://${absInput}`, { waitUntil: 'load' });
    await page.evaluate(() => document.fonts.ready);

    const slideCount = await page.evaluate(() => document.querySelectorAll('.slide').length);
    if (slideCount === 0) {
      console.error('NG: .slide 要素が1つも見つかりません');
      process.exitCode = 1;
      return;
    }

    await page.emulateMediaType('print');
    await page.pdf({
      path: absOutput,
      printBackground: true,
      preferCSSPageSize: true,
    });

    console.log(`OK: ${slideCount}枚のスライドを書き出しました: ${absOutput}`);
  } finally {
    await browser.close();
  }
}

main().catch((err) => {
  console.error('NG: PDF書き出し中にエラーが発生しました');
  console.error(err);
  process.exitCode = 1;
});
