export function randomArr<T>(arr: T[]) {
  return arr[Math.floor(Math.random() * arr.length)];
}

export function randomHexColor() {
  return '#' + (((1 << 24) * Math.random()) | 0).toString(16);
}

export function normalizeDate(date: Date) {
  const d = new Date(date);
  const month = d.getMonth() + 1;
  const day = d.getDate();
  const year = d.getFullYear();

  return `${day}/${month}/${year}`;
}
