import { hc } from 'hono/client';
import type { RpcAppType } from '../../worker/rpc';

const client = hc<RpcAppType>('/rpc');

export function createBackendClient() {
  return {
    submitContact: async (message: { name: string; email: string; text: string; imRobot: boolean }) => {
      const response = await client.contact.submit.$post({ json: message });
      if (!response.ok) throw new Error('Contact submission failed');
      return response.json();
    }
  };
}
