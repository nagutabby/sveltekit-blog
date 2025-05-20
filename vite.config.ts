import { defineConfig } from "vitest/config";
import { sveltekit } from "@sveltejs/kit/vite";
import { svelteTesting } from '@testing-library/svelte/vite';
import path from 'path';
import fs from 'fs/promises';
import sharp from 'sharp';

// 画像最適化プラグイン
const imageOptimizer = () => {
  return {
    name: 'sharp-image-optimizer',
    async buildStart() {
      console.log('画像の最適化を開始します...');
      // staticディレクトリの画像を処理
      const staticDir = path.resolve('static');
      await processDirectory(staticDir, staticDir);
      console.log('画像の最適化が完了しました');
    }
  };
};

// ディレクトリを再帰的に処理する関数
async function processDirectory(baseDir: string, currentDir: string) {
  try {
    const entries = await fs.readdir(currentDir, { withFileTypes: true });

    for (const entry of entries) {
      const fullPath = path.join(currentDir, entry.name);

      if (entry.isDirectory()) {
        // ディレクトリの場合は再帰的に処理
        await processDirectory(baseDir, fullPath);
      } else if (isImageFile(entry.name)) {
        // 画像ファイルの場合は最適化
        await optimizeImage(fullPath, baseDir);
      }
    }
  } catch (error) {
    console.error(`ディレクトリ処理中にエラーが発生しました: ${currentDir}`, error);
  }
}

// 画像ファイルかどうか判定する関数
function isImageFile(filename: string) {
  const ext = path.extname(filename).toLowerCase();
  return ['.jpg', '.jpeg', '.png', '.gif'].includes(ext);
}

// 画像を最適化する関数
async function optimizeImage(filePath: string, baseDir: string) {
  try {
    const relativePath = path.relative(baseDir, filePath);
    const dir = path.dirname(filePath);
    const filename = path.basename(filePath, path.extname(filePath));
    const outputPath = path.join(dir, `${filename}.webp`);

    // すでに最適化されたファイルがあるか確認
    try {
      await fs.access(outputPath);
      // ソースファイルの更新日時を取得
      const srcStat = await fs.stat(filePath);
      const destStat = await fs.stat(outputPath);

      // ソースが最適化済みファイルより新しい場合のみ最適化
      if (srcStat.mtime <= destStat.mtime) {
        console.log(`スキップ (既に最適化済み): ${relativePath}`);
        return;
      }
    } catch (e) {
      // ファイルが存在しない場合は続行
    }

    // sharpで画像を最適化
    await sharp(filePath)
      .webp({ quality: 80 }) // 品質調整
      .toFile(outputPath);

    console.log(`最適化完了: ${relativePath} → ${path.relative(baseDir, outputPath)}`);
  } catch (error) {
    console.error(`画像最適化エラーが発生しました: ${filePath}`, error);
  }
}


export default defineConfig({
  plugins: [sveltekit(), svelteTesting(), imageOptimizer()],

  build: {
    target: "ES2020",
  },
  test: {
    include: ['src/**/*.{test,spec}.{js,ts}'],
    environment: 'jsdom'
  }
});
