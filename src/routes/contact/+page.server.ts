import { createTransport } from "nodemailer";
import {
  SMTP_USERNAME,
  SMTP_TOKEN,
  SMTP_SERVER,
  SMTP_PORT,
  MY_EMAIL_ADDRESS,
  MY_BCC_EMAIL_ADDRESS
} from "$env/static/private";
import type { Actions } from "./$types"
import { fail } from "@sveltejs/kit";

export const actions: Actions = {
  default: async ({ request }) => {
    const data = await request.formData();
    const transporter = createTransport({
      host: SMTP_SERVER,
      port: Number(SMTP_PORT),
      requireTLS: true,
      secure: false,
      auth: {
        user: SMTP_USERNAME,
        pass: SMTP_TOKEN
      }
    });
    try {
      await transporter.sendMail({
        from: `"Hiroto Sasagawa" <${MY_EMAIL_ADDRESS}>`,
        to: data.get("email") as string,
        bcc: MY_BCC_EMAIL_ADDRESS,
        subject: `お問い合わせを受け付けました`,
        html: `<!DOCTYPE HTML><html><p>お問い合わせ内容は以下の通りです。</p><ul><li>氏名: ${data.get("name")}</li><li>メールアドレス: ${data.get("email")}</li><li>本文: ${data.get("text")}</li></ul><p>返信まで数日かかる場合がございます。予めご了承ください。</p></html>`,
        list: {
          help: `${SMTP_USERNAME}?subject=help`,
          unsubscribe: `${SMTP_USERNAME}?subject=unsubscribe`,
          subscribe: `${SMTP_USERNAME}?subject=subscribe`
        }
      })
      return {
        name: data.get("name"),
        email: data.get("email"),
        text: data.get("text"),
      }
    } catch {
      return fail(500);
    }
  }
};

export const prerender = false;
