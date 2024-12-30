import { error } from '@sveltejs/kit';
import { S3Client, GetObjectCommand } from '@aws-sdk/client-s3';
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

export async function GET({ params }) {
  try {
    const command = new GetObjectCommand({
      Bucket: 'slides',
      Key: `${params.name}.pdf`
    });

    const object = await r2.send(command);

    if (!object.Body) {
      throw error(404, 'PDF not found');
    }

    const stream = object.Body.transformToWebStream();

    return new Response(stream, {
      headers: {
        'content-type': 'application/pdf',
      }
    });
  } catch (err) {
    throw error(500, 'Failed to fetch PDF');
  }
}
