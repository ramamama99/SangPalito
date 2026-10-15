/**
 * Helper utilities for date management in Klinik Buah Batu
 * Timezone: Asia/Jakarta (WIB)
 */

export function getJakartaDateParts(offsetDays: number = 0): {
  year: number;
  month: number;
  day: number;
  dateStr: string;
} {
  const now = new Date();
  try {
    const formatter = new Intl.DateTimeFormat("en-US", {
      timeZone: "Asia/Jakarta",
      year: "numeric",
      month: "2-digit",
      day: "2-digit",
    });
    const parts = formatter.formatToParts(now);
    const y = parseInt(parts.find((p) => p.type === "year")?.value || String(now.getFullYear()), 10);
    const m = parseInt(parts.find((p) => p.type === "month")?.value || String(now.getMonth() + 1), 10);
    const d = parseInt(parts.find((p) => p.type === "day")?.value || String(now.getDate()), 10);

    const targetDate = new Date(Date.UTC(y, m - 1, d + offsetDays));
    const tY = targetDate.getUTCFullYear();
    const tM = String(targetDate.getUTCMonth() + 1).padStart(2, "0");
    const tD = String(targetDate.getUTCDate()).padStart(2, "0");
    return {
      year: tY,
      month: parseInt(tM, 10),
      day: parseInt(tD, 10),
      dateStr: `${tY}-${tM}-${tD}`,
    };
  } catch {
    const d = new Date();
    d.setDate(d.getDate() + offsetDays);
    const year = d.getFullYear();
    const month = String(d.getMonth() + 1).padStart(2, "0");
    const day = String(d.getDate()).padStart(2, "0");
    const dateStr = `${year}-${month}-${day}`;
    return { year, month: d.getMonth() + 1, day: d.getDate(), dateStr };
  }
}

export function getRelativeDateStr(offsetDays: number = 0): string {
  return getJakartaDateParts(offsetDays).dateStr;
}

export function formatShortIndoDate(dateStr: string): string {
  if (!dateStr) return "";
  try {
    const [y, m, d] = dateStr.slice(0, 10).split("-").map(Number);
    const obj = new Date(y, m - 1, d);
    return obj.toLocaleDateString("id-ID", { day: "numeric", month: "short" });
  } catch {
    return dateStr;
  }
}

export function formatDayAndDate(dateStr: string): string {
  if (!dateStr) return "";
  try {
    const [y, m, d] = dateStr.slice(0, 10).split("-").map(Number);
    const obj = new Date(y, m - 1, d);
    const dayName = obj.toLocaleDateString("id-ID", { weekday: "short" });
    const dayMonth = obj.toLocaleDateString("id-ID", { day: "numeric", month: "short" });
    return `${dayName}, ${dayMonth}`;
  } catch {
    return dateStr;
  }
}

export function normalizeDateStr(dateStr?: string): string {
  if (!dateStr) return "";
  const s = dateStr.trim();
  if (/^\d{4}-\d{2}-\d{2}/.test(s)) {
    return s.slice(0, 10);
  }
  const m = s.match(/^(\d{1,2})[-/.](\d{1,2})[-/.](\d{4})/);
  if (m) {
    return `${m[3]}-${m[2].padStart(2, "0")}-${m[1].padStart(2, "0")}`;
  }
  return s;
}

export function isSameDate(d1?: string, d2?: string): boolean {
  if (!d1 || !d2) return false;
  return normalizeDateStr(d1) === normalizeDateStr(d2);
}
