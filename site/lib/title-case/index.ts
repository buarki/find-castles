export function toTitleCase(inputString: string) {
  let words = inputString.split(/[\s_]+/);
  let capitalizedWords = words.map(word => {
    return word.charAt(0).toUpperCase() + word.slice(1).toLowerCase();
  });
  return capitalizedWords.join(' ');
}
