import { json } from '@sveltejs/kit';
import type { RequestHandler } from './$types';
import { S3Client, ListObjectsV2Command } from '@aws-sdk/client-s3';
import { PUBLIC_R2_URL } from '$env/static/public';
import { R2_ACCESS_KEY_ID, R2_SECRET_ACCESS_KEY } from '$env/static/private';

const r2 = new S3Client({
  region: 'auto',
  endpoint: PUBLIC_R2_URL,
  credentials: {
    accessKeyId: R2_ACCESS_KEY_ID,
    secretAccessKey: R2_SECRET_ACCESS_KEY
  }
});

export const GET: RequestHandler = async () => {
  try {
      const command = new ListObjectsV2Command({
          Bucket: 'slides'
      });

      const response = await r2.send(command);

      const files = (response.Contents || []).map(file => ({
          name: file.Key,
          size: file.Size,
          lastModified: file.LastModified,
      }));

      return json(files);

  } catch (error) {
      return json(
          { error: 'Failed to fetch slides' },
          { status: 500 }
      );
  }
};
