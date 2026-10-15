import type { Patient } from "../types";

export function downloadTicketImage(ticket: Patient, qrDataUrl: string) {
  if (!ticket) return;
  const canvas = document.createElement("canvas");
  canvas.width = 720;
  canvas.height = 940;
  const ctx = canvas.getContext("2d");
  if (!ctx) return;

  ctx.fillStyle = "#F3EDE6";
  ctx.fillRect(0, 0, 720, 940);
  ctx.fillStyle = "#195A5F";
  ctx.textAlign = "center";
  ctx.font = "bold 36px sans-serif";
  ctx.fillText("AntriSehat · E-Tiket", 360, 80);
  ctx.font = "bold 96px sans-serif";
  ctx.fillText(ticket.number, 360, 210);

  const img = new Image();
  img.onload = () => {
    // Draw white quiet-zone background behind QR code for maximum contrast and universal optical scanning
    ctx.fillStyle = "#ffffff";
    ctx.fillRect(230, 260, 260, 260);
    ctx.drawImage(img, 240, 270, 240, 240);
    ctx.font = "24px sans-serif";
    [
      ticket.name,
      ticket.poli,
      ticket.date,
      `Status: ${ticket.status}`,
      "Bawa KTP dan kartu BPJS saat kunjungan",
      "Scan QR di lokasi untuk check-in",
    ].forEach((line, i) => ctx.fillText(line, 360, 575 + i * 48));

    const link = document.createElement("a");
    link.download = `AntriSehat-${ticket.number}.png`;
    link.href = canvas.toDataURL();
    link.click();
  };
  img.src = qrDataUrl;
}
