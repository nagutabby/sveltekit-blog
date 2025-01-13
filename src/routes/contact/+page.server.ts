import {
  EMAIL_API_TOKEN,
  FROM_ADDRESS,
  BCC_ADDRESS
} from "$env/static/private";
import type { Actions } from "./$types";
import { fail } from "@sveltejs/kit";
import type { PageServerLoad } from "./$types";

export const load: PageServerLoad = async () => {
  const blogData: Blog = {
    image: {
      url: "https://images.microcms-assets.io/assets/99c53a99ae2b4682938f6c435d83e3d9/ca63de19468e45b2833ebf325dbfd56c/Microsoft-Fluentui-Emoji-3d-Cat-3d.1024.png"
    },
    title: "お問い合わせ",
    description: "お気軽にお問い合わせください"
  };
  return blogData;
};

export const actions: Actions = {
  default: async ({ request }) => {
    const data = await request.formData();
    const imRobot = data.get("im-robot") === "true";
    const name = data.get("name") as string;
    const email = data.get("email") as string;
    const text = data.get("text") as string;
    const errors: Record<string, string> = {};

    const emailRegex = /^[a-zA-Z0-9_+-]+(\.[a-zA-Z0-9_+-]+)*@([a-zA-Z0-9][a-zA-Z0-9-]*[a-zA-Z0-9]\.)+[a-zA-Z]{2,}$/;

    if (imRobot) {
      errors.imRobot = "Botによるメッセージ送信はできません";
    }

    if (!name) {
      errors.name = "氏名は必須です";
    }

    if (!email) {
      errors.email = "メールアドレスは必須です";
    } else if (!emailRegex.test(email)) {
      errors.email = "メールアドレスの形式が不適切です";
    }

    if (!text) {
      errors.text = "本文は必須です";
    }

    if (Object.keys(errors).length > 0) {
      return fail(400, {
        errors,
        values: {
          name,
          email,
          text
        }
      });
    } else {
      type EmailPayload = {
        from: {
          email: string;
          name: string;
        };
        to: {
          email: string;
          name: string;
        }[];
        bcc: {
          email: string;
          name: string;
        }[];
        subject: string;
        html: string;
      };
      const payload: EmailPayload = {
        from: {
          email: FROM_ADDRESS,
          name: "Hiroto Sasagawa"
        },
        to: [{
          email,
          name,
        }],
        bcc: [{
          email: BCC_ADDRESS,
          name: "Hiroto Sasagawa",
        }],
        subject: "お問い合わせを受け付けました",
        html: `<!DOCTYPE HTML><html><p>お問い合わせ内容は以下の通りです。</p><ul><li>氏名: ${name}</li><li>メールアドレス: ${email}</li><li>本文: ${text}</li></ul><p>返信まで数日かかる場合がございます。予めご了承ください。</p></html>`,
      };
      try {
        const response = await fetch("https://send.api.mailtrap.io/api/send", {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            "Api-Token": EMAIL_API_TOKEN
          },
          body: JSON.stringify(payload)
        });
        if (!response.ok) {
          const error = await response.json();
          throw new Error(`Mailtrap API error: ${JSON.stringify(error)}`);
        }
        return {
          name,
          email,
          text,
        };
      } catch {
        return fail(500);
      }
    }
  }
};
