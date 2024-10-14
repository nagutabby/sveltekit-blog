/*
 * Modal
 *
 * Pico.css - https://picocss.com
 * Copyright 2019-2023 - Licensed under MIT
 */

const isOpenClass = "modal-is-open";
const openingClass = "modal-is-opening";
const closingClass = "modal-is-closing";
const animationDuration = 400; // ms
let visibleModal: Element | null = null;

export const toggleModal = (event: Event) => {

  event.preventDefault();

  const modal = document.getElementById(
    (event.currentTarget as any).getAttribute(
      "data-target"
    )!
  );


  typeof modal !== "undefined" && modal !== null && isModalOpen(modal)
    ? closeModal(modal)
    : openModal(modal!);
}

const isModalOpen = (modal: Element) => {
  return modal.hasAttribute("open") && modal.getAttribute("open") !== "false"
    ? true
    : false;
};

const openModal = (modal: Element) => {
  if (isScrollbarVisible()) {
    document.documentElement.style.setProperty(
      "--scrollbar-width",
      `${getScrollbarWidth()}px`
    );
  }
  document.documentElement.classList.add(isOpenClass, openingClass);
  setTimeout(() => {
    visibleModal = modal;
    document.documentElement.classList.remove(openingClass);
  }, animationDuration);
  modal.setAttribute("open", "true");
};

const closeModal = (modal: Element) => {
  visibleModal = null;
  document.documentElement.classList.add(closingClass);
  setTimeout(() => {
    document.documentElement.classList.remove(closingClass, isOpenClass);
    document.documentElement.style.removeProperty("--scrollbar-width");
    modal.removeAttribute("open");
  }, animationDuration);
};

const getScrollbarWidth = () => {
  // Creating invisible container
  const outer = document.createElement("div");
  outer.style.visibility = "hidden";
  outer.style.overflow = "scroll"; // forcing scrollbar to appear
  outer.style.msOverflowStyle = "scrollbar"; // needed for WinJS apps
  document.body.appendChild(outer);

  // Creating inner element and placing it in the container
  const inner = document.createElement("div");
  outer.appendChild(inner);

  // Calculating difference between container's full width and the child width
  const scrollbarWidth = outer.offsetWidth - inner.offsetWidth;

  // Removing temporary elements from the DOM
  outer!.parentNode!.removeChild(outer);

  return scrollbarWidth;
};

const isScrollbarVisible = () => {
  return document.body.scrollHeight > screen.height;
};


export const closeWithClickOutside = () => {
  document.addEventListener("click", (event) => {
    if (visibleModal !== null) {
      const modalContent = visibleModal.querySelector("article");
      const isClickInside = modalContent?.contains(
        event.target as HTMLInputElement
      );
      !isClickInside && closeModal(visibleModal);
    }
  });
}

export const closeWithEscapeKey = () => {
  document.addEventListener("keydown", (event) => {
    if (event.key === "Escape" && visibleModal !== null) {
      closeModal(visibleModal);
    }
  });
}
