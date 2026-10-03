import { readdir, readFile, writeFile } from 'node:fs/promises';
import path from 'node:path';
import { fileURLToPath } from 'node:url';
import matter from 'gray-matter';

const scriptDir = path.dirname(fileURLToPath(import.meta.url));
const repoRoot = path.resolve(scriptDir, '../..');
const outputPath = path.resolve(scriptDir, '../src/worker/articles.generated.json');

/** @param {string} filename */
function dateFromFilename(filename) {
  const match = /^(\d{4}-\d{2}-\d{2})(?:-\d+)?\.md$/.exec(filename);
  if (!match) {
    throw new Error(`Content filename must use YYYY-MM-DD[-N].md: ${filename}`);
  }

  const date = new Date(`${match[1]}T00:00:00.000Z`);
  if (Number.isNaN(date.valueOf()) || date.toISOString().slice(0, 10) !== match[1]) {
    throw new Error(`Content filename has an invalid date: ${filename}`);
  }
  return date.toISOString();
}

/** @param {unknown} value */
function stringField(value) {
  if (typeof value === 'string') return value;
  if (value === null || value === undefined) return '';
  return String(value);
}

async function loadArticles() {
  const directory = path.join(repoRoot, 'backend/content/articles');
  const filenames = (await readdir(directory)).filter((name) => name.endsWith('.md')).sort();
  const articles = [];
  const articleIDs = new Set();

  for (const filename of filenames) {
    const source = await readFile(path.join(directory, filename), 'utf8');
    const { data } = matter(source);
    // Match Astro's fail-closed default: a missing or malformed is_draft
    // field must never make an article public by accident.
    if (data.is_draft !== false) continue;

    const id = stringField(data.id);
    const title = stringField(data.title);
    if (!id.trim() || !title.trim()) {
      throw new Error(`Published article must have non-empty id and title: ${filename}`);
    }
    if (articleIDs.has(id)) throw new Error(`Duplicate published article id "${id}": ${filename}`);
    articleIDs.add(id);

    articles.push({
      id,
      title,
      publishedAt: dateFromFilename(filename)
    });
  }

  // Filenames are read in lexical order before the stable newest-first date
  // sort. Keep that ordering for same-day suffixes.
  articles.sort((a, b) => b.publishedAt.localeCompare(a.publishedAt));
  return articles;
}

const articles = await loadArticles();
await writeFile(outputPath, `${JSON.stringify(articles, null, 2)}\n`);
console.log(`Generated ActivityPub metadata for ${articles.length} published articles.`);
