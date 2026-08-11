import type { Actions } from "./$types";
import { fail } from "@sveltejs/kit";
import { createBackendClient } from "$lib/server/backendClient";
import { ContactService } from "$lib/gen/blog/contact/v1/contact_pb";

const contactClient = createBackendClient(ContactService);

export const actions: Actions = {
  default: async ({ request }) => {
    const data = await request.formData();
    const imRobot = data.get("im-robot") === "true";
    const name = data.get("name") as string;
    const email = data.get("email") as string;
    const text = data.get("text") as string;

    try {
      const response = await contactClient.submitContact({ name, email, text, imRobot });

      if (Object.keys(response.errors).length > 0) {
        return fail(400, {
          errors: response.errors,
          values: {
            name,
            email,
            text
          }
        });
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
};
