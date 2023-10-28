import { createTransport } from "nodemailer";
import {
  SMTP_USERNAME,
  SMTP_TOKEN,
  SMTP_SERVER,
  SMTP_PORT,
  MY_CC_EMAIL_ADDRESS
} from "$env/static/private";
import type { Actions } from "./$types"

export const actions: Actions = {
  default: async ({ request }) => {
    const data = await request.formData();
    let isFailed = false;
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
    transporter.sendMail({
      from: `"Hiroto Sasagawa" <${SMTP_USERNAME}>`,
      to: data.get("email") as string,
      cc: MY_CC_EMAIL_ADDRESS,
      subject: `お問い合わせを受け付けました`,
      html: `<!DOCTYPE HTML><html><p>お問い合わせ内容は以下の通りです。</p><ul><li>氏名: ${data.get("name")}</li><li>メールアドレス: ${data.get("email")}</li><li>本文: ${data.get("text")}</li></ul><p>返信まで数日かかる場合がございます。予めご了承ください。</p></html>`,
      list: {
        help: `${SMTP_USERNAME}?subject=help`,
        unsubscribe: `${SMTP_USERNAME}?subject=unsubscribe`,
        subscribe: `${SMTP_USERNAME}?subject=subscribe`
      }
    }).catch(() => {
      isFailed = true;
    })
    return {
      name: data.get("name"),
      email: data.get("email"),
      text: data.get("text"),
      isFailed: isFailed
    }
  }
};
