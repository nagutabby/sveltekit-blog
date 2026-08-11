<script lang="ts">
  import { onMount } from "svelte";

  const { headingList } = $props();
  let rawToc = $state("");

  onMount(async () => {
    const toc = headingList.map((heading: Element) => ({
      text: heading.textContent,
      id: heading.getAttribute("id"),
      name: heading.tagName,
    }));
    let previousTag: string = "";
    for (let i = 0; i < toc.length; i++) {
      if (previousTag === "") {
        rawToc += `<ul><li><a href="#${toc[i].id}" class="link !link-secondary">${toc[i].text}</a>`;
      } else if (toc[i].name === previousTag) {
        rawToc += `</li><li><a href="#${toc[i].id}" class="link !link-secondary">${toc[i].text}</a>`;
      } else if (previousTag === "H1") {
        if (toc[i].name === "H2" || "H3") {
          rawToc += `<ul><li><a href="#${toc[i].id}" class="link !link-secondary">${toc[i].text}</a>`;
        }
      } else if (previousTag === "H2") {
        if (toc[i].name === "H3") {
          rawToc += `<ul><li><a href="#${toc[i].id}" class="link !link-secondary">${toc[i].text}</a>`;
        } else if (toc[i].name === "H1") {
          rawToc += `</li></ul></li><li><a href="#${toc[i].id}" class="link !link-secondary">${toc[i].text}</a>`;
        }
      } else if (previousTag === "H3") {
        if (toc[i].name === "H2") {
          rawToc += `</li></ul></li><li><a href="#${toc[i].id}" class="link !link-secondary">${toc[i].text}</a></li>`;
        } else if (toc[i].name === "H1") {
          rawToc += `</li></ul></li></ul></li><li><a href="#${toc[i].id}" class="link !link-secondary">${toc[i].text}</a>`;
        }
      }
      previousTag = toc[i].name;
    }

    rawToc += "</li>";

    if (previousTag === "H2") {
      rawToc += "</ul></li>";
    } else if (previousTag === "H3") {
      rawToc += "</ul></li></ul></li>";
    }

    rawToc += "</ul>";

    let activeHeading: any;
    const updateElements = (entries: any) => {
      entries.forEach((entry: any) => {
        if (entry.isIntersecting) {
          if (activeHeading !== undefined && activeHeading !== null) {
            activeHeading.classList.remove("active");
          }
          activeHeading = document.querySelector(
            `a[href="#${entry.target.id}"]`,
          );

          if (activeHeading !== null) {
            activeHeading.classList.add("active");
          }
        }
      });
    };
    const headings = document.querySelectorAll("h1, h2, h3");
    const options = {
      rootMargin: "0% 0% -100% 0%",
      threshold: 0,
    };
    const observer = new IntersectionObserver(updateElements, options);
    headings.forEach((heading: any) => {
      observer.observe(heading);
    });
  });
</script>

{@html rawToc}
