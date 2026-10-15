import type { Patient } from "../types";
import {
  calculateEstimatedArrival,
  calculateRecommendedArrivalTime,
  SHIFT_TIMES,
} from "../constants/doctorShifts";

export function downloadTicketImage(ticket: Patient, qrDataUrl: string) {
  if (!ticket) return;
  const canvas = document.createElement("canvas");
  canvas.width = 720;
  canvas.height = 1040;
  const ctx = canvas.getContext("2d");
  if (!ctx) return;

  // Background
  ctx.fillStyle = "#F8FAFB";
  ctx.fillRect(0, 0, 720, 1040);

  // Border frame
  ctx.strokeStyle = "#12B8B0";
  ctx.lineWidth = 4;
  ctx.strokeRect(20, 20, 680, 1000);

  // Header
  ctx.fillStyle = "#0F766E";
  ctx.textAlign = "center";
  ctx.font = "bold 32px sans-serif";
  ctx.fillText("AntriSehat · E-Tiket Pasien", 360, 75);

  // Queue Number
  ctx.fillStyle = "#0F172A";
  ctx.font = "bold 92px sans-serif";
  ctx.fillText(ticket.number, 360, 185);

  const img = new Image();
  img.onload = () => {
    // White background for QR code
    ctx.fillStyle = "#ffffff";
    ctx.fillRect(230, 220, 260, 260);
    ctx.strokeStyle = "#E2E8F0";
    ctx.lineWidth = 2;
    ctx.strokeRect(230, 220, 260, 260);
    ctx.drawImage(img, 240, 230, 240, 240);

    // Patient Name & Polyclinic
    ctx.fillStyle = "#0F172A";
    ctx.font = "bold 26px sans-serif";
    ctx.fillText(ticket.name, 360, 525);

    const shiftKey = (ticket.shift as "pagi" | "siang" | "sore") || "pagi";
    const shiftInfo = SHIFT_TIMES[shiftKey] || SHIFT_TIMES.pagi;
    const shiftLabel = `Shift ${shiftInfo.label}`;

    ctx.fillStyle = "#475569";
    ctx.font = "20px sans-serif";
    ctx.fillText(`${ticket.poli} · ${shiftLabel}`, 360, 560);
    ctx.fillText(`Tanggal: ${ticket.date} · Status: ${ticket.status}`, 360, 595);

    // Estimasi Kedatangan Box
    const estArrival = ticket.estimatedArrival || calculateEstimatedArrival(shiftKey, 10, 1);
    const recTime = calculateRecommendedArrivalTime(estArrival);

    ctx.fillStyle = "#ECFDF5";
    ctx.fillRect(80, 630, 560, 75);
    ctx.strokeStyle = "#10B981";
    ctx.lineWidth = 2;
    ctx.strokeRect(80, 630, 560, 75);

    ctx.fillStyle = "#065F46";
    ctx.font = "bold 22px sans-serif";
    ctx.fillText(`ESTIMASI KEDATANGAN: ${estArrival}`, 360, 675);

    // Peringatan 15 Menit Lebih Awal Box
    ctx.fillStyle = "#FFFBEB";
    ctx.fillRect(80, 725, 560, 100);
    ctx.strokeStyle = "#F59E0B";
    ctx.lineWidth = 2;
    ctx.strokeRect(80, 725, 560, 100);

    ctx.fillStyle = "#92400E";
    ctx.font = "bold 20px sans-serif";
    ctx.fillText(`⚠️ PERINGATAN KEDATANGAN`, 360, 760);
    ctx.font = "bold 18px sans-serif";
    ctx.fillText(`Wajib hadir pukul ${recTime}`, 360, 790);
    ctx.font = "15px sans-serif";
    ctx.fillText("(15 menit sebelum estimasi untuk verifikasi berkas & tensi)", 360, 812);

    // Footer info
    ctx.fillStyle = "#64748B";
    ctx.font = "16px sans-serif";
    ctx.fillText("Bawa e-KTP fisik dan Kartu BPJS / Asuransi saat kunjungan.", 360, 875);
    ctx.fillText("Tunjukkan QR ini ke scanner loket untuk check-in otomatis.", 360, 905);

    ctx.fillStyle = "#94A3B8";
    ctx.font = "italic 14px sans-serif";
    ctx.fillText("Klinik Buah Batu · Sistem Antrean Cerdas", 360, 955);

    const link = document.createElement("a");
    link.download = `Tiket-Antrean-${ticket.number}.png`;
    link.href = canvas.toDataURL();
    link.click();
  };
  img.src = qrDataUrl;
}
