export function hideElement(element) {
  addCss(element, 'hide');
}

export function unHideElement(element) {
  removeCss(element, 'hide');
}

export function toggleHideElement(element) {
  element.classList.toggle('hide');
}

export function addCss(element, cssClassName) {
  element.classList.add(cssClassName);
}

export function removeCss(element, cssClassName) {
  element.classList.remove(cssClassName);
}
