export function toTitleCase(s: string): string {
  return s.split(" ").map((part: string) => part.replace(
    /\w\S*/g,
    (s) => s.charAt(0).toUpperCase() + s.substring(1).toLowerCase()
  )).join(' ');
}
