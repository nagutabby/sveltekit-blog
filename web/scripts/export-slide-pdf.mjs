// スライドHTMLを1920x1080pxのPDFとして書き出すスクリプト。
// puppeteer-core(ローカルのGoogle Chrome/Chromiumを操作、バンドルDLなし)でヘッドレス印刷する。
// web/static/content/templates/slide.html の @media print / @page { size: 1920px 1080px; margin: 0 }
// を前提にしており、preferCSSPageSize: true でそのページサイズをそのまま採用するため、
// 縮小・再エンコードによる画質劣化なしに1920x1080px相当のPDFページを生成する
// (.slideごとに break-after: page で改ページされ、1スライド=1ページになる)。
//
// Chromiumの印刷パイプラインは、CSSのpx(1/96インチ)をPDFのpt(1/72インチ)に変換する際
// 必ず96→72換算を行うため、Puppeteer単体では/MediaBoxが[0 0 1440 810]になる
// (1920x1080ではない。1920px×72/96=1440pt、1080px×72/96=810ptで20x11.25インチ相当、
// 換算として正しい値であり画質劣化ではない)。
// Speaker Deck等、PDFの/MediaBoxをそのまま解像度とみなして判定するサービスに
// アップロードする際にこの数値が問題になるため、Ghostscript(gs)で後処理し、
// /MediaBoxを[0 0 1920 1080]に、ページ内容をそれに合わせて拡大するスケール変換
// (-dPDFFitPage -dFIXEDMEDIA)を焼き込む。ベクター文字・図形はそのまま拡大されるだけで
// ラスタライズ・再エンコードは発生しないため画質は劣化しない(確認済み)。
// ローカルにGhostscriptが無い場合は、この後処理をスキップしてPuppeteer出力
// (MediaBox 1440x810、内容自体は正しい)をそのまま採用する。
//
// 使い方:
//   node web/scripts/export-slide-pdf.mjs <path/to/slide.html> [path/to/output.pdf]
//   (出力先省略時は入力と同じディレクトリに拡張子だけ.pdfへ変えたパスに書き出す。
//    web/static/content/slides/<id>.html → web/static/content/slides/<id>.pdf)
//
// 前提: ローカルにGoogle Chrome(またはChromium)がインストールされていること。
//   別の場所にある場合は PUPPETEER_EXECUTABLE_PATH で明示する。
//   Ghostscript(`brew install ghostscript`)が無くても書き出し自体は可能(注記参照)。
// レイアウトのはみ出し検証は行わない(web/scripts/validate-slide-layout.mjsの責務)。
// 事前にそちらでOKになっていることを確認してから実行する。

import puppeteer from 'puppeteer-core';
import { existsSync, renameSync, unlinkSync } from 'node:fs';
import { resolve, dirname, basename, extname, join } from 'node:path';
import { spawnSync } from 'node:child_process';

const SLIDE_WIDTH_PX = 1920;
const SLIDE_HEIGHT_PX = 1080;

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

function findGhostscriptExecutable() {
  const candidates = [
    process.env.GHOSTSCRIPT_EXECUTABLE_PATH,
    '/opt/homebrew/bin/gs',
    '/usr/local/bin/gs',
    '/usr/bin/gs',
  ].filter(Boolean);
  const found = candidates.find((path) => existsSync(path));
  if (found) return found;
  const probe = spawnSync('gs', ['--version']);
  return probe.error ? null : 'gs';
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

  const rawOutput = `${absOutput}.raw.pdf`;

  const browser = await puppeteer.launch({ executablePath, headless: true });
  try {
    const page = await browser.newPage();
    await page.setViewport({ width: SLIDE_WIDTH_PX, height: SLIDE_HEIGHT_PX });
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
      path: rawOutput,
      printBackground: true,
      preferCSSPageSize: true,
    });

    const gsPath = findGhostscriptExecutable();
    if (!gsPath) {
      renameSync(rawOutput, absOutput);
      console.log(`OK: ${slideCount}枚のスライドを書き出しました: ${absOutput}`);
      console.log(
        'NG: Ghostscriptが見つからないため、/MediaBoxは[0 0 1440 810]のままです' +
          '(内容自体は正しく1920x1080px相当)。`brew install ghostscript`後に再実行すると' +
          '/MediaBoxを[0 0 1920 1080]に補正できます。'
      );
      return;
    }

    const gsResult = spawnSync(gsPath, [
      '-o', absOutput,
      '-sDEVICE=pdfwrite',
      `-dDEVICEWIDTHPOINTS=${SLIDE_WIDTH_PX}`,
      `-dDEVICEHEIGHTPOINTS=${SLIDE_HEIGHT_PX}`,
      '-dPDFFitPage',
      '-dFIXEDMEDIA',
      rawOutput,
    ]);
    unlinkSync(rawOutput);

    if (gsResult.status !== 0) {
      console.error('NG: GhostscriptによるMediaBox補正に失敗しました');
      console.error(gsResult.stderr?.toString() ?? gsResult.error);
      process.exitCode = 1;
      return;
    }

    console.log(`OK: ${slideCount}枚のスライドを書き出しました: ${absOutput}`);
    console.log(`OK: /MediaBoxを[0 0 ${SLIDE_WIDTH_PX} ${SLIDE_HEIGHT_PX}]に補正しました`);
  } finally {
    await browser.close();
  }
}

main().catch((err) => {
  console.error('NG: PDF書き出し中にエラーが発生しました');
  console.error(err);
  process.exitCode = 1;
});
