<script lang="ts">
  import { onMount } from "svelte";

  let rawToc: string = "";

  onMount(async () => {
    let headings: any = Array.from(
      document.getElementById("content")!.querySelectorAll("h1, h2, h3"),
    );
    const toc = headings.map((heading: Element) => ({
      text: heading.textContent,
      id: heading.getAttribute("id"),
      name: heading.tagName,
    }));
    let previousTag: string = "";
    for (let i = 0; i < toc.length; i++) {
      if (previousTag === "") {
        rawToc += `<ul><li><a href="#${toc[i].id}">${toc[i].text}</a></li>`;
      } else if (toc[i].name === previousTag) {
        rawToc += `<li><a href="#${toc[i].id}">${toc[i].text}</a></li>`;
      } else if (previousTag === "H1") {
        if (toc[i].name === "H2" || "H3") {
          rawToc += `<ul><li><a href="#${toc[i].id}">${toc[i].text}</a></li>`;
        }
      } else if (previousTag === "H2") {
        if (toc[i].name === "H3") {
          rawToc += `<ul><li><a href="#${toc[i].id}">${toc[i].text}</a></li>`;
        } else if (toc[i].name === "H1") {
          rawToc += `</ul><li><a href="#${toc[i].id}">${toc[i].text}</a></li>`;
        }
      } else if (previousTag === "H3") {
        if (toc[i].name === "H2") {
          rawToc += `</ul><li><a href="#${toc[i].id}">${toc[i].text}</a></li>`;
        } else if (toc[i].name === "H1") {
          rawToc += `</ul></ul><li><a href="#${toc[i].id}">${toc[i].text}</a></li>`;
        }
      }
      previousTag = toc[i].name;
    }
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
    headings = document.querySelectorAll("h1, h2, h3");
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

<article>
  {@html rawToc}
</article>
